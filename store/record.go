package store

import (
	"fmt"
	//"reflect"
)

// EntityData represents a database row from the entity's database
type Record []Field

// EntityDatum is a representation of a field in a database row of data from the database
type Field struct {
	ID       string
	Value    interface{}
	Hydrated bool
}

func (r *Record) GetField(elementID string) (*Field, error) {

	for _, field := range r {
		if field.ID == elementID {
			return *field, nil
		}
	}
	return nil, fmt.Errorf("Did not find elementID \"%s\" in list of fields", elementID)
}


//func NewEntityData(elements Elements) (elementsData EntityData) {
//	for _, element := range elements {
//		elementsData = append(elementsData, EntityDatum{ElementID: element.ID})
//	}
//	return elementsData
//}

//// HydrateFromRecord hydrates entity with record data (record data is usually marshalled from JSON to ClientRecord struct)
//func (d EntityData) HydrateFromRecord(record *ClientRecord, action string) error {
//	for i, _ := range d {
//		element := d[i].Element
//
//		for _, keyValue := range record.KeyValues {
//
//			if action == HYDRATE_FROM_RECORD_ACTION_POST && element.PrimaryKey == true {
//				continue
//			}
//
//			if keyValue.Key == element.ID {
//				d[i].Value = keyValue.Value
//				element.Hydrated = true
//				break
//			}
//		}
//	}
//	return nil
//}
//
//func (d EntityData) Validate(entityData EntityData, action string) error {
//
//	errors := make([]string, 0)
//	var primaryKey ElementLabel
//
//	for i, _ := range d {
//		element := d[i].Element
//
//		//if err := d.validateDataType(element); err != nil {
//		//	errors = append(errors, fmt.Sprintf(`"%s" (%s) has invalid data type: %s`, element.Label, element.ID, err))
//		//}
//
//		// This is useful to see if value was provided and whether a string is empty or not.  Use "Min" and "Max" for integers.
//		// Don't use anything for boolean because it'll either be true or false (or "nil" and be classed as not provided).
//		if element.Validation.Required && (element.Hydrated == false || d[i].Value == nil || d[i].Value == "") {
//			errors = append(errors, fmt.Sprintf(`"%s" (%s) is required and cannot be empty`, element.Label, element.ID))
//		}
//
//		if element.Validation.MustProvide == true && element.Hydrated == false {
//			errors = append(errors, fmt.Sprintf(`"%s" (%s) must be provided`, element.Label, element.ID))
//		}
//
//		if element.PrimaryKey == true {
//			if primaryKey != "" {
//				errors = append(errors, fmt.Sprintf(`"%s" (%s) cannot be a primary key because "%s" is already one`, element.Label, element.ID, primaryKey))
//			} else {
//				primaryKey = element.Label
//			}
//		}
//
//		if action != HYDRATE_FROM_RECORD_ACTION_PATCH && element.PrimaryKey != true && element.Hydrated == false {
//			errors = append(errors, fmt.Sprintf(`"%s" (%s) was not supplied on "%s"`, element.Label, element.ID, action))
//		}
//	}
//
//	if primaryKey == "" {
//		errors = append(errors, fmt.Sprintf(`Missing a primary key element`))
//	}
//
//	if len(errors) > 0 {
//		return fmt.Errorf("Validation errors: %v", errors)
//	}
//
//	return nil
//}

// validateDataType
// Unmarshal stores one of these in the interface value: "bool" for JSON booleans, "float64" for JSON numbers,
// "string" for JSON strings, "[]interface{}" for JSON arrays, "map[string]interface{}" for JSON objects,  "nil" for JSON null
//func (d *EntityData) validateDataType(element Element) error {
//	if element.Value == nil {
//		return nil
//	}
//
//	// Todo Move out of here so it's only created once!
//	dataTypes := make(map[string]string)
//	dataTypes[ELEMENT_DATA_TYPE_STRING] = "string"
//	dataTypes[ELEMENT_DATA_TYPE_NUMBER] = "float64"
//	dataTypes[ELEMENT_DATA_TYPE_BOOLEAN] = "bool"
//
//	if _, ok := dataTypes[element.DataType]; !ok {
//		return fmt.Errorf(`undefined data type "%s"`, element.DataType)
//	}
//
//	actualType := reflect.TypeOf(element.Value).String()
//	expectedType := dataTypes[element.DataType]
//	if actualType != expectedType {
//		return fmt.Errorf(`expected type to be "%s" but got "%s" with value "%v"`, expectedType, actualType, element.Value)
//	}
//
//	return nil
//}
