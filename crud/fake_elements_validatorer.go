package crud

// NewFakeElementsValidatorers returns array of elementsValidators
func NewFakeElementsValidatorers() []elementsValidatorer {
	return []elementsValidatorer{
		&FakeElementsValidatorer{},
	}
}

// NewFakeEmptyElementsValidatorers returns array of elementsValidators
func NewFakeEmptyElementsValidatorers() []elementsValidatorer {
	return []elementsValidatorer{}
}

// NewFakeElementsValidatorers returns a elementsValidator
func NewFakeElementsValidatorer() elementsValidatorer {
	return &FakeElementsValidatorer{}
}

type FakeElementsValidatorer struct {
	Success      bool
	ClientErrors *ClientErrors
}

func (f *FakeElementsValidatorer) Validate(entity *Entity, record StoreRecord, action string) (success bool, clientErrors *ClientErrors) {
	return f.Success, f.ClientErrors
}
