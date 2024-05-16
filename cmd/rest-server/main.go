package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	esv7 "github.com/elastic/go-elasticsearch/v7"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	rv8 "github.com/go-redis/redis/v8"

  
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riandyrn/otelchi"
	"go.uber.org/zap"

	"github.com/MarioCarrion/todo-api/cmd/internal"
	internaldomain "github.com/MarioCarrion/todo-api/internal"
	"github.com/MarioCarrion/todo-api/internal/elasticsearch"
	"github.com/MarioCarrion/todo-api/internal/envvar"
	"github.com/MarioCarrion/todo-api/internal/memcached"
	"github.com/MarioCarrion/todo-api/internal/postgresql"
	"github.com/MarioCarrion/todo-api/internal/redis"
	"github.com/MarioCarrion/todo-api/internal/rest"
	"github.com/MarioCarrion/todo-api/internal/service"
)

//go:embed static
var content embed.FS

func main() {
	var env, address string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.StringVar(&address, "address", ":9234", "HTTP Server Address")
	flag.Parse()

	errC, err := run(env, address)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run(env, address string) (<-chan error, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "zap.NewProduction")
	}

	if err := envvar.Load(env); err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "envvar.Load")
	}

	vault, err := internal.NewVaultProvider()
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewVaultProvider")
	}

	conf := envvar.New(vault)

	//-

	pool, err := internal.NewPostgreSQL(conf)
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewPostgreSQL")
	}

	esClient, err := internal.NewElasticSearch(conf)
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewElasticSearch")
	}

	memcached, err := internal.NewMemcached(conf)
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewMemcached")
	}

	// rmq, err := internal.NewRabbitMQ(conf)
	// if err != nil {
	// 	return nil, fmt.Errorf("internal.NewRabbitMQ %w", err)
	// }

	// kafka, err := internal.NewKafkaProducer(conf)
	// if err != nil {
	// 	return nil, fmt.Errorf("internal.NewKafka %w", err)
	// }

	rdb, err := internal.NewRedis(conf)
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewRedis")
	}

	//-

	promExporter, err := internal.NewOTExporter(conf)
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewOTExporter")
	}


	logging := func(c *gin.Context) {
		logger.Info(c.Request.Method,
			zap.Time("time", time.Now()),
			zap.String("url", c.Request.URL.String()),
		)
	}

	//-

	srv, err := newServer(serverConfig{
		Address:       address,
		DB:            pool,
		ElasticSearch: esClient,
		Metrics:       promExporter,


		Middlewares:   []func(next http.Handler) http.Handler{otelchi.Middleware("todo-api-server"), logging},
		Redis:         rdb,
		Logger:        logger,
		Memcached:     memcached,
		// RabbitMQ:      rmq,
		// Kafka:         kafka,
	})
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "newServer")
	}

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			_ = logger.Sync()

			pool.Close()
			// rmq.Close()
			rdb.Close()
			stop()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil { //nolint: contextcheck
			errC <- err
		}

		logger.Info("Shutdown completed")
	}()

	go func() {
		logger.Info("Listening and serving", zap.String("address", address))

		// "ListenAndServe always returns a non-nil error. After Shutdown or Close, the returned error is
		// ErrServerClosed."
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errC <- err
		}
	}()

	return errC, nil
}

type serverConfig struct {
	Address       string
	DB            *pgxpool.Pool
	ElasticSearch *esv7.Client
	Kafka         *internal.KafkaProducer
	RabbitMQ      *internal.RabbitMQ
	Redis         *rv8.Client
	Memcached     *memcache.Client
	Metrics       http.Handler


	Middlewares   []func(next http.Handler) http.Handler
	Logger        *zap.Logger
}

func newServer(conf serverConfig) (*http.Server, error) {


	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))

	for _, mw := range conf.Middlewares {
		router.Use(mw)
	}

	//-

	repo := postgresql.NewTask(conf.DB)
	mrepo := memcached.NewTask(conf.Memcached, repo, conf.Logger)

	search := elasticsearch.NewTask(conf.ElasticSearch)
	msearch := memcached.NewSearchableTask(conf.Memcached, search)

	// XXX mclient := memcached.NewSearchableTask(conf.Memcached, search, conf.Logger)
	// msgBroker, err := rabbitmq.NewTask(conf.RabbitMQ.Channel)
	// if err != nil {
	// 	return nil, fmt.Errorf("rabbitmq.NewTask %w", err)
	// }

	// msgBroker := kafka.NewTask(conf.Kafka.Producer, conf.Kafka.Topic)

	msgBroker := redis.NewTask(conf.Redis)

	svc := service.NewTask(conf.Logger, mrepo, msearch, msgBroker)

	rest.RegisterOpenAPI(router)
	rest.NewTaskHandler(svc).Register(router)

	//-

	fsys, _ := fs.Sub(content, "static")


	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(fsys))))

	router.GET("/metrics", gin.WrapH(conf.Metrics))

	//-

	lmt := tollbooth.NewLimiter(3, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Second})

	lmtmw := tollbooth.LimitHandler(lmt, router)

	//-

	return &http.Server{
		Handler:           lmtmw,
		Addr:              conf.Address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}, nil
}
