package fpgo

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
	"time"
)

// SimpleHTTP

const (
	// DefaultTimeoutMillisecond Default Request TimeoutMillisecond
	DefaultTimeoutMillisecond = 30 * time.Second
)

// Interceptor Interceptor functions
type Interceptor func(*http.Request) (*http.Response, error)

// SimpleHTTPDef SimpleHTTP inspired by Retrofits
type SimpleHTTPDef struct {
	client       *http.Client
	interceptors StreamDef

	TimeoutMillisecond int64

	clientTransport http.RoundTripper
	lastTransport   http.RoundTripper
}

// ResponseWithError Response with Error
type ResponseWithError struct {
	TargetObject interface{}
	Request      *http.Request
	Response     *http.Response

	Err error
}

// NewSimpleHTTP New a SimpleHTTP instance
func NewSimpleHTTP() *SimpleHTTPDef {
	return NewSimpleHTTPWithClientAndInterceptors(&http.Client{})
}

// NewSimpleHTTPWithClientAndInterceptors a SimpleHTTP instance with a given client & interceptor(s)
func NewSimpleHTTPWithClientAndInterceptors(client *http.Client, interceptors ...Interceptor) *SimpleHTTPDef {
	interceptorsInternal := make([]interface{}, len(interceptors))
	for i, interceptor := range interceptors {
		interceptorsInternal[i] = interceptor
	}
	newOne := &SimpleHTTPDef{
		client:       client,
		interceptors: StreamDef(interceptorsInternal),
	}
	newOne.SetHTTPClient(client)
	return newOne
}

// RemoveInterceptor Remove the interceptor(s)
func (simpleHTTPSelf *SimpleHTTPDef) RemoveInterceptor(interceptors ...Interceptor) {
	for _, interceptor := range interceptors {
		simpleHTTPSelf.interceptors.RemoveItem(interceptor)
	}
}

// AddInterceptor Add the interceptor(s)
func (simpleHTTPSelf *SimpleHTTPDef) AddInterceptor(interceptors ...Interceptor) {
	for _, interceptor := range interceptors {
		simpleHTTPSelf.interceptors.Append(interceptor)
	}
}

