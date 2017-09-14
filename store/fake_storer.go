package store

import "github.com/brettscott/gocrud/entity"

func NewFakeStorer() *FakeStorer {
	return &FakeStorer{}
}

// FakeStorer is a faked out storer
type FakeStorer struct{}

// List
func (f *FakeStorer) List() {}

// Get
func (f *FakeStorer) Get(e entity.Entity, recordID string) (entity.Record, error) {
	return entity.Record{}, nil
}

// Post
func (f *FakeStorer) Post(entity entity.Entity) (string, error) {
	return "", nil
}

// Put
func (f *FakeStorer) Put(entity entity.Entity) (string, error) {
	return "", nil
}

// Patch
func (f *FakeStorer) Patch() {}

// Delete
func (f *FakeStorer) Delete() {}
