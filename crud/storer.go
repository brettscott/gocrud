package crud

type Storer interface {
	List(entity *Entity) ([]StoreRecord, error)
	Get(e *Entity, recordID string) (StoreRecord, error)
	Post(entity *Entity, storeRecord StoreRecord) (string, error)
	Put(entity *Entity, storeRecord StoreRecord, recordID string) error
	Patch(entity *Entity, storeRecord StoreRecord, recordID string) error
	Delete(entity *Entity, recordID string) error
}
