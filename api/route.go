package api

import (
	"encoding/json"
	"fmt"
	"github.com/brettscott/gocrud/entity"
	"github.com/pressly/chi"
	"io/ioutil"
	"net/http"
	"github.com/brettscott/gocrud/store"
)

const ACTION_POST = "post"

type APIRoute struct {
	entities entity.Entities
	store store.Storer
	log      Logger
	statsd   StatsDer
}

// NewRoute prepares the routes for this package
func NewRoute(entities entity.Entities, store store.Storer, log Logger, statsd StatsDer) func(chi.Router) {

	apiRoute := &APIRoute{
		entities: entities,
		store: store,
		log:      log,
		statsd:   statsd,
	}

	return func(r chi.Router) {

		r.Get("/", apiRoute.root)

		// List
		// eg GET http://localhost:8080/gocrud/api/user
		r.Get("/:entityID", apiRoute.list)

		// Get record
		// eg GET http://localhost:8080/gocrud/api/user/12345
		// TODO check content-type header on POST
		// TODO validation
		// TODO get DB
		r.Get("/:entityID/:recordID", apiRoute.get)

		// Post/Create
		// eg POST http://localhost:8080/gocrud/api/user
		// TODO check content-type header on POST
		// TODO validation
		// TODO persist to DB
		r.Post("/:entityID", apiRoute.post)

		// Put/Update
		// eg PUT http://localhost:8080/gocrud/api/user/1234
		// TODO check content-type header on POST
		// TODO validation
		// TODO persist to DB
		r.Put("/:entityID/:recordID", apiRoute.put)

	}
}

func (a *APIRoute) root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the API"))
}

func (a *APIRoute) list(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityID")
	w.Write([]byte(fmt.Sprintf("List - entityID: %s", entityID)))
}

func (a *APIRoute) post(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityID")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Bad request - %v", err)))
		return
	}

	record, err := marshalBodyToRecord(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed to convert JSON - %v", err)))
		return
	}

	entity := a.entities[entityID]
	fmt.Println("Entity: %+v", entity)

	err = entity.HydrateFromRecord(record, ACTION_POST)
	fmt.Println("Hydrated Entity: %+v", entity)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to hydrate Entity from Record - %v", err)))
		return
	}
	err = entity.Validate(ACTION_POST)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed validation - %v", err)))
		return
	}

	dbID, err := a.store.Post(entity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed post entity %v.  Error: %v", entity, err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("Post\nrecordID: %s\nentityID: %s\nbody: %s\nentity: %+v\ndbID: %s\n", record.ID, entityID, body, entity, dbID)))
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

		w.Write([]byte(fmt.Sprintf("Post\nrecordID: %s\nentityID: %s\nrecord: %+v", recordID, entityID, record)))
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("Invalid entityID %s", entityID)))
}

// marshalBodyToRecord converts JSON to entity.Record
func marshalBodyToRecord(body []byte) (*entity.Record, error) {
	record := entity.Record{}
	err := json.Unmarshal(body, &record)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal body: %v", err)
	}
	return &record, nil
}

func (a *APIRoute) put(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityID")
	recordID := chi.URLParam(r, "recordID")
	w.Write([]byte(fmt.Sprintf("Put - entityID: %v, recordID: %v", entityID, recordID)))
}
