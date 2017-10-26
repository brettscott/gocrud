package crud

// NewFakeMutatorer returns array of mutators
func NewFakeMutatorers() []mutatorer {
	return []mutatorer{
		&fakeMutatorer{},
	}
}

// NewFakeElementsValidatorers returns a elementsValidator
func NewFakeMutatorer() mutatorer {
	return &fakeMutatorer{}
}

type fakeMutatorer struct {
}

func (f *fakeMutatorer) mutate(entity *Entity, storeRecord StoreRecord, action string) (mutatedStoreRecord StoreRecord, elementsErrors map[string][]string, globalErrors []string) {
	return mutatedStoreRecord, nil, nil
}
