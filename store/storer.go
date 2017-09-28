package store

import "github.com/brettscott/gocrud/entity"

type Storer interface {
	List(entity entity.Entity) (entity.List, error)
	Get(e entity.Entity, recordID string) (entity.ClientRecord, error)
	Post(entity entity.Entity) (string, error)
	Put(entity entity.Entity, recordID string) error
	Patch(entity entity.Entity, recordID string) error
	Delete(entity entity.Entity, recordID string) error
}
