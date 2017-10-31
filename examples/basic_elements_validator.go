package examples

import "github.com/brettscott/gocrud/crud"

type basicElementsValidator struct {
}

// validate will mark each element with a validation failure
func (m *basicElementsValidator) Validate(entity *crud.Entity, record crud.StoreRecord, action string) (success bool, clientErrors *crud.ClientErrors) {
	elementsErrors := map[string][]string{}
	globalErrors := []string{}
	for _, element := range entity.Elements {
		elementsErrors[element.ID] = []string{
			"I'm going fail for the sake of it",
		}
	}
	globalErrors = append(globalErrors, "a non-element specific error was identified")

	clientErrors = &crud.ClientErrors{} // instantiate only when there is an error
	clientErrors.ElementsErrors = elementsErrors
	clientErrors.GlobalErrors = globalErrors

	return !clientErrors.HasErrors(), clientErrors
}
