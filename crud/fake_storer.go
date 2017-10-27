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
	ListResponse     []StoreRecord
	ListError        error
	GetResponse      StoreRecord
	GetError         error
	PostCalled       int
	PostStoreRecord  StoreRecord
	PostError        error
	PostRecordID     string
	PutStoreRecord   StoreRecord
	PutError         error
	PatchStoreRecord StoreRecord
	PatchError       error
}

// Mode
func (f *FakeStorer) Mode(entity *Entity) *StoreMode {
	return &StoreMode{
		Read:   true,
		Write:  true,
		Delete: true,
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
	f.PostCalled++
	f.PostStoreRecord = storeRecord
	return f.PostRecordID, f.PostError
}

// Put
func (f *FakeStorer) Put(entity *Entity, storeRecord StoreRecord, recordID string) error {
	f.PutStoreRecord = storeRecord
	return f.PostError
}

// Patch
func (f *FakeStorer) Patch(entity *Entity, storeRecord StoreRecord, recordID string) error {
	f.PatchStoreRecord = storeRecord
	return f.PostError
}

// Delete
func (f *FakeStorer) Delete(entity *Entity, recordID string) error {
	return nil
}
