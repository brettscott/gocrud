package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mergermarket/gotools"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/brettscott/gocrud/api"
)

func TestInternal_Route_uses_db_handler(t *testing.T) {

	testRouter := routerWithTestHandlers(t)

	var routeTests = []struct {
		route  string
		result string
	}{
		{"/internal/healthcheck", "test response health check"},
	}

	for _, testCase := range routeTests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, testCase.route, nil)

		testRouter.ServeHTTP(w, req)
		if http.StatusOK != w.Code {
			t.Error("bad status expected 200 got", w.Code)
		}
		if testCase.result != w.Body.String() {
			t.Error("bad response expected", testCase.result, "got", w.Body.String())
		}
	}
}

func routerWithTestHandlers(t *testing.T) http.Handler {
	testLogger := &tools.TestLogger{T: t}
	tsdConfig := tools.NewStatsDConfig(false, testLogger)
	testStatsD, _ := tools.NewStatsD(tsdConfig)

	healthcheck := newHandler("test response health check")
	//apiGateway := &api.Gateway{}
	apiRouter := mux.NewRouter()

	return newRouter(testLogger, testStatsD, healthcheck, apiRouter)
}

func createNewHandler(prefix string) api.CreateHandlerWithPrefix
func newHandler(body string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})
}
