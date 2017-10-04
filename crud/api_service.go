package crud

import (
	"encoding/json"
	"fmt"
	"github.com/brettscott/gocrud/model"
	"github.com/brettscott/gocrud/store"
)

type apiService struct {
	store store.Storer
}

func newApiService(store store.Storer) apiService {
	return apiService{
		store: store,
	}
}

func (a *apiService) list(entity model.Entity) (jsonResponse []byte, err error) {
	storeRecords, err := a.store.List(entity)
	if err != nil {
		return nil, fmt.Errorf(`Store query failed for entity "%s" - %s`, entity.Label, err)
	}

	clientRecords := []Record{}
	for _, storeRecord := range storeRecords {
		clientRecord := marshalStoreRecordToClientRecord(storeRecord)
		clientRecords = append(clientRecords, clientRecord)
	}

	jsonResponse, err = json.Marshal(clientRecords)
	if err != nil {
		return nil, fmt.Errorf(`Failed to marshal record into JSON for entity "%s" - %s`, entity.Label, err)
	}

	return jsonResponse, nil
}

func (a *apiService) get(entity model.Entity, recordID string) (jsonResponse []byte, err error) {
	storeRecord, err := a.store.Get(entity, recordID)
	if err != nil {
		return nil, fmt.Errorf(`Store query failed for entity "%s" recordID "%s" - %s`, entity.Label, recordID, err)
	}

	clientRecord := marshalStoreRecordToClientRecord(storeRecord)

	jsonResponse, err = json.Marshal(clientRecord)
	if err != nil {
		return nil, fmt.Errorf(`Failed to marshal record into JSON for entity "%s" - %s`, entity.Label, err)
	}

	return jsonResponse, nil
}

// MarshalRecordToEntityData marshals the data received from client
func marshalRecordToEntityData(entity model.Entity, clientRecord *Record, action string) (data store.Record, err error) {

	for i, _ := range entity.Elements {
		element := &entity.Elements[i]

		datum := store.Field{
			ID: element.ID,
		}

		for _, keyValue := range clientRecord.KeyValues {

			if action == ACTION_POST && element.PrimaryKey == true {
				continue
			}

			if keyValue.Key == element.ID {
				datum.Value = keyValue.Value
				datum.Hydrated = true
				break
			}
		}

		data = append(data, datum)

	}
	return data, nil
}

func marshalStoreRecordToClientRecord(storeRecord store.Record) Record {
	clientRecord := Record{}
	kvs := KeyValues{}
	for _, field := range storeRecord {
		kv := KeyValue{Key: field.ID, Value: field.Value}
		kvs = append(kvs, kv)
	}
	clientRecord.KeyValues = kvs
	return clientRecord
}
