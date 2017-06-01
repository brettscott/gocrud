package examples

import (
	"fmt"
	"github.com/brettscott/gocrud/crud"
	"github.com/brettscott/gocrud/entity"
	"github.com/brettscott/gocrud/store"
	"github.com/mergermarket/gotools"
	"github.com/pressly/chi"
	"log"
	"net/http"
	"os"
)

// BasicExample should illustrate basic functionality of the CRUD
func BasicExample() {
	config, log, statsd := toolup()

	// TODO: Define schema
	// TODO: Build database connector - MySQL, Mongo
	// TODO: Pre/post hooks and override actions
	// TODO: Flexibility with rendering templates (custom head/foot/style)

	users := entity.Entity{
		ID:     "users",
		Label:  "User",
		Labels: "Users",
		Elements: entity.Elements{
			{
				ID:       "name",
				Label:    "Name",
				FormType: entity.ELEMENT_FORM_TYPE_TEXT,
				DataType: entity.ELEMENT_DATA_TYPE_STRING,
				Value:    "",
			},
			{
				ID:           "age",
				Label:        "Age",
				FormType:     entity.ELEMENT_FORM_TYPE_TEXT,
				DataType:     entity.ELEMENT_DATA_TYPE_INTEGER,
				Value:        "",
				DefaultValue: 22,
			},
		},
	}

	myConfig := &crud.Config{}

	myCrud := crud.NewCrud(myConfig, log, statsd)

	myStore, err := store.NewMongoStore("", "", "", statsd, log)
	if err != nil {
		fmt.Errorf("Error with store: %v", err)
	}
	myCrud.Store(myStore)

	// Register Entity
	myCrud.AddEntity(users)
	//myCrud.AddEntity(computers)

	// Add Sample data to DB
	myStore.Post()



	// Two ways to mount route in your application:
	// 1. Mount CRUD routes to /gocrud (using Chi)
	router := chi.NewRouter()
	router.Mount("/gocrud", myCrud.Handler())
	// 2. Simple approach to mount CRUD routes
	//router := Crud.Handler()

	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), router)
	if err != nil {
		log.Error("Problem starting server", err.Error())
		os.Exit(1)
	}
}

func toolup() (*appConfig, crud.Logger, crud.StatsDer) {
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
