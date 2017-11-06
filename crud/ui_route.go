package crud

import (
	"github.com/pressly/chi"
	"net/http"
	"encoding/json"
	"fmt"
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

		// Display entities
		r.Get("/", uiRoute.root)

		// List results for a given entity
		//r.Get("/:entityID", uiRoute.list)

		// React???
	}
}

func (u *UIRoute) root(w http.ResponseWriter, r *http.Request) {

	for _, entity := range u.entities {

	}
	records, err := a.apiService.list(entity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	jsonResponse, err := json.Marshal(records)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	return

}
