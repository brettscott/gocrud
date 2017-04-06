package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mergermarket/gotools"
	"net/http"
)

type CreateHandler func(string) http.Handler

func NewCreateHandler(logger tools.Logger, statsd tools.StatsD) CreateHandler {

	return func(prefix string) http.Handler {
		router := mux.NewRouter()

		//  /api/here
		router.HandleFunc(prefix + "/here", func(w http.ResponseWriter, r *http.Request) {
			logger.Info("/here/ hit")
			fmt.Fprint(w, "API HERE")
		}).Name("API HERE")

		router.HandleFunc(prefix + "/xxx/", func(w http.ResponseWriter, r *http.Request) {
			logger.Info("/api/ hit")
			fmt.Fprint(w, "API root")
		}).Name("API Root")

		router.HandleFunc(prefix + "/yyy/blah", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Blah is rendered")
		}).Name("API BLAH")

		return router
	}
}
