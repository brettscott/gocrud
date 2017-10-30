package crud

import (
	"encoding/json"
	"fmt"
	"github.com/pressly/chi"
	"io/ioutil"
	"net/http"
)

const ACTION_POST = "post"
const ACTION_PUT = "put"
const ACTION_PATCH = "patch"

type APIRoute struct {
	entities   Entities
	log        Logger
	statsd     StatsDer
	apiService apiServicer
}

type apiServicer interface {
	list(entity *Entity) (clientRecords ClientRecords, err error)
	get(entity *Entity, recordID string) (clientRecord ClientRecord, err error)
	save(entity *Entity, action string, clientRecord *ClientRecord, recordID string) (savedClientRecord ClientRecord, clientErrors *ClientErrors, err error)
	delete(entity *Entity, recordID string) error
}

// NewRoute prepares the routes for this package
func NewApiRoute(entities Entities, apiService apiServicer, log Logger, statsd StatsDer) func(chi.Router) {

	apiRoute := &APIRoute{
		entities:   entities,
		log:        log,
		statsd:     statsd,
		apiService: apiService,
	}

	return func(r chi.Router) {

		r.Get("/", apiRoute.root)

		// List
		// eg GET http://localhost:8080/gocrud/api/user
		r.Get("/:entityID", apiRoute.list)

		// Get record
		// eg GET http://localhost:8080/gocrud/api/user/12345
		r.Get("/:entityID/:recordID", apiRoute.get)

		// Post/Create
		// eg POST http://localhost:8080/gocrud/api/user
		r.Post("/:entityID", apiRoute.save(true, false))

		// Put/Update
		// eg PUT http://localhost:8080/gocrud/api/user/1234
		r.Put("/:entityID/:recordID", apiRoute.save(false, false))

		// Patch/Update partial
		// eg PATCH http://localhost:8080/gocrud/api/user/1234
		r.Patch("/:entityID/:recordID", apiRoute.save(false, true))

		// Delete
		// eg DELETE http://localhost:8080/gocrud/api/user/1234
		r.Delete("/:entityID/:recordID", apiRoute.delete)
	}
}

func (a *APIRoute) root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the API"))
}

func (a *APIRoute) list(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityID")

	if entity, ok := a.entities[entityID]; ok {
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
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("Invalid entityID: %s", entityID)))
}

// get returns a record from the database for the given recordID in given entityID
func (a *APIRoute) get(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityID")
	recordID := chi.URLParam(r, "recordID")

	if entity, ok := a.entities[entityID]; ok {
		record, err := a.apiService.get(entity, recordID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		jsonResponse, err := json.Marshal(record)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("Invalid entityID: %s", entityID)))
}

func (a *APIRoute) save(isRecordNew bool, isPartialPayload bool) func(w http.ResponseWriter, r *http.Request) {
	action := ACTION_PUT
	if isRecordNew && !isPartialPayload {
		action = ACTION_POST
	} else if isRecordNew && isPartialPayload {
		action = ACTION_PATCH
	}

	return func(w http.ResponseWriter, r *http.Request) {
		entityID := chi.URLParam(r, "entityID")
		recordID := chi.URLParam(r, "recordID")

		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`Bad request - missing Content Type header "application/json"`))
			return
		}

		if (action == ACTION_PUT || action == ACTION_PATCH) && len(recordID) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad request - missing recordID"))
			return
		}

		if entity, ok := a.entities[entityID]; ok {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("Bad request - %v", err)))
				return
			}
			record := &ClientRecord{}
			err = json.Unmarshal(body, record)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			savedRecord, clientErrors, err := a.apiService.save(entity, action, record, recordID)
			if err != nil {
				a.log.Error(err)
				fmt.Printf("api_route: do something with clientErrors %+v", clientErrors) // todo something with this
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			jsonResponse, err := json.Marshal(savedRecord)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			if action == ACTION_POST {
				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusOK)
			}
			w.Write(jsonResponse)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Invalid entityID: %s", entityID)))
	}
}

// delete removes a record from the database
func (a *APIRoute) delete(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityID")
	recordID := chi.URLParam(r, "recordID")

	if entity, ok := a.entities[entityID]; ok {
		err := a.apiService.delete(entity, recordID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("Invalid entityID: %s", entityID)))
}
