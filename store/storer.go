package store

import (
	"github.com/brettscott/gocrud/model"
)

type Storer interface {
	List(entity model.Entity) ([]Record, error)
	Get(e model.Entity, recordID string) (Record, error)
	Post(entity model.Entity, storeRecord Record) (string, error)
	Put(entity model.Entity, storeRecord Record, recordID string) error
	Patch(entity model.Entity, storeRecord Record, recordID string) error
	Delete(entity model.Entity, recordID string) error
}
