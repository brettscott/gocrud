package entity

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestEntity_HydrateFromRecord(t *testing.T) {

	t.Run("Hydrate", func(t *testing.T) {

		entity := &Entity{
			Elements: Elements{
				{
					ID:       "name",
					Label:    "Name",
					FormType: ELEMENT_FORM_TYPE_TEXT,
					DataType: ELEMENT_DATA_TYPE_STRING,
					Value:    "",
				},
			},
		}

		record := &Record{
			ID: "users",
			KeyValues: KeyValues{
				{
					Key: "id",
					Value: "1234",
				},
				{
					Key: "name",
					Value: "Brett Scott",
				},
			},
		}

		entity.HydrateFromRecord(record)

		nameElement, err := entity.GetElement("name")
		if err != nil {
			t.Fatalf("Failed to get element \"name\" with error: %v", err)
		}
		assert.Equal(t, "Brett Scott", nameElement.Value)
	})
}

func TestEntity_ValidateSchema(t *testing.T) {

	t.Run("Validate POST", func(t *testing.T) {

		entity := &Entity{
			Elements: Elements{
				{
					ID:         "id",
					Label:      "Identifier",
					FormType:   ELEMENT_FORM_TYPE_HIDDEN,
					DataType:   ELEMENT_DATA_TYPE_STRING,
					Value:      "",
					Immutable:  true,
					PrimaryKey: true,
					Validation: ElementValidation{
					},
				},
				{
					ID:       "name",
					Label:    "Name",
					FormType: ELEMENT_FORM_TYPE_TEXT,
					DataType: ELEMENT_DATA_TYPE_STRING,
					Value:    "",
					Validation: ElementValidation{
						Required: true,
					},
				},
			},
		}

		t.Run("Passes when required fields provided", func(t *testing.T) {
			record := &Record{
				ID: "users",
				KeyValues: KeyValues{
					{
						Key: "name",
						Value: "Brett Scott",
					},
				},
			}


			entity.HydrateFromRecord(record)
			err := entity.Validate(VALIDATE_ACTION_POST)

			assert.NoError(t, err, "Should not have failed because required fields provided")
		})


		t.Run("Fails when required fields NOT provided", func(t *testing.T) {
			record := &Record{
				ID: "users",
				KeyValues: KeyValues{
					{
						Key: "name",
						Value: "",
					},
				},
			}

			entity.HydrateFromRecord(record)
			err := entity.Validate(VALIDATE_ACTION_POST)

			assert.Error(t, err, "Should have failed because required fields not provided")
		})

		t.Run("Fails when field receives the wrong type of data", func(t *testing.T) {
			record := &Record{
				ID: "users",
				KeyValues: KeyValues{
					{
						Key: "name",
						Value: 12345,
					},
				},
			}

			entity.HydrateFromRecord(record)
			err := entity.Validate(VALIDATE_ACTION_POST)

			assert.Error(t, err, `Should have failed because data type is "string" and "integer" supplied`)
		})
	})

}