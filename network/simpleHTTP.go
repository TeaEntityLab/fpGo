package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	fpgo "github.com/TeaEntityLab/fpGo"
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
	if index >= simpleHTTPSelf.interceptors.Len() && simpleHTTPSelf.clientTransport != nil {
		return simpleHTTPSelf.clientTransport.RoundTrip(request)
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
		return context.WithTimeout(context.Background(), time.Duration(simpleHTTPSelf.TimeoutMillisecond))
	}

	return context.WithTimeout(context.Background(), DefaultTimeoutMillisecond)
}

// Get HTTP Method Get
func (simpleHTTPSelf *SimpleHTTPDef) Get(givenURL string) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	defer cancel()
	return simpleHTTPSelf.DoNewRequest(ctx, nil, http.MethodGet, givenURL)
}

// Head HTTP Method Head
func (simpleHTTPSelf *SimpleHTTPDef) Head(givenURL string) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	defer cancel()
	return simpleHTTPSelf.DoNewRequest(ctx, nil, http.MethodHead, givenURL)
}

// Options HTTP Method Options
func (simpleHTTPSelf *SimpleHTTPDef) Options(givenURL string) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	defer cancel()
	return simpleHTTPSelf.DoNewRequest(ctx, nil, http.MethodOptions, givenURL)
}

// Delete HTTP Method Delete
func (simpleHTTPSelf *SimpleHTTPDef) Delete(givenURL string) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	defer cancel()
	return simpleHTTPSelf.DoNewRequest(ctx, nil, http.MethodDelete, givenURL)
}

// Post HTTP Method Post
func (simpleHTTPSelf *SimpleHTTPDef) Post(givenURL, contentType string, body io.Reader) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	defer cancel()
	return simpleHTTPSelf.DoNewRequestWithBodyOptions(ctx, nil, http.MethodPost, givenURL, body, contentType)
}

// Put HTTP Method Put
func (simpleHTTPSelf *SimpleHTTPDef) Put(givenURL, contentType string, body io.Reader) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	defer cancel()
	return simpleHTTPSelf.DoNewRequestWithBodyOptions(ctx, nil, http.MethodPut, givenURL, body, contentType)
}

// Patch HTTP Method Patch
func (simpleHTTPSelf *SimpleHTTPDef) Patch(givenURL, contentType string, body io.Reader) *ResponseWithError {
	ctx, cancel := simpleHTTPSelf.GetContextTimeout()
	defer cancel()
	return simpleHTTPSelf.DoNewRequestWithBodyOptions(ctx, nil, http.MethodPatch, givenURL, body, contentType)
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
			defer file.Close()

			io.Copy(part, file)
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
	urlInstance, _ := url.Parse(baseURL)

	return &SimpleAPIDef{
		BaseURL: urlInstance.String(),

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
	return APIMakeDoNewRequestWithBodySerializer[T, R](simpleAPISelf, http.MethodPost, relativeURL, "application/json", simpleAPISelf.RequestSerializerForJSON)
}

// APIMakePatchJSONBody Make a Patch API with json Body
func APIMakePatchJSONBody[T any, R any](simpleAPISelf *SimpleAPIDef, relativeURL string) APIHasBody[T, R] {
	return APIMakeDoNewRequestWithBodySerializer[T, R](simpleAPISelf, http.MethodPost, relativeURL, "application/json", simpleAPISelf.RequestSerializerForJSON)
}

// APIMakePostMultipartBody Make a Post API with multipart Body
func APIMakePostMultipartBody[R any](simpleAPISelf *SimpleAPIDef, relativeURL string) APIMultipart[R] {
	return APIMakeDoNewRequestWithMultipartSerializer[R](simpleAPISelf, http.MethodPost, relativeURL, simpleAPISelf.RequestSerializerForMultipart)
}

// APIMakePutMultipartBody Make a Put API with multipart Body
func APIMakePutMultipartBody[R any](simpleAPISelf *SimpleAPIDef, relativeURL string) APIMultipart[R] {
	return APIMakeDoNewRequestWithMultipartSerializer[R](simpleAPISelf, http.MethodPost, relativeURL, simpleAPISelf.RequestSerializerForMultipart)
}

// APIMakePatchMultipartBody Make a Patch API with multipart Body
func APIMakePatchMultipartBody[R any](simpleAPISelf *SimpleAPIDef, relativeURL string) APIMultipart[R] {
	return APIMakeDoNewRequestWithMultipartSerializer[R](simpleAPISelf, http.MethodPost, relativeURL, simpleAPISelf.RequestSerializerForMultipart)
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

// decodeResponseBody Make a API without body options
func decodeResponseBody[R any](simpleAPISelf *SimpleAPIDef, response *APIResponse[R], target *R) *APIResponse[R] {
	responseBody, readResponseErr := ioutil.ReadAll(response.Response.Body)
	if readResponseErr != nil {
		response.Err = readResponseErr
		return response
	}

	var tempTarget interface{}
	tempTarget, response.Err = simpleAPISelf.ResponseDeserializer(responseBody, target)
	response.TargetObject = tempTarget.(*R)
	return response
}

func (simpleAPISelf *SimpleAPIDef) replacePathParams(relativeURL string, pathParam PathParam) string {
	finalURL := relativeURL
	for k, v := range pathParam {
		finalURL = strings.ReplaceAll(relativeURL, fmt.Sprintf("{%s}", k), fmt.Sprintf("%v", v))
	}
	return simpleAPISelf.BaseURL + "/" + finalURL
}
