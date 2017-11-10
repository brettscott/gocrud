package crud

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientErrors(t *testing.T) {

	t.Run("Has errors", func(t *testing.T) {

		t.Run("returns true when there is at least one elements error", func(t *testing.T) {
			elementsErrors := map[string][]string{}
			elementsErrors["id"] = []string{
				"must exist",
			}
			clientErrors := ClientErrors{
				ElementsErrors: elementsErrors,
			}

			assert.Equal(t, true, clientErrors.HasErrors())
		})

		t.Run("returns true when there is at least one global error", func(t *testing.T) {
			globalErrors := []string{
				"a non-element error occurred",
			}
			clientErrors := ClientErrors{
				GlobalErrors: globalErrors,
			}

			assert.Equal(t, true, clientErrors.HasErrors())
		})

		t.Run("returns true when there is at least one global error and one elements error", func(t *testing.T) {
			globalErrors := []string{
				"a non-element error occurred",
			}
			elementsErrors := map[string][]string{}
			elementsErrors["id"] = []string{
				"must exist",
			}
			clientErrors := ClientErrors{
				ElementsErrors: elementsErrors,
				GlobalErrors:   globalErrors,
			}

			assert.Equal(t, true, clientErrors.HasErrors())
		})
	})
}
