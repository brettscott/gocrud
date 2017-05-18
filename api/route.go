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

		// List
		// eg http://localhost:8080/gocrud/api/user
		r.Get("/:entityID", list)
	}
}

func list(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityID")
	w.Write([]byte(fmt.Sprintf("List: %v", entityID)))



}
