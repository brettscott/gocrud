package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mergermarket/gotools"
	"net/http"
)

// NewRouter adds handlers to routes
func NewRouter(router *mux.Router, log tools.Logger, statsd tools.StatsD) *mux.Router {
	//router := mux.NewRouter()

	log.Info("inside api_router")

	router.HandleFunc("/here", func(w http.ResponseWriter, r *http.Request) {
		log.Info("/here/ hit")
		fmt.Fprint(w, "API HERE")
	}).Name("API HERE")

	router.HandleFunc("/xxx/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("/api/ hit")
		fmt.Fprint(w, "API root")
	}).Name("API Root")

	router.HandleFunc("/yyy/blah", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Blah is rendered")
	}).Name("API BLAH")

	return router
}
