package main

import (
	"github.com/brettscott/gocrud/entity"
	"net/http"
	"fmt"
	"github.com/brettscott/gocrud/api"
)

type crud struct {
	entities []entity.Entity
	log Logger
	statsd StatsDer
}

func NewCrud() *crud {
	return &crud{}
}

func (c *crud) AddEntity(entity entity.Entity) {
	c.entities = append(c.entities, entity)
}

func (c *crud) Handler() http.Handler {
	healthcheckHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Healthy")
	})

	apiRouteHandler := api.NewRoute(c.log, c.statsd)

	return newRouter(c.log, c.statsd, healthcheckHandler, apiRouteHandler)
}

// Logger interface for logging
type Logger interface {
	Info(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
	Warn(msg ...interface{})
}

// StatsD is the interface for the all the DataDog StatsD methods
type StatsDer interface {
	Histogram(name string, value float64, tags ...string)
	Gauge(name string, value float64, tags ...string)
	Incr(name string, tags ...string)
}