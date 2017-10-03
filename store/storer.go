package store

import (
	"github.com/brettscott/gocrud/crud"
)

type Storer interface {
	List(entity crud.Entity) ([]Record, error)
	Get(e crud.Entity, recordID string) (Record, error)
	Post(entity crud.Entity, storeRecord Record) (string, error)
	Put(entity crud.Entity, storeRecord Record, recordID string) error
	Patch(entity crud.Entity, storeRecord Record, recordID string) error
	Delete(entity crud.Entity, recordID string) error
}
