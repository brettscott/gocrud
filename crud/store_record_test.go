package crud

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecord(t *testing.T) {

	idField := Field{
		ID:    "id",
		Value: 12345,
	}
	nameField := Field{
		ID:    "name",
		Value: "Hulk Hogan",
	}
	record := StoreRecord{}
	record = append(record, idField)
	record = append(record, nameField)

	t.Run("Access field properties directly", func(t *testing.T) {
		assert.Equal(t, 2, len(record))
		assert.Equal(t, "id", record[0].ID)
		assert.Equal(t, 12345, record[0].Value)
		assert.Equal(t, "name", record[1].ID)
		assert.Equal(t, "Hulk Hogan", record[1].Value)
	})

	t.Run("GetField returns a specific field", func(t *testing.T) {
		id, _ := record.GetField("id")
		name, _ := record.GetField("name")

		assert.Equal(t, "id", id.ID)
		assert.Equal(t, 12345, id.Value)
		assert.Equal(t, "name", name.ID)
		assert.Equal(t, "Hulk Hogan", name.Value)
	})

	t.Run("GetField returns error when field not found", func(t *testing.T) {
		_, err := record.GetField("i-do-not-exist")
		if err == nil {
			t.Error("Should have failed")
		}
	})

}
