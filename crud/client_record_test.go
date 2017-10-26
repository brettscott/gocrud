package crud

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientRecord(t *testing.T) {

	t.Run("ClientRecords", func(t *testing.T) {

		t.Run("GetClientRecordByKeyValue returns client record", func(t *testing.T) {
			johnSmith := ClientRecord{
				KeyValues: KeyValues{
					{
						Key:   "id",
						Value: "1234",
					},
					{
						Key:   "name",
						Value: "John Smith",
					},
				},
			}

			jillJones := ClientRecord{
				KeyValues: KeyValues{
					{
						Key:   "id",
						Value: "5678",
					},
					{
						Key:   "name",
						Value: "Jill Jones",
					},
				},
			}

			clientRecords := ClientRecords{}
			clientRecords = append(clientRecords, johnSmith)
			clientRecords = append(clientRecords, jillJones)

			assert.Equal(t, 2, len(clientRecords), "Number of records")

			john, err := clientRecords.GetClientRecordByKeyValue("id", "1234")
			assert.NoError(t, err)
			johnName, err := john.KeyValues.GetValue("name")
			assert.NoError(t, err)
			assert.Equal(t, "John Smith", johnName)

			jill, err := clientRecords.GetClientRecordByKeyValue("id", "5678")
			assert.NoError(t, err)
			jillName, err := jill.KeyValues.GetValue("name")
			assert.NoError(t, err)
			assert.Equal(t, "Jill Jones", jillName)
		})
	})

	t.Run("ClientRecord", func(t *testing.T) {

		t.Run("GetKeyValue returns key value", func(t *testing.T) {
			johnSmith := ClientRecord{
				KeyValues: KeyValues{
					{
						Key:   "id",
						Value: "1234",
					},
					{
						Key:   "name",
						Value: "John Smith",
					},
				},
			}

			id, err := johnSmith.GetKeyValue("id")
			assert.NoError(t, err)
			assert.Equal(t, "id", id.Key)
			assert.Equal(t, "1234", id.Value)

			name, err := johnSmith.GetKeyValue("name")
			assert.NoError(t, err)
			assert.Equal(t, "name", name.Key)
			assert.Equal(t, "John Smith", name.Value)
		})

		t.Run("GetValue returns value", func(t *testing.T) {
			johnSmith := ClientRecord{
				KeyValues: KeyValues{
					{
						Key:   "id",
						Value: "1234",
					},
					{
						Key:   "name",
						Value: "John Smith",
					},
				},
			}

			id, err := johnSmith.GetValue("id")
			assert.NoError(t, err)
			assert.Equal(t, "1234", id)

			name, err := johnSmith.GetValue("name")
			assert.NoError(t, err)
			assert.Equal(t, "John Smith", name)
		})
	})

	t.Run("KeyValues", func(t *testing.T) {

		t.Run("GetKeyValue returns key value", func(t *testing.T) {
			kvs := KeyValues{
				{
					Key:   "id",
					Value: "1234",
				},
				{
					Key:   "name",
					Value: "John Smith",
				},
			}

			assert.Equal(t, 2, len(kvs), "Number of kvs")

			id, err := kvs.GetKeyValue("id")
			assert.NoError(t, err)
			assert.Equal(t, "id", id.Key)
			assert.Equal(t, "1234", id.Value)

			name, err := kvs.GetKeyValue("name")
			assert.NoError(t, err)
			assert.Equal(t, "name", name.Key)
			assert.Equal(t, "John Smith", name.Value)
		})

		t.Run("GetValue returns value", func(t *testing.T) {
			kvs := KeyValues{
				{
					Key:   "id",
					Value: "1234",
				},
				{
					Key:   "name",
					Value: "John Smith",
				},
			}

			assert.Equal(t, 2, len(kvs), "Number of kvs")

			id, err := kvs.GetValue("id")
			assert.NoError(t, err)
			assert.Equal(t, "1234", id)

			name, err := kvs.GetValue("name")
			assert.NoError(t, err)
			assert.Equal(t, "John Smith", name)
		})
	})

}
