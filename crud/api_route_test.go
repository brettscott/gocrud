package crud

import (
	"bytes"
	"errors"
	"github.com/brettscott/gocrud/model"
	"github.com/mergermarket/gotools"
	"github.com/pressly/chi"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPIRoute(t *testing.T) {

	t.Run("GET /<entity> returns 200", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{
			listResponseBody: []byte("the-test-response"),
		}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		testRouter.ServeHTTP(w, req)
		if http.StatusOK != w.Code {
			t.Error("bad status: expected", http.StatusOK, "got", w.Code, "body:", w.Body.String())
		}
		if strings.Contains(w.Body.String(), "the-test-response") == false {
			t.Error("body doesn't contain expected string.  Body: ", w.Body.String())
		}
	})

	t.Run("GET /<entity> returns 500 when error response from service", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{
			listResponseError: errors.New("the service failed"),
		}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		testRouter.ServeHTTP(w, req)
		if http.StatusInternalServerError != w.Code {
			t.Error("bad status: expected", http.StatusInternalServerError, "got", w.Code, "body:", w.Body.String())
		}
	})

	t.Run("GET /<entity> returns 400 when invalid entity provided", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/doesnt-exist", nil)
		testRouter.ServeHTTP(w, req)
		if http.StatusBadRequest != w.Code {
			t.Error("bad status: expected", http.StatusBadRequest, "got", w.Code, "body:", w.Body.String())
		}
	})

	t.Run("GET /<entity>/<recordID> returns 200", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{
			getResponseBody: []byte("the-test-response"),
		}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/59df7716ed16fc4a2ca31c08", nil)
		testRouter.ServeHTTP(w, req)
		if http.StatusOK != w.Code {
			t.Error("bad status: expected", http.StatusOK, "got", w.Code, "body:", w.Body.String())
		}
		if strings.Contains(w.Body.String(), "the-test-response") == false {
			t.Error("body doesn't contain expected string.  Body: ", w.Body.String())
		}
	})

	t.Run("GET /<entity>/<recordID> returns 500 when error response from service", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{
			getResponseError: errors.New("the service failed"),
		}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/59df7716ed16fc4a2ca31c08", nil)
		testRouter.ServeHTTP(w, req)
		if http.StatusInternalServerError != w.Code {
			t.Error("bad status: expected", http.StatusInternalServerError, "got", w.Code)
		}
	})

	t.Run("GET /<entity>/<recordID> returns 400 when invalid entity provided", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/doesnt-exist/59df7716ed16fc4a2ca31c08", nil)
		testRouter.ServeHTTP(w, req)
		if http.StatusBadRequest != w.Code {
			t.Error("bad status: expected", http.StatusBadRequest, "got", w.Code, "body:", w.Body.String())
		}
	})

	t.Run("POST /<entity> returns 200", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{
			saveResponseBody: []byte("the-test-response"),
		}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(""))
		req.Header.Set("content-type", "application/json")
		testRouter.ServeHTTP(w, req)
		if http.StatusCreated != w.Code {
			t.Error("bad status: expected", http.StatusCreated, "got", w.Code)
		}
		if strings.Contains(w.Body.String(), "the-test-response") == false {
			t.Error("body doesn't contain expected string.  Body: ", w.Body.String())
		}
	})

	t.Run("POST /<entity> returns 500 when error response from service", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{
			saveResponseError: errors.New("the service failed"),
		}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(""))
		req.Header.Set("content-type", "application/json")
		testRouter.ServeHTTP(w, req)
		if http.StatusInternalServerError != w.Code {
			t.Error("bad status: expected", http.StatusInternalServerError, "got", w.Code, "body:", w.Body.String())
		}
	})

	t.Run("POST /<entity> returns 400 when invalid entity provided", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/doesnt-exist", bytes.NewBufferString(""))
		req.Header.Set("content-type", "application/json")
		testRouter.ServeHTTP(w, req)
		if http.StatusBadRequest != w.Code {
			t.Error("bad status: expected", http.StatusBadRequest, "got", w.Code, "body:", w.Body.String())
		}
	})

	t.Run("POST /<entity> returns 400 if invalid or missing content-type", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(""))
		testRouter.ServeHTTP(w, req)
		if http.StatusBadRequest != w.Code {
			t.Error("bad status: expected", http.StatusBadRequest, "got", w.Code, "body:", w.Body.String())
		}
		if strings.Contains(w.Body.String(), "missing Content Type") == false {
			t.Error("body doesn't contain expected string.  Body: ", w.Body.String())
		}
	})

	t.Run("PUT /<entity>/<recordID> returns 200", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{
			saveResponseBody: []byte("the-test-response"),
		}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/users/59df7716ed16fc4a2ca31c08", bytes.NewBufferString(""))
		req.Header.Set("content-type", "application/json")
		testRouter.ServeHTTP(w, req)
		if http.StatusOK != w.Code {
			t.Error("bad status: expected", http.StatusOK, "got", w.Code, "body:", w.Body.String())
		}
		if strings.Contains(w.Body.String(), "the-test-response") == false {
			t.Error("body doesn't contain expected string.  Body: ", w.Body.String())
		}
	})

	t.Run("PUT /<entity> returns 405 when missing recordID", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{
			saveResponseBody: []byte("the-test-response"),
		}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/users", bytes.NewBufferString(""))
		req.Header.Set("content-type", "application/json")
		testRouter.ServeHTTP(w, req)
		if http.StatusMethodNotAllowed != w.Code {
			t.Error("bad status: expected", http.StatusMethodNotAllowed, "got", w.Code, "body:", w.Body.String())
		}
	})

	t.Run("PATCH /<entity>/<recordID> returns 200", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{
			saveResponseBody: []byte("the-test-response"),
		}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPatch, "/users/59df7716ed16fc4a2ca31c08", bytes.NewBufferString(""))
		req.Header.Set("content-type", "application/json")
		testRouter.ServeHTTP(w, req)
		if http.StatusOK != w.Code {
			t.Error("bad status: expected", http.StatusOK, "got", w.Code, "body:", w.Body.String())
		}
		if strings.Contains(w.Body.String(), "the-test-response") == false {
			t.Error("body doesn't contain expected string.  Body: ", w.Body.String())
		}
	})

	t.Run("PATCH /<entity> returns 405 when missing recordID", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{
			saveResponseBody: []byte("the-test-response"),
		}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPatch, "/users", bytes.NewBufferString(""))
		req.Header.Set("content-type", "application/json")
		testRouter.ServeHTTP(w, req)
		if http.StatusMethodNotAllowed != w.Code {
			t.Error("bad status: expected", http.StatusMethodNotAllowed, "got", w.Code, "body:", w.Body.String())
		}
	})

	t.Run("DELETE /<entity>/<recordID> returns 204", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/users/59df7716ed16fc4a2ca31c08", nil)
		testRouter.ServeHTTP(w, req)
		if http.StatusNoContent != w.Code {
			t.Error("bad status: expected", http.StatusNoContent, "got", w.Code, "body:", w.Body.String())
		}
	})

	t.Run("DELETE /<entity>/<recordID> returns 400 when invalid entity provided", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/doesnt-exist/59df7716ed16fc4a2ca31c08", nil)
		testRouter.ServeHTTP(w, req)
		if http.StatusBadRequest != w.Code {
			t.Error("bad status: expected", http.StatusBadRequest, "got", w.Code, "body:", w.Body.String())
		}
	})

	t.Run("DELETE /<entity> returns 405 when missing recordID", func(t *testing.T) {
		fakeApiService := &fakeApiServicer{
			saveResponseBody: []byte("the-test-response"),
		}
		testRouter := makeTestRouter(t, makeEntities(), fakeApiService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/users", bytes.NewBufferString(""))
		testRouter.ServeHTTP(w, req)
		if http.StatusMethodNotAllowed != w.Code {
			t.Error("bad status: expected", http.StatusMethodNotAllowed, "got", w.Code, "body:", w.Body.String())
		}
	})

}

func makeTestRouter(t *testing.T, entities model.Entities, apiService apiServicer) chi.Router {
	testLogger := &tools.TestLogger{T: t}
	testConfig := tools.NewStatsDConfig(false, testLogger)
	testStatsD, _ := tools.NewStatsD(testConfig)

	testAPIRoute := NewApiRoute(entities, apiService, testLogger, testStatsD)
	testRouter := chi.NewRouter()
	testRouter.Route("/", testAPIRoute)
	return testRouter
}

func makeEntities() model.Entities {
	return model.Entities{
		"users":     model.Entity{},
		"computers": model.Entity{},
	}
}
