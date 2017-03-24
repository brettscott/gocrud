package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mergermarket/gotools"
	"net/http"
)

// NewRouter adds handlers to routes
func NewRouter(log tools.Logger, statsd tools.StatsD) *mux.Router {
	apiRouter := mux.NewRouter()

	apiRouter.HandleFunc("/blah", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Blah is rendered")
	})

	return apiRouter
}
