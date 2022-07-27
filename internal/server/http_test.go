package httpserver

import (
	"net/http/httptest"
	"testing"
	"strings"
	"io"
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
}


func TestNs(t *testing.T) {
	server := httpServer()

	testRequest(t, server, "PUT", "/null/foo", "ola", 404)
	testRequest(t, server, "POST", "/null", "", 201)
	testRequest(t, server, "POST", "/null", "", 409)
	testRequest(t, server, "PUT", "/null/foo", "ola", 200)
	assert.Equal(t, testRequest(t, server, "GET", "/null/foo", "", 200), "ola")
	

}
