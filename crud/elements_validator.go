package crud

import (
	"fmt"
	"reflect"
)

func NewElementsValidator() *elementsValidator {
	dataTypes := make(map[string]string)
	dataTypes[ELEMENT_DATA_TYPE_STRING] = "string"
	dataTypes[ELEMENT_DATA_TYPE_NUMBER] = "float64"
	dataTypes[ELEMENT_DATA_TYPE_BOOLEAN] = "bool"

	return &elementsValidator{
		dataTypes: dataTypes,
	}
}

type elementsValidator struct {
	dataTypes map[string]string
}

// validate ensures values supplied by users are valid
func (e *elementsValidator) validate(entity *Entity, record StoreRecord, action string) (success bool, elementsErrors map[string][]string, globalErrors []string) {
	success = true
	var primaryKey ElementLabel
	elementsErrors = map[string][]string{}

	for _, element := range entity.Elements {
		elementErrors := make([]string, 0)
		userData, ok := record[element.ID]
		if !ok {
			elementErrors = append(elementErrors, "is missing")
		}

		if err := e.validateDataType(element, userData.Value); err != nil {
			elementErrors = append(elementErrors, fmt.Sprintf(`has invalid data type: %s`, err))
		}

		// This is useful to see if value was provided and whether a string is empty or not.  Use "Min" and "Max" for integers.
		// Don't use anything for boolean because it'll either be true or false (or "nil" and be classed as not provided).
		if element.Validation.Required && element.PrimaryKey == false && (action == ACTION_POST || userData.Hydrated == true) && (userData.Hydrated == false || userData.Value == nil || userData.Value == "") {
			elementErrors = append(elementErrors, "is required and cannot be empty")
		}

		if element.Validation.MustProvide == true && userData.Hydrated == false {
			elementErrors = append(elementErrors, "must be provided")
		} else if action == ACTION_POST && element.Validation.MustProvideOnPost == true && userData.Hydrated == false {
			elementErrors = append(elementErrors, "must be provided on POST")
		} else if action == ACTION_PUT && element.Validation.MustProvideOnPut == true && userData.Hydrated == false {
			elementErrors = append(elementErrors, "must be provided on PUT")
		} else if action == ACTION_PATCH && element.Validation.MustProvideOnPatch == true && userData.Hydrated == false {
			elementErrors = append(elementErrors, "must be provided on PATCH")
		}

		if element.PrimaryKey == true {
			if primaryKey != "" {
				elementErrors = append(elementErrors, fmt.Sprintf(`cannot be a primary key because "%s" is already one`, primaryKey))
			} else {
				primaryKey = element.Label
			}
		}

		if len(elementErrors) > 0 {
			elementsErrors[element.ID] = []string{}
			elementsErrors[element.ID] = elementErrors
		}
	}

	if primaryKey == "" {
		globalErrors = append(globalErrors, "Missing a primary key")
	}

	if len(elementsErrors) > 0 || len(globalErrors) > 0 {
		success = false
	}

	return
}

// validateDataType
// Unmarshal stores one of these in the interface value: "bool" for JSON booleans, "float64" for JSON numbers,
// "string" for JSON strings, "[]interface{}" for JSON arrays, "map[string]interface{}" for JSON objects,  "nil" for JSON null
func (e *elementsValidator) validateDataType(element Element, value interface{}) error {
	if value == nil {
		return nil
	}
	if _, ok := e.dataTypes[element.DataType]; !ok {
		return fmt.Errorf(`undefined data type "%s"`, element.DataType)
	}

	actualType := reflect.TypeOf(value).String()
	expectedType := e.dataTypes[element.DataType]
	if actualType != expectedType {
		return fmt.Errorf(`expected type to be "%s" but got "%s" with value "%v"`, expectedType, actualType, value)
	}
	return nil
}
