package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	fpgo "github.com/TeaEntityLab/fpGo/v2"
)

// SimpleHTTP

const (
	// DefaultTimeoutMillisecond Default Request TimeoutMillisecond
	DefaultTimeoutMillisecond = 30 * time.Second
)

// Interceptor Interceptor functions
type Interceptor func(*http.Request) error

// SimpleHTTPDef SimpleHTTP inspired by Retrofits
type SimpleHTTPDef struct {
	client       *http.Client
	interceptors fpgo.StreamDef[*Interceptor]

	TimeoutMillisecond int64

	clientTransport http.RoundTripper
	lastTransport   http.RoundTripper
}

// ResponseWithError Response with Error
type ResponseWithError struct {
	Request  *http.Request
	Response *http.Response

	Err error
}

// NewSimpleHTTP New a SimpleHTTP instance
func NewSimpleHTTP() *SimpleHTTPDef {
	return NewSimpleHTTPWithClientAndInterceptors(&http.Client{})
}

// NewSimpleHTTPWithClientAndInterceptors a SimpleHTTP instance with a given client & interceptor(s)
func NewSimpleHTTPWithClientAndInterceptors(client *http.Client, interceptors ...*Interceptor) *SimpleHTTPDef {
	newOne := &SimpleHTTPDef{
		client:       client,
		interceptors: fpgo.StreamDef[*Interceptor](interceptors),
	}
	newOne.SetHTTPClient(client)
	return newOne
}

// RemoveInterceptor Remove the interceptor(s)
func (simpleHTTPSelf *SimpleHTTPDef) RemoveInterceptor(interceptors ...*Interceptor) {
	for _, interceptor := range interceptors {
		simpleHTTPSelf.interceptors = *simpleHTTPSelf.interceptors.RemoveItem(interceptor)
	}
}

// AddInterceptor Add the interceptor(s)
func (simpleHTTPSelf *SimpleHTTPDef) AddInterceptor(interceptors ...*Interceptor) {
	for _, interceptor := range interceptors {
		simpleHTTPSelf.interceptors = *simpleHTTPSelf.interceptors.Append(interceptor)
	}
}

// ClearInterceptor Clear the interceptor(s)
func (simpleHTTPSelf *SimpleHTTPDef) ClearInterceptor() {
	simpleHTTPSelf.interceptors = fpgo.StreamDef[*Interceptor]{}
}

// GetHTTPClient Get the http client
func (simpleHTTPSelf *SimpleHTTPDef) GetHTTPClient() *http.Client {
	return simpleHTTPSelf.client
}

// SetHTTPClient Get the http client and setup interceptors
func (simpleHTTPSelf *SimpleHTTPDef) SetHTTPClient(client *http.Client) {
	if client.Transport == nil {
		client.Transport = http.DefaultTransport
	}
	// Avoid setting up again next time
	if client.Transport != simpleHTTPSelf.lastTransport {
		// Keep old one
		simpleHTTPSelf.clientTransport = client.Transport

		// Setup myself
		client.Transport = simpleHTTPSelf
		// Avoid setting up again next time
		simpleHTTPSelf.lastTransport = client.Transport
	}

	simpleHTTPSelf.client = client
}

// RoundTrip Do RoundTrip things(interceptors)
func (simpleHTTPSelf *SimpleHTTPDef) RoundTrip(request *http.Request) (*http.Response, error) {
	return simpleHTTPSelf.recursiveVisit(request, 0)
}

func (simpleHTTPSelf *SimpleHTTPDef) recursiveVisit(request *http.Request, index int) (*http.Response, error) {
	if index >= simpleHTTPSelf.interceptors.Len() {
		if simpleHTTPSelf.clientTransport != nil {
			return simpleHTTPSelf.clientTransport.RoundTrip(request)
		}
		// A RoundTripper must not panic; surface a real error.
		return nil, fmt.Errorf("SimpleHTTP has no transport: call SetHTTPClient or NewSimpleHTTP first")
	}

	err := (*simpleHTTPSelf.interceptors[index])(request)
	if err != nil {
		return nil, err
	}
	return simpleHTTPSelf.recursiveVisit(request, index+1)
}

