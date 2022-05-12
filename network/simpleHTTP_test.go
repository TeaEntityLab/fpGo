package network

import (
	"bytes"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
