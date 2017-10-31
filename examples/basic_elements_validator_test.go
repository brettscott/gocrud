package examples

import (
	"github.com/brettscott/gocrud/crud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicElementsValidation(t *testing.T) {

	t.Run("Fail regardless of what you pass in", func(t *testing.T) {
		testUsersEntity := &crud.Entity{
			ID:     "users",
			Label:  "User",
			Labels: "Users",
			Elements: crud.Elements{
				{
					ID:       "name",
					Label:    "Name",
					FormType: crud.ELEMENT_FORM_TYPE_TEXT,
					DataType: crud.ELEMENT_DATA_TYPE_BOOLEAN,
				},
			},
		}

		nameField := &crud.Field{
			ID:       "name",
			Value:    "Jack Daniels",
			Hydrated: true,
		}
		storeRecord := crud.StoreRecord{}
		storeRecord["name"] = nameField

		validator := &basicElementsValidator{}
		isValid, clientErrors := validator.Validate(testUsersEntity, storeRecord, crud.ACTION_POST)

		assert.Equal(t, false, isValid)
		assert.NotNil(t, clientErrors)
		assert.Equal(t, 1, len(clientErrors.ElementsErrors["name"]))
		assert.Equal(t, "I'm going fail for the sake of it", clientErrors.ElementsErrors["name"][0])
		assert.Equal(t, 1, len(clientErrors.GlobalErrors))
	})

}
