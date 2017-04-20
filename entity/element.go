package entity


type Element struct {
	ID string
	Label ElementLabel
	FormType ElementFormType
	Value interface{}
	ValueType string  // type of "Value"
	DefaultValue interface{}

}

type ElementLabel string
type ElementFormType string

const ELEMENT_FORM_TYPE_TEXT = "text"
const ELEMENT_FORM_TYPE_SELECT = "select"

const ELEMENT_VALUE_TYPE_STRING = "string"
const ELEMENT_VALUE_TYPE_INTEGER = "integer"
const ELEMENT_VALUE_TYPE_BOOLEAN = "boolean"

//ElementTypes := []string {
//	"text",
//	"select",
//}

