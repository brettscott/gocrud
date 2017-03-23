[![Build Status](https://travis-ci.org/mergermarket/gotools.svg?branch=master)](https://travis-ci.org/mergermarket/gotools)

# gotools

Gotools provides a few simple building blocks for Mergermarket apps.

## internal

Internal endpoints for http services. 

Example usage:

```
	router := http.NewServeMux()
	router.HandleFunc("/internal/healthcheck", tools.InternalHealthCheck)
	router.HandleFunc("/internal/log-config", tools.NewInternalLogConfig(config, log))

```

## logger

Logger logs messages in a structured format in AWS and pretty colours in local.

Example usage:

```
  logger := tools.NewLogger(config.IsLocal())
	logger.Info("Hello!")

```

## statsd

StatsD provides an interface for all of the DataDog StatsD methods used by mergermarket. 

Example usage:

```
	statsdConfig := tools.NewStatsDConfig(!config.IsLocal(), logger)
	statsd, err := tools.NewStatsD(statsdConfig)
	if err != nil {
		logger.Error("Error connecting to StatsD. Stats will only be logged. Error: ", err.Error())
	}
  statsd.Histogram("important_action", .0001, "tag1:tag1value", "tag2:tag2value")

```

## http-handler-with-stats

HTTPHandlerWithStats takes an http.Handler and adds the sending of response time metrics to DataDog, and debug logging of request details

Example usage:

```
  router := http.NewServeMux()
  helloHandlerPattern := "/hello"
  helloHandlerWithStats := tools.HTTPHandlerWithStats(helloHandlerPattern, helloHandler, log, statsd)
  router.Handle(helloHandlerPattern, helloHandlerWithStats)
```

## http-client-with-stats

HTTPClientWithStats takes an http.Client and adds the sending of metrics to DataDog.
Optionally pass in tags with the requests.

Example usage:

```
    logger := tools.NewLogger(config.IsLocal())
	statsd, _ := tools.NewStatsD(tools.NewStatsDConfig(!config.IsLocal(), logger))
	httpClient := NewHTTPClientWithStats(http.DefaultClient, statsd)
    resp, err := httpClient.Get(ts.URL, "callee:my-remote-service", "operation:getstuff")
```

## test tools

```
testLogger, testStatsd := tools.NewTestTools(t)
```

## httputil.ValidateParamsHandler

Example:

```
type MainHandler struct{}

func (mh MainHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Hello world\n")
}

func main() {
	http.Handle("/", httputil.ValidateParamsHandler(MainHandler{}, "X-Component", "X-User-ID"))
	http.ListenAndServe(":8080", nil)
}
```

## Contributing

TODO
