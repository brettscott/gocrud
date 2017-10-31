package examples

import (
	"fmt"
	"github.com/brettscott/gocrud/crud"
	"strings"
)

type basicMutator struct {
}

// mutate will trim whitespace from beginning and end of all element values which have the data type of a "string"
func (m *basicMutator) Mutate(entity *crud.Entity, storeRecord *crud.StoreRecord, action string) (clientErrors *crud.ClientErrors, err error) {
	for id, field := range *storeRecord {
		element, err := entity.GetElement(id)
		if err != nil {
			return clientErrors, fmt.Errorf("Element not found: %s", id)
		}

		if field.Hydrated == true && element.DataType == crud.ELEMENT_DATA_TYPE_STRING {
			field.Value = strings.TrimSpace(field.Value.(string))
		}
	}
	return clientErrors, nil
}
