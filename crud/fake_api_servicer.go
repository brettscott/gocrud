package crud

func NewFakeApiServicer() *fakeApiServicer {
	return &fakeApiServicer{}
}

type fakeApiServicer struct {
	listResponseBody    []ClientRecord
	listResponseError   error
	getResponseBody     ClientRecord
	getResponseError    error
	saveResponseBody    ClientRecord
	saveResponseError   error
	deleteResponseError error
}

func (f *fakeApiServicer) list(entity Entity) (clientRecords []ClientRecord, err error) {
	return f.listResponseBody, f.listResponseError
}

func (f *fakeApiServicer) get(entity Entity, recordID string) (clientRecord ClientRecord, err error) {
	return f.getResponseBody, f.getResponseError
}

func (f *fakeApiServicer) save(entity Entity, action string, clientRecord *ClientRecord, recordID string) (savedClientRecord ClientRecord, err error) {
	return f.saveResponseBody, f.saveResponseError
}

func (f *fakeApiServicer) delete(entity Entity, recordID string) error {
	return f.deleteResponseError
}
