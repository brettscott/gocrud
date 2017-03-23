package api

import (
	"github.com/mergermarket/gotools"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
)

// newRouter adds handlers to routes
func NewRouter(log tools.Logger, statsd tools.StatsD) *mux.Router {
	apiRouter := mux.NewRouter()

	apiRouter.HandleFunc("/blah", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Blah is rendered")
	})

	return apiRouter
}