package crud

import (
	"github.com/mergermarket/gotools"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"net/http"
	"time"
)

type ChiRouteHandler func(chi.Router)

// newRouter adds handlers to routes
func newRouter(log tools.Logger, statsd tools.StatsD, healthcheckHandlerFunc http.HandlerFunc, apiRouteHandler ChiRouteHandler, uiRouteHandler ChiRouteHandler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Timeout(60 * time.Second))
	//
	//router.Get("/healthcheck", healthcheckHandlerFunc)

	//router.Get("/", func(w http.ResponseWriter, r *http.Request) {
	//	http.Redirect(w, r, "/", http.StatusFound)
	//	return
	//})

	//ui := UIRoute{}
	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	}))

	//healthcheckHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	w.WriteHeader(http.StatusOK)
	//	fmt.Fprint(w, "Healthy")
	//})

	//uiRouteHandler(router)  // mount to the root of this route

	//router.Route("/api", apiRouteHandler)

	return router
}
