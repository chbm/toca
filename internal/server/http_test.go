package httpserver

import (
	"net/http/httptest"
	"testing"
	"strings"
	"bytes"
	"io"
	"gotest.tools/v3/assert"

	"github.com/go-chi/chi/v5"
	
	"github.com/chbm/toca/internal/storage"
)

func httpServer() *chi.Mux {
	clerkCh := storage.Start()
	return Start(clerkCh)
}


func TestBasic(t *testing.T) {
	server := httpServer()
	recorder := httptest.NewRecorder()

	putReq := httptest.NewRequest("PUT", "/default/foo", strings.NewReader("abcdefg"))
	server.ServeHTTP(recorder, putReq)
	res := recorder.Result()
	assert.Equal(t, res.StatusCode, 201) 

	readReq := httptest.NewRequest("GET", "/default/foo", bytes.NewReader([]byte{}))
	recorder = httptest.NewRecorder()
	server.ServeHTTP(recorder, readReq)
	res = recorder.Result()
	assert.Equal(t, res.StatusCode, 200)
	b, _ := io.ReadAll(res.Body)
	assert.Equal(t, string(b), "abcdefg")

	delReq := httptest.NewRequest("DELETE", "/default/foo", bytes.NewReader([]byte{}))
	recorder = httptest.NewRecorder()
	server.ServeHTTP(recorder, delReq)
	res = recorder.Result()
	assert.Equal(t, res.StatusCode, 204) 

	rereadReq := httptest.NewRequest("GET", "/default/foo", bytes.NewReader([]byte{}))
	recorder = httptest.NewRecorder()
	server.ServeHTTP(recorder, rereadReq)
	res = recorder.Result()
	assert.Equal(t, res.StatusCode, 404)
}

