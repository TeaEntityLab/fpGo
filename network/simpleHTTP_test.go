package network

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type errReadCloser struct{}

func (e errReadCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New("read failed")
}
func (e errReadCloser) Close() error { return nil }

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type PostListResponse struct {
	Data []Post `json:"data"`
}

func TestSimpleAPI(t *testing.T) {
	var actualPath string
	var actualRequest *http.Request
	var actualRequestBody []byte
	var actualContentType string

	postsHandler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		actualRequestBody, _ = ioutil.ReadAll(req.Body)

		// auth := req.Header.Get("Auth")
		_, err := writer.Write([]byte(`
{
	"data": [
	  {
	    "userId": 1,
	    "id": 1,
	    "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
	    "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
	  },
	  {
	    "userId": 1,
	    "id": 2,
	    "title": "qui est esse",
	    "body": "est rerum tempore vitae\nsequi sint nihil reprehenderit dolor beatae ea dolores neque\nfugiat blanditiis voluptate porro vel nihil molestiae ut reiciendis\nqui aperiam non debitis possimus qui neque nisi nulla"
	  }
	]
}
			`))
		assert.NoError(t, err)
	})

	server := httptest.NewServer(postsHandler)
	defer server.Close()
	// router := httprouter.New()
	// router.GET("/posts", postsHandler)
	// recorder := httptest.NewRecorder()

	var response *ResponseWithError

	client := NewSimpleHTTP()

	interceptorForTest := Interceptor(func(request *http.Request) error {
		actualPath = request.URL.Path
		actualRequest = request
		actualContentType = actualRequest.Header.Get("Content-Type")
		return nil
	})
	client.AddInterceptor(&interceptorForTest)

	response = client.Get(server.URL + "/posts")
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "/posts", actualPath)

	response = client.Options(server.URL + "/posts")
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "/posts", actualPath)

	response = client.Head(server.URL + "/posts")
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "/posts", actualPath)

	response = client.Delete(server.URL + "/posts/1")
	assert.Equal(t, "/posts/1", actualPath)
	assert.Equal(t, nil, response.Err)

	actualContentType = ""
	response = client.Post(server.URL+"/posts", "application/json", bytes.NewReader([]byte(`{"userId":0,"id":5,"title":"aa","body":""}`)))
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "application/json", actualContentType)
	assert.Equal(t, `{"userId":0,"id":5,"title":"aa","body":""}`, string(actualRequestBody))

	actualContentType = ""
	response = client.Put(server.URL+"/posts", "application/json", bytes.NewReader([]byte(`{"userId":0,"id":4,"title":"bb","body":""}`)))
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "application/json", actualContentType)
	assert.Equal(t, `{"userId":0,"id":4,"title":"bb","body":""}`, string(actualRequestBody))

	actualContentType = ""
	response = client.Patch(server.URL+"/posts", "application/json", bytes.NewReader([]byte(`{"userId":0,"id":3,"title":"cc","body":""}`)))
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "application/json", actualContentType)
	assert.Equal(t, `{"userId":0,"id":3,"title":"cc","body":""}`, string(actualRequestBody))

	// Test RemoveInterceptor
	client.RemoveInterceptor(&interceptorForTest)
	actualContentType = ""
	response = client.Patch(server.URL+"/posts", "application/json", bytes.NewReader([]byte(`{"userId":0,"id":3,"title":"cc","body":""}`)))
	assert.Equal(t, "", actualContentType)

	// api := NewSimpleAPI("https://jsonplaceholder.typicode.com")
	api := NewSimpleAPI(server.URL)
	api.GetSimpleHTTP().AddInterceptor(&interceptorForTest)

	var apiResponse *APIResponse[PostListResponse]

	postsGet := APIMakeGet[PostListResponse](api, "posts")
	apiResponse = postsGet(nil, &PostListResponse{}).Eval()
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "/posts", actualPath)
	assert.Equal(t, 2, len(apiResponse.TargetObject.Data))

	postsGetOne := APIMakeGet[PostListResponse](api, "posts/{id}")
	apiResponse = postsGetOne(PathParam{"id": 1}, &PostListResponse{}).Eval()
	assert.Equal(t, "/posts/1", actualPath)
	assert.Equal(t, nil, response.Err)

	postsDeleteOne := APIMakeDelete[PostListResponse](api, "posts/{id}")
	apiResponse = postsDeleteOne(PathParam{"id": 1}, &PostListResponse{}).Eval()
	assert.Equal(t, "/posts/1", actualPath)
	assert.Equal(t, nil, response.Err)

	actualContentType = ""
	postsPost := APIMakePostJSONBody[Post, PostListResponse](api, "posts")
	apiResponse = postsPost(nil, Post{ID: 5, Title: "aa"}, &PostListResponse{}).Eval()
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "application/json", actualContentType)
	assert.Equal(t, `{"userId":0,"id":5,"title":"aa","body":""}`, string(actualRequestBody))

	actualContentType = ""
	postsPut := APIMakePutJSONBody[Post, PostListResponse](api, "posts")
	apiResponse = postsPut(nil, Post{ID: 4, Title: "bb"}, &PostListResponse{}).Eval()
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "application/json", actualContentType)
	assert.Equal(t, `{"userId":0,"id":4,"title":"bb","body":""}`, string(actualRequestBody))

	actualContentType = ""
	postsPatch := APIMakePatchJSONBody[Post, PostListResponse](api, "posts")
	apiResponse = postsPatch(nil, Post{ID: 3, Title: "cc"}, &PostListResponse{}).Eval()
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "application/json", actualContentType)
	assert.Equal(t, `{"userId":0,"id":3,"title":"cc","body":""}`, string(actualRequestBody))

	// Test ClearInterceptor
	api.GetSimpleHTTP().ClearInterceptor()
	actualContentType = ""
	apiResponse = postsPatch(nil, Post{ID: 3, Title: "cc"}, &PostListResponse{}).Eval()
	assert.Equal(t, "", actualContentType)
}

