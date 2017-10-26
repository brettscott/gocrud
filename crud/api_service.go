package crud

import (
	"fmt"
)

type elementsValidatorer interface {
	validate(entity *Entity, record StoreRecord, action string) (success bool, elementsErrors map[string][]string, globalErrors []string)
}

type mutatorer interface {
	mutate(entity *Entity, storeRecord StoreRecord, action string) (mutatedStoreRecord StoreRecord, elementsErrors map[string][]string, globalErrors []string)
}

func newApiService(store Storer, elementsValidators []elementsValidatorer, mutators []mutatorer) apiService {
	return apiService{
		store:              store,
		elementsValidators: elementsValidators,
		mutators:           mutators,
	}
}

type apiService struct {
	store              Storer
	elementsValidators []elementsValidatorer
	mutators           []mutatorer
}

func (a *apiService) list(entity *Entity) (clientRecords []ClientRecord, err error) {
	storeRecords, err := a.store.List(entity)
	if err != nil {
		return nil, fmt.Errorf(`Store query failed for entity "%s" - %s`, entity.Label, err)
	}

	for _, storeRecord := range storeRecords {
		clientRecord := marshalStoreRecordToClientRecord(storeRecord)
		clientRecords = append(clientRecords, clientRecord)
	}

	return
}

func (a *apiService) get(entity *Entity, recordID string) (clientRecord ClientRecord, err error) {
	storeRecord, err := a.store.Get(entity, recordID)
	if err != nil {
		return clientRecord, fmt.Errorf(`Store query failed for entity "%s" recordID "%s" - %s`, entity.Label, recordID, err)
	}

	// Not found in database
	if storeRecord.IsHydrated() == false {
		return
	}

	clientRecord = marshalStoreRecordToClientRecord(storeRecord)

	return
}

func (a *apiService) save(entity *Entity, action string, clientRecord *ClientRecord, recordID string) (savedClientRecord ClientRecord, err error) {
	storeRecord, err := marshalClientRecordToStoreRecord(entity, clientRecord, action)
	if err != nil {
		return savedClientRecord, fmt.Errorf(`Failed to marshal client record to store record for entity "%s" - %s`, entity.Label, err)
	}

	mergedElementsValidators := append(a.elementsValidators, entity.ElementsValidators...)
	for _, validator := range mergedElementsValidators {
		// TODO Goroutine in order to run through each validator and report all issues that each validator finds
		isValid, elementsErrors, globalErrors := validator.validate(entity, storeRecord, action)
		if !isValid {
			return savedClientRecord, fmt.Errorf(`Failed validation for entity "%s" - %v %v`, entity.Label, elementsErrors, globalErrors)
		}
	}

	// TODO allow users to specify order of mutators
	mergedMutators := append(a.mutators, entity.Mutators...)
	for _, mutator := range mergedMutators {
		// TODO Goroutine in order to run through each validator and report all issues that each validator finds
		_, elementsErrors, globalErrors := mutator.mutate(entity, storeRecord, action)
		if len(elementsErrors) > 0 || len(globalErrors) > 0 {
			return savedClientRecord, fmt.Errorf(`Failed validation for entity "%s" - %v %v`, entity.Label, elementsErrors, globalErrors)
		}
	}

	switch action {
	case ACTION_POST:
		recordID, err = a.store.Post(entity, storeRecord)
		break
	case ACTION_PUT:
		if recordID == "" {
			return savedClientRecord, fmt.Errorf(`Missing record ID for entity "%s"`, entity.Label)
		}
		err = a.store.Put(entity, storeRecord, recordID)
		break
	case ACTION_PATCH:
		if recordID == "" {
			return savedClientRecord, fmt.Errorf(`Missing record ID for entity "%s"`, entity.Label)
		}
		err = a.store.Patch(entity, storeRecord, recordID)
		break
	default:
		return savedClientRecord, fmt.Errorf(`Invalid action "%s" for entity "%s"`, action, entity.Label)
		break
	}
	if err != nil {
		return savedClientRecord, fmt.Errorf(`Failed to "%s" for entity "%s" - %s`, action, entity.Label, err)
	}

	savedStoreRecord, err := a.store.Get(entity, recordID)
	if err != nil {
		return savedClientRecord, fmt.Errorf(`Failed to get newly created DB record for entity "%s" - %s`, entity.Label, err)
	}

	if savedStoreRecord.IsHydrated() == false {
		return savedClientRecord, fmt.Errorf(`New created DB record was not found in database for entity "%s" - %s`, entity.Label, err)
	}

	savedClientRecord = marshalStoreRecordToClientRecord(savedStoreRecord)
	return savedClientRecord, nil
}

func (a *apiService) delete(entity *Entity, recordID string) error {
	err := a.store.Delete(entity, recordID)
	if err != nil {
		return fmt.Errorf(`Store delete failed for entity "%s" recordID "%s" - %s`, entity.Label, recordID, err)
	}
	return nil
}

func marshalClientRecordToStoreRecord(entity *Entity, clientRecord *ClientRecord, action string) (data StoreRecord, err error) {

	for i, _ := range entity.Elements {
		element := &entity.Elements[i]

		datum := Field{
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

func marshalStoreRecordToClientRecord(storeRecord StoreRecord) ClientRecord {
	clientRecord := ClientRecord{}
	kvs := KeyValues{}
	for _, field := range storeRecord {
		kv := KeyValue{Key: field.ID, Value: field.Value}
		kvs = append(kvs, kv)
	}
	clientRecord.KeyValues = kvs
	return clientRecord
}
