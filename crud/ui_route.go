package crud

import (
	"github.com/pressly/chi"
	"net/http"
)

type UIRoute struct {
	entities   Entities
	log        Logger
	statsd     StatsDer
	apiService apiServicer
}

func NewUiRoute(entities Entities, apiService apiServicer, log Logger, statsd StatsDer) func(chi.Router) {

	uiRoute := &UIRoute{
		entities:   entities,
		log:        log,
		statsd:     statsd,
		apiService: apiService,
	}

	return func(r chi.Router) {

		r.Get("/", uiRoute.root)

	}
}

func (u *UIRoute) root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the CRUD"))
}
