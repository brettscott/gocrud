package api

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
)

func TestAPI_ServeHTTP(t *testing.T) {

	t.Run("can call API and get a response", func(t *testing.T) {

		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		w := httptest.NewRecorder()

		apiGateway := Gateway{}
		apiGateway.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
