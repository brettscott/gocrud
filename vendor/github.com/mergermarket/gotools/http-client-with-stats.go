package tools

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClientWithStats takes an http.Client and adds the sending of response time metrics to DataDog
type HTTPClientWithStats interface {
	Do(r *http.Request, tags ...string) (*http.Response, error)
	Get(url string, tags ...string) (*http.Response, error)
	Post(url string, bodyType string, body io.Reader, tags ...string) (*http.Response, error)
}

type clock interface {
	Now() time.Time
}

type httpClientWithStats struct {
	httpClient *http.Client
	statsd     StatsD
	clock      clock
}

func (thc *httpClientWithStats) Do(r *http.Request, tags ...string) (*http.Response, error) {
	tags = append(tags, fmt.Sprintf("method:%s", r.Method))
	start := thc.clock.Now()
	resp, err := thc.httpClient.Do(r)
	if err != nil {
		thc.statsd.Incr(HttpClientResponseErrorKey, tags...)
	} else {
		finish := thc.clock.Now()
		duration := (finish.Nanosecond() - start.Nanosecond()) / 1000000
		tags = append(tags, fmt.Sprintf("resp_status:%d", resp.StatusCode))
		thc.statsd.Histogram(HttpClientResponseTimeKey, float64(duration), tags...)
		thc.statsd.Incr(HttpClientResponseSuccessKey, tags...)
		responseCodeKey := fmt.Sprintf(HttpClientResponseCodeFormatKey, resp.StatusCode)
		thc.statsd.Incr(responseCodeKey, tags...)
		thc.statsd.Incr(HttpClientResponseCodeAllKey, tags...)
	}
	return resp, err
}

func (thc *httpClientWithStats) Get(url string, tags ...string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return resp, err
	}
	return thc.Do(req, tags...)
}

func (thc *httpClientWithStats) Post(url string, bodyType string, body io.Reader, tags ...string) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return resp, err
	}
	req.Header.Set("Content-Type", bodyType)
	return thc.Do(req, tags...)
}

type timeClock struct{}

func (c *timeClock) Now() time.Time {
	return time.Now()
}

func NewHTTPClientWithStats(client *http.Client, statsd StatsD) HTTPClientWithStats {
	return &httpClientWithStats{statsd: statsd, httpClient: client, clock: &timeClock{}}
}
