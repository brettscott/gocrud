package api

import (
	"fmt"
	"github.com/pressly/chi"
	"net/http"
)

// NewRoute prepares the routes for this package
func NewRoute(logger Logger, statsd StatsDer) func(chi.Router) {

	return func(r chi.Router) {
		r.Get("/here", getHere)
	}
}

func getHere(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprint("Hello world")))
}
