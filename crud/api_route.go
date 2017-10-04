package crud

import (
	"encoding/json"
	"fmt"
	"github.com/brettscott/gocrud/model"
	"github.com/brettscott/gocrud/store"
	"github.com/pressly/chi"
	"io/ioutil"
	"net/http"
	"reflect"
)

const ACTION_POST = "post"
const ACTION_PUT = "put"
const ACTION_PATCH = "patch"

type APIRoute struct {
	entities model.Entities
	store    store.Storer
	log      Logger
	statsd   StatsDer
	apiService apiService
}

// NewRoute prepares the routes for this package
func NewApiRoute(entities model.Entities, store store.Storer, apiService apiService, log Logger, statsd StatsDer) func(chi.Router) {

	apiRoute := &APIRoute{
		entities: entities,
		store:    store,
		log:      log,
		statsd:   statsd,
		apiService: apiService,
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
		r.Post("/:entityID", apiRoute.save(true, false))

		// Put/Update
		// eg PUT http://localhost:8080/gocrud/api/user/1234
		// TODO check content-type header on PUT
		// TODO validation
		r.Put("/:entityID/:recordID", apiRoute.save(false, false))

		// Patch/Update partial
		// eg PATCH http://localhost:8080/gocrud/api/user/1234
		// TODO check content-type header on PUT
		// TODO validation
		r.Patch("/:entityID/:recordID", apiRoute.save(false, true))

	}
}

func (a *APIRoute) root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the API"))
}

func (a *APIRoute) list(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityID")

	if entity, ok := a.entities[entityID]; ok {
		jsonResponse, err := a.apiService.list(entity)
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
		jsonResponse, err := a.apiService.get(entity, recordID)
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Bad request - %v", err)))
			return
		}

		record := &Record{}
		err = record.UnmarshalJSON(body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Failed to convert JSON - %v", err)))
			return
		}

		entity := a.entities[entityID]
		entityData, err := marshalRecordToEntityData(entity, record, action)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to hydrate Entity from ClientRecord - %v", err)))
			return
		}
		err = validate(entity, entityData, action)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Failed validation - %v", err)))
			return
		}

		var recordID string
		switch action {
		case ACTION_POST:
			recordID, err = a.store.Post(entity, entityData)
			break
		case ACTION_PUT:
			recordID = chi.URLParam(r, "recordID")
			if recordID == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("Missing recordID - %v", err)))
			}
			err = a.store.Put(entity, entityData, recordID)
			break
		case ACTION_PATCH:
			recordID = chi.URLParam(r, "recordID")
			if recordID == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("Missing recordID - %v", err)))
			}
			err = a.store.Patch(entity, entityData, recordID)
			break
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Invalid action - %s", action)))
			break
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed post/put e %v.  Error: %v", entity, err)))
			return
		}

		storeRecord, err := a.store.Get(entity, recordID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to get newly created DB record. entityID: %s, recordID: %s.  Error: %v", entityID, recordID, err)))
			return
		}

		clientRecord := marshalStoreRecordToClientRecord(storeRecord)
		jsonResponse, err := json.Marshal(clientRecord)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to convert DB record to json.  Error: %v", err)))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}
}


func validate(entity model.Entity, record store.Record, action string) error {

	errors := make([]string, 0)
	var primaryKey model.ElementLabel

	for _, element := range entity.Elements {

		userData, err := record.GetField(element.ID)
		if err != nil {
			errors = append(errors, fmt.Sprintf(`Missing element "%s" - %v`, element.ID, err))
		}

		if err := validateDataType(element, userData.Value); err != nil {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) has invalid data type: %s`, element.Label, element.ID, err))
		}

		// This is useful to see if value was provided and whether a string is empty or not.  Use "Min" and "Max" for integers.
		// Don't use anything for boolean because it'll either be true or false (or "nil" and be classed as not provided).
		if element.Validation.Required && (userData.Hydrated == false || userData.Value == nil || userData.Value == "") {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) is required and cannot be empty`, element.Label, element.ID))
		}

		if element.Validation.MustProvide == true && userData.Hydrated == false {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) must be provided`, element.Label, element.ID))
		}

		if element.PrimaryKey == true {
			if primaryKey != "" {
				errors = append(errors, fmt.Sprintf(`"%s" (%s) cannot be a primary key because "%s" is already one`, element.Label, element.ID, primaryKey))
			} else {
				primaryKey = element.Label
			}
		}

		if action != ACTION_PATCH && element.PrimaryKey != true && userData.Hydrated == false {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) was not supplied on "%s"`, element.Label, element.ID, action))
		}
	}

	if primaryKey == "" {
		errors = append(errors, fmt.Sprintf(`Missing a primary key element`))
	}

	if len(errors) > 0 {
		return fmt.Errorf("Validation errors: %v", errors)
	}

	return nil
}

// validateDataType
// Unmarshal stores one of these in the interface value: "bool" for JSON booleans, "float64" for JSON numbers,
// "string" for JSON strings, "[]interface{}" for JSON arrays, "map[string]interface{}" for JSON objects,  "nil" for JSON null
func validateDataType(element model.Element, value interface{}) error {
	if value == nil {
		return nil
	}

	// Todo Move out of here so it's only created once!
	dataTypes := make(map[string]string)
	dataTypes[model.ELEMENT_DATA_TYPE_STRING] = "string"
	dataTypes[model.ELEMENT_DATA_TYPE_NUMBER] = "float64"
	dataTypes[model.ELEMENT_DATA_TYPE_BOOLEAN] = "bool"

	if _, ok := dataTypes[element.DataType]; !ok {
		return fmt.Errorf(`undefined data type "%s"`, element.DataType)
	}

	actualType := reflect.TypeOf(value).String()
	expectedType := dataTypes[element.DataType]
	if actualType != expectedType {
		return fmt.Errorf(`expected type to be "%s" but got "%s" with value "%v"`, expectedType, actualType, value)
	}

	return nil
}
