package crud

type ElementsErrors map[string][]string
type GlobalErrors []string

func newClientErrors(elementsErrors ElementsErrors, globalErrors GlobalErrors) *ClientErrors {
	return &ClientErrors{
		ElementsErrors: elementsErrors,
		GlobalErrors:   globalErrors,
	}
}

type ClientErrors struct {
	ElementsErrors map[string][]string `json:"elementsErrors"`
	GlobalErrors   []string            `json:"globalErrors"`
}

// HasErrors returns true when there are element or global errors present
func (c *ClientErrors) HasErrors() bool {
	if len(c.ElementsErrors) > 0 || len(c.GlobalErrors) > 0 {
		return true
	}
	return false
}