// GetContextTimeout Get Context by TimeoutMillisecond
func (simpleHTTPSelf *SimpleHTTPDef) GetContextTimeout() (context.Context, context.CancelFunc) {
	if simpleHTTPSelf.TimeoutMillisecond > 0 {
		return context.WithTimeout(context.Background(), time.Duration(simpleHTTPSelf.TimeoutMillisecond)*time.Millisecond)
	}

	return context.WithTimeout(context.Background(), DefaultTimeoutMillisecond)
}

// cancelOnCloseBody cancels the request context when the body is closed.
// Cancelling the context as soon as client.Do returns (before the caller has
// read the body) aborts body reads with context.Canceled; cancelling on Close
// bounds the whole request lifetime, matching the standard http.Client
// "caller must close Body" contract.
type cancelOnCloseBody struct {
	io.ReadCloser
	cancel context.CancelFunc
	once   sync.Once
}

func (b *cancelOnCloseBody) Close() error {
	err := b.ReadCloser.Close()
	b.once.Do(b.cancel)
	return err
}

// bindCancelToBody ties context cancellation to the response body lifetime.
// On request failure (no response) the context is cancelled immediately.
func (simpleHTTPSelf *SimpleHTTPDef) bindCancelToBody(result *ResponseWithError, cancel context.CancelFunc) *ResponseWithError {
	if result == nil || result.Response == nil || result.Response.Body == nil {
		cancel()
		return result
	}
	result.Response.Body = &cancelOnCloseBody{ReadCloser: result.Response.Body, cancel: cancel}
	return result
}

// Get HTTP Method Get
func (simpleHTTPSelf *SimpleHTTPDef) Get(givenURL string) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	return simpleHTTPSelf.bindCancelToBody(simpleHTTPSelf.DoNewRequest(ctx, nil, http.MethodGet, givenURL), cancel)
}

// Head HTTP Method Head
func (simpleHTTPSelf *SimpleHTTPDef) Head(givenURL string) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	return simpleHTTPSelf.bindCancelToBody(simpleHTTPSelf.DoNewRequest(ctx, nil, http.MethodHead, givenURL), cancel)
}

// Options HTTP Method Options
func (simpleHTTPSelf *SimpleHTTPDef) Options(givenURL string) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	return simpleHTTPSelf.bindCancelToBody(simpleHTTPSelf.DoNewRequest(ctx, nil, http.MethodOptions, givenURL), cancel)
}

// Delete HTTP Method Delete
func (simpleHTTPSelf *SimpleHTTPDef) Delete(givenURL string) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	return simpleHTTPSelf.bindCancelToBody(simpleHTTPSelf.DoNewRequest(ctx, nil, http.MethodDelete, givenURL), cancel)
}

// Post HTTP Method Post
func (simpleHTTPSelf *SimpleHTTPDef) Post(givenURL, contentType string, body io.Reader) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	return simpleHTTPSelf.bindCancelToBody(simpleHTTPSelf.DoNewRequestWithBodyOptions(ctx, nil, http.MethodPost, givenURL, body, contentType), cancel)
}

// Put HTTP Method Put
func (simpleHTTPSelf *SimpleHTTPDef) Put(givenURL, contentType string, body io.Reader) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	return simpleHTTPSelf.bindCancelToBody(simpleHTTPSelf.DoNewRequestWithBodyOptions(ctx, nil, http.MethodPut, givenURL, body, contentType), cancel)
}

// Patch HTTP Method Patch
func (simpleHTTPSelf *SimpleHTTPDef) Patch(givenURL, contentType string, body io.Reader) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	return simpleHTTPSelf.bindCancelToBody(simpleHTTPSelf.DoNewRequestWithBodyOptions(ctx, nil, http.MethodPatch, givenURL, body, contentType), cancel)
}

