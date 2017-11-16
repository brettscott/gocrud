package crud

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/pressly/chi"
	"io/ioutil"
	"net/http"
)

const TEMPLATE_PATH string = "crud/templates/%s.hbs"

var templateNames []string = []string{
	"root",
	"list",
	"form",
}

type templateList map[string]*raymond.Template

func NewUiRoute(entities Entities, apiService apiServicer, log Logger, statsd StatsDer) func(chi.Router) {

	uiRoute := &UIRoute{
		entities:   entities,
		log:        log,
		statsd:     statsd,
		apiService: apiService,
		templates:  templates(),
	}

	return func(r chi.Router) {

		// Display entities
		r.Get("/", uiRoute.root)

		// List results for a given entity
		r.Get("/{entityID}", uiRoute.list)

		// List results for a given entity
		r.Get("/{entityID}/create", uiRoute.form(true))

		// List results for a given entity
		r.Get("/{entityID}/{recordID}", uiRoute.form(false))

		// React SPA ??
	}
}

func templates() (tmpls templateList) {
	tmpls = templateList{}

	for _, name := range templateNames {
		filename := fmt.Sprintf(TEMPLATE_PATH, name)
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(fmt.Sprintf("Missing name: %s - %+v", name, err))
		}
		tpl, err := raymond.Parse(string(contents))
		if err != nil {
			panic(err)
		}
		tpl.RegisterHelper("listColumnHeadings", ListColumnHeadings)
		tpl.RegisterHelper("listRows", ListRows)
		tpl.RegisterHelper("listCells", ListCells)
		tmpls[name] = tpl
	}

	return tmpls
}

type UIRoute struct {
	entities   Entities
	log        Logger
	statsd     StatsDer
	apiService apiServicer
	templates  templateList
}

func (u *UIRoute) root(w http.ResponseWriter, r *http.Request) {
	ctx := map[string]interface{}{
		"entities": u.entities,
	}
	html, err := u.templates["root"].Exec(ctx)
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
	html, err := u.templates["list"].Exec(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
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

		// Get record from MongoDB

		// Convert record into "row" struct
		row := row{}
		//for _, field := range records {
		//	record[field.ID] = field.Value
		//}
		row["id"] = "blah"

		ctx := map[string]interface{}{
			"entity":   entity,
			"recordID": recordID,
			"row":      row,
		}
		html, err := u.templates["form"].Exec(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))

	}
}
