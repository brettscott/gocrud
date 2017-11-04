package examples

import (
	"github.com/brettscott/gocrud/crud"
)

func NewExampleStaticStore() *ExampleStaticStore {
	supermanStoreRecord := crud.StoreRecord{}
	supermanStoreRecord["id"] = &crud.Field{
		ID:       "id",
		Value:    "the-superman-id",
		Hydrated: true,
	}
	supermanStoreRecord["name"] = &crud.Field{
		ID:       "name",
		Value:    "Superman",
		Hydrated: true,
	}
	supermanStoreRecord["age"] = &crud.Field{
		ID:       "age",
		Value:    11,
		Hydrated: true,
	}

	catwomanStoreRecord := crud.StoreRecord{}
	catwomanStoreRecord["id"] = &crud.Field{
		ID:       "id",
		Value:    "the-catwoman-id",
		Hydrated: true,
	}
	catwomanStoreRecord["name"] = &crud.Field{
		ID:       "name",
		Value:    "Catwoman",
		Hydrated: true,
	}
	catwomanStoreRecord["age"] = &crud.Field{
		ID:       "age",
		Value:    11,
		Hydrated: true,
	}

	storeRecords := []crud.StoreRecord{
		catwomanStoreRecord,
		supermanStoreRecord,
	}

	return &ExampleStaticStore{
		fakeDatabase: storeRecords,
	}
}

type ExampleStaticStore struct {
	fakeDatabase []crud.StoreRecord
}

// Mode
func (e *ExampleStaticStore) Mode(entity *crud.Entity) *crud.StoreMode {
	return &crud.StoreMode{
		Read:   true,
		Write:  true,
		Delete: true,
	}
}

// List
func (e *ExampleStaticStore) List(entity *crud.Entity) ([]crud.StoreRecord, error) {
	return e.fakeDatabase, nil
}

// Get
func (e *ExampleStaticStore) Get(entity *crud.Entity, recordID string) (crud.StoreRecord, error) {
	for _, storeRecord := range e.fakeDatabase {
		if storeRecord["id"].Value == recordID {
		return storeRecord, nil
		}
	}
	return crud.StoreRecord{}, nil
}

// Post
func (e *ExampleStaticStore) Post(entity *crud.Entity, storeRecord crud.StoreRecord) (string, error) {
	// Does nothing!
	return "", nil
}

// Put
func (e *ExampleStaticStore) Put(entity *crud.Entity, storeRecord crud.StoreRecord, recordID string) error {
	// Does nothing!
	return nil
}

// Patch
func (e *ExampleStaticStore) Patch(entity *crud.Entity, storeRecord crud.StoreRecord, recordID string) error {
	// Does nothing!
	return nil
}

// Delete
func (e *ExampleStaticStore) Delete(entity *crud.Entity, recordID string) error {
	// Does nothing!
	return nil
}
