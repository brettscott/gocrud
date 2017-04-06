package main

import (
	"github.com/gorilla/mux"
	"github.com/mergermarket/gotools"
	"net/http"
	"github.com/brettscott/gocrud/api"
)

// newRouter adds handlers to routes
func newRouter(log tools.Logger, statsd tools.StatsD, healthcheckHandler http.Handler) http.Handler {
	router := mux.NewRouter()

	router.Handle("/internal/healthcheck", healthcheckHandler)

	apiRouter := router.PathPrefix("/api").Subrouter()
	router.NewRoute().Handler(api.NewRouter(apiRouter, log, statsd))
	return router
}
