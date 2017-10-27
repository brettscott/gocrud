package crud

import (
	"fmt"
	"strings"
)

// NewBasicMutator creates the basic mutator.
// This is mutator is here only as an example
func NewBasicMutator() mutatorer {
	return &basicMutator{}
}

type basicMutator struct {
}

// mutate will trim whitespace from beginning and end of all element values which have the data type of a "string"
func (m *basicMutator) mutate(entity *Entity, storeRecord *StoreRecord, action string) (clientErrors *ClientErrors, err error) {
	for id, field := range *storeRecord {
		element, err := entity.GetElement(id)
		if err != nil {
			return clientErrors, fmt.Errorf("Element not found: %s", id)
		}

		if field.Hydrated == true && element.DataType == ELEMENT_DATA_TYPE_STRING {
			field.Value = strings.TrimSpace(field.Value.(string))
		}
	}
	return clientErrors, nil
}