func TestSimpleAPIMultipart(t *testing.T) {
	// var actualPath string
	var actualRequest *http.Request
	var actualRequestBody []byte
	var actualContentType string

	postsHandler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		actualRequestBody, _ = ioutil.ReadAll(req.Body)

		// auth := req.Header.Get("Auth")
		_, err := writer.Write([]byte(`{}`))
		assert.NoError(t, err)
	})

	server := httptest.NewServer(postsHandler)
	defer server.Close()
	// router := httprouter.New()
	// router.GET("/posts", postsHandler)
	// recorder := httptest.NewRecorder()

	var apiResponse *APIResponse[PostListResponse]

	interceptorForTest := Interceptor(func(request *http.Request) error {
		// actualPath = request.URL.Path
		actualRequest = request
		actualContentType = actualRequest.Header.Get("Content-Type")
		return nil
	})

	// api := NewSimpleAPI("https://jsonplaceholder.typicode.com")
	api := NewSimpleAPI(server.URL)
	api.GetSimpleHTTP().AddInterceptor(&interceptorForTest)

	var multipartReader *multipart.Reader
	var params map[string]string

	var actualForm *multipart.Form
	var sentValues map[string][]string

	fileDir, _ := os.Getwd()
	fileName := "simpleHTTP_test.go"
	filePath := path.Join(fileDir, fileName)

	actualContentType = ""
	postsPost := APIMakePostMultipartBody[PostListResponse](api, "posts")
	sentValues = map[string][]string{"userId": {"0"}, "id": {"5"}, "title": {"aa"}, "body": {""}}
	sentFiles := map[string][]string{"file": {filePath}}
	apiResponse = postsPost(nil, &MultipartForm{Value: sentValues, File: sentFiles}, &PostListResponse{}).Eval()
	assert.Equal(t, nil, apiResponse.Err)
	_, params, _ = mime.ParseMediaType(actualContentType)
	multipartReader = multipart.NewReader(bytes.NewReader(actualRequestBody), params["boundary"])
	actualForm, _ = multipartReader.ReadForm(1024)
	assert.Equal(t, sentValues, actualForm.Value)
	assert.Equal(t, 1, len(actualForm.File["file"]))

	actualContentType = ""
	postsPut := APIMakePutMultipartBody[PostListResponse](api, "posts")
	sentValues = map[string][]string{"userId": {"0"}, "id": {"4"}, "title": {"bb"}, "body": {""}}
	apiResponse = postsPut(nil, &MultipartForm{Value: sentValues}, &PostListResponse{}).Eval()
	assert.Equal(t, nil, apiResponse.Err)
	_, params, _ = mime.ParseMediaType(actualContentType)
	multipartReader = multipart.NewReader(bytes.NewReader(actualRequestBody), params["boundary"])
	actualForm, _ = multipartReader.ReadForm(1024)
	assert.Equal(t, sentValues, actualForm.Value)
	assert.Equal(t, 0, len(actualForm.File["file"]))

	actualContentType = ""
	postsPatch := APIMakePatchMultipartBody[PostListResponse](api, "posts")
	sentValues = map[string][]string{"userId": {"0"}, "id": {"3"}, "title": {"cc"}, "body": {""}}
	apiResponse = postsPatch(nil, &MultipartForm{Value: sentValues}, &PostListResponse{}).Eval()
	assert.Equal(t, nil, apiResponse.Err)
	_, params, _ = mime.ParseMediaType(actualContentType)
	multipartReader = multipart.NewReader(bytes.NewReader(actualRequestBody), params["boundary"])
	actualForm, _ = multipartReader.ReadForm(1024)
	assert.Equal(t, sentValues, actualForm.Value)
	assert.Equal(t, 0, len(actualForm.File["file"]))
}

