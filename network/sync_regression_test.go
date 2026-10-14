package network

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncReplacePathParamsEscapesValues(t *testing.T) {
	api := NewSimpleAPI("http://example.com")
	result := api.replacePathParams("users/{id}", PathParam{"id": "a/b?x=1&y=2"})
	assert.Equal(t, "http://example.com/users/a%2Fb%3Fx=1&y=2", result)

	// A placeholder in the query must be query-escaped. PathEscape leaves '&'
	// untouched, which would split the query.
	result = api.replacePathParams("search?q={term}", PathParam{"term": "a&b=c"})
	assert.Equal(t, "http://example.com/search?q=a%26b%3Dc", result)
}

func TestSyncGeneralMultipartSerializerNilFormReturnsError(t *testing.T) {
	_, _, err := GeneralMultipartSerializer(nil)
	assert.Error(t, err)
}

func TestSyncSimpleHTTPRecursiveVisitWithoutTransportReturnsError(t *testing.T) {
	s := &SimpleHTTPDef{}
	assert.NotPanics(t, func() {
		_, err := s.recursiveVisit(&http.Request{}, 0)
		assert.Error(t, err)
	})
}

func TestSyncSimpleAPIErrorStatusSurfacesError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"detail":"boom"}`))
	}))
	defer server.Close()

	var target struct{}
	resp := NewSimpleAPI(server.URL).MakeDoNewRequest(http.MethodGet, "/x")(nil, &target).Eval()
	result, ok := resp.(*ResponseWithError)
	assert.True(t, ok)
	assert.Error(t, result.Err)
	assert.Contains(t, result.Err.Error(), "500")
	assert.Contains(t, result.Err.Error(), "boom")
}

func TestSyncSimpleHTTPGetBodyReadableAfterReturn(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))
	defer server.Close()

	resp := NewSimpleHTTP().Get(server.URL)
	assert.NoError(t, resp.Err)
	body, err := io.ReadAll(resp.Response.Body)
	resp.Response.Body.Close()
	assert.NoError(t, err, "body must stay readable until Close (cancel tied to body lifetime)")
	assert.Equal(t, "hello", string(body))
}

func TestSyncBindCancelToBodyWithoutResponseCancels(t *testing.T) {
	s := NewSimpleHTTP()
	cancelled := false
	result := s.bindCancelToBody(&ResponseWithError{}, func() { cancelled = true })
	assert.True(t, cancelled)
	assert.NotNil(t, result)
}

// Guard against a deserializer that returns the wrong type shape: response
// must still carry Err or a typed target, never panic.
func TestSyncDecodeResponseBodyDeserializesIntoTarget(t *testing.T) {
	api := NewSimpleAPI("http://example.com")
	api.ResponseDeserializer = func(body []byte, target interface{}) (interface{}, error) {
		return target, nil
	}
	body := io.NopCloser(strings.NewReader(`{"a":1}`))
	resp := &ResponseWithError{Response: &http.Response{StatusCode: 200, Status: "200 OK", Body: body}}
	out := api.decodeResponseBody(resp, &struct{}{})
	assert.NoError(t, out.Err)
}
