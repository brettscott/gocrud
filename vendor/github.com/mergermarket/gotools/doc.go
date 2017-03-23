/*

This package contains a few simple tools for building Go web applications.

InternalHealthCheck handles healthcheck requests from load balancers:

	serveMux.HandleFunc("/internal/healthcheck", tools.InternalHealthCheck)

InternalLogConfig dumps your application's configuration to your logger:

	serveMux.HandleFunc("/internal/log-config", tools.InternalLogConfig(config, log))

Logger logs messages in a structured format in AWS (for forwarding to LogEntries) and pretty colours in local:

	logger := tools.NewLogger(config.IsLocal())
	logger.Info("Hello!")

StatsD sends metrics to DataDog:

	statsd, err := tools.NewStatsD(tools.NewStatsDConfig(!config.IsLocal(), logger))
	if err != nil {
		logger.Error("Error connecting to StatsD. Stats will only be logged. Error: ", err.Error())
	}
	statsd.Histogram("important_action", .0001, "tag1:tag1value", "tag2:tag2value")


HTTPHandlerWithStats takes an http.Handler and adds the sending of response time metrics to DataDog, and debug logging of request details:

  	serveMux.Handle("/my-important-endpoint", importantHandler)
  	tools.HTTPHandlerWithStats("/", serveMux, logger, statsd)


HTTPClientWithStats takes an http.Client and adds the sending of response time metrics to DataDog:

	httpClient := tools.NewHTTPClientWithStats(http.DefaultClient, statsd, "foo-api")
	resp, err := httpClient.Get("http://foo-api.com/important-data")

*/
package tools
