package crud

// newFakeMutatorers returns a slice of mutators
func newFakeMutatorers() []mutatorer {
	return []mutatorer{
		&FakeMutatorer{},
	}
}

// newFakeEmptyMutatorers returns an empty slice of mutators
func newFakeEmptyMutatorers() []mutatorer {
	return []mutatorer{}
}

// newFakeElementsValidatorers returns a elementsValidator
func newFakeMutatorer() mutatorer {
	return &FakeMutatorer{}
}

type FakeMutatorer struct {
	Err          error
	ClientErrors *ClientErrors
	StoreRecord  *StoreRecord
}

func (f *FakeMutatorer) Mutate(entity *Entity, storeRecord *StoreRecord, action string) (clientErrors *ClientErrors, err error) {
	if f.StoreRecord != nil {
		storeRecord = f.StoreRecord
	}
	return f.ClientErrors, f.Err
}
