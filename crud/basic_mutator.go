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

// mutate will trim whitespace from beginning and end of all element values which are strings
func (m *basicMutator) mutate(entity *Entity, storeRecord *StoreRecord, action string) (err error, elementsErrors map[string][]string, globalErrors []string) {

	for id, field := range *storeRecord {
		fmt.Println("id", id)
		element, err := entity.GetElement(id)
		if err != nil {
			return fmt.Errorf("Element not found: %s", id), elementsErrors, globalErrors
		}

		if field.Hydrated == true && element.DataType == ELEMENT_DATA_TYPE_STRING {
			field.Value = strings.TrimSpace(field.Value.(string))
		}
	}

	return err, elementsErrors, globalErrors
}
