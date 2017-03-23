package tools

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type fakeClock struct {
	currentTime time.Time
}

func (f *fakeClock) Now() time.Time {
	f.currentTime = f.currentTime.Add(100 * time.Millisecond)
	return f.currentTime
}

func TestHTTPClientWithStats_Do(t *testing.T) {

	fc := &fakeClock{time.Now()}
	msd := &MockStatsD{}
	hc := http.DefaultClient
	wc := &httpClientWithStats{statsd: msd, httpClient: hc, clock: fc}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	}))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	assert.Len(t, msd.Calls, 0)

	resp, _ := wc.Do(req, "http_callee:my-remote-service", "operation:my-operation")

	expectedTags := []string{"http_callee:my-remote-service", "operation:my-operation", "method:GET", "resp_status:200"}

	assert.Len(t, msd.Calls, 4)
	assert.NotNil(t, resp)
	assert.Equal(t, "Histogram", msd.Calls[0].Method)
	assert.Equal(t, 100.0, msd.Calls[0].Args.Value)
	assert.Equal(t, expectedTags, msd.Calls[0].Args.Tags)
	assert.Equal(t, "Incr", msd.Calls[1].Method)
	assert.Equal(t, "http_client.response_success", msd.Calls[1].Args.Name)
	assert.Equal(t, "Incr", msd.Calls[2].Method)
	assert.Equal(t, "http_client.response_code.200", msd.Calls[2].Args.Name)
	assert.Equal(t, "Incr", msd.Calls[3].Method)
	assert.Equal(t, "http_client.response_code.all", msd.Calls[3].Args.Name)
	assert.Equal(t, expectedTags, msd.Calls[1].Args.Tags)
}

func TestHTTPClientWithStats_Do_Error(t *testing.T) {

	msd := &MockStatsD{}
	hc := http.DefaultClient
	wc := NewHTTPClientWithStats(hc, msd)

	req, _ := http.NewRequest("GET", "http://not-a-domain.zyzyuyziuy", nil)

	assert.Len(t, msd.Calls, 0)

	resp, err := wc.Do(req, "http_callee:my-remote-service", "operation:my-operation")

	assert.Len(t, msd.Calls, 1)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "Incr", msd.Calls[0].Method)
	assert.Equal(t, "http_client.response_error", msd.Calls[0].Args.Name)
	assert.Equal(t, []string{"http_callee:my-remote-service", "operation:my-operation", "method:GET"}, msd.Calls[0].Args.Tags)
}

func TestHTTPClientWithStats_Get(t *testing.T) {

	msd := &MockStatsD{}
	hc := http.DefaultClient
	wc := NewHTTPClientWithStats(hc, msd)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	}))
	defer ts.Close()

	assert.Len(t, msd.Calls, 0)

	resp, _ := wc.Get(ts.URL, "http_callee:my-remote-service", "operation:my-operation")

	assert.Len(t, msd.Calls, 4)
	assert.NotNil(t, resp)
	assert.Equal(t, "Histogram", msd.Calls[0].Method)
	assert.Equal(t, []string{"http_callee:my-remote-service", "operation:my-operation", "method:GET", "resp_status:200"}, msd.Calls[0].Args.Tags)
	assert.Equal(t, "Incr", msd.Calls[2].Method)
	assert.Equal(t, "http_client.response_code.200", msd.Calls[2].Args.Name)
	assert.Equal(t, "Incr", msd.Calls[3].Method)
	assert.Equal(t, "http_client.response_code.all", msd.Calls[3].Args.Name)

}

func TestHttpClientWithStats_Post(t *testing.T) {
	msd := &MockStatsD{}
	hc := http.DefaultClient
	wc := NewHTTPClientWithStats(hc, msd)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	}))
	defer ts.Close()
	assert.Len(t, msd.Calls, 0)

	resp, _ := wc.Post(ts.URL, "application/json", strings.NewReader(`{"hello":"world"}`), "http_callee:my-remote-service", "operation:my-operation")

	assert.Len(t, msd.Calls, 4)
	assert.NotNil(t, resp)
	assert.Equal(t, "Histogram", msd.Calls[0].Method)
	assert.Equal(t, []string{"http_callee:my-remote-service", "operation:my-operation", "method:POST", "resp_status:200"}, msd.Calls[0].Args.Tags)
	assert.Equal(t, "Incr", msd.Calls[2].Method)
	assert.Equal(t, "http_client.response_code.200", msd.Calls[2].Args.Name)
	assert.Equal(t, "Incr", msd.Calls[3].Method)
	assert.Equal(t, "http_client.response_code.all", msd.Calls[3].Args.Name)
}
