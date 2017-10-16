package crud

import (
	"github.com/brettscott/gocrud/model"
)

func NewFakeApiServicer() *fakeApiServicer {
	return &fakeApiServicer{}
}

type fakeApiServicer struct {
	listResponseBody    []Record
	listResponseError   error
	getResponseBody     Record
	getResponseError    error
	saveResponseBody    Record
	saveResponseError   error
	deleteResponseError error
}

func (f *fakeApiServicer) list(entity model.Entity) (clientRecords []Record, err error) {
	return f.listResponseBody, f.listResponseError
}

func (f *fakeApiServicer) get(entity model.Entity, recordID string) (clientRecord Record, err error) {
	return f.getResponseBody, f.getResponseError
}

func (f *fakeApiServicer) save(entity model.Entity, action string, clientRecord *Record, recordID string) (savedClientRecord Record, err error) {
	return f.saveResponseBody, f.saveResponseError
}

func (f *fakeApiServicer) delete(entity model.Entity, recordID string) error {
	return f.deleteResponseError
}
