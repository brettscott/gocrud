package crud

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStoreMode(t *testing.T) {
	t.Run("IsReadable", func(t *testing.T) {
		storeMode := StoreMode{Read: true}
		assert.Equal(t, true, storeMode.IsReadable())

		storeMode = StoreMode{Read: false}
		assert.Equal(t, false, storeMode.IsReadable())
	})

	t.Run("IsWritable", func(t *testing.T) {
		storeMode := StoreMode{Write: true}
		assert.Equal(t, true, storeMode.IsWritable())

		storeMode = StoreMode{Write: false}
		assert.Equal(t, false, storeMode.IsWritable())
	})

	t.Run("IsDeletable", func(t *testing.T) {
		storeMode := StoreMode{Delete: true}
		assert.Equal(t, true, storeMode.IsDeletable())

		storeMode = StoreMode{Delete: false}
		assert.Equal(t, false, storeMode.IsDeletable())
	})

}
