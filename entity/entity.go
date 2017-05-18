package entity

// Entity eg User
type Entity struct {
	ID       string
	Label    string
	Labels   string
	Elements []Element
	Form Form
	List List
}

// Entities eg Users
type Entities []Entity