// DoNewRequest Do New HTTP Request
func (simpleHTTPSelf *SimpleHTTPDef) DoNewRequest(context context.Context, header http.Header, method string, givenURL string) *ResponseWithError {
	request, newRequestErr := http.NewRequestWithContext(context, method, givenURL, nil)
	if newRequestErr != nil {
		return &ResponseWithError{
			Request: request,
			Err:     newRequestErr,
		}
	}

	if header != nil {
		request.Header = header
	}

	return simpleHTTPSelf.DoRequest(request)
}

// DoNewRequestWithBodyOptions Do New HTTP Request with body options
func (simpleHTTPSelf *SimpleHTTPDef) DoNewRequestWithBodyOptions(context context.Context, header http.Header, method string, givenURL string, body io.Reader, contentType string) *ResponseWithError {
	request, newRequestErr := http.NewRequestWithContext(context, method, givenURL, body)
	if newRequestErr != nil {
		return &ResponseWithError{
			Request: request,
			Err:     newRequestErr,
		}
	}

	if header != nil {
		request.Header = header
	}
	if contentType != "" {
		request.Header.Add("Content-Type", contentType)
	}

	return simpleHTTPSelf.DoRequest(request)
}

// DoRequest Do HTTP Request with interceptors
func (simpleHTTPSelf *SimpleHTTPDef) DoRequest(request *http.Request) *ResponseWithError {
	response, err := simpleHTTPSelf.client.Do(request)

	return &ResponseWithError{
		Request:  request,
		Response: response,
		Err:      err,
	}
}

// // SimpleHTTP SimpleHTTP utils instance
// var SimpleHTTP SimpleHTTPDef

// SimpleAPI

// PathParam Path params for API usages
type PathParam map[string]interface{}

// MultipartForm Path params for API usages
type MultipartForm struct {
	Value map[string][]string
	// File The absolute paths of files
	File map[string][]string
}

// APIResponse Response with Error & Type
type APIResponse[R any] struct {
	ResponseWithError
	TargetObject *R
}

// APINoBody API without request body options
type APINoBody[R any] func(pathParam PathParam, target *R) *fpgo.MonadIODef[*APIResponse[R]]

// APIHasBody API with request body options
type APIHasBody[T any, R any] func(pathParam PathParam, body T, target *R) *fpgo.MonadIODef[*APIResponse[R]]

// APIMultipart API with request body options
type APIMultipart[R any] func(pathParam PathParam, body *MultipartForm, target *R) *fpgo.MonadIODef[*APIResponse[R]]

// // APIResponseOnly API with only response options
// type APIResponseOnly[R any] func(target *R) *fpgo.MonadIODef[*APIResponse[R]]

// BodySerializer Serialize the body (for put/post/patch etc)
type BodySerializer func(body interface{}) (io.Reader, error)

// BodyDeserializer Deserialize the body (for response)
type BodyDeserializer func(body []byte, target interface{}) (interface{}, error)

// MultipartSerializer Serialize the multipart body (for put/post/patch etc)
type MultipartSerializer func(body *MultipartForm) (io.Reader, string, error)

// JSONBodyDeserializer Default JSON Body deserializer
func JSONBodyDeserializer(body []byte, target interface{}) (interface{}, error) {
	err := json.Unmarshal(body, target)
	return target, err
}

// JSONBodySerializer Default JSON Body serializer
func JSONBodySerializer(body interface{}) (io.Reader, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonBytes), err
}

