package crud

// Element (eg id) is an attribute of Entity (eg users)
type Element struct {
	// Identifier eg id, description
	ID string

	// Label eg Name, Description
	Label ElementLabel

	// FormType eg hidden, text, select
	FormType ElementFormType

	// DataType of "Value" eg string,integer,boolean,keyValues
	DataType string

	// DefaultValue eg "1"
	DefaultValue interface{}

	// Validation rules
	Validation ElementValidation

	// Immutability of element in record.  When "true", cannot be changed after creation/POST
	Immutable bool

	// PrimaryKey is set to "true" when primary key for record
	PrimaryKey bool
}

type Elements []Element

type ElementLabel string
type ElementFormType string

const ELEMENT_FORM_TYPE_HIDDEN = "hidden"
const ELEMENT_FORM_TYPE_TEXT = "text"
const ELEMENT_FORM_TYPE_SELECT = "select"

const ELEMENT_DATA_TYPE_STRING = "string"
const ELEMENT_DATA_TYPE_NUMBER = "number"
const ELEMENT_DATA_TYPE_BOOLEAN = "boolean"