func TestGeneralMultipartSerializerMultipleFilesAndValues(t *testing.T) {
	file1, err := os.CreateTemp("", "fpgo-multipart-1-*.txt")
	assert.NoError(t, err)
	defer os.Remove(file1.Name())
	_, err = file1.WriteString("first")
	assert.NoError(t, err)
	assert.NoError(t, file1.Close())

	file2, err := os.CreateTemp("", "fpgo-multipart-2-*.txt")
	assert.NoError(t, err)
	defer os.Remove(file2.Name())
	_, err = file2.WriteString("second")
	assert.NoError(t, err)
	assert.NoError(t, file2.Close())

	reader, contentType, err := GeneralMultipartSerializer(&MultipartForm{
		Value: map[string][]string{
			"alpha": {"1", "2"},
			"beta":  {"x"},
		},
		File: map[string][]string{
			"file": {file1.Name(), file2.Name()},
		},
	})
	assert.NoError(t, err)
	assert.Contains(t, contentType, "multipart/form-data")

	body, err := io.ReadAll(reader)
	assert.NoError(t, err)
	_, params, err := mime.ParseMediaType(contentType)
	assert.NoError(t, err)

	multipartReader := multipart.NewReader(bytes.NewReader(body), params["boundary"])
	form, err := multipartReader.ReadForm(1024 * 1024)
	assert.NoError(t, err)
	assert.Equal(t, []string{"1", "2"}, form.Value["alpha"])
	assert.Equal(t, []string{"x"}, form.Value["beta"])
	assert.Len(t, form.File["file"], 2)
}

func TestGeneralMultipartSerializerMissingFile(t *testing.T) {
	reader, contentType, err := GeneralMultipartSerializer(&MultipartForm{
		Value: map[string][]string{"alpha": {"1"}},
		File:  map[string][]string{"file": {"/definitely/missing/file.txt"}},
	})
	assert.Nil(t, reader)
	assert.Equal(t, "", contentType)
	assert.Error(t, err)
}

func TestGeneralMultipartSerializerNilFormPanics(t *testing.T) {
	assert.Panics(t, func() {
		_, _, _ = GeneralMultipartSerializer(nil)
	})
}

func TestSimpleHTTPRecursiveVisitWithoutTransportPanics(t *testing.T) {
	client := &SimpleHTTPDef{}
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	assert.Panics(t, func() {
		_, _ = client.RoundTrip(req)
	})
}

