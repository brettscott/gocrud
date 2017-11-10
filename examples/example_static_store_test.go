package examples

import (
	"github.com/brettscott/gocrud/crud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExampleStaticStore(t *testing.T) {

	t.Run("List", func(t *testing.T) {
		usersEntity := &crud.Entity{
			ID:     "users",
			Label:  "User",
			Labels: "Users",
		}

		store := NewExampleStaticStore()
		results, err := store.List(usersEntity)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(results))
	})

	t.Run("Get", func(t *testing.T) {
		usersEntity := &crud.Entity{
			ID:     "users",
			Label:  "User",
			Labels: "Users",
		}

		store := NewExampleStaticStore()
		superman, err := store.Get(usersEntity, "the-superman-id")
		assert.NoError(t, err)
		name, err := superman.GetValue("name")
		assert.NoError(t, err)
		assert.Equal(t, "Superman", name)
	})
}
