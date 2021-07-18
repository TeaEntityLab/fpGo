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

	postsGet := api.MakeGet("posts")
	response = postsGet(nil, &PostListResponse{}).Eval().(*ResponseWithError)
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "/posts", actualPath)
	assert.Equal(t, 2, len(response.TargetObject.(*PostListResponse).Data))

	postsGetOne := api.MakeGet("posts/{id}")
	response = postsGetOne(PathParam{"id": 1}, &struct{}{}).Eval().(*ResponseWithError)
	assert.Equal(t, "/posts/1", actualPath)
	assert.Equal(t, nil, response.Err)

	postsDeleteOne := api.MakeDelete("posts/{id}")
	response = postsDeleteOne(PathParam{"id": 1}, &struct{}{}).Eval().(*ResponseWithError)
	assert.Equal(t, "/posts/1", actualPath)
	assert.Equal(t, nil, response.Err)

	actualContentType = ""
	postsPost := api.MakePostJSONBody("posts")
	response = postsPost(nil, Post{ID: 5, Title: "aa"}, &struct{}{}).Eval().(*ResponseWithError)
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "application/json", actualContentType)
	assert.Equal(t, `{"userId":0,"id":5,"title":"aa","body":""}`, string(actualRequestBody))

	actualContentType = ""
	postsPut := api.MakePutJSONBody("posts")
	response = postsPut(nil, Post{ID: 4, Title: "bb"}, &struct{}{}).Eval().(*ResponseWithError)
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "application/json", actualContentType)
	assert.Equal(t, `{"userId":0,"id":4,"title":"bb","body":""}`, string(actualRequestBody))

	actualContentType = ""
	postsPatch := api.MakePatchJSONBody("posts")
	response = postsPatch(nil, Post{ID: 3, Title: "cc"}, &struct{}{}).Eval().(*ResponseWithError)
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "application/json", actualContentType)
	assert.Equal(t, `{"userId":0,"id":3,"title":"cc","body":""}`, string(actualRequestBody))

	// Test ClearInterceptor
	api.GetSimpleHTTP().ClearInterceptor()
	actualContentType = ""
	response = postsPatch(nil, Post{ID: 3, Title: "cc"}, &struct{}{}).Eval().(*ResponseWithError)
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

	var response *ResponseWithError

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
	postsPost := api.MakePostMultipartBody("posts")
	sentValues = map[string][]string{"userId": []string{"0"}, "id": []string{"5"}, "title": []string{"aa"}, "body": []string{""}}
	sentFiles := map[string][]string{"file": []string{filePath}}
	response = postsPost(nil, &MultipartForm{Value: sentValues, File: sentFiles}, &struct{}{}).Eval().(*ResponseWithError)
	assert.Equal(t, nil, response.Err)
	_, params, _ = mime.ParseMediaType(actualContentType)
	multipartReader = multipart.NewReader(bytes.NewReader(actualRequestBody), params["boundary"])
	actualForm, _ = multipartReader.ReadForm(1024)
	assert.Equal(t, sentValues, actualForm.Value)
	assert.Equal(t, 1, len(actualForm.File["file"]))

	actualContentType = ""
	postsPut := api.MakePutMultipartBody("posts")
	sentValues = map[string][]string{"userId": []string{"0"}, "id": []string{"4"}, "title": []string{"bb"}, "body": []string{""}}
	response = postsPut(nil, &MultipartForm{Value: sentValues}, &struct{}{}).Eval().(*ResponseWithError)
	assert.Equal(t, nil, response.Err)
	_, params, _ = mime.ParseMediaType(actualContentType)
	multipartReader = multipart.NewReader(bytes.NewReader(actualRequestBody), params["boundary"])
	actualForm, _ = multipartReader.ReadForm(1024)
	assert.Equal(t, sentValues, actualForm.Value)
	assert.Equal(t, 0, len(actualForm.File["file"]))

	actualContentType = ""
	postsPatch := api.MakePatchMultipartBody("posts")
	sentValues = map[string][]string{"userId": []string{"0"}, "id": []string{"3"}, "title": []string{"cc"}, "body": []string{""}}
	response = postsPatch(nil, &MultipartForm{Value: sentValues}, &struct{}{}).Eval().(*ResponseWithError)
	assert.Equal(t, nil, response.Err)
	_, params, _ = mime.ParseMediaType(actualContentType)
	multipartReader = multipart.NewReader(bytes.NewReader(actualRequestBody), params["boundary"])
	actualForm, _ = multipartReader.ReadForm(1024)
	assert.Equal(t, sentValues, actualForm.Value)
	assert.Equal(t, 0, len(actualForm.File["file"]))
}
