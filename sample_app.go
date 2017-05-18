package main

import (
	"github.com/brettscott/gocrud/entity"
	"fmt"
	"os"
	"net/http"
	"github.com/mergermarket/gotools"
	"log"
	"github.com/pressly/chi"
)

func main() {
	config, log, statsd := toolup()

	// TODO: Define schema
	// TODO: Build database connector - MySQL, Mongo
	// TODO: Pre/post hooks and override actions
	// TODO: http.ListenAndServe(<port>) - return a route/router for sample app to "listen and serve"
	// TODO: Flexibility with rendering templates (custom head/foot/style)

	users := entity.Entity{
		ID: "users",
		Label: "User",
		Labels: "Users",
		Elements: []entity.Element{
			{
				ID: "name",
				Label: "Name",
				FormType: entity.ELEMENT_FORM_TYPE_TEXT,
				ValueType: entity.ELEMENT_VALUE_TYPE_STRING,
				Value: "",
			},
			{
				ID: "age",
				Label: "Age",
				FormType: entity.ELEMENT_FORM_TYPE_TEXT,
				ValueType: entity.ELEMENT_VALUE_TYPE_INTEGER,
				Value: "",
				DefaultValue: 22,
			},
		},
	}


	crud := NewCrud(log, statsd)

	crud.AddEntity(users)

	// Two ways to mount route in your application:
	// 1. Mount CRUD routes to /gocrud (using Chi)
	router := chi.NewRouter()
	router.Mount("/gocrud", crud.Handler())
	// 2. Simple approach to mount CRUD routes
	//router := crud.Handler()

	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), router)
	if err != nil {
		log.Error("Problem starting server", err.Error())
		os.Exit(1)
	}
}

func toolup() (*appConfig, Logger, StatsDer) {
	config, err := loadAppConfig()
	if err != nil {
		log.Fatal("Error loading config", err.Error())
	}

	log := tools.NewLogger(config.IsLocal())
	log.Info(fmt.Sprintf("Application config - %+v", config))

	statsdConfig := tools.NewStatsDConfig(!config.IsLocal(), log)
	statsd, err := tools.NewStatsD(statsdConfig)
	if err != nil {
		log.Error("Error connecting to StatsD - defaulting to logging stats: ", err.Error())
	}

	return config, log, statsd
}
