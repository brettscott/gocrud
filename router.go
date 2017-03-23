package main

import (
	"github.com/mergermarket/gotools"
	"net/http"
)

// newRouter adds handlers to routes
func newRouter(log tools.Logger, statsd tools.StatsD, healthcheckHandler, apiGateway http.Handler) http.Handler {
	router := http.NewServeMux()

	router.Handle("/internal/healthcheck", healthcheckHandler)

	router.Handle("/api", apiGateway)
	return router
}
