package main

import (
	"github.com/brettscott/gocrud/entity"
	"fmt"
	"github.com/kyoh86/richgo/config"
	"os"
	"net/http"
	"log"
)

func main() {


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


	crud := NewCrud()

	crud.AddEntity(users)

	err := http.ListenAndServe("8080", crud.Handler())
	if err != nil {
		fmt.Print("Problem starting server", err.Error())
		os.Exit(1)
	}


}
