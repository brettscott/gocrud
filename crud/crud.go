package crud

import (
	"fmt"
	"github.com/brettscott/gocrud/model"
	"github.com/brettscott/gocrud/store"
	"net/http"
)

type Crud struct {
	entities   model.Entities
	config     *Config
	log        Logger
	statsd     StatsDer
	store      store.Storer
	apiService apiService
}

// NewCrud creates a new CRUD instance
func NewCrud(config *Config, log Logger, statsd StatsDer) *Crud {
	return &Crud{
		config:   config,
		log:      log,
		statsd:   statsd,
		entities: make(map[string]model.Entity),
	}
}

// Store defines which database to use
func (c *Crud) Store(store store.Storer) {
	c.store = store
}

// AddEntity for each entity type (eg User)
func (c *Crud) AddEntity(entity model.Entity) {
	c.entities[entity.ID] = entity
}

// Handler for mounting routes for CRUD
func (c *Crud) Handler() http.Handler {

	elementsValidator := NewElementsValidator()
	c.apiService = newApiService(c.store, elementsValidator) // TODO  change to &c.store

	healthcheckHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Healthy")
	})

	apiRouteHandler := NewApiRoute(c.entities, &c.apiService, c.log, c.statsd)

	return newRouter(c.log, c.statsd, healthcheckHandler, apiRouteHandler)
}