// GeneralMultipartSerializer Default Multipart Body serializer
func GeneralMultipartSerializer(form *MultipartForm) (io.Reader, string, error) {
	if form == nil {
		return nil, "", fmt.Errorf("multipart form is nil")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for fieldName, values := range form.Value {
		for _, value := range values {
			writeFieldErr := writer.WriteField(fieldName, value)
			if writeFieldErr != nil {
				return nil, "", writeFieldErr
			}
		}
	}
	for fieldName, filePaths := range form.File {
		for _, filePath := range filePaths {
			part, createFormFileErr := writer.CreateFormFile(fieldName, filepath.Base(filePath))
			if createFormFileErr != nil {
				return nil, "", createFormFileErr
			}
			file, openFileErr := os.Open(filePath)
			if openFileErr != nil {
				return nil, "", openFileErr
			}

			// Close per file, not deferred: a large form would hold every
			// descriptor open until serialization finished.
			_, copyErr := io.Copy(part, file)
			closeErr := file.Close()
			if copyErr != nil {
				return nil, "", copyErr
			}
			if closeErr != nil {
				return nil, "", closeErr
			}
		}
	}
	writer.Close()

	return body, writer.FormDataContentType(), nil
}

// SimpleAPIDef SimpleAPIDef inspired by Retrofits
type SimpleAPIDef struct {
	simpleHTTP    *SimpleHTTPDef
	BaseURL       string
	DefaultHeader http.Header

	RequestSerializerForMultipart MultipartSerializer
	RequestSerializerForJSON      BodySerializer
	ResponseDeserializer          BodyDeserializer
}

// NewSimpleAPI New a NewSimpleAPI instance
func NewSimpleAPI(baseURL string) *SimpleAPIDef {
	return NewSimpleAPIWithSimpleHTTP(baseURL, NewSimpleHTTP())
}

// NewSimpleAPIWithSimpleHTTP a SimpleAPIDef instance with a SimpleHTTP
func NewSimpleAPIWithSimpleHTTP(baseURL string, simpleHTTP *SimpleHTTPDef) *SimpleAPIDef {
	// url.Parse can fail and return nil; keep the raw base so the request
	// fails with a real error later instead of dereferencing a nil *URL.
	base := baseURL
	if urlInstance, err := url.Parse(baseURL); err == nil && urlInstance != nil {
		base = urlInstance.String()
	}

	return &SimpleAPIDef{
		BaseURL: base,

		RequestSerializerForMultipart: GeneralMultipartSerializer,
		RequestSerializerForJSON:      JSONBodySerializer,
		ResponseDeserializer:          JSONBodyDeserializer,

		simpleHTTP: simpleHTTP,
	}
}

// GetSimpleHTTP Get the SimpleHTTP
func (simpleAPISelf *SimpleAPIDef) GetSimpleHTTP() *SimpleHTTPDef {
	return simpleAPISelf.simpleHTTP
}

// APIMakeGet Make a Get API
func APIMakeGet[R any](simpleAPISelf *SimpleAPIDef, relativeURL string) APINoBody[R] {
	return APIMakeDoNewRequest[R](simpleAPISelf, http.MethodGet, relativeURL)
}

// APIMakeDelete Make a Delete API
func APIMakeDelete[R any](simpleAPISelf *SimpleAPIDef, relativeURL string) APINoBody[R] {
	return APIMakeDoNewRequest[R](simpleAPISelf, http.MethodDelete, relativeURL)
}

// APIMakePostJSONBody Make a Post API with json Body
func APIMakePostJSONBody[T any, R any](simpleAPISelf *SimpleAPIDef, relativeURL string) APIHasBody[T, R] {
	return APIMakeDoNewRequestWithBodySerializer[T, R](simpleAPISelf, http.MethodPost, relativeURL, "application/json", simpleAPISelf.RequestSerializerForJSON)
}

// APIMakePutJSONBody Make a Put API with json Body
func APIMakePutJSONBody[T any, R any](simpleAPISelf *SimpleAPIDef, relativeURL string) APIHasBody[T, R] {
	return APIMakeDoNewRequestWithBodySerializer[T, R](simpleAPISelf, http.MethodPut, relativeURL, "application/json", simpleAPISelf.RequestSerializerForJSON)
}

// APIMakePatchJSONBody Make a Patch API with json Body
func APIMakePatchJSONBody[T any, R any](simpleAPISelf *SimpleAPIDef, relativeURL string) APIHasBody[T, R] {
	return APIMakeDoNewRequestWithBodySerializer[T, R](simpleAPISelf, http.MethodPatch, relativeURL, "application/json", simpleAPISelf.RequestSerializerForJSON)
}

// APIMakePostMultipartBody Make a Post API with multipart Body
func APIMakePostMultipartBody[R any](simpleAPISelf *SimpleAPIDef, relativeURL string) APIMultipart[R] {
	return APIMakeDoNewRequestWithMultipartSerializer[R](simpleAPISelf, http.MethodPost, relativeURL, simpleAPISelf.RequestSerializerForMultipart)
}

// APIMakePutMultipartBody Make a Put API with multipart Body
func APIMakePutMultipartBody[R any](simpleAPISelf *SimpleAPIDef, relativeURL string) APIMultipart[R] {
	return APIMakeDoNewRequestWithMultipartSerializer[R](simpleAPISelf, http.MethodPut, relativeURL, simpleAPISelf.RequestSerializerForMultipart)
}

// APIMakePatchMultipartBody Make a Patch API with multipart Body
func APIMakePatchMultipartBody[R any](simpleAPISelf *SimpleAPIDef, relativeURL string) APIMultipart[R] {
	return APIMakeDoNewRequestWithMultipartSerializer[R](simpleAPISelf, http.MethodPatch, relativeURL, simpleAPISelf.RequestSerializerForMultipart)
}

// APIMakeDoNewRequestWithBodySerializer Make a API with request body options
func APIMakeDoNewRequestWithBodySerializer[T any, R any](simpleAPISelf *SimpleAPIDef, method string, relativeURL string, contentType string, bodySerializer BodySerializer) APIHasBody[T, R] {
	return APIHasBody[T, R](func(pathParam PathParam, body T, target *R) *fpgo.MonadIODef[*APIResponse[R]] {
		return fpgo.MonadIONewGenerics[*APIResponse[R]](func() *APIResponse[R] {
			var bodyReader io.Reader
			if !fpgo.IsNil(body) {
				var newBodyReaderErr error
				bodyReader, newBodyReaderErr = bodySerializer(body)
				if newBodyReaderErr != nil {
					return &APIResponse[R]{
						// Request: request,
						ResponseWithError: ResponseWithError{
							Err: newBodyReaderErr,
						},
					}
				}
			}

			ctx, cancel := simpleAPISelf.GetSimpleHTTP().GetContextTimeout()
			defer cancel()
			response := simpleAPISelf.simpleHTTP.DoNewRequestWithBodyOptions(ctx, simpleAPISelf.DefaultHeader.Clone(), method, simpleAPISelf.replacePathParams(relativeURL, pathParam), bodyReader, contentType)
			if response.Err != nil {
				return &APIResponse[R]{
					ResponseWithError: *response,
				}
			}
			return decodeResponseBody[R](simpleAPISelf, &APIResponse[R]{
				ResponseWithError: *response,
			}, target)
		})
	})
}

// APIMakeDoNewRequestWithMultipartSerializer Make a API with request body options
func APIMakeDoNewRequestWithMultipartSerializer[R any](simpleAPISelf *SimpleAPIDef, method string, relativeURL string, multipartSerializer MultipartSerializer) APIMultipart[R] {
	return APIMultipart[R](func(pathParam PathParam, body *MultipartForm, target *R) *fpgo.MonadIODef[*APIResponse[R]] {
		return fpgo.MonadIONewGenerics[*APIResponse[R]](func() *APIResponse[R] {
			var bodyReader io.Reader
			var contentType string
			if !fpgo.IsNil(body) {
				var newBodyReaderErr error
				bodyReader, contentType, newBodyReaderErr = multipartSerializer(body)
				if newBodyReaderErr != nil {
					return &APIResponse[R]{
						// Request: request,
						ResponseWithError: ResponseWithError{
							Err: newBodyReaderErr,
						},
					}
				}
			}

			ctx, cancel := simpleAPISelf.GetSimpleHTTP().GetContextTimeout()
			defer cancel()
			response := simpleAPISelf.simpleHTTP.DoNewRequestWithBodyOptions(ctx, simpleAPISelf.DefaultHeader.Clone(), method, simpleAPISelf.replacePathParams(relativeURL, pathParam), bodyReader, contentType)
			if response.Err != nil {
				return &APIResponse[R]{
					ResponseWithError: *response,
				}
			}
			return decodeResponseBody[R](simpleAPISelf, &APIResponse[R]{
				ResponseWithError: *response,
			}, target)
		})
	})
}

// APIMakeDoNewRequest Make a API without body options
func APIMakeDoNewRequest[R any](simpleAPISelf *SimpleAPIDef, method string, relativeURL string) APINoBody[R] {
	return APINoBody[R](func(pathParam PathParam, target *R) *fpgo.MonadIODef[*APIResponse[R]] {
		return fpgo.MonadIONewGenerics[*APIResponse[R]](func() *APIResponse[R] {
			ctx, cancel := simpleAPISelf.GetSimpleHTTP().GetContextTimeout()
			defer cancel()
			response := simpleAPISelf.simpleHTTP.DoNewRequest(ctx, simpleAPISelf.DefaultHeader.Clone(), method, simpleAPISelf.replacePathParams(relativeURL, pathParam))
			if response.Err != nil {
				return &APIResponse[R]{
					ResponseWithError: *response,
				}
			}
			return decodeResponseBody[R](simpleAPISelf, &APIResponse[R]{
				ResponseWithError: *response,
			}, target)
		})
	})
}

// decodeResponseBody Close the body and decode it into target. Non-2xx
// responses surface an error carrying the status instead of attempting to
// deserialize an error page into the target type.
func decodeResponseBody[R any](simpleAPISelf *SimpleAPIDef, response *APIResponse[R], target *R) *APIResponse[R] {
	defer response.Response.Body.Close()

	if response.Response.StatusCode >= http.StatusBadRequest {
		// Keep a short slice of the error body. The status alone hides the
		// payload the server sent (validation JSON, plain-text reason).
		snippet, readErr := io.ReadAll(io.LimitReader(response.Response.Body, 512))
		_, _ = io.Copy(io.Discard, response.Response.Body)
		if readErr != nil {
			response.Err = readErr
			return response
		}
		if len(snippet) == 0 {
			response.Err = fmt.Errorf("http error: %s", response.Response.Status)
		} else {
			response.Err = fmt.Errorf("http error: %s: %s", response.Response.Status, snippet)
		}
		return response
	}

	responseBody, readResponseErr := io.ReadAll(response.Response.Body)
	if readResponseErr != nil {
		response.Err = readResponseErr
		return response
	}

	var tempTarget interface{}
	tempTarget, response.Err = simpleAPISelf.ResponseDeserializer(responseBody, target)
	typed, ok := tempTarget.(*R)
	if !ok {
		response.Err = fmt.Errorf("response deserializer returned %T, not *R", tempTarget)
		return response
	}
	response.TargetObject = typed
	return response
}

func (simpleAPISelf *SimpleAPIDef) replacePathParams(relativeURL string, pathParam PathParam) string {
	// Split on the first '?'. PathEscape leaves '&' and '=' alone, which is
	// legal inside a path segment but splits a query. Values in the query use
	// QueryEscape; values in the path use PathEscape.
	pathPart, queryPart, hasQuery := strings.Cut(relativeURL, "?")
	for k, v := range pathParam {
		token := fmt.Sprintf("{%s}", k)
		raw := fmt.Sprintf("%v", v)
		pathPart = strings.ReplaceAll(pathPart, token, url.PathEscape(raw))
		if hasQuery {
			queryPart = strings.ReplaceAll(queryPart, token, url.QueryEscape(raw))
		}
	}
	finalURL := pathPart
	if hasQuery {
		finalURL += "?" + queryPart
	}
	return simpleAPISelf.BaseURL + "/" + finalURL
}
