package crud

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecord(t *testing.T) {

	idField := &Field{
		ID:    "id",
		Value: 12345,
	}
	nameField := &Field{
		ID:    "name",
		Value: "Hulk Hogan",
	}
	record := StoreRecord{}
	record["id"] = idField
	record["name"] = nameField

	t.Run("Access field properties directly", func(t *testing.T) {
		assert.Equal(t, 2, len(record))
		assert.Equal(t, "id", record["id"].ID)
		assert.Equal(t, 12345, record["id"].Value)
		assert.Equal(t, "name", record["name"].ID)
		assert.Equal(t, "Hulk Hogan", record["name"].Value)
	})

}
