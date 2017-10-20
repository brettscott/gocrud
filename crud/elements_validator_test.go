package crud

import (
	"testing"
	"github.com/brettscott/gocrud/model"
	"github.com/brettscott/gocrud/store"
	"github.com/stretchr/testify/assert"
)

func TestElementsValidator(t *testing.T) {

	elementsValidator := NewElementsValidator()

	t.Run("Passes when user data is valid", func(t *testing.T) {
		testEntity := model.Entity{
			ID: "test",
			Elements: []model.Element{
				{
					ID:"id",
					Label: "Identifier",
					DataType: model.ELEMENT_DATA_TYPE_STRING,
					PrimaryKey: true,
				},
				{
					ID:"name",
					Label: "Name",
					DataType: model.ELEMENT_DATA_TYPE_STRING,
					Validation: model.ElementValidation{
						Required: true,
					},
				},
			},
		}

		userData := store.Record{
			{
				ID: "id",
				Value: "12345",
				Hydrated: true,
			},
			{
				ID: "name",
				Value: "John Smith",
				Hydrated: true,
			},
		}

		success, elementsErrors, globalErrors := elementsValidator.validate(testEntity, userData, ACTION_POST)

		assert.Equal(t, true, success)
		assert.Equal(t, 0, len(elementsErrors))
		assert.Equal(t, 0, len(globalErrors))
	})


	t.Run("Fails when posting user data without a required field", func(t *testing.T) {
		testEntity := model.Entity{
			ID: "test",
			Elements: []model.Element{
				{
					ID:"id",
					Label: "Identifier",
					DataType: model.ELEMENT_DATA_TYPE_STRING,
					PrimaryKey: true,
				},
				{
					ID:"name",
					Label: "Name",
					DataType: model.ELEMENT_DATA_TYPE_STRING,
				},
			},
		}

		userData := store.Record{
			{
				ID: "id",
				Value: "12345",
				Hydrated: true,
			},
			{
				ID: "name",
				Value: "",
				Hydrated: false,
			},
		}

		success, elementsErrors, globalErrors := elementsValidator.validate(testEntity, userData, ACTION_POST)

		assert.Equal(t, false, success)
		assert.Equal(t, 1, len(elementsErrors))
		assert.Equal(t, 0, len(globalErrors))
	})
}


