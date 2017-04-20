package api

import (
	"github.com/mergermarket/gotools"
	"github.com/pressly/chi"
	"net/http"
	"fmt"
)

// NewRoute prepares the routes for this package
func NewRoute(logger tools.Logger, statsd tools.StatsD) func(chi.Router) {

	return func(r chi.Router) {
		r.Get("/here", getHere)
	}
}

func getHere(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprint("Hello world")))
}
