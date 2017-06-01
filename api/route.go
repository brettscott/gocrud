package api

import (
	"fmt"
	"github.com/brettscott/gocrud/entity"
	"github.com/pressly/chi"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

// NewRoute prepares the routes for this package
func NewRoute(entities entity.Entities, logger Logger, statsd StatsDer) func(chi.Router) {

	return func(r chi.Router) {

		r.Get("/", root)

		// List
		// eg GET http://localhost:8080/gocrud/api/user
		r.Get("/:entityID", list)

		// Post/Create
		// eg POST http://localhost:8080/gocrud/api/user
		// TODO check content-type header on POST
		// TODO validation
		// TODO persist to DB
		r.Post("/:entityID", post)

		// Put/Update
		// eg PUT http://localhost:8080/gocrud/api/user/1234
		// TODO check content-type header on POST
		// TODO validation
		// TODO persist to DB
		r.Put("/:entityID/:recordID", put)


	}
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the API"))
}

func list(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityID")
	w.Write([]byte(fmt.Sprintf("List - entityID: %s", entityID)))
}

func post(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityID")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Bad request - %v", err)))
		return

	}
	//decoder := json.NewDecoder(r.Body)
	//var t test_struct
	//err := decoder.Decode(&t)
	//if err != nil {
	//	// TODO failure
	//
	//}

	recordID, err := convertPostBodyToRecord(body)
	if err != nil {
		// TODO failure
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to convert JSON - %v", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("Post\nrecordID: %s\nentityID: %s\nbody: %s", recordID, entityID, body)))
}

func convertPostBodyToRecord(body []byte) (string, error) {
	record := entity.Record{}
	err := json.Unmarshal(body, &record)
	if err != nil {
		return "", fmt.Errorf("Unable to unmarshal body: %v", err)
	}

	return record.ID, nil  // Todo remove .ID
}

func put(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityID")
	recordID := chi.URLParam(r, "recordID")
	w.Write([]byte(fmt.Sprintf("Put - entityID: %v, recordID: %v", entityID, recordID)))
}

