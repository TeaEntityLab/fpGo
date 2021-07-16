package fpgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// SimpleHTTP

// Interceptor Interceptor functions
type Interceptor func(*http.Request) (*http.Response, error)

// SimpleHTTPDef SimpleHTTP inspired by Retrofits
type SimpleHTTPDef struct {
	client       *http.Client
	interceptors StreamDef

	clientTransport http.RoundTripper
	lastTransport   http.RoundTripper
}

// ResponseWithError Response with Error
type ResponseWithError struct {
	TargetObject interface{}
	Response     *http.Response
	Err          error
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

// Get HTTP Get
func (simpleHTTPSelf *SimpleHTTPDef) Get(theURL string) *ResponseWithError {
	response, err := simpleHTTPSelf.client.Get(theURL)

	return &ResponseWithError{
		Response: response,
		Err:      err,
	}
}

// Post HTTP Post
func (simpleHTTPSelf *SimpleHTTPDef) Post(theURL, contentType string, body io.Reader) *ResponseWithError {
	response, err := simpleHTTPSelf.client.Post(theURL, contentType, body)

	return &ResponseWithError{
		Response: response,
		Err:      err,
	}
}

// // SimpleHTTP SimpleHTTP utils instance
// var SimpleHTTP SimpleHTTPDef

// SimpleAPI

type apiNoBody func(pathParam map[string]interface{}, target interface{}) *MonadIODef
type apiHasBody func(pathParam map[string]interface{}, body interface{}, target interface{}) *MonadIODef

// APIGet API Get function definition
type APIGet apiNoBody

// APIOption API Option function definition
type APIOption apiNoBody

// APIDelete API Delete function definition
type APIDelete apiNoBody

// APIPost API Post function definition
type APIPost apiHasBody

// APIPut API Put function definition
type APIPut apiHasBody

// BodySerializer Serialize the body (for put/post etc)
type BodySerializer func(body interface{}) (io.Reader, error)

// BodyDeserializer Deserialize the body (for response)
type BodyDeserializer func(body []byte, target interface{}) (interface{}, error)

// JSONBodyDeserializer Default JSON Body deserializer
func JSONBodyDeserializer(body []byte, target interface{}) (interface{}, error) {
	err := json.Unmarshal(body, target)
	return target, err
}

// SimpleAPIDef SimpleAPIDef inspired by Retrofits
type SimpleAPIDef struct {
	simpleHTTP *SimpleHTTPDef
	BaseURL    string

	ResponseDeserializer BodyDeserializer
}

// NewSimpleAPI New a NewSimpleAPI instance
func NewSimpleAPI(baseURL string) *SimpleAPIDef {
	return NewSimpleAPIWithSimpleHTTP(baseURL, NewSimpleHTTP())
}

// NewSimpleAPIWithSimpleHTTP a SimpleAPIDef instance with a SimpleHTTP
func NewSimpleAPIWithSimpleHTTP(baseURL string, simpleHTTP *SimpleHTTPDef) *SimpleAPIDef {
	theURL, _ := url.Parse(baseURL)

	return &SimpleAPIDef{
		BaseURL:              theURL.String(),
		ResponseDeserializer: JSONBodyDeserializer,

		simpleHTTP: simpleHTTP,
	}
}

// GetSimpleHTTP Get the SimpleHTTP
func (simpleAPISelf *SimpleAPIDef) GetSimpleHTTP() *SimpleHTTPDef {
	return simpleAPISelf.simpleHTTP
}

// MakeGet Make a Get API
func (simpleAPISelf *SimpleAPIDef) MakeGet(theURL string) APIGet {
	return APIGet(func(pathParam map[string]interface{}, target interface{}) *MonadIODef {
		return MonadIO.New(func() interface{} {

			response := simpleAPISelf.simpleHTTP.Get(simpleAPISelf.replacePathParams(theURL, pathParam))
			if response.Err != nil {
				return ResponseWithError{
					Err: response.Err,
				}
			}
			responseBody, readResponseErr := ioutil.ReadAll(response.Response.Body)
			if readResponseErr != nil {
				return ResponseWithError{
					Err: readResponseErr,
				}
			}
			response.TargetObject, response.Err = simpleAPISelf.ResponseDeserializer(responseBody, target)
			return response
		})
	})
}

// MakePostJSONBody Make a Post API with json Body
func (simpleAPISelf *SimpleAPIDef) MakePostJSONBody(theURL string) APIPost {
	return simpleAPISelf.MakePostWithBodySerializer(theURL, "application/json", BodySerializer(func(body interface{}) (io.Reader, error) {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(jsonBytes), err
	}))
}

// MakePostWithBodySerializer Make a Post API with Body serializer
func (simpleAPISelf *SimpleAPIDef) MakePostWithBodySerializer(theURL string, contentType string, bodySerializer BodySerializer) APIPost {
	return APIPost(func(pathParam map[string]interface{}, body interface{}, target interface{}) *MonadIODef {
		return MonadIO.New(func() interface{} {
			bodyReader, err := bodySerializer(body)
			if err != nil {
				return ResponseWithError{
					Err: err,
				}
			}

			response := simpleAPISelf.simpleHTTP.Post(simpleAPISelf.replacePathParams(theURL, pathParam), contentType, bodyReader)
			if response.Err != nil {
				return ResponseWithError{
					Err: response.Err,
				}
			}
			responseBody, readResponseErr := ioutil.ReadAll(response.Response.Body)
			if readResponseErr != nil {
				return ResponseWithError{
					Err: readResponseErr,
				}
			}
			response.TargetObject, response.Err = simpleAPISelf.ResponseDeserializer(responseBody, target)
			return response
		})
	})
}

func (simpleAPISelf *SimpleAPIDef) replacePathParams(theURL string, pathParam map[string]interface{}) string {
	finalURL := theURL
	for k, v := range pathParam {
		finalURL = strings.ReplaceAll(theURL, fmt.Sprintf("{%s}", k), fmt.Sprintf("%v", v))
	}
	return simpleAPISelf.BaseURL + "/" + finalURL
}