// ClearInterceptor Clear the interceptor(s)
func (simpleHTTPSelf *SimpleHTTPDef) ClearInterceptor() {
	simpleHTTPSelf.interceptors = StreamDef{}
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

	simpleHTTPSelf.interceptors[index].(Interceptor)(request)
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

	response, err := simpleHTTPSelf.client.Do(request)

	return &ResponseWithError{
		Request:  request,
		Response: response,
		Err:      err,
	}
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

// APINoBody API without request body options
type APINoBody func(pathParam map[string]interface{}, target interface{}) *MonadIODef

// APIHasBody API with request body options
type APIHasBody func(pathParam map[string]interface{}, body interface{}, target interface{}) *MonadIODef

// BodySerializer Serialize the body (for put/post etc)
type BodySerializer func(body interface{}) (io.Reader, error)

// BodyDeserializer Deserialize the body (for response)
type BodyDeserializer func(body []byte, target interface{}) (interface{}, error)

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

// SimpleAPIDef SimpleAPIDef inspired by Retrofits
type SimpleAPIDef struct {
	simpleHTTP    *SimpleHTTPDef
	BaseURL       string
	DefaultHeader http.Header

	RequestSerializerForJSON BodySerializer
	ResponseDeserializer     BodyDeserializer
}

// NewSimpleAPI New a NewSimpleAPI instance
func NewSimpleAPI(baseURL string) *SimpleAPIDef {
	return NewSimpleAPIWithSimpleHTTP(baseURL, NewSimpleHTTP())
}

// NewSimpleAPIWithSimpleHTTP a SimpleAPIDef instance with a SimpleHTTP
func NewSimpleAPIWithSimpleHTTP(baseURL string, simpleHTTP *SimpleHTTPDef) *SimpleAPIDef {
	urlInstance, _ := url.Parse(baseURL)

	return &SimpleAPIDef{
		BaseURL:              urlInstance.String(),
		ResponseDeserializer: JSONBodyDeserializer,

		simpleHTTP: simpleHTTP,
	}
}

// GetSimpleHTTP Get the SimpleHTTP
func (simpleAPISelf *SimpleAPIDef) GetSimpleHTTP() *SimpleHTTPDef {
	return simpleAPISelf.simpleHTTP
}

// MakeGet Make a Get API
func (simpleAPISelf *SimpleAPIDef) MakeGet(relativeURL string) APINoBody {
	return simpleAPISelf.MakeDoNewRequest(http.MethodGet, relativeURL)
}

// MakeDelete Make a Delete API
func (simpleAPISelf *SimpleAPIDef) MakeDelete(relativeURL string) APINoBody {
	return simpleAPISelf.MakeDoNewRequest(http.MethodDelete, relativeURL)
}

// MakePostJSONBody Make a Post API with json Body
func (simpleAPISelf *SimpleAPIDef) MakePostJSONBody(relativeURL string) APIHasBody {
	return simpleAPISelf.MakeDoNewRequestWithBodyOptions(http.MethodPost, relativeURL, "application/json", simpleAPISelf.RequestSerializerForJSON)
}

// MakePutJSONBody Make a Put API with json Body
func (simpleAPISelf *SimpleAPIDef) MakePutJSONBody(relativeURL string) APIHasBody {
	return simpleAPISelf.MakeDoNewRequestWithBodyOptions(http.MethodPost, relativeURL, "application/json", simpleAPISelf.RequestSerializerForJSON)
}

// MakePatchJSONBody Make a Patch API with json Body
func (simpleAPISelf *SimpleAPIDef) MakePatchJSONBody(relativeURL string) APIHasBody {
	return simpleAPISelf.MakeDoNewRequestWithBodyOptions(http.MethodPost, relativeURL, "application/json", simpleAPISelf.RequestSerializerForJSON)
}

// MakeDoNewRequestWithBodyOptions Make a API with request body options
func (simpleAPISelf *SimpleAPIDef) MakeDoNewRequestWithBodyOptions(method string, relativeURL string, contentType string, bodySerializer BodySerializer) APIHasBody {
	return APIHasBody(func(pathParam map[string]interface{}, body interface{}, target interface{}) *MonadIODef {
		return MonadIO.New(func() interface{} {
			var bodyReader io.Reader
			if !IsNil(body) {
				var newBodyReaderErr error
				bodyReader, newBodyReaderErr = bodySerializer(body)
				if newBodyReaderErr != nil {
					return &ResponseWithError{
						// Request: request,
						Err: newBodyReaderErr,
					}
				}
			}

			ctx, cancel := simpleAPISelf.GetSimpleHTTP().GetContextTimeout()
			defer cancel()
			response := simpleAPISelf.simpleHTTP.DoNewRequestWithBodyOptions(ctx, simpleAPISelf.DefaultHeader, method, simpleAPISelf.replacePathParams(relativeURL, pathParam), bodyReader, contentType)
			if response.Err != nil {
				return response
			}
			responseBody, readResponseErr := ioutil.ReadAll(response.Response.Body)
			if readResponseErr != nil {
				return &ResponseWithError{
					Request:  response.Request,
					Response: response.Response,
					Err:      readResponseErr,
				}
			}
			response.TargetObject, response.Err = simpleAPISelf.ResponseDeserializer(responseBody, target)
			return response
		})
	})
}

// MakeDoNewRequest Make a API without body options
func (simpleAPISelf *SimpleAPIDef) MakeDoNewRequest(method string, relativeURL string) APINoBody {
	return APINoBody(func(pathParam map[string]interface{}, target interface{}) *MonadIODef {
		return MonadIO.New(func() interface{} {

			ctx, cancel := simpleAPISelf.GetSimpleHTTP().GetContextTimeout()
			defer cancel()
			response := simpleAPISelf.simpleHTTP.DoNewRequest(ctx, simpleAPISelf.DefaultHeader, method, simpleAPISelf.replacePathParams(relativeURL, pathParam))
			if response.Err != nil {
				return response
			}
			responseBody, readResponseErr := ioutil.ReadAll(response.Response.Body)
			if readResponseErr != nil {
				return &ResponseWithError{
					Request:  response.Request,
					Response: response.Response,
					Err:      readResponseErr,
				}
			}
			response.TargetObject, response.Err = simpleAPISelf.ResponseDeserializer(responseBody, target)
			return response
		})
	})
}

func (simpleAPISelf *SimpleAPIDef) replacePathParams(relativeURL string, pathParam map[string]interface{}) string {
	finalURL := relativeURL
	for k, v := range pathParam {
		finalURL = strings.ReplaceAll(relativeURL, fmt.Sprintf("{%s}", k), fmt.Sprintf("%v", v))
	}
	return simpleAPISelf.BaseURL + "/" + finalURL
}
