package crud

import (
	"testing"
	"github.com/brettscott/gocrud/store"
	"github.com/brettscott/gocrud/model"
	"github.com/stretchr/testify/assert"
)

func TestAPIService(t *testing.T) {

	testUsersEntity := model.Entity{
		ID:     "users",
		Label:  "User",
		Labels: "Users",
		Elements: model.Elements{
			{
				ID:         "id",
				Label:      "ID",
				PrimaryKey: true,
				FormType:   model.ELEMENT_FORM_TYPE_HIDDEN,
				DataType:   model.ELEMENT_DATA_TYPE_STRING,
			},
			{
				ID:       "name",
				Label:    "Name",
				FormType: model.ELEMENT_FORM_TYPE_TEXT,
				DataType: model.ELEMENT_DATA_TYPE_STRING,
			},
			{
				ID:           "age",
				Label:        "Age",
				FormType:     model.ELEMENT_FORM_TYPE_TEXT,
				DataType:     model.ELEMENT_DATA_TYPE_NUMBER,
				DefaultValue: 22,
			},
		},
	}

	t.Run("List returns records from database as JSON response", func(t *testing.T) {
		fakeStore := store.NewFakeStorer()
		fakeStore.ListResponse = []store.Record{
			{
				store.Field{
					ID: "id",
					Value: "1",
					Hydrated: true,
				},
				store.Field{
					ID: "name",
					Value: "Superman",
					Hydrated: true,
				},
				store.Field{
					ID: "age",
					Value: 11,
					Hydrated: true,
				},
			},
			{
				store.Field{
					ID: "id",
					Value: "2",
					Hydrated: true,
				},
				store.Field{
					ID: "name",
					Value: "Catwoman",
					Hydrated: true,
				},
				store.Field{
					ID: "age",
					Value: 22,
					Hydrated: true,
				},
			},

		}
		fakeStore.ListError = nil
		apiService := newApiService(fakeStore)

		clientRecords, err := apiService.list(testUsersEntity)

		assert.NoError(t, err)

		assert.Equal(t, 2, len(clientRecords), "Should be 2 records returned")
		assert.Equal(t, 3, len(clientRecords[0].KeyValues), "First record should have 3 fields")

		assert.Equal(t, "id", clientRecords[0].KeyValues[0].Key, "First record's first field ID is wrong")
		assert.Equal(t, "1", clientRecords[0].KeyValues[0].Value, "First record's first field Value is wrong")
		assert.Equal(t, "name", clientRecords[0].KeyValues[1].Key, "First record's second field ID is wrong")
		assert.Equal(t, "Superman", clientRecords[0].KeyValues[1].Value, "First record's second field Value is wrong")

		assert.Equal(t, "id", clientRecords[1].KeyValues[0].Key, "Second record's first field ID is wrong")
		assert.Equal(t, "2", clientRecords[1].KeyValues[0].Value, "Second record's first field Value is wrong")
		assert.Equal(t, "name", clientRecords[1].KeyValues[1].Key, "Second record's second field ID is wrong")
		assert.Equal(t, "Catwoman", clientRecords[1].KeyValues[1].Value, "Second record's second field Value is wrong")

	})
}


