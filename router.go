package main

import (
	"github.com/mergermarket/gotools"
	"net/http"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"time"
)

type ChiRouteHandler func(chi.Router)

// newRouter adds handlers to routes
func newRouter(log tools.Logger, statsd tools.StatsD, healthcheckHandlerFunc http.HandlerFunc, apiRouteHandler ChiRouteHandler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/healthcheck", healthcheckHandlerFunc)

	router.Route("/api", apiRouteHandler)

	return router
}
