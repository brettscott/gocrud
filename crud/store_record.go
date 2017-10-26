package crud

import (
	"fmt"
	//"reflect"
)

// EntityData represents a database row from the entity's database
type StoreRecord []Field

// EntityDatum is a representation of a field in a database row of data from the database
type Field struct {
	ID       string
	Value    interface{}
	Hydrated bool
}

// GetField returns a particular field's (key's) value
func (r *StoreRecord) GetField(elementID string) (*Field, error) {

	for _, field := range *r {
		if field.ID == elementID {
			return &field, nil
		}
	}
	return nil, fmt.Errorf("Did not find elementID \"%s\" in list of fields", elementID)
}

// IsHydrated lets you know if any data (key-values) are attached to record
func (r *StoreRecord) IsHydrated() bool {
	if len(*r) == 0 {
		return false
	}
	return true
}
