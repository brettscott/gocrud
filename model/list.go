package model

import "github.com/brettscott/gocrud/api"

// List contains the results attributes for a given entity (eg User)
type List struct {
	Records []api.Record `json:"records"`
}
