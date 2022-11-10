// Package openapi3 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package openapi3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// SearchTask request with any body
	SearchTaskWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	SearchTask(ctx context.Context, body SearchTaskJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateTask request with any body
	CreateTaskWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateTask(ctx context.Context, body CreateTaskJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteTask request
	DeleteTask(ctx context.Context, taskId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ReadTask request
	ReadTask(ctx context.Context, taskId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateTask request with any body
	UpdateTaskWithBody(ctx context.Context, taskId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateTask(ctx context.Context, taskId openapi_types.UUID, body UpdateTaskJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) SearchTaskWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSearchTaskRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SearchTask(ctx context.Context, body SearchTaskJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSearchTaskRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateTaskWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateTaskRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateTask(ctx context.Context, body CreateTaskJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateTaskRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteTask(ctx context.Context, taskId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteTaskRequest(c.Server, taskId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ReadTask(ctx context.Context, taskId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewReadTaskRequest(c.Server, taskId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateTaskWithBody(ctx context.Context, taskId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateTaskRequestWithBody(c.Server, taskId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateTask(ctx context.Context, taskId openapi_types.UUID, body UpdateTaskJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateTaskRequest(c.Server, taskId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewSearchTaskRequest calls the generic SearchTask builder with application/json body
func NewSearchTaskRequest(server string, body SearchTaskJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewSearchTaskRequestWithBody(server, "application/json", bodyReader)
}

// NewSearchTaskRequestWithBody generates requests for SearchTask with any type of body
func NewSearchTaskRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/search/tasks")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewCreateTaskRequest calls the generic CreateTask builder with application/json body
func NewCreateTaskRequest(server string, body CreateTaskJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateTaskRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateTaskRequestWithBody generates requests for CreateTask with any type of body
func NewCreateTaskRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/tasks")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteTaskRequest generates requests for DeleteTask
func NewDeleteTaskRequest(server string, taskId openapi_types.UUID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "taskId", runtime.ParamLocationPath, taskId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/tasks/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewReadTaskRequest generates requests for ReadTask
func NewReadTaskRequest(server string, taskId openapi_types.UUID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "taskId", runtime.ParamLocationPath, taskId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/tasks/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewUpdateTaskRequest calls the generic UpdateTask builder with application/json body
func NewUpdateTaskRequest(server string, taskId openapi_types.UUID, body UpdateTaskJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateTaskRequestWithBody(server, taskId, "application/json", bodyReader)
}

// NewUpdateTaskRequestWithBody generates requests for UpdateTask with any type of body
func NewUpdateTaskRequestWithBody(server string, taskId openapi_types.UUID, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "taskId", runtime.ParamLocationPath, taskId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/tasks/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// SearchTask request with any body
	SearchTaskWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SearchTaskResponse, error)

	SearchTaskWithResponse(ctx context.Context, body SearchTaskJSONRequestBody, reqEditors ...RequestEditorFn) (*SearchTaskResponse, error)

	// CreateTask request with any body
	CreateTaskWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTaskResponse, error)

	CreateTaskWithResponse(ctx context.Context, body CreateTaskJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTaskResponse, error)

	// DeleteTask request
	DeleteTaskWithResponse(ctx context.Context, taskId openapi_types.UUID, reqEditors ...RequestEditorFn) (*DeleteTaskResponse, error)

	// ReadTask request
	ReadTaskWithResponse(ctx context.Context, taskId openapi_types.UUID, reqEditors ...RequestEditorFn) (*ReadTaskResponse, error)

	// UpdateTask request with any body
	UpdateTaskWithBodyWithResponse(ctx context.Context, taskId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTaskResponse, error)

	UpdateTaskWithResponse(ctx context.Context, taskId openapi_types.UUID, body UpdateTaskJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTaskResponse, error)
}

type SearchTaskResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Tasks *[]Task `json:"tasks,omitempty"`
		Total *int64  `json:"total,omitempty"`
	}
	JSON400 *struct {
		Error *string `json:"error,omitempty"`
	}
	JSON500 *struct {
		Error *string `json:"error,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r SearchTaskResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SearchTaskResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateTaskResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *struct {
		Task *Task `json:"task,omitempty"`
	}
	JSON400 *struct {
		Error *string `json:"error,omitempty"`
	}
	JSON500 *struct {
		Error *string `json:"error,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r CreateTaskResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateTaskResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteTaskResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *struct {
		Error *string `json:"error,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r DeleteTaskResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteTaskResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ReadTaskResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Task *Task `json:"task,omitempty"`
	}
	JSON500 *struct {
		Error *string `json:"error,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r ReadTaskResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ReadTaskResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateTaskResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *struct {
		Error *string `json:"error,omitempty"`
	}
	JSON500 *struct {
		Error *string `json:"error,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r UpdateTaskResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateTaskResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// SearchTaskWithBodyWithResponse request with arbitrary body returning *SearchTaskResponse
func (c *ClientWithResponses) SearchTaskWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SearchTaskResponse, error) {
	rsp, err := c.SearchTaskWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSearchTaskResponse(rsp)
}

func (c *ClientWithResponses) SearchTaskWithResponse(ctx context.Context, body SearchTaskJSONRequestBody, reqEditors ...RequestEditorFn) (*SearchTaskResponse, error) {
	rsp, err := c.SearchTask(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSearchTaskResponse(rsp)
}

// CreateTaskWithBodyWithResponse request with arbitrary body returning *CreateTaskResponse
func (c *ClientWithResponses) CreateTaskWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTaskResponse, error) {
	rsp, err := c.CreateTaskWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTaskResponse(rsp)
}

func (c *ClientWithResponses) CreateTaskWithResponse(ctx context.Context, body CreateTaskJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTaskResponse, error) {
	rsp, err := c.CreateTask(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTaskResponse(rsp)
}

// DeleteTaskWithResponse request returning *DeleteTaskResponse
func (c *ClientWithResponses) DeleteTaskWithResponse(ctx context.Context, taskId openapi_types.UUID, reqEditors ...RequestEditorFn) (*DeleteTaskResponse, error) {
	rsp, err := c.DeleteTask(ctx, taskId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteTaskResponse(rsp)
}

// ReadTaskWithResponse request returning *ReadTaskResponse
func (c *ClientWithResponses) ReadTaskWithResponse(ctx context.Context, taskId openapi_types.UUID, reqEditors ...RequestEditorFn) (*ReadTaskResponse, error) {
	rsp, err := c.ReadTask(ctx, taskId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseReadTaskResponse(rsp)
}

// UpdateTaskWithBodyWithResponse request with arbitrary body returning *UpdateTaskResponse
func (c *ClientWithResponses) UpdateTaskWithBodyWithResponse(ctx context.Context, taskId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTaskResponse, error) {
	rsp, err := c.UpdateTaskWithBody(ctx, taskId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTaskResponse(rsp)
}

func (c *ClientWithResponses) UpdateTaskWithResponse(ctx context.Context, taskId openapi_types.UUID, body UpdateTaskJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTaskResponse, error) {
	rsp, err := c.UpdateTask(ctx, taskId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTaskResponse(rsp)
}

// ParseSearchTaskResponse parses an HTTP response from a SearchTaskWithResponse call
func ParseSearchTaskResponse(rsp *http.Response) (*SearchTaskResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SearchTaskResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Tasks *[]Task `json:"tasks,omitempty"`
			Total *int64  `json:"total,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Error *string `json:"error,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest struct {
			Error *string `json:"error,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseCreateTaskResponse parses an HTTP response from a CreateTaskWithResponse call
func ParseCreateTaskResponse(rsp *http.Response) (*CreateTaskResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateTaskResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest struct {
			Task *Task `json:"task,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Error *string `json:"error,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest struct {
			Error *string `json:"error,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseDeleteTaskResponse parses an HTTP response from a DeleteTaskWithResponse call
func ParseDeleteTaskResponse(rsp *http.Response) (*DeleteTaskResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteTaskResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest struct {
			Error *string `json:"error,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseReadTaskResponse parses an HTTP response from a ReadTaskWithResponse call
func ParseReadTaskResponse(rsp *http.Response) (*ReadTaskResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ReadTaskResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Task *Task `json:"task,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest struct {
			Error *string `json:"error,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseUpdateTaskResponse parses an HTTP response from a UpdateTaskWithResponse call
func ParseUpdateTaskResponse(rsp *http.Response) (*UpdateTaskResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateTaskResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Error *string `json:"error,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest struct {
			Error *string `json:"error,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}
