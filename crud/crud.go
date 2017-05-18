package crud

import (
	"fmt"
	"github.com/brettscott/gocrud/api"
	"github.com/brettscott/gocrud/entity"
	"net/http"
)

type Crud struct {
	entities []entity.Entity
	config   *Config
	log      Logger
	statsd   StatsDer
}

// NewCrud creates a new CRUD instance
func NewCrud(config *Config, log Logger, statsd StatsDer) *Crud {
	return &Crud{
		config: config,
		log:    log,
		statsd: statsd,
	}
}

func (c *Crud) AddEntity(entity entity.Entity) {
	c.entities = append(c.entities, entity)
}

func (c *Crud) Handler() http.Handler {
	healthcheckHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Healthy")
	})

	apiRouteHandler := api.NewRoute(c.entities, c.log, c.statsd)

	return newRouter(c.log, c.statsd, healthcheckHandler, apiRouteHandler)
}
