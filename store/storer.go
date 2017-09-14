package store

import "github.com/brettscott/gocrud/entity"

type Storer interface {
	List()
	Get(e entity.Entity, recordID string) (entity.Record, error)
	Post(entity entity.Entity) (string, error)
	Put(entity entity.Entity, recordID string) error
	Patch()
	Delete()
}
