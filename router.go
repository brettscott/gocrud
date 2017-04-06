package main

import (
	"github.com/gorilla/mux"
	"github.com/mergermarket/gotools"
	"net/http"
	"github.com/brettscott/gocrud/api"
)

//type CreateHandlerWithPrefix func(string) http.Handler

// newRouter adds handlers to routes
func newRouter(logger tools.Logger, statsd tools.StatsD, healthcheckHandler http.Handler, createAPIHandler api.CreateHandlerWithPrefix) http.Handler {
	router := mux.NewRouter()

	router.Handle("/internal/healthcheck", healthcheckHandler)

	router.NewRoute().Handler(createAPIHandler("/api"))

	return router
}
