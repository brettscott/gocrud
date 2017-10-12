package crud

import (
	"github.com/brettscott/gocrud/model"
)

func NewFakeApiServicer() *fakeApiServicer {
	return &fakeApiServicer{}
}

type fakeApiServicer struct {
	listResponseBody []byte
	listResponseError error
	getResponseBody []byte
	getResponseError error
	saveResponseBody []byte
	saveResponseError error
	deleteResponseError error
}

func (f *fakeApiServicer) list(entity model.Entity) (jsonResponse []byte, err error) {
	return f.listResponseBody, f.listResponseError
}

func (f *fakeApiServicer) get(entity model.Entity, recordID string) (jsonResponse []byte, err error) {
	return f.getResponseBody, f.getResponseError
}

func (f *fakeApiServicer) save(entity model.Entity, action string, body []byte, recordID string) (jsonResponse []byte, err error) {
	return f.saveResponseBody, f.saveResponseError
}

func (f *fakeApiServicer) delete(entity model.Entity, recordID string) error {
	return f.deleteResponseError
}
