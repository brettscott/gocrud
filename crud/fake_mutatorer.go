package crud

// newFakeMutatorers returns a slice of mutators
func newFakeMutatorers() []mutatorer {
	return []mutatorer{
		&fakeMutatorer{},
	}
}

// newFakeEmptyMutatorers returns an empty slice of mutators
func newFakeEmptyMutatorers() []mutatorer {
	return []mutatorer{}
}

// newFakeElementsValidatorers returns a elementsValidator
func newFakeMutatorer() mutatorer {
	return &fakeMutatorer{}
}

type fakeMutatorer struct {
}

func (f *fakeMutatorer) mutate(entity *Entity, storeRecord *StoreRecord, action string) (err error, elementsErrors map[string][]string, globalErrors []string) {
	return nil, nil, nil
}
