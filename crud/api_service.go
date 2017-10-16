package crud

import (
	"encoding/json"
	"fmt"
	"github.com/brettscott/gocrud/model"
	"github.com/brettscott/gocrud/store"
)

func newApiService(store store.Storer) apiService {
	return apiService{
		store: store,
	}
}

type apiService struct {
	store store.Storer
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

	// Not found in database
	if storeRecord.IsHydrated() == false {
		return jsonResponse, nil
	}

	clientRecord := marshalStoreRecordToClientRecord(storeRecord)

	jsonResponse, err = json.Marshal(clientRecord)
	if err != nil {
		return nil, fmt.Errorf(`Failed to marshal record into JSON for entity "%s" - %s`, entity.Label, err)
	}

	return jsonResponse, nil
}

func (a *apiService) save(entity model.Entity, action string, body []byte, recordID string) (jsonResponse []byte, err error) {
	record := &Record{}
	err = record.UnmarshalJSON(body)
	if err != nil {
		return nil, fmt.Errorf(`Failed to unmarshal JSON for entity "%s" - %s`, entity.Label, err)
	}

	storeRecord, err := marshalClientRecordToStoreRecord(entity, record, action)
	if err != nil {
		return nil, fmt.Errorf(`Failed to marshal client record to store record for entity "%s" - %s`, entity.Label, err)
	}

	err = validate(entity, storeRecord, action)
	if err != nil {
		return nil, fmt.Errorf(`Failed validation for entity "%s" - %s`, entity.Label, err)
	}

	switch action {
	case ACTION_POST:
		recordID, err = a.store.Post(entity, storeRecord)
		break
	case ACTION_PUT:
		if recordID == "" {
			return nil, fmt.Errorf(`Missing record ID for entity "%s"`, entity.Label)
		}
		err = a.store.Put(entity, storeRecord, recordID)
		break
	case ACTION_PATCH:
		if recordID == "" {
			return nil, fmt.Errorf(`Missing record ID for entity "%s"`, entity.Label)
		}
		err = a.store.Patch(entity, storeRecord, recordID)
		break
	default:
		return nil, fmt.Errorf(`Invalid action "%s" for entity "%s"`, action, entity.Label)
		break
	}
	if err != nil {
		return nil, fmt.Errorf(`Failed to "%s" for entity "%s" - %s`, action, entity.Label, err)
	}

	savedStoreRecord, err := a.store.Get(entity, recordID)
	if err != nil {
		return nil, fmt.Errorf(`Failed to get newly created DB record for entity "%s" - %s`, entity.Label, err)
	}

	if savedStoreRecord.IsHydrated() == false {
		return nil, fmt.Errorf(`New created DB record was not found in database for entity "%s" - %s`, entity.Label, err)
	}

	clientRecord := marshalStoreRecordToClientRecord(savedStoreRecord)
	jsonResponse, err = json.Marshal(clientRecord)
	if err != nil {
		return nil, fmt.Errorf(`Failed to marshal record into JSON for entity "%s" - %s`, entity.Label, err)
	}
	return jsonResponse, nil
}

func (a *apiService) delete(entity model.Entity, recordID string) error {
	err := a.store.Delete(entity, recordID)
	if err != nil {
		return fmt.Errorf(`Store delete failed for entity "%s" recordID "%s" - %s`, entity.Label, recordID, err)
	}
	return nil
}

func marshalClientRecordToStoreRecord(entity model.Entity, clientRecord *Record, action string) (data store.Record, err error) {

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
