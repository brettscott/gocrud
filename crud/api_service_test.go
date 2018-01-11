package crud

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAPIService(t *testing.T) {

	testUsersEntity := &Entity{
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
			{
				ID:       "isActive",
				Label:    "Is Active",
				FormType: ELEMENT_FORM_TYPE_TEXT,
				DataType: ELEMENT_DATA_TYPE_BOOLEAN,
			},
		},
	}

	supermanStoreRecord := StoreRecord{}
	supermanStoreRecord["id"] = &Field{
		ID:       "id",
		Value:    "the-superman-id",
		Hydrated: true,
	}
	supermanStoreRecord["name"] = &Field{
		ID:       "name",
		Value:    "Superman",
		Hydrated: true,
	}
	supermanStoreRecord["age"] = &Field{
		ID:       "age",
		Value:    11,
		Hydrated: true,
	}

	catwomanStoreRecord := StoreRecord{}
	catwomanStoreRecord["id"] = &Field{
		ID:       "id",
		Value:    "the-catwoman-id",
		Hydrated: true,
	}
	catwomanStoreRecord["name"] = &Field{
		ID:       "name",
		Value:    "Catwoman",
		Hydrated: true,
	}
	catwomanStoreRecord["age"] = &Field{
		ID:       "age",
		Value:    11,
		Hydrated: true,
	}

	t.Run("List", func(t *testing.T) {

		t.Run("returns records from database and returns it in client record format", func(t *testing.T) {
			fakeStore := NewFakeStorer()
			fakeStore.ListResponse = []StoreRecord{
				supermanStoreRecord,
				catwomanStoreRecord,
			}
			fakeStore.ListError = nil
			fakeStores := NewFakeStorers(fakeStore)
			fakeElementsValidators := NewFakeElementsValidatorers()
			fakeMutators := newFakeEmptyMutatorers()
			apiService := newApiService(fakeStores, fakeElementsValidators, fakeMutators)

			clientRecords, err := apiService.list(testUsersEntity)
			assert.NoError(t, err)

			assert.Equal(t, 2, len(clientRecords), "Should be 2 records returned")

			superman, err := clientRecords.GetClientRecordByKeyValue("id", "the-superman-id")
			assert.NoError(t, err)
			catwoman, err := clientRecords.GetClientRecordByKeyValue("id", "the-catwoman-id")
			assert.NoError(t, err)

			id, err := superman.GetValue("id")
			assert.NoError(t, err)
			name, err := superman.GetValue("name")
			assert.NoError(t, err)
			assert.Equal(t, "the-superman-id", id, "First record's first field Value is wrong")
			assert.Equal(t, "Superman", name, "First record's second field Value is wrong")

			id, err = catwoman.GetValue("id")
			assert.NoError(t, err)
			name, err = catwoman.GetValue("name")
			assert.NoError(t, err)
			assert.Equal(t, "the-catwoman-id", id, "Second record's first field Value is wrong")
			assert.Equal(t, "Catwoman", name, "Second record's second field Value is wrong")
		})
	})

	t.Run("Get", func(t *testing.T) {

		t.Run("returns record from database and returns it in client record format", func(t *testing.T) {
			fakeStore := NewFakeStorer()
			fakeStore.GetResponse = supermanStoreRecord
			fakeStore.GetError = nil
			fakeStores := NewFakeStorers(fakeStore)
			fakeElementsValidators := NewFakeElementsValidatorers()
			fakeMutators := newFakeEmptyMutatorers()
			apiService := newApiService(fakeStores, fakeElementsValidators, fakeMutators)

			clientRecord, err := apiService.get(testUsersEntity, "the-superman-id")
			assert.NoError(t, err)
			id, err := clientRecord.GetValue("id")
			assert.NoError(t, err)
			name, err := clientRecord.GetValue("name")
			assert.NoError(t, err)
			assert.Equal(t, 3, len(clientRecord.KeyValues), "Should have 3 fields (key-values)")
			assert.Equal(t, "the-superman-id", id, "ID is wrong")
			assert.Equal(t, "Superman", name, "Name is wrong")
		})
	})

	t.Run("Save", func(t *testing.T) {

		t.Run("persists client record data in database and get back new record", func(t *testing.T) {
			recordID := "1234567"
			clientRecord := &ClientRecord{
				KeyValues: KeyValues{
					{
						Key:   "id",
						Value: recordID,
					},
					{
						Key:   "name",
						Value: "Jim Beam",
					},
				},
			}
			fakeStore := NewFakeStorer()
			fakeStore.GetResponse = supermanStoreRecord
			fakeStore.GetError = nil
			fakeStores := NewFakeStorers(fakeStore)
			fakeElementsValidators := NewFakeEmptyElementsValidatorers()
			fakeMutators := newFakeEmptyMutatorers()
			apiService := newApiService(fakeStores, fakeElementsValidators, fakeMutators)

			savedClientRecord, clientErrors, err := apiService.save(testUsersEntity, ACTION_POST, clientRecord, recordID)
			assert.NoError(t, err)

			assert.NotNil(t, clientErrors)
			assert.Equal(t, false, clientErrors.HasErrors(), "Should not have any errors")

			// Test data persisted to store
			assert.Equal(t, 1, fakeStore.PostCalled)
			storedID, err := fakeStore.PostStoreRecord.GetValue("id")
			assert.NoError(t, err)
			assert.Nil(t, storedID)
			storedName, err := fakeStore.PostStoreRecord.GetValue("name")
			assert.NoError(t, err)
			assert.Equal(t, "Jim Beam", storedName)

			// Test client record returned from call to save()
			name, err := savedClientRecord.GetValue("name")
			assert.NoError(t, err)
			assert.Equal(t, name, "Superman")
			id, err := savedClientRecord.GetValue("id")
			assert.NoError(t, err)
			assert.Equal(t, id, "the-superman-id")
		})

		t.Run("mutates client record before persisting in store", func(t *testing.T) {
			recordID := "1234567"
			clientRecord := &ClientRecord{
				KeyValues: KeyValues{
					{
						Key:   "id",
						Value: recordID,
					},
					{
						Key:   "name",
						Value: "Jim Beam",
					},
				},
			}
			fakeStore := NewFakeStorer()
			fakeStore.GetResponse = supermanStoreRecord
			fakeStore.GetError = nil
			fakeStores := NewFakeStorers(fakeStore)
			fakeElementsValidators := NewFakeEmptyElementsValidatorers()
			mutatedStoreRecord := StoreRecord{}
			mutatedStoreRecord["name"] = &Field{
				ID:       "name",
				Value:    "Jack Daniels",
				Hydrated: true,
			}
			fakeMutator := &FakeMutatorer{
				StoreRecord: &mutatedStoreRecord,
			}
			mutators := []mutatorer{
				fakeMutator,
			}
			apiService := newApiService(fakeStores, fakeElementsValidators, mutators)
			_, clientErrors, err := apiService.save(testUsersEntity, ACTION_POST, clientRecord, recordID)
			assert.NoError(t, err)
			assert.NotNil(t, clientErrors)
			assert.Equal(t, false, clientErrors.HasErrors(), "Should not have any errors")

			assert.Equal(t, 1, fakeStore.PostCalled, "Should post to store")
			storedName, err := fakeStore.PostStoreRecord.GetValue("name")
			assert.NoError(t, err)
			assert.Equal(t, "Jim Beam", storedName)
		})

		t.Run("basic element validator fails validation for each element", func(t *testing.T) {
			recordID := "1234567"
			clientRecord := &ClientRecord{
				KeyValues: KeyValues{
					{
						Key:   "id",
						Value: recordID,
					},
					{
						Key:   "name",
						Value: "John Smith",
					},
					{
						Key:   "isActive",
						Value: "true",
					},
				},
			}
			fakeStore := NewFakeStorer()
			fakeStore.GetResponse = supermanStoreRecord
			fakeStore.GetError = nil
			fakeStores := NewFakeStorers(fakeStore)

			fakeElementsErrors := map[string][]string{}
			fakeElementsErrors["id"] = append(fakeElementsErrors["id"], "id is invalid")
			fakeElementsErrors["name"] = append(fakeElementsErrors["name"], "name is invalid")
			fakeClientErrors := &ClientErrors{
				ElementsErrors: fakeElementsErrors,
			}
			fakeElementsValidator := &FakeElementsValidatorer{
				Success:      false,
				ClientErrors: fakeClientErrors,
			}
			elementsValidators := []elementsValidatorer{
				fakeElementsValidator,
			}

			fakeMutators := newFakeEmptyMutatorers()

			apiService := newApiService(fakeStores, elementsValidators, fakeMutators)

			_, clientErrors, err := apiService.save(testUsersEntity, ACTION_POST, clientRecord, recordID)
			assert.Error(t, err)
			assert.NotNil(t, clientErrors)
			assert.Equal(t, true, clientErrors.HasErrors(), "Should have error because validator fails each element")
			assert.Equal(t, 1, len(clientErrors.ElementsErrors["id"]))
			assert.Equal(t, 1, len(clientErrors.ElementsErrors["name"]))
			//assert.Equal(t, "name is invalid", clientErrors.ElementsErrors["name"][0])
			assert.Equal(t, 0, fakeStore.PostCalled, "Should not try to save invalid record")
		})
	})
}
