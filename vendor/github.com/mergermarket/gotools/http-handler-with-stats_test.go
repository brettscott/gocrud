package tools

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func checkMetricsCalled(t *testing.T, statsd *MockStatsD, routeName string, response int, statusCode string) {
	call := statsd.Calls

	expectedTags := []string{"route:" + routeName, "response:" + strconv.Itoa(response)}

	assert.Len(t, call, 3)
	assert.Equal(t, "Histogram", call[0].Method)
	assert.Equal(t, expectedTags, call[0].Args.Tags)
	assert.Equal(t, "Incr", call[1].Method)
	assert.Equal(t, fmt.Sprintf("web.response_code.%s", statusCode), call[1].Args.Name)
	assert.Equal(t, expectedTags, call[1].Args.Tags)
	assert.Equal(t, "Incr", call[2].Method)
	assert.Equal(t, "web.response_code.all", call[2].Args.Name)
	assert.Equal(t, expectedTags, call[2].Args.Tags)
}

type MockHandler struct {
	response int
}

func (h MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(h.response)
}

func TestHTTPHandlerWithStats(t *testing.T) {
	statsd := &MockStatsD{}
	logger := &MockLogger{}
	router := &MockHandler{response: http.StatusOK}
	httpHandler := HTTPHandlerWithStats("route", router, logger, statsd)

	httpHandler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "http://example.com", nil))

	lastLogCall := logger.LastCall()

	assert.NotNil(t, lastLogCall, "Expected call to be made")
	assert.Equal(t, "Debug", lastLogCall.Method)
	assert.Equal(t, "[Request to http://example.com had response code 200]", lastLogCall.Args.Msg)

	checkMetricsCalled(t, statsd, "route", http.StatusOK, "200")

}

func TestHTTPHandlerWithStats_Error(t *testing.T) {
	statsd := &MockStatsD{}
	logger := &MockLogger{}
	router := &MockHandler{response: http.StatusInternalServerError}
	httpHandler := HTTPHandlerWithStats("route", router, logger, statsd)

	httpHandler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "http://example.com", nil))

	lastLogCall := logger.LastCall()

	assert.NotNil(t, lastLogCall, "Expected call to be made")
	assert.Equal(t, "Debug", lastLogCall.Method)
	assert.Equal(t, "[Request to http://example.com had response code 500]", lastLogCall.Args.Msg)

	checkMetricsCalled(t, statsd, "route", http.StatusInternalServerError, "500")
}
