package crud

import (
	"fmt"
	"github.com/brettscott/gocrud/model"
	"github.com/brettscott/gocrud/store"
	"reflect"
)

func validate(entity model.Entity, record store.Record, action string) error {

	errors := make([]string, 0)
	var primaryKey model.ElementLabel

	for _, element := range entity.Elements {

		userData, err := record.GetField(element.ID)
		if err != nil {
			errors = append(errors, fmt.Sprintf(`Missing element "%s" - %v`, element.ID, err))
		}

		if err := validateDataType(element, userData.Value); err != nil {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) has invalid data type: %s`, element.Label, element.ID, err))
		}

		// This is useful to see if value was provided and whether a string is empty or not.  Use "Min" and "Max" for integers.
		// Don't use anything for boolean because it'll either be true or false (or "nil" and be classed as not provided).
		if element.Validation.Required && (userData.Hydrated == false || userData.Value == nil || userData.Value == "") {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) is required and cannot be empty`, element.Label, element.ID))
		}

		if element.Validation.MustProvide == true && userData.Hydrated == false {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) must be provided`, element.Label, element.ID))
		}

		if element.PrimaryKey == true {
			if primaryKey != "" {
				errors = append(errors, fmt.Sprintf(`"%s" (%s) cannot be a primary key because "%s" is already one`, element.Label, element.ID, primaryKey))
			} else {
				primaryKey = element.Label
			}
		}

		if action != ACTION_PATCH && element.PrimaryKey != true && userData.Hydrated == false {
			errors = append(errors, fmt.Sprintf(`"%s" (%s) was not supplied on "%s"`, element.Label, element.ID, action))
		}
	}

	if primaryKey == "" {
		errors = append(errors, fmt.Sprintf(`Missing a primary key element`))
	}

	if len(errors) > 0 {
		return fmt.Errorf("Validation errors: %v", errors)
	}

	return nil
}

// validateDataType
// Unmarshal stores one of these in the interface value: "bool" for JSON booleans, "float64" for JSON numbers,
// "string" for JSON strings, "[]interface{}" for JSON arrays, "map[string]interface{}" for JSON objects,  "nil" for JSON null
func validateDataType(element model.Element, value interface{}) error {
	if value == nil {
		return nil
	}

	// Todo Move out of here so it's only created once!
	dataTypes := make(map[string]string)
	dataTypes[model.ELEMENT_DATA_TYPE_STRING] = "string"
	dataTypes[model.ELEMENT_DATA_TYPE_NUMBER] = "float64"
	dataTypes[model.ELEMENT_DATA_TYPE_BOOLEAN] = "bool"

	if _, ok := dataTypes[element.DataType]; !ok {
		return fmt.Errorf(`undefined data type "%s"`, element.DataType)
	}

	actualType := reflect.TypeOf(value).String()
	expectedType := dataTypes[element.DataType]
	if actualType != expectedType {
		return fmt.Errorf(`expected type to be "%s" but got "%s" with value "%v"`, expectedType, actualType, value)
	}

	return nil
}
