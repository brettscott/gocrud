package crud

func NewFakeApiServicer() *fakeApiServicer {
	return &fakeApiServicer{
		listCalled: false,
		getCalled: false,
	}
}

type fakeApiServicer struct {
	listResponseBody         ClientRecords
	listResponseError        error
	listCalled				 bool
	getResponseBody          ClientRecord
	getResponseError         error
	getRequestEntity	     *Entity
	getRequestRecordID	     string
	getCalled				 bool
	saveResponseBody         ClientRecord
	saveResponseClientErrors *ClientErrors
	saveResponseError        error
	deleteResponseError      error
}

func (f *fakeApiServicer) list(entity *Entity) (clientRecords ClientRecords, err error) {
	f.listCalled = true
	return f.listResponseBody, f.listResponseError
}

func (f *fakeApiServicer) get(entity *Entity, recordID string) (clientRecord ClientRecord, err error) {
	f.getCalled = true
	f.getRequestEntity = entity
	f.getRequestRecordID = recordID
	return f.getResponseBody, f.getResponseError
}

func (f *fakeApiServicer) save(entity *Entity, action string, clientRecord *ClientRecord, recordID string) (savedClientRecord ClientRecord, clientErrors *ClientErrors, err error) {
	return f.saveResponseBody, f.saveResponseClientErrors, f.saveResponseError
}

func (f *fakeApiServicer) delete(entity *Entity, recordID string) error {
	return f.deleteResponseError
}
