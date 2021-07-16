package fpgo

import (
	"net/http"
	"net/http/httptest"
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

	postsHandler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		actualPath = req.URL.Path

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

	// api := NewSimpleAPI("https://jsonplaceholder.typicode.com")
	api := NewSimpleAPI(server.URL)
	postsGet := api.MakeGet("posts")
	response := postsGet(nil, &PostListResponse{}).Eval().(*ResponseWithError)
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "/posts", actualPath)

	result := response.TargetObject.(*PostListResponse)
	assert.Equal(t, 2, len(result.Data))

}
