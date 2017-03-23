package tools

import (
	"fmt"
	"github.com/felixge/httpsnoop"
	"net/http"
)

type logger interface {
	Debug(...interface{})
	Info(...interface{})
	Error(...interface{})
}

// HTTPHandlerWithStats takes an http.Handler and adds the sending of response time metrics to DataDog, and debug logging of request details
func HTTPHandlerWithStats(routeName string, router http.Handler, logger logger, statsd StatsD) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(r.Method, "at", r.URL.String())
		metrics := httpsnoop.CaptureMetrics(router, w, r)

		logResult(routeName, metrics, statsd, logger, r)
	})
}

func logResult(routeName string, metrics httpsnoop.Metrics, statsd StatsD, logger logger, req *http.Request) {
	responseTag := fmt.Sprintf("response:%d", metrics.Code)
	tags := []string{"route:" + routeName, responseTag}
	if caller := req.Header.Get("X-Component"); caller != "" {
		tags = append(tags, "caller:"+caller)
	}
	statsd.Histogram(WebResponseTimeKey, float64(metrics.Duration.Nanoseconds())/1000000, tags...)
	statsd.Incr(fmt.Sprintf(WebResponseCodeFormatKey, metrics.Code), tags...)
	statsd.Incr(WebResponseCodeAllKey, tags...)
	logger.Debug("Request to ", req.URL.String(), " had response code ", metrics.Code)
}
