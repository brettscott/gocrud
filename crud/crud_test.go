package crud

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCrud(t *testing.T) {

	testLog := NewTestLog(t)
	testStatsD := NewTestStatsD(t)

	t.Run("AddStore", func(t *testing.T) {

		t.Run("appends store to slice", func(t *testing.T) {
			crud := NewCrud(&Config{}, testLog, testStatsD)
			fakeStore1 := NewFakeStorer()
			crud.AddStore(fakeStore1)
			fakeStore2 := NewFakeStorer()
			crud.AddStore(fakeStore2)
			stores := crud.GetStores()
			assert.Equal(t, 2, len(stores))
		})
	})

	t.Run("AddEntity", func(t *testing.T) {
		usersEntity := &Entity{
			ID:     "users",
			Label:  "User",
			Labels: "Users",
		}
		computersEntity := &Entity{
			ID:     "computers",
			Label:  "Computer",
			Labels: "Computers",
		}

		t.Run("appends entity to slice", func(t *testing.T) {
			crud := NewCrud(&Config{}, testLog, testStatsD)
			crud.AddEntity(usersEntity)
			crud.AddEntity(computersEntity)
			entities := crud.GetEntities()
			assert.Equal(t, 2, len(entities))
		})
	})


	t.Run("AddElementsValidator", func(t *testing.T) {
		elementsValidator1 := NewFakeElementsValidatorer()
		elementsValidator2 := NewFakeElementsValidatorer()

		t.Run("appends elements validators to slice", func(t *testing.T) {
			crud := NewCrud(&Config{}, testLog, testStatsD)
			crud.AddElementsValidator(elementsValidator1)
			crud.AddElementsValidator(elementsValidator2)
			elementsValidators := crud.GetElementsValidators()
			assert.Equal(t, 2, len(elementsValidators))
		})
	})

	t.Run("AddEntityElementsValidator", func(t *testing.T) {
		elementsValidator1 := NewFakeElementsValidatorer()
		elementsValidator2 := NewFakeElementsValidatorer()
		usersEntity := &Entity{
			ID:     "users",
			Label:  "User",
			Labels: "Users",
		}

		t.Run("appends elements validator to an entity's elements validators", func(t *testing.T) {
			crud := NewCrud(&Config{}, testLog, testStatsD)
			crud.AddEntity(usersEntity)
			crud.AddEntityElementsValidator("users", elementsValidator1)
			crud.AddEntityElementsValidator("users", elementsValidator2)
			elementsValidators, err := crud.GetEntityElementsValidators("users")
			assert.NoError(t, err)
			assert.Equal(t, 2, len(elementsValidators))
		})

		t.Run("panics when entity not registered", func(t *testing.T) {

			assert.Panics(t, func() {
				crud := NewCrud(&Config{}, testLog, testStatsD)
				crud.AddEntityElementsValidator("users", elementsValidator1)
			})
		})
	})

	t.Run("GetEntityElementsValidators", func(t *testing.T) {

		t.Run("throws an error when entity not registered", func(t *testing.T) {
			crud := NewCrud(&Config{}, testLog, testStatsD)
			_, err := crud.GetEntityElementsValidators("users")
			assert.Error(t, err)
		})
	})

	t.Run("AddMutator", func(t *testing.T) {
		mutator1 := newFakeMutatorer()
		mutator2 := newFakeMutatorer()

		t.Run("appends mutator to slice", func(t *testing.T) {
			crud := NewCrud(&Config{}, testLog, testStatsD)
			crud.AddMutator(mutator1)
			crud.AddMutator(mutator2)
			mutators := crud.GetMutators()
			assert.Equal(t, 2, len(mutators))
		})
	})

	t.Run("AddEntityMutator", func(t *testing.T) {
		mutator1 := newFakeMutatorer()
		mutator2 := newFakeMutatorer()
		usersEntity := &Entity{
			ID:     "users",
			Label:  "User",
			Labels: "Users",
		}

		t.Run("appends mutator to entity's slice", func(t *testing.T) {
			crud := NewCrud(&Config{}, testLog, testStatsD)
			crud.AddEntity(usersEntity)
			crud.AddEntityMutator("users", mutator1)
			crud.AddEntityMutator("users", mutator2)
			mutators, err := crud.GetEntityMutators("users")
			assert.NoError(t, err)
			assert.Equal(t, 2, len(mutators))
		})
	})
}
