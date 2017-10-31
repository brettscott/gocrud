package crud

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElementsValidator(t *testing.T) {

	elementsValidator := NewElementsValidator()

	t.Run("Basic", func(t *testing.T) {

		t.Run("Passes when user data has no validation rules", func(t *testing.T) {
			testEntity := &Entity{
				ID: "test",
				Elements: []Element{
					{
						ID:         "id",
						Label:      "Identifier",
						DataType:   ELEMENT_DATA_TYPE_STRING,
						PrimaryKey: true,
					},
					{
						ID:       "name",
						Label:    "Name",
						DataType: ELEMENT_DATA_TYPE_STRING,
					},
				},
			}

			userData := StoreRecord{}
			userData["id"] = &Field{
				ID:       "id",
				Value:    "12345",
				Hydrated: true,
			}
			userData["name"] = &Field{
				ID:       "name",
				Value:    "John Smith",
				Hydrated: true,
			}

			success, clientErrors := elementsValidator.Validate(testEntity, userData, ACTION_POST)

			assert.Equal(t, true, success)
			assert.Nil(t, clientErrors)
		})

		t.Run("Passes when user data is valid", func(t *testing.T) {
			testEntity := &Entity{
				ID: "test",
				Elements: []Element{
					{
						ID:         "id",
						Label:      "Identifier",
						DataType:   ELEMENT_DATA_TYPE_STRING,
						PrimaryKey: true,
					},
					{
						ID:       "name",
						Label:    "Name",
						DataType: ELEMENT_DATA_TYPE_STRING,
						Validation: ElementValidation{
							Required: true,
						},
					},
				},
			}

			userData := StoreRecord{}
			userData["id"] = &Field{
				ID:       "id",
				Value:    "12345",
				Hydrated: true,
			}
			userData["name"] = &Field{
				ID:       "name",
				Value:    "John Smith",
				Hydrated: true,
			}

			success, clientErrors := elementsValidator.Validate(testEntity, userData, ACTION_POST)

			assert.Equal(t, true, success)
			assert.Nil(t, clientErrors)
			//assert.Equal(t, 0, len(clientErrors.ElementsErrors))
			//assert.Equal(t, 0, len(clientErrors.GlobalErrors))
		})

	})

	t.Run("Required", func(t *testing.T) {

		t.Run("Passes when posting user data with a required field", func(t *testing.T) {
			testEntity := &Entity{
				ID: "test",
				Elements: []Element{
					{
						ID:         "id",
						Label:      "Identifier",
						DataType:   ELEMENT_DATA_TYPE_STRING,
						PrimaryKey: true,
					},
					{
						ID:       "name",
						Label:    "Name",
						DataType: ELEMENT_DATA_TYPE_STRING,
						Validation: ElementValidation{
							Required: true,
						},
					},
				},
			}

			userData := StoreRecord{}
			userData["id"] = &Field{
				ID:       "id",
				Value:    "12345",
				Hydrated: true,
			}
			userData["name"] = &Field{
				ID:       "name",
				Value:    "John Smith",
				Hydrated: true,
			}

			success, clientErrors := elementsValidator.Validate(testEntity, userData, ACTION_POST)

			assert.Equal(t, true, success, "Should not be valid")
			assert.Nil(t, clientErrors)
			//assert.Equal(t, 0, len(clientErrors.ElementsErrors), "Element error")
			//assert.Equal(t, 0, len(clientErrors.GlobalErrors), "Global error")
		})

		t.Run("Fails when posting user data without a required field being provided", func(t *testing.T) {
			testEntity := &Entity{
				ID: "test",
				Elements: []Element{
					{
						ID:         "id",
						Label:      "Identifier",
						DataType:   ELEMENT_DATA_TYPE_STRING,
						PrimaryKey: true,
					},
					{
						ID:       "name",
						Label:    "Name",
						DataType: ELEMENT_DATA_TYPE_STRING,
						Validation: ElementValidation{
							Required: true,
						},
					},
				},
			}

			userData := StoreRecord{}
			userData["id"] = &Field{
				ID:       "id",
				Value:    "12345",
				Hydrated: true,
			}
			userData["name"] = &Field{
				ID:       "name",
				Value:    "",
				Hydrated: false,
			}

			success, clientErrors := elementsValidator.Validate(testEntity, userData, ACTION_POST)

			assert.Equal(t, false, success, "Should not be valid")
			assert.NotNil(t, clientErrors)
			assert.Equal(t, 1, len(clientErrors.ElementsErrors), "Element error")
			assert.Equal(t, 0, len(clientErrors.GlobalErrors), "Global error")
		})

		t.Run("Passes when putting and patching user data when a required field is not provided", func(t *testing.T) {
			for _, action := range []string{ACTION_PUT, ACTION_PATCH} {
				testEntity := &Entity{
					ID: "test",
					Elements: []Element{
						{
							ID:         "id",
							Label:      "Identifier",
							DataType:   ELEMENT_DATA_TYPE_STRING,
							PrimaryKey: true,
						},
						{
							ID:       "name",
							Label:    "Name",
							DataType: ELEMENT_DATA_TYPE_STRING,
							Validation: ElementValidation{
								Required: true,
							},
						},
					},
				}

				userData := StoreRecord{}
				userData["id"] = &Field{
					ID:       "id",
					Value:    "12345",
					Hydrated: true,
				}
				userData["name"] = &Field{
					ID:       "name",
					Value:    "",
					Hydrated: false,
				}

				success, clientErrors := elementsValidator.Validate(testEntity, userData, action)

				assert.Equal(t, true, success, "Should not be valid")
				assert.Nil(t, clientErrors)
				//assert.Equal(t, 0, len(clientErrors.ElementsErrors), "Element error")
				//assert.Equal(t, 0, len(clientErrors.GlobalErrors), "Global error")
			}
		})

		t.Run("Fails when posting, putting and patching user data when a required field is empty", func(t *testing.T) {
			for _, action := range []string{ACTION_POST, ACTION_PUT, ACTION_PATCH} {
				testEntity := &Entity{
					ID: "test",
					Elements: []Element{
						{
							ID:         "id",
							Label:      "Identifier",
							DataType:   ELEMENT_DATA_TYPE_STRING,
							PrimaryKey: true,
						},
						{
							ID:       "name",
							Label:    "Name",
							DataType: ELEMENT_DATA_TYPE_STRING,
							Validation: ElementValidation{
								Required: true,
							},
						},
					},
				}

				userData := StoreRecord{}
				userData["id"] = &Field{
					ID:       "id",
					Value:    "12345",
					Hydrated: true,
				}
				userData["name"] = &Field{
					ID:       "name",
					Value:    "",
					Hydrated: true,
				}

				success, clientErrors := elementsValidator.Validate(testEntity, userData, action)

				assert.Equal(t, false, success, fmt.Sprintf("Should not be valid on %s", action))
				assert.NotNil(t, clientErrors)
				assert.Equal(t, 1, len(clientErrors.ElementsErrors), fmt.Sprintf("Element error on %s", action))
				assert.Equal(t, 0, len(clientErrors.GlobalErrors), fmt.Sprintf("Global error on %s", action))
			}
		})
	})

	t.Run("MustProvide", func(t *testing.T) {

		t.Run("Passes when posting, putting and patching user data when a 'must be provided' field is provided", func(t *testing.T) {
			for _, action := range []string{ACTION_POST, ACTION_PUT, ACTION_PATCH} {
				testEntity := &Entity{
					ID: "test",
					Elements: []Element{
						{
							ID:         "id",
							Label:      "Identifier",
							DataType:   ELEMENT_DATA_TYPE_STRING,
							PrimaryKey: true,
						},
						{
							ID:       "name",
							Label:    "Name",
							DataType: ELEMENT_DATA_TYPE_STRING,
							Validation: ElementValidation{
								MustProvide: true,
							},
						},
					},
				}

				userData := StoreRecord{}
				userData["id"] = &Field{
					ID:       "id",
					Value:    "12345",
					Hydrated: true,
				}
				userData["name"] = &Field{
					ID:       "name",
					Value:    "John Smith",
					Hydrated: true,
				}

				success, clientErrors := elementsValidator.Validate(testEntity, userData, action)

				assert.Equal(t, true, success, fmt.Sprintf("Should not be valid on %s", action))
				assert.Nil(t, clientErrors)
				//assert.Equal(t, 0, len(clientErrors.ElementsErrors), fmt.Sprintf("Element error on %s", action))
				//assert.Equal(t, 0, len(clientErrors.GlobalErrors), fmt.Sprintf("Global error on %s", action))
			}
		})

		t.Run("Fails when posting, putting and patching user data when a 'must be provided' field is missing", func(t *testing.T) {
			for _, action := range []string{ACTION_POST, ACTION_PUT, ACTION_PATCH} {
				testEntity := &Entity{
					ID: "test",
					Elements: []Element{
						{
							ID:         "id",
							Label:      "Identifier",
							DataType:   ELEMENT_DATA_TYPE_STRING,
							PrimaryKey: true,
						},
						{
							ID:       "name",
							Label:    "Name",
							DataType: ELEMENT_DATA_TYPE_STRING,
							Validation: ElementValidation{
								MustProvide: true,
							},
						},
					},
				}

				userData := StoreRecord{}
				userData["id"] = &Field{
					ID:       "id",
					Value:    "12345",
					Hydrated: true,
				}
				userData["name"] = &Field{
					ID:       "name",
					Value:    "",
					Hydrated: false,
				}

				success, clientErrors := elementsValidator.Validate(testEntity, userData, action)

				assert.Equal(t, false, success, fmt.Sprintf("Should not be valid on %s", action))
				assert.NotNil(t, clientErrors)
				assert.Equal(t, 1, len(clientErrors.ElementsErrors), fmt.Sprintf("Element error on %s", action))
				assert.Equal(t, 0, len(clientErrors.GlobalErrors), fmt.Sprintf("Global error on %s", action))
			}
		})

		t.Run("Fails when posting user data when a 'must be provided on posting' field is missing", func(t *testing.T) {
			testEntity := &Entity{
				ID: "test",
				Elements: []Element{
					{
						ID:         "id",
						Label:      "Identifier",
						DataType:   ELEMENT_DATA_TYPE_STRING,
						PrimaryKey: true,
					},
					{
						ID:       "name",
						Label:    "Name",
						DataType: ELEMENT_DATA_TYPE_STRING,
						Validation: ElementValidation{
							MustProvideOnPost: true,
						},
					},
				},
			}

			userData := StoreRecord{}
			userData["id"] = &Field{
				ID:       "id",
				Value:    "12345",
				Hydrated: true,
			}
			userData["name"] = &Field{
				ID:       "name",
				Value:    "",
				Hydrated: false,
			}

			success, clientErrors := elementsValidator.Validate(testEntity, userData, ACTION_POST)

			assert.Equal(t, false, success, "Should not be valid")
			assert.NotNil(t, clientErrors)
			assert.Equal(t, 1, len(clientErrors.ElementsErrors), "Element error")
			assert.Equal(t, 0, len(clientErrors.GlobalErrors), "Global error")
		})

		t.Run("Fails when posting user data when a 'must be provided on putting' field is missing", func(t *testing.T) {
			testEntity := &Entity{
				ID: "test",
				Elements: []Element{
					{
						ID:         "id",
						Label:      "Identifier",
						DataType:   ELEMENT_DATA_TYPE_STRING,
						PrimaryKey: true,
					},
					{
						ID:       "name",
						Label:    "Name",
						DataType: ELEMENT_DATA_TYPE_STRING,
						Validation: ElementValidation{
							MustProvideOnPut: true,
						},
					},
				},
			}

			userData := StoreRecord{}
			userData["id"] = &Field{
				ID:       "id",
				Value:    "12345",
				Hydrated: true,
			}
			userData["name"] = &Field{
				ID:       "name",
				Value:    "",
				Hydrated: false,
			}

			success, clientErrors := elementsValidator.Validate(testEntity, userData, ACTION_PUT)

			assert.Equal(t, false, success, "Should not be valid")
			assert.NotNil(t, clientErrors)
			assert.Equal(t, 1, len(clientErrors.ElementsErrors), "Element error")
			assert.Equal(t, 0, len(clientErrors.GlobalErrors), "Global error")
		})

		t.Run("Fails when posting user data when a 'must be provided on patching' field is missing", func(t *testing.T) {
			testEntity := &Entity{
				ID: "test",
				Elements: []Element{
					{
						ID:         "id",
						Label:      "Identifier",
						DataType:   ELEMENT_DATA_TYPE_STRING,
						PrimaryKey: true,
					},
					{
						ID:       "name",
						Label:    "Name",
						DataType: ELEMENT_DATA_TYPE_STRING,
						Validation: ElementValidation{
							MustProvideOnPatch: true,
						},
					},
				},
			}

			userData := StoreRecord{}
			userData["id"] = &Field{
				ID:       "id",
				Value:    "12345",
				Hydrated: true,
			}
			userData["name"] = &Field{
				ID:       "name",
				Value:    "",
				Hydrated: false,
			}

			success, clientErrors := elementsValidator.Validate(testEntity, userData, ACTION_PATCH)

			assert.Equal(t, false, success, "Should not be valid")
			assert.NotNil(t, clientErrors)
			assert.Equal(t, 1, len(clientErrors.ElementsErrors), "Element error")
			assert.Equal(t, 0, len(clientErrors.GlobalErrors), "Global error")
		})
	})
}
