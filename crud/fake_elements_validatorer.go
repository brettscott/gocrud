package crud

type fakeElementsValidatorer struct {
}

func (f *fakeElementsValidatorer) validate(entity Entity, record StoreRecord, action string) (success bool, elementsErrors map[string][]string, globalErrors []string) {
	return true, nil, nil
}
