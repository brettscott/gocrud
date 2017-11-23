package crud

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplateService(t *testing.T) {

	t.Run("should render registered template", func(t *testing.T) {
		template := newTemplateService()
		ctx := map[string]interface{}{}

		html, err := template.exec("root", ctx)

		assert.NoError(t, err)
		assert.NotEmpty(t, html)
	})

	t.Run("should error to render non-registered template", func(t *testing.T) {
		template := newTemplateService()
		ctx := map[string]interface{}{}

		_, err := template.exec("i-do-not-exist-in-register", ctx)

		assert.Error(t, err)
	})
}
