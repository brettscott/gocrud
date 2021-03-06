package crud

import (
	"fmt"
	"strconv"
)

type elementsValidatorer interface {
	Validate(entity *Entity, record StoreRecord, action string) (success bool, clientErrors *ClientErrors)
}

type mutatorer interface {
	Mutate(entity *Entity, storeRecord *StoreRecord, action string) (clientErrors *ClientErrors, err error)
}

func newApiService(stores []Storer, elementsValidators []elementsValidatorer, mutators []mutatorer) apiService {
	return apiService{
		stores:             stores,
		elementsValidators: elementsValidators,
		mutators:           mutators,
	}
}

type apiService struct {
	stores             []Storer
	elementsValidators []elementsValidatorer
	mutators           []mutatorer
}

func (a *apiService) list(entity *Entity) (clientRecords ClientRecords, err error) {
	clientRecords = ClientRecords{}
	// TODO reads from first readable database.  Add ability to order databases?
	for _, store := range a.stores {
		if store.Mode(entity).IsReadable() {
			storeRecords, err := store.List(entity)
			if err != nil {
				return nil, fmt.Errorf(`Store query failed for entity "%s" - %s`, entity.Label, err)
			}

			for _, storeRecord := range storeRecords {
				clientRecord := marshalStoreRecordToClientRecord(storeRecord)
				clientRecords = append(clientRecords, clientRecord)
			}
			return clientRecords, nil
		}
	}

	return
}

func (a *apiService) get(entity *Entity, recordID string) (clientRecord ClientRecord, err error) {
	// TODO reads from first readable database.  Add ability to order databases?
	for _, store := range a.stores {
		if store.Mode(entity).IsReadable() {
			storeRecord, err := store.Get(entity, recordID)
			if err != nil {
				return clientRecord, fmt.Errorf(`Store query failed for entity "%s" recordID "%s" - %s`, entity.Label, recordID, err)
			}

			// Not found in database
			if storeRecord.IsHydrated() == false {
				return clientRecord, nil
			}

			clientRecord = marshalStoreRecordToClientRecord(storeRecord)
		}
	}
	return
}

func (a *apiService) save(entity *Entity, action string, clientRecord *ClientRecord, recordID string) (savedClientRecord ClientRecord, clientErrors *ClientErrors, err error) {
	clientErrors = &ClientErrors{}
	storeRecord, err := marshalClientRecordToStoreRecord(entity, clientRecord, action)
	if err != nil {
		return savedClientRecord, clientErrors, nil
	}

	mergedElementsValidators := append(a.elementsValidators, entity.ElementsValidators...)
	for _, validator := range mergedElementsValidators {
		// TODO Goroutine in order to run through each validator and report all issues that each validator finds
		var isValid bool
		isValid, validateClientErrors := validator.Validate(entity, storeRecord, action)
		if validateClientErrors == nil {
			validateClientErrors = clientErrors
		}
		if !isValid {
			err = fmt.Errorf(`validation failure for entity "%s"`, entity.Label)
			return savedClientRecord, validateClientErrors, err
		}
	}

	// TODO allow users to specify order of mutators
	mergedMutators := append(a.mutators, entity.Mutators...)
	for _, mutator := range mergedMutators {
		// TODO Goroutine in order to run through each validator and report all issues that each validator finds
		mutateClientErrors, err := mutator.Mutate(entity, &storeRecord, action)
		if mutateClientErrors == nil {
			mutateClientErrors = clientErrors
		}
		if err != nil {
			err = fmt.Errorf(`mutation error for entity "%s" with error: %v`, entity.Label, err)
			return savedClientRecord, mutateClientErrors, err
		}
	}

	// TODO support multiple stores - at the moment, once written to first DB, doesn't write to subsequent.  Add Goroutine
	for _, store := range a.stores {
		if store.Mode(entity).IsWritable() == false {
			continue
		}
		switch action {
		case ACTION_POST:
			recordID, err = store.Post(entity, storeRecord)
			break
		case ACTION_PUT:
			if recordID == "" {
				err = fmt.Errorf(`Missing record ID for entity "%s"`, entity.Label)
				return savedClientRecord, clientErrors, err
			}
			err = store.Put(entity, storeRecord, recordID)
			break
		case ACTION_PATCH:
			if recordID == "" {
				err = fmt.Errorf(`Missing record ID for entity "%s"`, entity.Label)
				return savedClientRecord, clientErrors, err
			}
			err = store.Patch(entity, storeRecord, recordID)
			break
		default:
			err = fmt.Errorf(`Invalid action "%s" for entity "%s"`, action, entity.Label)
			return savedClientRecord, clientErrors, err
			break
		}
		if err != nil {
			err = fmt.Errorf(`Failed to "%s" for entity "%s" - %s`, action, entity.Label, err)
			return savedClientRecord, clientErrors, err
		}

		savedStoreRecord, err := store.Get(entity, recordID)
		if err != nil {
			err = fmt.Errorf(`Failed to get newly created record from the database for entity "%s" - %s`, entity.Label, err)
			return savedClientRecord, clientErrors, err
		}

		if savedStoreRecord.IsHydrated() == false {
			err = fmt.Errorf(`Newly created record was not found in the database (save didn't work) for entity "%s" recordID: "%s"`, entity.Label, recordID)
			return savedClientRecord, clientErrors, err
		}

		savedClientRecord = marshalStoreRecordToClientRecord(savedStoreRecord)
		return savedClientRecord, clientErrors, err
	}
	err = fmt.Errorf("Failed to find a writeable database")
	return savedClientRecord, clientErrors, err
}

func (a *apiService) delete(entity *Entity, recordID string) error {
	// TODO transaction in case delete works to one DB but not the rest
	for _, store := range a.stores {
		if store.Mode(entity).IsDeletable() {
			err := store.Delete(entity, recordID)
			if err != nil {
				return fmt.Errorf(`Store delete failed for entity "%s" recordID "%s" - %s`, entity.Label, recordID, err)
			}

			return nil
		}
	}

	return fmt.Errorf("Could not find a deletable database")
}

func marshalClientRecordToStoreRecord(entity *Entity, clientRecord *ClientRecord, action string) (storeRecord StoreRecord, err error) {
	storeRecord = StoreRecord{}
	for i, _ := range entity.Elements {
		element := &entity.Elements[i]
		field := &Field{
			ID: element.ID,
		}

		for _, keyValue := range clientRecord.KeyValues {
			if action == ACTION_POST && element.PrimaryKey == true {
				continue
			}
			val := keyValue.Value
			if element.DataType == ELEMENT_DATA_TYPE_NUMBER {
				f, err := strconv.ParseFloat(val.(string), 64)
				if err == nil {
					val = f
				}
			}
			if keyValue.Key == element.ID {
				field.Value = val
				field.Hydrated = true
				break
			}
		}
		storeRecord[element.ID] = field
	}
	return storeRecord, nil
}

func marshalStoreRecordToClientRecord(storeRecord StoreRecord) ClientRecord {
	clientRecord := ClientRecord{}
	kvs := KeyValues{}
	for id, field := range storeRecord {
		kv := KeyValue{Key: id, Value: field.Value}
		kvs = append(kvs, kv)
	}
	clientRecord.KeyValues = kvs
	return clientRecord
}