func TestAPIMakeJSONAndMultipartMethodSelection(t *testing.T) {
	var methods []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		methods = append(methods, req.Method)
		_, err := w.Write([]byte(`{"data":[]}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	api := NewSimpleAPI(server.URL)

	putJSON := APIMakePutJSONBody[Post, PostListResponse](api, "posts")
	patchJSON := APIMakePatchJSONBody[Post, PostListResponse](api, "posts")
	putMultipart := APIMakePutMultipartBody[PostListResponse](api, "posts")
	patchMultipart := APIMakePatchMultipartBody[PostListResponse](api, "posts")

	assert.NoError(t, putJSON(nil, Post{ID: 1}, &PostListResponse{}).Eval().Err)
	assert.NoError(t, patchJSON(nil, Post{ID: 2}, &PostListResponse{}).Eval().Err)
	assert.NoError(t, putMultipart(nil, &MultipartForm{Value: map[string][]string{"k": {"v"}}}, &PostListResponse{}).Eval().Err)
	assert.NoError(t, patchMultipart(nil, &MultipartForm{Value: map[string][]string{"k": {"v"}}}, &PostListResponse{}).Eval().Err)

	assert.Equal(t, []string{http.MethodPut, http.MethodPatch, http.MethodPut, http.MethodPatch}, methods)
}

func TestSimpleHTTPInterceptor(t *testing.T) {
	client := NewSimpleHTTP()

	interceptor := Interceptor(func(request *http.Request) error {
		request.Header.Set("X-Custom-Header", "test")
		return nil
	})

	client.AddInterceptor(&interceptor)
	assert.Equal(t, 1, len(client.interceptors))

	client.RemoveInterceptor(&interceptor)
	assert.Equal(t, 0, len(client.interceptors))

	client.AddInterceptor(&interceptor)
	client.ClearInterceptor()
	assert.Equal(t, 0, len(client.interceptors))
}

func TestNewSimpleHTTPWithClient(t *testing.T) {
	httpClient := &http.Client{}
	client := NewSimpleHTTPWithClientAndInterceptors(httpClient)
	assert.NotNil(t, client)
	assert.Equal(t, httpClient, client.client)
}

func TestSimpleHTTPClearInterceptor(t *testing.T) {
	client := NewSimpleHTTP()

	interceptor := Interceptor(func(request *http.Request) error {
		return nil
	})

	client.AddInterceptor(&interceptor)
	client.ClearInterceptor()
	assert.Equal(t, 0, len(client.interceptors))
}

func TestSimpleHTTPClientGetterSetter(t *testing.T) {
	client := NewSimpleHTTP()

	httpClient := client.GetHTTPClient()
	assert.NotNil(t, httpClient)

	newClient := &http.Client{
		Timeout: 5 * 1000000000,
	}
	client.SetHTTPClient(newClient)
	assert.Equal(t, newClient, client.GetHTTPClient())
}

func TestNewSimpleAPIWithSimpleHTTP(t *testing.T) {
	httpClient := &http.Client{}
	simpleHTTP := NewSimpleHTTPWithClientAndInterceptors(httpClient)

	api := NewSimpleAPIWithSimpleHTTP("http://example.com", simpleHTTP)
	assert.NotNil(t, api)
	assert.Equal(t, simpleHTTP, api.GetSimpleHTTP())
}

func TestJSONBodySerializer(t *testing.T) {
	body := map[string]string{"key": "value"}
	reader, err := JSONBodySerializer(body)
	assert.Nil(t, err)
	assert.NotNil(t, reader)
}

func TestJSONBodyDeserializer(t *testing.T) {
	body := []byte(`{"key": "value"}`)
	var target map[string]string
	result, err := JSONBodyDeserializer(body, &target)
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

func TestGeneralMultipartSerializer(t *testing.T) {
	form := &MultipartForm{
		Value: map[string][]string{
			"key": {"value"},
		},
	}
	reader, contentType, err := GeneralMultipartSerializer(form)
	assert.Nil(t, err)
	assert.NotNil(t, reader)
	assert.NotEmpty(t, contentType)
}

func TestGetContextTimeoutBranches(t *testing.T) {
	client := NewSimpleHTTP()
	client.TimeoutMillisecond = 1
	ctx, cancel := client.GetContextTimeout()
	defer cancel()
	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.True(t, time.Until(deadline) <= 5*time.Millisecond)

	client.TimeoutMillisecond = 0
	ctx2, cancel2 := client.GetContextTimeout()
	defer cancel2()
	_, ok2 := ctx2.Deadline()
	assert.True(t, ok2)
}

func TestDoNewRequestInvalidURL(t *testing.T) {
	client := NewSimpleHTTP()
	resp := client.DoNewRequest(context.Background(), nil, http.MethodGet, "://bad-url")
	assert.Error(t, resp.Err)
	assert.Nil(t, resp.Response)
}

func TestDoNewRequestWithBodyOptionsEmptyContentType(t *testing.T) {
	var gotContentType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := NewSimpleHTTP()
	resp := client.DoNewRequestWithBodyOptions(context.Background(), nil, http.MethodPost, srv.URL, strings.NewReader("{}"), "")
	assert.NoError(t, resp.Err)
	assert.Equal(t, "", gotContentType)
}

func TestRecursiveVisitInterceptorError(t *testing.T) {
	client := NewSimpleHTTP()
	expected := errors.New("stop chain")
	badInterceptor := Interceptor(func(request *http.Request) error {
		return expected
	})
	client.AddInterceptor(&badInterceptor)

	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	_, err := client.RoundTrip(req)
	assert.ErrorIs(t, err, expected)
}

func TestDecodeResponseBodyErrorBranchesAndReplacePathParams(t *testing.T) {
	api := NewSimpleAPI("http://example.com")

	readErrResp := &APIResponse[PostListResponse]{
		ResponseWithError: ResponseWithError{
			Response: &http.Response{Body: errReadCloser{}},
		},
	}
	decodeResponseBody(api, readErrResp, &PostListResponse{})
	assert.Error(t, readErrResp.Err)

	expected := errors.New("deserialize failed")
	api.ResponseDeserializer = func(body []byte, target interface{}) (interface{}, error) {
		return target, expected
	}
	resp := &APIResponse[PostListResponse]{
		ResponseWithError: ResponseWithError{
			Response: &http.Response{Body: io.NopCloser(strings.NewReader(`{}`))},
		},
	}
	out := decodeResponseBody(api, resp, &PostListResponse{})
	assert.ErrorIs(t, out.Err, expected)
	assert.NotNil(t, out.TargetObject)

	finalURL := api.replacePathParams("posts/{id}", PathParam{"id": 7})
	assert.Equal(t, "http://example.com/posts/7", finalURL)
}

func TestJSONBodySerializerError(t *testing.T) {
	ch := make(chan int)
	reader, err := JSONBodySerializer(ch)
	assert.Error(t, err)
	assert.Nil(t, reader)
}

func TestGeneralMultipartSerializerWriteFieldError(t *testing.T) {
	form := &MultipartForm{
		Value: map[string][]string{
			"field1": {"value1"},
			"field2": {"value2"},
		},
	}
	reader, contentType, err := GeneralMultipartSerializer(form)
	assert.Nil(t, err)
	assert.NotNil(t, reader)
	assert.NotEmpty(t, contentType)
	assert.Contains(t, contentType, "multipart/form-data")
}

func TestDoNewRequestWithHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := NewSimpleHTTP()
	header := http.Header{}
	header.Set("X-Custom", "value")
	resp := client.DoNewRequest(context.Background(), header, http.MethodGet, srv.URL)
	assert.NoError(t, resp.Err)
	assert.NotNil(t, resp.Response)
	assert.Equal(t, "value", resp.Request.Header.Get("X-Custom"))
}

func TestDoNewRequestWithBodyOptionsBadURL(t *testing.T) {
	client := NewSimpleHTTP()
	resp := client.DoNewRequestWithBodyOptions(context.Background(), nil, http.MethodPost, "://bad-url", nil, "text/plain")
	assert.Error(t, resp.Err)
	assert.Nil(t, resp.Response)
}

func TestDoNewRequestWithBodyOptionsWithHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := NewSimpleHTTP()
	header := http.Header{}
	header.Set("X-Custom", "value")
	resp := client.DoNewRequestWithBodyOptions(context.Background(), header, http.MethodPost, srv.URL, strings.NewReader("body"), "application/json")
	assert.NoError(t, resp.Err)
	assert.NotNil(t, resp.Response)
	assert.Equal(t, "value", resp.Request.Header.Get("X-Custom"))
	assert.Equal(t, "application/json", resp.Request.Header.Get("Content-Type"))
}

func TestGeneralMultipartSerializerWithFiles(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-upload-*.txt")
	assert.NoError(t, err)
	_, err = tmpFile.WriteString("test content")
	assert.NoError(t, err)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	form := &MultipartForm{
		Value: map[string][]string{"field1": {"val1"}},
		File:  map[string][]string{"file1": {tmpFile.Name()}},
	}
	reader, contentType, err := GeneralMultipartSerializer(form)
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.NotEmpty(t, contentType)
	assert.Contains(t, contentType, "multipart/form-data")
}

func TestGeneralMultipartSerializerOpenFileError(t *testing.T) {
	form := &MultipartForm{
		File: map[string][]string{"bad": {"/nonexistent/file.txt"}},
	}
	_, _, err := GeneralMultipartSerializer(form)
	assert.Error(t, err)
}

func TestGeneralMultipartSerializerMultipleValuesAndFiles(t *testing.T) {
	tmpFile1, err := os.CreateTemp("", "test-upload-1-*.txt")
	assert.NoError(t, err)
	_, err = tmpFile1.WriteString("first")
	assert.NoError(t, err)
	assert.NoError(t, tmpFile1.Close())
	defer os.Remove(tmpFile1.Name())

	tmpFile2, err := os.CreateTemp("", "test-upload-2-*.txt")
	assert.NoError(t, err)
	_, err = tmpFile2.WriteString("second")
	assert.NoError(t, err)
	assert.NoError(t, tmpFile2.Close())
	defer os.Remove(tmpFile2.Name())

	form := &MultipartForm{
		Value: map[string][]string{"field": {"v1", "v2"}},
		File:  map[string][]string{"file": {tmpFile1.Name(), tmpFile2.Name()}},
	}

	reader, contentType, err := GeneralMultipartSerializer(form)
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Contains(t, contentType, "multipart/form-data")

	bodyBytes, err := io.ReadAll(reader)
	assert.NoError(t, err)

	_, params, err := mime.ParseMediaType(contentType)
	assert.NoError(t, err)

	multipartReader := multipart.NewReader(bytes.NewReader(bodyBytes), params["boundary"])
	actualForm, err := multipartReader.ReadForm(4096)
	assert.NoError(t, err)
	assert.Equal(t, []string{"v1", "v2"}, actualForm.Value["field"])
	assert.Len(t, actualForm.File["file"], 2)
}

func TestAPIMakePostJSONBodyNilBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer srv.Close()

	api := NewSimpleAPI(srv.URL)
	postsPost := APIMakePostJSONBody[*Post, PostListResponse](api, "posts")
	apiResponse := postsPost(nil, nil, &PostListResponse{}).Eval()
	assert.NoError(t, apiResponse.Err)
	assert.NotNil(t, apiResponse.TargetObject)
}

func TestAPIMakePostJSONBodySerializerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer srv.Close()

	api := NewSimpleAPI(srv.URL)
	postsPost := APIMakePostJSONBody[chan int, PostListResponse](api, "posts")
	apiResponse := postsPost(nil, make(chan int), &PostListResponse{}).Eval()
	assert.Error(t, apiResponse.Err)
}

func TestAPIMakePostMultipartBodyNilBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer srv.Close()

	api := NewSimpleAPI(srv.URL)
	postsPost := APIMakePostMultipartBody[PostListResponse](api, "posts")
	apiResponse := postsPost(nil, nil, &PostListResponse{}).Eval()
	assert.NoError(t, apiResponse.Err)
	assert.NotNil(t, apiResponse.TargetObject)
}

func TestAPIMakePostMultipartBodySerializerError(t *testing.T) {
	api := NewSimpleAPI("http://example.com")
	postsPost := APIMakePostMultipartBody[PostListResponse](api, "posts")
	form := &MultipartForm{
		File: map[string][]string{"bad": {"/nonexistent/file.txt"}},
	}
	apiResponse := postsPost(nil, form, &PostListResponse{}).Eval()
	assert.Error(t, apiResponse.Err)
}

func TestAPIMakeDoNewRequestWithHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer srv.Close()

	api := NewSimpleAPI(srv.URL)
	api.DefaultHeader = http.Header{}
	api.DefaultHeader.Set("X-Auth", "token123")

	var gotHeader string
	interceptor := Interceptor(func(req *http.Request) error {
		gotHeader = req.Header.Get("X-Auth")
		return nil
	})
	api.GetSimpleHTTP().AddInterceptor(&interceptor)

	getPosts := APIMakeGet[PostListResponse](api, "posts")
	apiResponse := getPosts(nil, &PostListResponse{}).Eval()
	assert.NoError(t, apiResponse.Err)
	assert.Equal(t, "token123", gotHeader)
}

func TestAPIMakeDoNewRequestWithBodyOptionsResponseError(t *testing.T) {
	api := NewSimpleAPI("http://localhost:1")
	api.GetSimpleHTTP().TimeoutMillisecond = int64(10 * time.Millisecond)
	api.DefaultHeader = http.Header{}
	api.DefaultHeader.Set("X-Auth", "token123")

	postsPost := APIMakePostJSONBody[*Post, PostListResponse](api, "posts")
	apiResponse := postsPost(nil, nil, &PostListResponse{}).Eval()
	assert.Error(t, apiResponse.Err)
}

func TestAPIMakeDoNewRequestResponseError(t *testing.T) {
	api := NewSimpleAPI("http://localhost:1")
	api.GetSimpleHTTP().TimeoutMillisecond = int64(10 * time.Millisecond)

	getPosts := APIMakeGet[PostListResponse](api, "posts")
	apiResponse := getPosts(nil, &PostListResponse{}).Eval()
	assert.Error(t, apiResponse.Err)
}

func TestGeneralMultipartSerializerWriteError(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*")
	assert.NoError(t, err)
	_, err = tmpFile.WriteString("data")
	assert.NoError(t, err)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	form := &MultipartForm{
		Value: map[string][]string{"name": {"test"}},
		File:  map[string][]string{"upload": {tmpFile.Name()}},
	}
	reader, contentType, err := GeneralMultipartSerializer(form)
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Contains(t, contentType, "multipart/form-data")
}

func TestGeneralMultipartSerializerWriteFieldMultipleFields(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*")
	assert.NoError(t, err)
	_, err = tmpFile.WriteString("test content")
	assert.NoError(t, err)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	form := &MultipartForm{
		Value: map[string][]string{
			"a": {"1"},
			"b": {"2"},
			"c": {"3"},
		},
		File: map[string][]string{
			"f": {tmpFile.Name()},
		},
	}
	reader, ct, err := GeneralMultipartSerializer(form)
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Contains(t, ct, "multipart/form-data")
}

func TestJSONBodySerializerErrorNonMarshalable(t *testing.T) {
	ch := make(chan int)
	reader, err := JSONBodySerializer(ch)
	assert.Error(t, err)
	assert.Nil(t, reader)
}

func TestAPIMakePostMultipartBodyResponseErr(t *testing.T) {
	api := NewSimpleAPI("http://0.0.0.0:1")
	postMultipart := APIMakePostMultipartBody[struct{}](api, "test")
	apiResponse := postMultipart(nil, nil, &struct{}{}).Eval()
	assert.Error(t, apiResponse.Err)
}

func TestGeneralMultipartSerializerMultipleFiles(t *testing.T) {
	tmpFile1, err := os.CreateTemp("", "test-*")
	assert.NoError(t, err)
	_, err = tmpFile1.WriteString("content1")
	assert.NoError(t, err)
	tmpFile1.Close()
	defer os.Remove(tmpFile1.Name())

	tmpFile2, err := os.CreateTemp("", "test2-*")
	assert.NoError(t, err)
	_, err = tmpFile2.WriteString("content2")
	assert.NoError(t, err)
	tmpFile2.Close()
	defer os.Remove(tmpFile2.Name())

	form := &MultipartForm{
		Value: map[string][]string{"field1": {"val1"}},
		File: map[string][]string{
			"file1": {tmpFile1.Name()},
			"file2": {tmpFile2.Name()},
		},
	}
	reader, ct, err := GeneralMultipartSerializer(form)
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Contains(t, ct, "multipart/form-data")
}

func TestGeneralMultipartSerializerEmptyForm(t *testing.T) {
	reader, contentType, err := GeneralMultipartSerializer(&MultipartForm{})
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Contains(t, contentType, "multipart/form-data")
}

func TestGeneralMultipartSerializerValueAndFileRoundTrip(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "multipart-roundtrip-*.txt")
	assert.NoError(t, err)
	_, err = tmpFile.WriteString("payload")
	assert.NoError(t, err)
	assert.NoError(t, tmpFile.Close())
	defer os.Remove(tmpFile.Name())

	reader, contentType, err := GeneralMultipartSerializer(&MultipartForm{
		Value: map[string][]string{"field": {"value"}},
		File:  map[string][]string{"upload": {tmpFile.Name()}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Contains(t, contentType, "multipart/form-data")

	body, err := io.ReadAll(reader)
	assert.NoError(t, err)
	_, params, err := mime.ParseMediaType(contentType)
	assert.NoError(t, err)

	multipartReader := multipart.NewReader(bytes.NewReader(body), params["boundary"])
	form, err := multipartReader.ReadForm(1024 * 1024)
	assert.NoError(t, err)
	assert.Equal(t, []string{"value"}, form.Value["field"])
	assert.Len(t, form.File["upload"], 1)
}

func TestGeneralMultipartSerializerNilMapsAndCloseErrorFree(t *testing.T) {
	reader, contentType, err := GeneralMultipartSerializer(&MultipartForm{})
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Contains(t, contentType, "multipart/form-data")
}

// TestSimpleAPIHTTPMethods guards against the copy-paste defect where
// APIMakePut*/APIMakePatch* issued POST instead of PUT/PATCH.
func TestSimpleAPIHTTPMethods(t *testing.T) {
	var actualMethod string

	handler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		actualMethod = req.Method
		_, err := writer.Write([]byte(`{"data":[]}`))
		assert.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	api := NewSimpleAPI(server.URL)

	postJSON := APIMakePostJSONBody[Post, PostListResponse](api, "posts")
	postJSON(nil, Post{ID: 1}, &PostListResponse{}).Eval()
	assert.Equal(t, http.MethodPost, actualMethod)

	putJSON := APIMakePutJSONBody[Post, PostListResponse](api, "posts")
	putJSON(nil, Post{ID: 2}, &PostListResponse{}).Eval()
	assert.Equal(t, http.MethodPut, actualMethod)

	patchJSON := APIMakePatchJSONBody[Post, PostListResponse](api, "posts")
	patchJSON(nil, Post{ID: 3}, &PostListResponse{}).Eval()
	assert.Equal(t, http.MethodPatch, actualMethod)

	postMultipart := APIMakePostMultipartBody[PostListResponse](api, "posts")
	postMultipart(nil, &MultipartForm{Value: map[string][]string{"id": {"1"}}}, &PostListResponse{}).Eval()
	assert.Equal(t, http.MethodPost, actualMethod)

	putMultipart := APIMakePutMultipartBody[PostListResponse](api, "posts")
	putMultipart(nil, &MultipartForm{Value: map[string][]string{"id": {"2"}}}, &PostListResponse{}).Eval()
	assert.Equal(t, http.MethodPut, actualMethod)

	patchMultipart := APIMakePatchMultipartBody[PostListResponse](api, "posts")
	patchMultipart(nil, &MultipartForm{Value: map[string][]string{"id": {"3"}}}, &PostListResponse{}).Eval()
	assert.Equal(t, http.MethodPatch, actualMethod)
}

type trackingReadCloser struct {
	data    []byte
	pos     int
	closed  bool
	onClose func()
}

func (t *trackingReadCloser) Read(p []byte) (n int, err error) {
	if t.pos >= len(t.data) {
		return 0, io.EOF
	}
	n = copy(p, t.data[t.pos:])
	t.pos += n
	return n, nil
}

func (t *trackingReadCloser) Close() error {
	t.closed = true
	if t.onClose != nil {
		t.onClose()
	}
	return nil
}

func TestReplacePathParamsMultipleParams(t *testing.T) {
	api := NewSimpleAPI("http://example.com")
	result := api.replacePathParams("users/{userId}/posts/{postId}", PathParam{
		"userId": "123",
		"postId": "456",
	})
	assert.Equal(t, "http://example.com/users/123/posts/456", result)
}

func TestGetContextTimeoutMillisecondUnit(t *testing.T) {
	client := NewSimpleHTTP()
	client.TimeoutMillisecond = 100 // 100ms
	ctx, cancel := client.GetContextTimeout()
	defer cancel()
	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	expected := time.Now().Add(100 * time.Millisecond)
	diff := deadline.Sub(expected)
	if diff < 0 {
		diff = -diff
	}
	assert.Less(t, diff, 50*time.Millisecond,
		"TimeoutMillisecond=100 should produce ~100ms deadline, not ~100ns")
}

func TestDecodeResponseBodyClosesBody(t *testing.T) {
	closed := false
	body := &trackingReadCloser{
		data:    []byte(`{"value":42}`),
		onClose: func() { closed = true },
	}
	api := NewSimpleAPI("http://example.com")
	api.ResponseDeserializer = func(b []byte, target interface{}) (interface{}, error) {
		return target, nil
	}
	resp := &APIResponse[int]{
		ResponseWithError: ResponseWithError{
			Response: &http.Response{Body: body},
		},
	}
	target := 0
	decodeResponseBody[int](api, resp, &target)
	assert.True(t, closed, "response body should be closed after decoding")
}
