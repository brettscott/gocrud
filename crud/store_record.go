package crud

import "fmt"

// EntityData represents a database row from the entity's database
//type StoreRecord []Field
type StoreRecord map[string]*Field

// EntityDatum is a representation of a field in a database row of data from the database
type Field struct {
	ID       string
	Value    interface{}
	Hydrated bool
}

// GetValue
func (r *StoreRecord) GetValue(key string) (interface{}, error) {
	if field, ok := (*r)[key]; ok {
		//fmt.Printf("\n\nField: %+v\n\n", field)
		return (*field).Value, nil
	}
	return nil, fmt.Errorf("Key \"%s\" not found in store record", key)
}

// IsHydrated lets you know if any data (key-values) are attached to record
func (r *StoreRecord) IsHydrated() bool {
	if len(*r) == 0 {
		return false
	}
	return true
}
