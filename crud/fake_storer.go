package crud

func NewFakeStorers(store Storer) []Storer {
	return []Storer{
		store,
	}
}

func NewFakeStorer() *FakeStorer {
	return &FakeStorer{}
}

// FakeStorer is a faked out storer
type FakeStorer struct {
	ListResponse []StoreRecord
	ListError    error
	GetResponse  StoreRecord
	GetError     error
}

func (f *FakeStorer) Mode(entity *Entity) *StoreMode {
	return &StoreMode{
		Read:  true,
		Write: true,
	}
}

// List
func (f *FakeStorer) List(entity *Entity) ([]StoreRecord, error) {
	return f.ListResponse, f.ListError
}

// Get
func (f *FakeStorer) Get(entity *Entity, recordID string) (StoreRecord, error) {
	return f.GetResponse, f.GetError
}

// Post
func (f *FakeStorer) Post(entity *Entity, storeRecord StoreRecord) (string, error) {
	return "", nil
}

// Put
func (f *FakeStorer) Put(entity *Entity, storeRecord StoreRecord, recordID string) error {
	return nil
}

// Patch
func (f *FakeStorer) Patch(entity *Entity, storeRecord StoreRecord, recordID string) error {
	return nil
}

// Delete
func (f *FakeStorer) Delete(entity *Entity, recordID string) error {
	return nil
}
