package crud

// NewFakeElementsValidatorers returns array of elementsValidators
func NewFakeElementsValidatorers() []elementsValidatorer {
	return []elementsValidatorer{
		&fakeElementsValidatorer{},
	}
}

// NewFakeEmptyElementsValidatorers returns array of elementsValidators
func NewFakeEmptyElementsValidatorers() []elementsValidatorer {
	return []elementsValidatorer{}
}

// NewFakeElementsValidatorers returns a elementsValidator
func NewFakeElementsValidatorer() elementsValidatorer {
	return &fakeElementsValidatorer{}
}

type fakeElementsValidatorer struct {
	success bool
}

func (f *fakeElementsValidatorer) validate(entity *Entity, record StoreRecord, action string) (success bool, clientErrors *ClientErrors) {
	return f.success, clientErrors
}
