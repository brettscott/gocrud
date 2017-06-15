package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
				},
			},
		}

		record := &Record{
			ID: "users",
			KeyValues: KeyValues{
				{
					Key:   "id",
					Value: "1234",
				},
				{
					Key:   "name",
					Value: "Brett Scott",
				},
			},
		}

		entity.HydrateFromRecord(record, HYDRATE_FROM_RECORD_ACTION_POST)

		nameElement, err := entity.GetElement("name")
		if err != nil {
			t.Fatalf("Failed to get element \"name\" with error: %v", err)
		}
		assert.Equal(t, "Brett Scott", nameElement.Value)
	})
}

func TestEntity_Validate(t *testing.T) {

	t.Run("Validate POST", func(t *testing.T) {

		entity := &Entity{
			Elements: Elements{
				{
					ID:         "id",
					Label:      "Identifier",
					FormType:   ELEMENT_FORM_TYPE_HIDDEN,
					DataType:   ELEMENT_DATA_TYPE_STRING,
					Immutable:  true,
					PrimaryKey: true,
					Validation: ElementValidation{},
				},
				{
					ID:       "name",
					Label:    "Name",
					FormType: ELEMENT_FORM_TYPE_TEXT,
					DataType: ELEMENT_DATA_TYPE_STRING,
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
						Key:   "name",
						Value: "Brett Scott",
					},
				},
			}

			entity.HydrateFromRecord(record, HYDRATE_FROM_RECORD_ACTION_POST)
			err := entity.Validate(VALIDATE_ACTION_POST)

			assert.NoError(t, err, "Should not have failed because required fields provided")
		})

		t.Run("Fails when required field is empty", func(t *testing.T) {
			record := &Record{
				ID: "users",
				KeyValues: KeyValues{
					{
						Key:   "name",
						Value: "",
					},
				},
			}

			entity.HydrateFromRecord(record, HYDRATE_FROM_RECORD_ACTION_POST)
			err := entity.Validate(VALIDATE_ACTION_POST)

			assert.Error(t, err, "Should have failed because required fields not provided")
		})

		t.Run("Fails when field receives the wrong type of data", func(t *testing.T) {
			record := &Record{
				ID: "users",
				KeyValues: KeyValues{
					{
						Key:   "name",
						Value: 12345,
					},
				},
			}

			entity.HydrateFromRecord(record, HYDRATE_FROM_RECORD_ACTION_POST)
			err := entity.Validate(VALIDATE_ACTION_POST)

			assert.Error(t, err, `Should have failed because data type is "string" and "integer" supplied`)
		})

		t.Run("Passes when field is present when it must be provided", func(t *testing.T) {
			entity := &Entity{
				Elements: Elements{
					{
						ID:       "name",
						Label:    "Name",
						FormType: ELEMENT_FORM_TYPE_TEXT,
						DataType: ELEMENT_DATA_TYPE_STRING,
						Validation: ElementValidation{
							MustProvide: true,
						},
					},
				},
			}

			record := &Record{
				ID: "users",
				KeyValues: KeyValues{
					{
						Key:   "name",
						Value: "Brett",
					},
				},
			}

			entity.HydrateFromRecord(record, HYDRATE_FROM_RECORD_ACTION_POST)
			err := entity.Validate(VALIDATE_ACTION_POST)

			assert.NoError(t, err, `Should not have failed because field was provided`)
		})

		t.Run("Fails when field is missing when it must be provided", func(t *testing.T) {
			entity := &Entity{
				Elements: Elements{
					{
						ID:       "name",
						Label:    "Name",
						FormType: ELEMENT_FORM_TYPE_TEXT,
						DataType: ELEMENT_DATA_TYPE_STRING,
						Validation: ElementValidation{
							MustProvide: true,
						},
					},
				},
			}

			record := &Record{
				ID:        "users",
				KeyValues: KeyValues{},
			}

			entity.HydrateFromRecord(record, HYDRATE_FROM_RECORD_ACTION_POST)
			err := entity.Validate(VALIDATE_ACTION_POST)

			assert.Error(t, err, `Should have failed because field must be provided when "MustProvide" is set`)
		})

	})

}
