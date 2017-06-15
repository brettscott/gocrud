package store

import "github.com/brettscott/gocrud/entity"

type Storer interface {
	List()
	Get()
	Post(entity entity.Entity) (string, error)
	Put()
	Patch()
	Delete()
}
