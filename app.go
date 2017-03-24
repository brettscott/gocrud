package main

import (
	"fmt"
	"github.com/brettscott/gocrud/api"
	"github.com/mergermarket/gotools"
	"log"
	"net/http"
	"os"
)

func main() {
	config, logger, statsd := toolup()

	healthcheckHandler := http.HandlerFunc(tools.InternalHealthCheck)

	//func InternalHealthCheck(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprint(w, "Healthy")
	//}

	//apiGateway := api.NewGateway()
	apiRouter := api.NewRouter(logger, statsd)

	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), newRouter(logger, statsd, healthcheckHandler, apiRouter))
	if err != nil {
		logger.Error("Problem starting server", err.Error())
		os.Exit(1)
	}
}

func toolup() (*appConfig, tools.Logger, tools.StatsD) {
	config, err := loadAppConfig()
	if err != nil {
		log.Fatal("Error loading config", err.Error())
	}

	logger := tools.NewLogger(config.IsLocal())
	logger.Info(fmt.Sprintf("Application config - %+v", config))

	statsdConfig := tools.NewStatsDConfig(!config.IsLocal(), logger)
	statsd, err := tools.NewStatsD(statsdConfig)
	if err != nil {
		logger.Error("Error connecting to StatsD - defaulting to logging stats: ", err.Error())
	}

	return config, logger, statsd
}
