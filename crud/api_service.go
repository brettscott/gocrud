package crud

import (
	"fmt"
)

type elementsValidatorer interface {
	validate(entity *Entity, record StoreRecord, action string) (success bool, elementsErrors map[string][]string, globalErrors []string)
}

type mutatorer interface {
	mutate(entity *Entity, storeRecord *StoreRecord, action string) (err error, elementsErrors map[string][]string, globalErrors []string)
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

func (a *apiService) save(entity *Entity, action string, clientRecord *ClientRecord, recordID string) (savedClientRecord ClientRecord, err error) {


	//fmt.Printf("%+v", clientRecord)

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
		err, elementsErrors, globalErrors := mutator.mutate(entity, &storeRecord, action)
		if err != nil {
			return savedClientRecord, fmt.Errorf(`Failed mutating for entity "%s" - %v`, entity.Label, err)
		}
		if len(elementsErrors) > 0 || len(globalErrors) > 0 {
			return savedClientRecord, fmt.Errorf(`Failed validation for entity "%s" - %v %v`, entity.Label, elementsErrors, globalErrors)
		}
	}

	// TODO support multiple stores - at the moment, once written to first DB, doesn't write to subsequent.  Add Goroutine
	for _, store := range a.stores {
		if store.Mode(entity).IsWritable() {
			switch action {
			case ACTION_POST:
				recordID, err = store.Post(entity, storeRecord)
				break
			case ACTION_PUT:
				if recordID == "" {
					return savedClientRecord, fmt.Errorf(`Missing record ID for entity "%s"`, entity.Label)
				}
				err = store.Put(entity, storeRecord, recordID)
				break
			case ACTION_PATCH:
				if recordID == "" {
					return savedClientRecord, fmt.Errorf(`Missing record ID for entity "%s"`, entity.Label)
				}
				err = store.Patch(entity, storeRecord, recordID)
				break
			default:
				return savedClientRecord, fmt.Errorf(`Invalid action "%s" for entity "%s"`, action, entity.Label)
				break
			}
			if err != nil {
				return savedClientRecord, fmt.Errorf(`Failed to "%s" for entity "%s" - %s`, action, entity.Label, err)
			}

			savedStoreRecord, err := store.Get(entity, recordID)
			if err != nil {
				return savedClientRecord, fmt.Errorf(`Failed to get newly created DB record for entity "%s" - %s`, entity.Label, err)
			}

			if savedStoreRecord.IsHydrated() == false {
				return savedClientRecord, fmt.Errorf(`New created DB record was not found in database for entity "%s" - %s`, entity.Label, err)
			}

			savedClientRecord = marshalStoreRecordToClientRecord(savedStoreRecord)
			return savedClientRecord, nil
		}
	}

	return savedClientRecord, fmt.Errorf("Failed to find a writeable database")
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
			if keyValue.Key == element.ID {
				field.Value = keyValue.Value
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
