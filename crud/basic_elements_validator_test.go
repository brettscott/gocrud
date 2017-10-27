package crud

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicElementsValidation(t *testing.T) {

	t.Run("Fail regardless of what you pass in", func(t *testing.T) {
		testUsersEntity := &Entity{
			ID:     "users",
			Label:  "User",
			Labels: "Users",
			Elements: Elements{
				{
					ID:       "name",
					Label:    "Name",
					FormType: ELEMENT_FORM_TYPE_TEXT,
					DataType: ELEMENT_DATA_TYPE_BOOLEAN,
				},
			},
		}

		nameField := &Field{
			ID:       "name",
			Value:    "Jack Daniels",
			Hydrated: true,
		}
		storeRecord := StoreRecord{}
		storeRecord["name"] = nameField

		validator := NewBasicElementsValidator()
		isValid, clientErrors := validator.validate(testUsersEntity, storeRecord, ACTION_POST)

		assert.Equal(t, false, isValid)
		assert.NotNil(t, clientErrors)
		assert.Equal(t, 1, len(clientErrors.ElementsErrors["name"]))
		assert.Equal(t, "I'm going fail for the sake of it", clientErrors.ElementsErrors["name"][0])
		assert.Equal(t, 1, len(clientErrors.GlobalErrors))
	})

}
