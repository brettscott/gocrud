package entity

import (
	"fmt"
	"reflect"
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

func (e *Entity) HydrateFromRecord(record *Record) error {
	for i, _ := range e.Elements {
		element := &e.Elements[i]
		for _, keyValue := range record.KeyValues {
			if keyValue.Key == element.ID {

				// TODO decide whether to do this:
				//if element.DataType != keyValue.DataType {
				//	return fmt.Errorf(`Record element "%s" has the wrong type - element: %s, record: %s`, element.name, element.DataType, keyValue.DataType)
				//}

				element.Value = keyValue.Value
				element.Hydrated = true
				break
			}
		}
	}
	return nil
}

func (e *Entity) Validate(action string) error {

	errors := make([]string, 0)
	for _, element := range e.Elements {

		err := validateDataType(element)
		if err != nil {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) has invalid data type: %s`, element.Label, element.ID, err))
		}

		if element.Validation.Required && element.Value == nil {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) is required and cannot be empty`, element.Label, element.ID))
		}
		if element.Validation.MustProvide == true && element.Hydrated == false {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) must be provided`, element.Label, element.ID))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("Validation errors: %v", errors)
	}

	return nil
}

func validateDataType(element Element) error {

	//var dataTypeMappings = []struct{
	//	name string
	//	golang        string
	//}{
	//	{ELEMENT_DATA_TYPE_STRING, "string"},
	//	{ELEMENT_DATA_TYPE_INTEGER, "int"},
	//}

	dataTypes := make(map[string]string)
	dataTypes[ELEMENT_DATA_TYPE_STRING] = "string"
	dataTypes[ELEMENT_DATA_TYPE_INTEGER] = "int"

	actualType := reflect.TypeOf(element.Value).String()
	expectedType := dataTypes[element.DataType]
	fmt.Printf("actual: %s, expected: %s", actualType, expectedType)
	if actualType != expectedType {
		return fmt.Errorf(`expected type to be "%s" but got "%s"`, expectedType, actualType)
	}

	//var ok bool
	//switch element.DataType {
	//case ELEMENT_DATA_TYPE_STRING:
	//	_, ok = element.Value.(string)
	//}
	//if !ok {
	//	return fmt.Errorf(`Not a "%s"`, element.DataType)
	//}
	//
	//for _, dataType := range dataTypeMappings {
	//	if element.DataType == dataType.name {
	//		_, ok := element.Value.(string)
	//		if !ok {
	//			return error(`Not a "string"`)
	//		}
	//	}
	//}
	return nil
}

const VALIDATE_ACTION_POST = "post"
const VALIDATE_ACTION_PUT = "put"

// CheckConfiguration makes sure the entity and its elements have a sensible configuration
// TODO this should be kicked off when application starts.  Create "NewEntity()".
func (e *Entity) CheckConfiguration() error {

	errors := make([]string, 0)
	for _, element := range e.Elements {

		if element.Validation.Required == true && element.Validation.MustProvide == false {
			errors = append(errors, fmt.Sprintf(`"%s"(%s) is "required" and "optional" which doesn't make sense.  Choose one.`, element.Label, element.ID))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("Configuration errors: %v", errors)
	}

	return nil

}