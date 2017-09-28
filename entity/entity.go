package entity

import (
	"fmt"
)

// Entity eg User
type Entity struct {
	ID       string
	Label    string
	Labels   string
	Elements Elements
	Form     Form
	List     List
}

// Entities eg Users
type Entities map[string]Entity

func (e *Entity) GetElement(elementID string) (*Element, error) {

	for _, element := range e.Elements {
		if element.ID == elementID {
			return &element, nil
		}
	}
	return nil, fmt.Errorf("Did not find elementID \"%s\" in list of elements", elementID)
}

const HYDRATE_FROM_RECORD_ACTION_POST = "post"
const HYDRATE_FROM_RECORD_ACTION_PUT = "put"
const HYDRATE_FROM_RECORD_ACTION_PATCH = "patch"

const VALIDATE_ACTION_POST = "post"
const VALIDATE_ACTION_PUT = "put"
const VALIDATE_ACTION_PATCH = "patch"

// CheckConfiguration makes sure the entity and its elements have a sensible configuration
// TODO this should be kicked off when application starts.  Create "NewEntity()".
func (e *Entity) CheckConfiguration() error {

	// Todo Move out of here so it's only created once!
	dataTypes := make(map[string]string)
	dataTypes[ELEMENT_DATA_TYPE_STRING] = "string"
	dataTypes[ELEMENT_DATA_TYPE_NUMBER] = "float64"
	dataTypes[ELEMENT_DATA_TYPE_BOOLEAN] = "bool"

	errors := make([]string, 0)
	for _, element := range e.Elements {

		if element.Validation.Required == true && element.Validation.MustProvide == false {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) is "required" and "optional" which doesn't make sense.  Choose one.`, element.Label, element.ID))
		}
		if _, exists := dataTypes[element.DataType]; !exists {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) has an invalid data type of "%s".`, element.Label, element.ID, element.DataType))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("Configuration errors: %v", errors)
	}

	return nil

}
