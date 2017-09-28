package entity

// List contains the results attributes for a given entity (eg User)
type List struct {
	Records []Record `json:"records"`
}
