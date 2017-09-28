package api

import (
	"encoding/json"
	"fmt"
	"github.com/brettscott/gocrud/entity"
	"github.com/brettscott/gocrud/store"
	"github.com/pressly/chi"
	"io/ioutil"
	"net/http"
)

const ACTION_POST = "post"
const ACTION_PUT = "put"

type APIRoute struct {
	entities entity.Entities
	store    store.Storer
	log      Logger
	statsd   StatsDer
}

// NewRoute prepares the routes for this package
func NewRoute(entities entity.Entities, store store.Storer, log Logger, statsd StatsDer) func(chi.Router) {

	apiRoute := &APIRoute{
		entities: entities,
		store:    store,
		log:      log,
		statsd:   statsd,
	}

	return func(r chi.Router) {

		r.Get("/", apiRoute.root)

		// List
		// eg GET http://localhost:8080/gocrud/api/user
		// TODO pagination
		r.Get("/:entityID", apiRoute.list)

		// Get record
		// eg GET http://localhost:8080/gocrud/api/user/12345
		r.Get("/:entityID/:recordID", apiRoute.get)

		// Post/Create
		// eg POST http://localhost:8080/gocrud/api/user
		// TODO check content-type header on POST
		r.Post("/:entityID", apiRoute.save(true))

		// Put/Update
		// eg PUT http://localhost:8080/gocrud/api/user/1234
		// TODO check content-type header on PUT
		// TODO validation
		// TODO persist to DB
		r.Put("/:entityID/:recordID", apiRoute.save(false))

	}
}

func (a *APIRoute) root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the API"))
}

func (a *APIRoute) list(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityID")

	if entity, ok := a.entities[entityID]; ok {

		records, err := a.store.List(entity)
		fmt.Printf("\nRecords: %+v", records)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Invalid entityID: %s", entityID)))
			return
		}

		jsonResponse, err := json.Marshal(records)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to convert record to json.  Error: %v", err)))
			return
		}
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
		fmt.Println("Entity: %+v", entity)

		record, err := a.store.Get(entity, recordID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to get entityID: %s, recordID: %s.  Error: %v", entityID, recordID, err)))
			return
		}

		jsonResponse, err := json.Marshal(record)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to convert record to json.  Error: %v", err)))
			return
		}
		w.Write(jsonResponse)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("Invalid entityID: %s", entityID)))
}

func (a *APIRoute) save(isRecordNew bool) func(w http.ResponseWriter, r *http.Request) {
	action := ACTION_PUT
	if isRecordNew {
		action = ACTION_POST
	}

	return func(w http.ResponseWriter, r *http.Request) {
		entityID := chi.URLParam(r, "entityID")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Bad request - %v", err)))
			return
		}

		record := &entity.Record{}
		err = record.UnmarshalJSON(body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Failed to convert JSON - %v", err)))
			return
		}

		entity := a.entities[entityID]
		err = entity.HydrateFromRecord(record, action)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to hydrate Entity from Record - %v", err)))
			return
		}
		err = entity.Validate(action)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Failed validation - %v", err)))
			return
		}

		var recordID string
		if action == ACTION_POST {
			recordID, err = a.store.Post(entity)
		} else {
			recordID = chi.URLParam(r, "recordID")
			if recordID == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("Missing recordID - %v", err)))
			}
			err = a.store.Put(entity, recordID)
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed post/put entity %v.  Error: %v", entity, err)))
			return
		}

		dbRecord, err := a.store.Get(entity, recordID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to get newly created DB record. entityID: %s, recordID: %s.  Error: %v", entityID, recordID, err)))
			return
		}

		jsonResponse, err := json.Marshal(dbRecord)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to convert DB record to json.  Error: %v", err)))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}
}
