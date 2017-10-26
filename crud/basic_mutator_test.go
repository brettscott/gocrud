package crud

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicMutator(t *testing.T) {

	t.Run("Simple mutator trims whitespace", func(t *testing.T) {

		testEntity := &Entity{
			ID: "test",
			Elements: []Element{
				{
					ID:         "id",
					Label:      "Identifier",
					DataType:   ELEMENT_DATA_TYPE_STRING,
					PrimaryKey: true,
				},
				{
					ID:       "name",
					Label:    "Name",
					DataType: ELEMENT_DATA_TYPE_STRING,
				},
			},
		}

		userData := StoreRecord{}
		userData["id"] = &Field{
			ID:       "id",
			Value:    "1234567",
			Hydrated: true,
		}
		userData["name"] = &Field{
			ID:       "name",
			Value:    "  John Smith  ",
			Hydrated: true,
		}

		basicMutator := basicMutator{}
		err, elementsErrors, globalErrors := basicMutator.mutate(testEntity, &userData, ACTION_POST)

		assert.NoError(t, err, "Should not error")
		assert.Equal(t, 0, len(elementsErrors), "Elements errors should be empty")
		assert.Equal(t, 0, len(globalErrors), "Global errors should be empty")

		assert.Equal(t, "John Smith", userData["name"].Value, "Name value is wrong")
	})
}
