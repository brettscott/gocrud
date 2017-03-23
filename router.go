package main

import (
	"github.com/mergermarket/gotools"
	"github.com/gorilla/mux"
	"net/http"
)

// newRouter adds handlers to routes
func newRouter(log tools.Logger, statsd tools.StatsD, healthcheckHandler http.Handler, apiRouter *mux.Router) http.Handler {
	router := mux.NewRouter()

	router.Handle("/internal/healthcheck", healthcheckHandler)
	router.Handle("/api", apiRouter)

	return router
}
