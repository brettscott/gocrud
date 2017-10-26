package crud

import (
	"fmt"
	"strings"
)

func NewBasicMutator() mutatorer {
	return &basicMutator{}
}

type basicMutator struct {
}

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
