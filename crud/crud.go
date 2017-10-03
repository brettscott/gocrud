package crud

import (
	"fmt"
	"github.com/brettscott/gocrud/api"
	"github.com/brettscott/gocrud/store"
	"net/http"
)

type Crud struct {
	entities Entities
	config   *Config
	log      Logger
	statsd   StatsDer
	store    store.Storer
}

// NewCrud creates a new CRUD instance
func NewCrud(config *Config, log Logger, statsd StatsDer) *Crud {
	return &Crud{
		config:   config,
		log:      log,
		statsd:   statsd,
		entities: make(map[string]Entity),
	}
}

// Store defines which database to use
func (c *Crud) Store(store store.Storer) {
	c.store = store
}

// AddEntity for each entity type (eg User)
func (c *Crud) AddEntity(entity Entity) {
	c.entities[entity.ID] = entity
}

// Handler for mounting routes for CRUD
func (c *Crud) Handler() http.Handler {
	healthcheckHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Healthy")
	})

	apiRouteHandler := api.NewRoute(c.entities, c.store, c.log, c.statsd)

	return newRouter(c.log, c.statsd, healthcheckHandler, apiRouteHandler)
}
