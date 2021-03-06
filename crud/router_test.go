package crud

import (
	"fmt"
	"github.com/mergermarket/gotools"
	"github.com/pressly/chi"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInternal_Route(t *testing.T) {

	testRouter := routerWithTestHandlers(t)

	var routeTests = []struct {
		route  string
		result string
	}{
		{"/healthcheck", "test response health check"},
		{"/api/test-url", "test response for API route"},
	}

	for _, testCase := range routeTests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, testCase.route, nil)

		testRouter.ServeHTTP(w, req)
		if http.StatusOK != w.Code {
			t.Error("bad status expected 200 got", w.Code)
		}
		if testCase.result != w.Body.String() {
			t.Error(testCase.route, "bad response expected", testCase.result, "got", w.Body.String())
		}
	}
}

func routerWithTestHandlers(t *testing.T) http.Handler {
	testLogger := &tools.TestLogger{T: t}
	tsdConfig := tools.NewStatsDConfig(false, testLogger)
	testStatsD, _ := tools.NewStatsD(tsdConfig)

	healthcheckHandlerFunc := newHandlerFunc("test response health check")
	apiRouteHandler := newChiRouteHandler("test response for API route")
	uiRouteHandler := newChiRouteHandler("test response for UI route")

	return newRouter(testLogger, testStatsD, healthcheckHandlerFunc, apiRouteHandler, uiRouteHandler)
}

func newHandlerFunc(body string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})
}

func newChiRouteHandler(body string) func(chi.Router) {
	return func(router chi.Router) {
		router.Get("/test-url", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, body)
		})
	}
}
