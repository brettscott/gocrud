package crud

import (
	"fmt"
	"github.com/pressly/chi"
	"net/http"
	"net/url"
)

type templateServicer interface {
	exec(tplName string, ctx map[string]interface{}) (html string, err error)
}

func NewUiRoute(entities Entities, apiService apiServicer, templateService templateServicer, log Logger, statsd StatsDer) func(chi.Router) {

	registerTemplateHelpers()

	uiRoute := &UIRoute{
		entities:        entities,
		log:             log,
		statsd:          statsd,
		apiService:      apiService,
		templateService: templateService,
	}

	return func(r chi.Router) {

		// Display entities
		r.Get("/", uiRoute.root)

		// List results for a given entity
		r.Get("/{entityID}", uiRoute.list)

		// Create a record
		r.Get("/{entityID}/create", uiRoute.form(true))

		// View record
		// TODO create route
		r.Get("/{entityID}/{recordID}/view", uiRoute.view)

		// Edit record
		r.Get("/{entityID}/{recordID}/edit", uiRoute.form(false))

		// Delete record
		// TODO create route
		r.Get("/{entityID}/{recordID}/delete", uiRoute.delete)

		// Save record (triggered by form submit)
		r.Post("/{entityID}/save", uiRoute.save(ACTION_POST))
		r.Post("/{entityID}/{recordID}/save", uiRoute.save(ACTION_PUT))

		// React SPA ??
	}
}

type UIRoute struct {
	entities        Entities
	log             Logger
	statsd          StatsDer
	apiService      apiServicer
	templateService templateServicer
}

func (u *UIRoute) root(w http.ResponseWriter, r *http.Request) {
	ctx := map[string]interface{}{
		"entities": u.entities,
	}
	html, err := u.templateService.exec("root", ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func marshalClientRecordsToRows(entity *Entity, clientRecords []ClientRecord) (rows []row, err error) {
	rows = []row{}

	for _, clientRecord := range clientRecords {
		row, err := marshalClientRecordToRow(entity, clientRecord)
		if err != nil {
			return rows, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func marshalClientRecordToRow(entity *Entity, clientRecord ClientRecord) (row, error) {
	r := row{}
	for _, element := range entity.Elements {
		val, err := clientRecord.GetValue(element.ID)
		if err != nil {
			return r, nil
		}
		r[element.ID] = val
	}
	return r, nil
}

func (u *UIRoute) list(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityID")
	entity, ok := u.entities[entityID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unknown entity: %s", entityID)))
		return
	}

	clientRecords, err := u.apiService.list(entity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Printf("clientRecords: %+v\n\n", clientRecords)
	rows, err := marshalClientRecordsToRows(entity, clientRecords)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ctx := map[string]interface{}{
		"entity": entity,
		"rows":   rows,
	}
	fmt.Printf("ctx: %+v", ctx)
	html, err := u.templateService.exec("list", ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

type ElementValue struct {
	Element
	Value interface{}
}

func marshalClientRecordToElementValues(entity *Entity, clientRecord ClientRecord) []ElementValue {
	var evs []ElementValue
	for _, element := range entity.Elements {
		ev := ElementValue{
			Element: element,
		}
		val, err := clientRecord.GetValue(element.ID)
		if err == nil {
			ev.Value = val
		}
		evs = append(evs, ev)
	}
	return evs
}

func (u *UIRoute) form(create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entityID := chi.URLParam(r, "entityID")
		recordID := chi.URLParam(r, "recordID")
		entity, ok := u.entities[entityID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Unknown entity: %s", entityID)))
			return
		}

		clientRecord := ClientRecord{}
		if !create {
			var err error
			clientRecord, err = u.apiService.get(entity, recordID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}
		evs := marshalClientRecordToElementValues(entity, clientRecord)

		ctx := map[string]interface{}{
			"create":        create,
			"entity":        entity,
			"recordID":      recordID,
			"elementValues": evs,
		}
		html, err := u.templateService.exec("form", ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}
}

func (u *UIRoute) view(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("TODO"))
}

func (u *UIRoute) delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("TODO"))
}

func (u *UIRoute) save(action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		entityID := chi.URLParam(r, "entityID")
		recordID := chi.URLParam(r, "recordID")
		entity, ok := u.entities[entityID]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Unknown entity: %s", entityID)))
			return
		}
		if action == ACTION_PUT && len(recordID) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Missing record ID: %s", recordID)))
			return
		}
		fmt.Printf("Form: %+v, Body: %+v", r.Form, r.Body)
		clientRecord := marshalFormToClientRecord(entity, r.Form)

		fmt.Printf("clientRecord: %+v", clientRecord)

		savedClientRecord, clientErrors, err := u.apiService.save(entity, action, clientRecord, recordID)
		if clientErrors.HasErrors() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Client Errors - elements: %+v, global: %+v", clientErrors.ElementsErrors, clientErrors.GlobalErrors)))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		fmt.Printf("SUCCESS: %+v", savedClientRecord)
	}
}

func marshalFormToClientRecord(entity *Entity, formValues url.Values) *ClientRecord {
	clientRecord := ClientRecord{}
	for _, element := range entity.Elements {
		val, ok := formValues[element.ID]
		if !ok {
			continue
		}
		var v interface{}
		if element.DataType != ELEMENT_DATA_TYPE_ARRAY {
			v = val[0]
		} else {
			v = val
		}
		kv := KeyValue{
			Key:   element.ID,
			Value: v,
		}
		clientRecord.KeyValues = append(clientRecord.KeyValues, kv)
	}
	return &clientRecord
}
