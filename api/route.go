package api

import (
	"fmt"
	"github.com/brettscott/gocrud/entity"
	"github.com/pressly/chi"
	"net/http"
)

// NewRoute prepares the routes for this package
func NewRoute(entities entity.Entities, logger Logger, statsd StatsDer) func(chi.Router) {

	return func(r chi.Router) {
		r.Get("/here", getHere)
	}
}

func getHere(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprint("Hello world")))
}
