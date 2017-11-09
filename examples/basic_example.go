package examples

import (
	"fmt"
	"github.com/brettscott/gocrud/crud"
	"github.com/mergermarket/gotools"
	"github.com/pressly/chi"
	"log"
	"net/http"
	"os"
)

// BasicExample should illustrate basic functionality of the CRUD
func BasicExample() {
	config, log, statsd := infra()

	// TODO: Define schema
	// TODO: Build database connector - MySQL, Mongo
	// TODO: Pre/post hooks and override actions
	// TODO: Flexibility with rendering templates (custom head/foot/style)

	users := &crud.Entity{
		ID:     "users",
		Label:  "User",
		Labels: "Users",
		Elements: crud.Elements{
			{
				ID:         "id",
				Label:      "ID",
				PrimaryKey: true,
				FormType:   crud.ELEMENT_FORM_TYPE_HIDDEN,
				DataType:   crud.ELEMENT_DATA_TYPE_STRING,
			},
			{
				ID:       "name",
				Label:    "Name",
				FormType: crud.ELEMENT_FORM_TYPE_TEXT,
				DataType: crud.ELEMENT_DATA_TYPE_STRING,
			},
			{
				ID:           "age",
				Label:        "Age",
				FormType:     crud.ELEMENT_FORM_TYPE_TEXT,
				DataType:     crud.ELEMENT_DATA_TYPE_NUMBER,
				DefaultValue: 22,
			},
		},
	}

	// Todo: should do NewEntity and not newing up Entity manually.
	err := users.CheckConfiguration()
	if err != nil {
		log.Error(fmt.Sprintf(`Error with "users" entity: %v`, err))
		os.Exit(1)
	}

	myConfig := &crud.Config{}

	myCrud := crud.NewCrud(myConfig, log, statsd)

	// Add store (database)
	myStore, err := crud.NewMongoStore(os.Getenv("MONGO_DB_CONNECTION"), "", os.Getenv("MONGO_DB_NAME"), statsd, log)
	if err != nil {
		log.Error(fmt.Sprintf("Error with store: %v", err))
		os.Exit(1)
	}
	myCrud.AddStore(myStore)

	// Register Entity
	myCrud.AddEntity(users)

	// Basic mutator added
	myCrud.AddMutator(&basicMutator{})

	// Basic validator added
	//myCrud.AddElementsValidator(&basicElementsValidator{})  // TODO add a non-destructive validator for illustrative purposes (this one fails on all requests)

	// Mount Route
	router := chi.NewRouter()
	// TODO remove this warning when Chi is fixed
	router.Get("/gocrud", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bug in chi.  Add / to end"))
	})
	router.Mount("/gocrud/", myCrud.Handler())

	// Serve HTTP endpoint
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), router)
	if err != nil {
		log.Error("Problem starting server", err.Error())
		os.Exit(1)
	}
}

func infra() (*appConfig, crud.Logger, crud.StatsDer) {
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
