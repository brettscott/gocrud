package store

import (
	"github.com/brettscott/gocrud/model"
)

func NewFakeStorer() *FakeStorer {
	return &FakeStorer{}
}

// FakeStorer is a faked out storer
type FakeStorer struct{}

// List
func (f *FakeStorer) List(entity model.Entity) ([]Record, error) {
	return []Record{}, nil
}

// Get
func (f *FakeStorer) Get(e model.Entity, recordID string) (Record, error) {
	return Record{}, nil
}

// Post
func (f *FakeStorer) Post(entity model.Entity, storeRecord Record) (string, error) {
	return "", nil
}

// Put
func (f *FakeStorer) Put(entity model.Entity, storeRecord Record, recordID string) error {
	return nil
}

// Patch
func (f *FakeStorer) Patch(entity model.Entity, storeRecord Record, recordID string) error {
	return nil
}

// Delete
func (f *FakeStorer) Delete(entity model.Entity, recordID string) error {
	return nil
}
