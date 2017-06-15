package api

import (
	"github.com/brettscott/gocrud/entity"
	"github.com/mergermarket/gotools"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pressly/chi"
)

func TestAPIRoute(t *testing.T) {

	testLogger := &tools.TestLogger{T: t}
	tsdConfig := tools.NewStatsDConfig(false, testLogger)
	testStatsD, _ := tools.NewStatsD(tsdConfig)

	testAPIRoute := NewRoute(makeEntities(), testLogger, testStatsD)
	testRouter := chi.NewRouter()
	testRouter.Route("/", testAPIRoute)

	var routeTests = []struct {
		route  string
		result string
	}{
		{"/", "Welcome to the API"},
	}

	for _, testCase := range routeTests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, testCase.route, nil)

		testRouter.ServeHTTP(w, req)
		if http.StatusOK != w.Code {
			t.Error("bad status, expected 200 got", w.Code)
		}
		if testCase.result != w.Body.String() {
			t.Error(testCase.route, "bad response, expected", testCase.result, "got", w.Body.String())
		}
	}
}

func makeEntities() entity.Entities {
	return entity.Entities{
		"users":     entity.Entity{},
		"computers": entity.Entity{},
	}

}
