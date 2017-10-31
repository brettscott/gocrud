package crud

// NewBasicElementsValidator creates the basic elements validator.
// This is elements validator is here only as an example
func NewBasicElementsValidator() elementsValidatorer {
	return &basicElementsValidator{}
}

type basicElementsValidator struct {
}

// validate will mark each element with a validation failure
func (m *basicElementsValidator) validate(entity *Entity, record StoreRecord, action string) (success bool, clientErrors *ClientErrors) {
	elementsErrors := map[string][]string{}
	globalErrors := []string{}
	for _, element := range entity.Elements {
		elementsErrors[element.ID] = []string{
			"I'm going fail for the sake of it",
		}
	}
	globalErrors = append(globalErrors, "a non-element specific error was identified")

	clientErrors = &ClientErrors{} // instantiate only when there is an error
	clientErrors.ElementsErrors = elementsErrors
	clientErrors.GlobalErrors = globalErrors

	return !clientErrors.HasErrors(), clientErrors
}
