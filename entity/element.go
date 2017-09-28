package entity

// Element (eg id) is an attribute of Entity (eg users)
type Element struct {
	// Identifier eg id, description
	ID string

	// Label eg Name, Description
	Label ElementLabel

	// FormType eg hidden, text, select
	FormType ElementFormType

	// Value of element in record
	Value interface{} // TODO should this be private as we don't want it specified on instantiation??

	// DataType of "Value" eg string,integer,boolean,keyValues
	DataType string

	// DefaultValue eg "1"
	DefaultValue interface{}

	// Validation rules
	Validation ElementValidation

	// Immutability of element in record.  When "true", cannot be changed after creation/POST
	Immutable bool

	// Hydrated indicates whether this record was populated by user/API
	Hydrated bool

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

//ElementTypes := []string {
//	"text",
//	"select",
//}

// ElementValidation
// Inspired by Joi https://github.com/hapijs/joi
type ElementValidation struct {
	// Required means the element must not be a zero value (0, false, "")
	Required bool

	// MustProvide means the element must be submitted on every POST and PUT
	MustProvide bool

	//Forbidden bool // must not be sent
	//Strip bool  // remove from output after validation
	//Any bool  // any data type
	//String bool
	//Min int
	//Max int
	//Allow string // whitelist
	//Disallow string // blacklist
	// TODO ... loads more from Joi ...
}
