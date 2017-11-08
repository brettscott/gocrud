package crud

import (
	"fmt"
	"github.com/pressly/chi"
	"io/ioutil"
	"net/http"
)

func NewUiRoute(entities Entities, apiService apiServicer, log Logger, statsd StatsDer) func(chi.Router) {

	uiRoute := &UIRoute{
		entities:   entities,
		log:        log,
		statsd:     statsd,
		apiService: apiService,
	}

	templates := []string{
		"root",
	}

	templateContents := map[string][]byte{}

	for _, template := range templates {
		templateContents["root"] = []byte{}
		filename := fmt.Sprintf("crud/templates/%s.hbs", template)
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(fmt.Sprintf("Missing template: %s - %+v", template, err))
		}
		templateContents["root"] = contents
	}
	uiRoute.templates = templateContents

	log.Info(templateContents)

	return func(r chi.Router) {

		// Display entities
		r.Get("/", uiRoute.root)

		// List results for a given entity
		//r.Get("/:entityID", uiRoute.list)

		// React???
	}
}

type rootContext struct {
	entities Entities
}

type UIRoute struct {
	entities   Entities
	log        Logger
	statsd     StatsDer
	apiService apiServicer
	templates  map[string][]byte
}

func (u *UIRoute) root(w http.ResponseWriter, r *http.Request) {

	fmt.Println("ui root called")
	//result, err := raymond.Render(tpl, ctx)

	w.Write([]byte("Hello"))

	//tpl :=
	//ctx := rootContext{
	//	entities: u.entities,
	//}
	//result, err := raymond.Render(tpl, ctx)
	//
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	//jsonResponse, err := json.Marshal(records)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	//w.WriteHeader(http.StatusOK)
	//w.Write(jsonResponse)
	return

}
