package crud

type mutatorer interface {
	mutate(entity Entity, storeRecord StoreRecord, action string) (mutatedStoreRecord StoreRecord, elementsErrors map[string][]string, globalErrors []string)
}

type mutator struct {
}

func (m *mutator) mutate(entity Entity, storeRecord StoreRecord, action string) (mutatedStoreRecord StoreRecord, elementsErrors map[string][]string, globalErrors []string) {
	return storeRecord, elementsErrors, globalErrors
}
