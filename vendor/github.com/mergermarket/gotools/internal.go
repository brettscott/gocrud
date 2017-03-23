package tools

import (
	"fmt"
	"net/http"
)

// InternalHealthCheck is used by our infrastructure to check the service is listening
func InternalHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Healthy")
}

// InternalLogConfig creates an http handler which logs out the app's config
func InternalLogConfig(config interface{}, logger logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("Application config - %+v", config))
		fmt.Fprint(w, "Logged the config")
	}
}
