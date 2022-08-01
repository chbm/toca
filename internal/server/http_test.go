package httpserver

import (
	"io"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/go-chi/chi/v5"

	"github.com/chbm/toca/internal/storage"
)

func httpServer() *chi.Mux {
	clerkCh := storage.Start()
	return Start(clerkCh)
}

func testRequest(t *testing.T, server *chi.Mux, how string, where string, what string, expect int) string {
	testReq := httptest.NewRequest(how, where, strings.NewReader(what))
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, testReq)
	res := recorder.Result()
	assert.Equal(t, res.StatusCode, expect)
	b, e := io.ReadAll(res.Body)
	if e != nil {
		return "" 
	}
	return string(b)
}

func TestBasic(t *testing.T) {
	server := httpServer()
	
	aKey := "/default/foo"
	aValue := "ipsum bacon"

	testRequest(t, server, "PUT", aKey, aValue, 201)
	assert.Equal(t, testRequest(t, server, "GET", aKey, "", 200), aValue)
	testRequest(t, server, "DELETE", aKey, "", 204)
	testRequest(t, server, "GET", aKey, "", 404)
	
	testRequest(t, server, "PUT", aKey, aValue, 201)
	assert.Equal(t, testRequest(t, server, "PUT", aKey, "", 200), aValue)
}


func TestNs(t *testing.T) {
	server := httpServer()

	testRequest(t, server, "PUT", "/null/foo", "ola", 404)
	testRequest(t, server, "POST", "/null", "", 201)
	testRequest(t, server, "POST", "/null", "", 409)
	testRequest(t, server, "PUT", "/null/foo", "ola", 201)
	assert.Equal(t, testRequest(t, server, "GET", "/null/foo", "", 200), "ola")
	
}

func TestLoadSave(t *testing.T) {
	server := httpServer()

	testRequest(t, server, "POST", "/null", "", 201)
	testRequest(t, server, "PUT", "/null/foo", "ola", 201)
	testRequest(t, server, "POST", "/null/_save", "", 200)

	otherserver := httpServer()
	testRequest(t, otherserver, "POST", "/null/_load", "", 200)
	assert.Equal(t, testRequest(t, otherserver, "GET", "/null/foo", "", 200), "ola")
}

func FuzzStore(f *testing.F) {
	server := httpServer()

	f.Fuzz(func(t *testing.T, key string, value string) {
		path := "/default/" + url.QueryEscape(key)	
		testReq := httptest.NewRequest("PUT", path, strings.NewReader(value))
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, testReq)
		res := recorder.Result()
		if res.StatusCode > 299 {
			t.Fatalf("PUT %v : %v", key, value)
		}
		assert.Equal(t, testRequest(t, server, "GET", path, "", 200), value)
	})
}
