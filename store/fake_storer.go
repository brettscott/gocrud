package store

import (
	"github.com/brettscott/gocrud/crud"
)

func NewFakeStorer() *FakeStorer {
	return &FakeStorer{}
}

// FakeStorer is a faked out storer
type FakeStorer struct{}

// List
func (f *FakeStorer) List(e crud.Entity) ([]Record, error) {
	return []Record{}, nil
}

// Get
func (f *FakeStorer) Get(e crud.Entity, recordID string) (Record, error) {
	return Record{}, nil
}

// Post
func (f *FakeStorer) Post(entity crud.Entity) (string, error) {
	return "", nil
}

// Put
func (f *FakeStorer) Put(entity crud.Entity, recordID string) error {
	return nil
}

// Patch
func (f *FakeStorer) Patch(entity crud.Entity, recordID string) error {
	return nil
}

// Delete
func (f *FakeStorer) Delete(entity crud.Entity, recordID string) error {
	return nil
}
