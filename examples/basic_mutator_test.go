package examples

import (
	"github.com/brettscott/gocrud/crud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicMutator(t *testing.T) {

	t.Run("Simple mutator trims whitespace", func(t *testing.T) {

		testEntity := &crud.Entity{
			ID: "test",
			Elements: []crud.Element{
				{
					ID:         "id",
					Label:      "Identifier",
					DataType:   crud.ELEMENT_DATA_TYPE_STRING,
					PrimaryKey: true,
				},
				{
					ID:       "name",
					Label:    "Name",
					DataType: crud.ELEMENT_DATA_TYPE_STRING,
				},
			},
		}

		userData := crud.StoreRecord{}
		userData["id"] = &crud.Field{
			ID:       "id",
			Value:    "1234567",
			Hydrated: true,
		}
		userData["name"] = &crud.Field{
			ID:       "name",
			Value:    "  John Smith  ",
			Hydrated: true,
		}

		basicMutator := basicMutator{}
		clientErrors, err := basicMutator.Mutate(testEntity, &userData, crud.ACTION_POST)

		assert.NoError(t, err, "Should not error")
		assert.Nil(t, clientErrors)
		assert.Equal(t, "John Smith", userData["name"].Value, "Name value is wrong")
	})
}
