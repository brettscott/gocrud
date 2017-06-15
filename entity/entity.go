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

const HYDRATE_FROM_RECORD_ACTION_POST = "post"
const HYDRATE_FROM_RECORD_ACTION_PUT = "put"
const HYDRATE_FROM_RECORD_ACTION_PATCH = "patch"

func (e *Entity) HydrateFromRecord(record *Record, action string) error {
	for i, _ := range e.Elements {
		element := &e.Elements[i]

		for _, keyValue := range record.KeyValues {

			if action == HYDRATE_FROM_RECORD_ACTION_POST && element.PrimaryKey == true {
				continue
			}

			if keyValue.Key == element.ID {
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

		if err := e.validateDataType(element); err != nil {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) has invalid data type: %s`, element.Label, element.ID, err))
		}

		// This is useful to see if value was provided and whether a string is empty or not.  Use "Min" and "Max" for integers.
		// Don't use anything for boolean because it'll either be true or false (or "nil" and be classed as not provided).
		if element.Validation.Required && (element.Hydrated == false || element.Value == nil || element.Value == "") {
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

func (e *Entity) validateDataType(element Element) error {

	if element.Value == nil {
		return nil
	}

	// Todo Move out of here so it's only created once!
	dataTypes := make(map[string]string)
	dataTypes[ELEMENT_DATA_TYPE_STRING] = "string"
	dataTypes[ELEMENT_DATA_TYPE_INTEGER] = "int"
	dataTypes[ELEMENT_DATA_TYPE_BOOLEAN] = "bool"

	if _, ok := dataTypes[element.DataType]; !ok {
		return fmt.Errorf(`undefined data type "%s"`, element.DataType)
	}

	actualType := reflect.TypeOf(element.Value).String()
	expectedType := dataTypes[element.DataType]
	if actualType != expectedType {
		return fmt.Errorf(`expected type to be "%s" but got "%s" with value "%v"`, expectedType, actualType, element.Value)
	}

	return nil
}

const VALIDATE_ACTION_POST = "post"
const VALIDATE_ACTION_PUT = "put"
const VALIDATE_ACTION_PATCH = "patch"

// CheckConfiguration makes sure the entity and its elements have a sensible configuration
// TODO this should be kicked off when application starts.  Create "NewEntity()".
func (e *Entity) CheckConfiguration() error {

	// Todo Move out of here so it's only created once!
	dataTypes := make(map[string]string)
	dataTypes[ELEMENT_DATA_TYPE_STRING] = "string"
	dataTypes[ELEMENT_DATA_TYPE_INTEGER] = "int"
	dataTypes[ELEMENT_DATA_TYPE_BOOLEAN] = "bool"

	errors := make([]string, 0)
	for _, element := range e.Elements {

		if element.Value != nil {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) should not have "Value" attribute specified.  Value: "%v"`, element.Label, element.ID, element.Value))
		}
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
