package crud

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAPIService(t *testing.T) {

	testUsersEntity := Entity{
		ID:     "users",
		Label:  "User",
		Labels: "Users",
		Elements: Elements{
			{
				ID:         "id",
				Label:      "ID",
				PrimaryKey: true,
				FormType:   ELEMENT_FORM_TYPE_HIDDEN,
				DataType:   ELEMENT_DATA_TYPE_STRING,
			},
			{
				ID:       "name",
				Label:    "Name",
				FormType: ELEMENT_FORM_TYPE_TEXT,
				DataType: ELEMENT_DATA_TYPE_STRING,
			},
			{
				ID:           "age",
				Label:        "Age",
				FormType:     ELEMENT_FORM_TYPE_TEXT,
				DataType:     ELEMENT_DATA_TYPE_NUMBER,
				DefaultValue: 22,
			},
		},
	}

	t.Run("List returns records from database and returns it in client record format", func(t *testing.T) {
		fakeStore := NewFakeStorer()
		fakeStore.ListResponse = []StoreRecord{
			{
				Field{
					ID:       "id",
					Value:    "1",
					Hydrated: true,
				},
				Field{
					ID:       "name",
					Value:    "Superman",
					Hydrated: true,
				},
				Field{
					ID:       "age",
					Value:    11,
					Hydrated: true,
				},
			},
			{
				Field{
					ID:       "id",
					Value:    "2",
					Hydrated: true,
				},
				Field{
					ID:       "name",
					Value:    "Catwoman",
					Hydrated: true,
				},
				Field{
					ID:       "age",
					Value:    22,
					Hydrated: true,
				},
			},
		}
		fakeStore.ListError = nil
		fakeElementsValidator := &fakeElementsValidatorer{}
		apiService := newApiService(fakeStore, fakeElementsValidator)

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

	t.Run("Get returns record from database and returns it in client record format", func(t *testing.T) {
		fakeStore := NewFakeStorer()
		fakeStore.GetResponse = StoreRecord{
			Field{
				ID:       "id",
				Value:    "1",
				Hydrated: true,
			},
			Field{
				ID:       "name",
				Value:    "Superman",
				Hydrated: true,
			},
			Field{
				ID:       "age",
				Value:    11,
				Hydrated: true,
			},
		}
		fakeStore.GetError = nil
		fakeElementsValidator := &fakeElementsValidatorer{}
		apiService := newApiService(fakeStore, fakeElementsValidator)

		clientRecord, err := apiService.get(testUsersEntity, "1")

		id, _ := clientRecord.KeyValues.GetKeyValue("id")
		name, _ := clientRecord.KeyValues.GetKeyValue("name")

		assert.NoError(t, err)
		assert.Equal(t, 3, len(clientRecord.KeyValues), "Should have 3 fields (key-values)")
		assert.Equal(t, "1", id.Value, "ID is wrong")
		assert.Equal(t, "Superman", name.Value, "Name is wrong")
	})
}
