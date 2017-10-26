package crud

// NewFakeElementsValidatorers returns array of elementsValidators
func NewFakeElementsValidatorers() []elementsValidatorer {
	return []elementsValidatorer{
		&fakeElementsValidatorer{},
	}
}

// NewFakeElementsValidatorers returns a elementsValidator
func NewFakeElementsValidatorer() elementsValidatorer {
	return &fakeElementsValidatorer{}
}

type fakeElementsValidatorer struct {
}

func (f *fakeElementsValidatorer) validate(entity *Entity, record StoreRecord, action string) (success bool, elementsErrors map[string][]string, globalErrors []string) {
	return true, nil, nil
}
