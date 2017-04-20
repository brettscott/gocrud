package main

import (
	"github.com/mergermarket/gotools"
	"net/http"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"time"
)

//type CreateHandlerWithPrefix func(string) http.Handler

type ChiRouteHandler func(chi.Router)

// newRouter adds handlers to routes
func newRouter(logger tools.Logger, statsd tools.StatsD, healthcheckHandler http.HandlerFunc, apiRouteHandler ChiRouteHandler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/_gocrud/healthcheck", healthcheckHandler)

	router.Route("/_gocrud/api", apiRouteHandler)

	return router
}
