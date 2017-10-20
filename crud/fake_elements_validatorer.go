package crud

import (
	"github.com/brettscott/gocrud/model"
	"github.com/brettscott/gocrud/store"
)

type fakeElementsValidatorer struct {
}

func (f *fakeElementsValidatorer) validate(entity model.Entity, record store.Record, action string) (success bool, elementsErrors map[string][]string, globalErrors []string) {
	return true, nil, nil
}
