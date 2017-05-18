package store

type Storer interface {
	List()
	Get()
	Post()
	Put()
	Patch()
	Delete()
}
