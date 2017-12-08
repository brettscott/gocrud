package crud

import (
	"github.com/pressly/chi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNewUiRoute(t *testing.T) {

	testLog := NewTestLog(t)
	testStatsD := NewTestStatsD(t)

	idElement := Element{
		ID:         "id",
		Label:      "Identifier",
		DataType:   ELEMENT_DATA_TYPE_STRING,
		PrimaryKey: true,
	}
	nameElement := Element{
		ID:       "name",
		Label:    "Name",
		DataType: ELEMENT_DATA_TYPE_STRING,
	}

	userEntity := &Entity{
		ID: "users",
		Elements: []Element{
			idElement,
			nameElement,
		},
	}
	entities := Entities{
		"users": userEntity,
	}

	t.Run("list", func(t *testing.T) {

		t.Run("returns 400 when unknown entity provided", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/this-is-unknown", nil)
			router := makeTestUiRouter(t, entities, &fakeApiServicer{}, &fakeTemplateServicer{})
			router.ServeHTTP(w, req)
			if http.StatusBadRequest != w.Code {
				t.Error("bad status: expected", http.StatusBadRequest, "got", w.Code, "body:", w.Body.String())
			}
		})

		t.Run("loads \"list\" template and sends correct context to template", func(t *testing.T) {
			fakeApiService := NewFakeApiServicer()
			fakeApiService.listResponseBody = ClientRecords{
				{
					KeyValues: KeyValues{
						{
							Key:   "id",
							Value: "12345",
						},
						{
							Key:   "name",
							Value: "Bruce Lee",
						},
					},
				},
				{
					KeyValues: KeyValues{
						{
							Key:   "id",
							Value: "67890",
						},
						{
							Key:   "name",
							Value: "Jackie Chan",
						},
					},
				},
			}

			fakeTemplateService := &fakeTemplateServicer{}
			uiRouteHandler := NewUiRoute(entities, fakeApiService, fakeTemplateService, testLog, testStatsD)
			router := chi.NewRouter()
			router.Route("/", uiRouteHandler)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/users", nil)
			router.ServeHTTP(w, req)
			if http.StatusOK != w.Code {
				t.Error("bad status: expected", http.StatusOK, "got", w.Code, "body:", w.Body.String())
			}
			expectedContext := map[string]interface{}{
				"entity": userEntity,
				"rows": []row{
					{
						"id":   "12345",
						"name": "Bruce Lee",
					},
					{
						"id":   "67890",
						"name": "Jackie Chan",
					},
				},
			}
			assert.True(t, fakeApiService.listCalled, "Should have requested list()")
			assert.Equal(t, "list", fakeTemplateService.execTmplName, "Incorrect template name")
			assert.Equal(t, expectedContext, fakeTemplateService.execContext, "Incorrect context")
		})
	})

	t.Run("create", func(t *testing.T) {

		t.Run("returns 400 when unknown entity provided", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/this-is-unknown/create", nil)
			router := makeTestUiRouter(t, entities, &fakeApiServicer{}, &fakeTemplateServicer{})
			router.ServeHTTP(w, req)
			if http.StatusBadRequest != w.Code {
				t.Error("bad status: expected", http.StatusBadRequest, "got", w.Code, "body:", w.Body.String())
			}
		})

		t.Run("loads \"form\" template and sends correct context to template", func(t *testing.T) {
			fakeApiService := NewFakeApiServicer()
			fakeTemplateService := &fakeTemplateServicer{}
			uiRouteHandler := NewUiRoute(entities, fakeApiService, fakeTemplateService, testLog, testStatsD)
			router := chi.NewRouter()
			router.Route("/", uiRouteHandler)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/users/create", nil)
			router.ServeHTTP(w, req)
			if http.StatusOK != w.Code {
				t.Error("bad status: expected", http.StatusOK, "got", w.Code, "body:", w.Body.String())
			}
			expectedContext := map[string]interface{}{
				"create":   true,
				"entity":   userEntity,
				"recordID": "",
				"elementValues": []ElementValue{
					{
						Element: idElement,
						Value:   nil,
					},
					{
						Element: nameElement,
						Value:   nil,
					},
				},
			}
			assert.Equal(t, "form", fakeTemplateService.execTmplName, "Incorrect template name")
			assert.Equal(t, expectedContext, fakeTemplateService.execContext, "Incorrect context")
			assert.False(t, fakeApiService.getCalled, "Should not have requested get()")
		})
	})

	t.Run("edit", func(t *testing.T) {

		t.Run("returns 400 when unknown entity provided", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/this-is-unknown/12345/edit", nil)
			router := makeTestUiRouter(t, entities, &fakeApiServicer{}, &fakeTemplateServicer{})
			router.ServeHTTP(w, req)
			if http.StatusBadRequest != w.Code {
				t.Error("bad status: expected", http.StatusBadRequest, "got", w.Code, "body:", w.Body.String())
			}
		})

		t.Run("loads \"form\" template and sends correct context to template", func(t *testing.T) {
			fakeApiService := NewFakeApiServicer()
			fakeApiService.getResponseBody = ClientRecord{
				KeyValues: KeyValues{
					{
						Key:   "id",
						Value: "12345",
					},
					{
						Key:   "name",
						Value: "Bruce Lee",
					},
				},
			}
			fakeTemplateService := &fakeTemplateServicer{}
			uiRouteHandler := NewUiRoute(entities, fakeApiService, fakeTemplateService, testLog, testStatsD)
			router := chi.NewRouter()
			router.Route("/", uiRouteHandler)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/users/12345/edit", nil)
			router.ServeHTTP(w, req)
			if http.StatusOK != w.Code {
				t.Error("bad status: expected", http.StatusOK, "got", w.Code, "body:", w.Body.String())
			}
			expectedContext := map[string]interface{}{
				"create":   false,
				"entity":   userEntity,
				"recordID": "12345",
				"elementValues": []ElementValue{
					{
						Element: idElement,
						Value:   "12345",
					},
					{
						Element: nameElement,
						Value:   "Bruce Lee",
					},
				},
			}
			assert.True(t, fakeApiService.getCalled, "Should have requested get()")
			assert.Equal(t, "form", fakeTemplateService.execTmplName, "Incorrect template name")
			assert.Equal(t, expectedContext, fakeTemplateService.execContext, "Incorrect context")
		})
	})

	t.Run("save", func(t *testing.T) {

		t.Run("returns 400 when unknown entity provided", func(t *testing.T) {
			fakeApiService := NewFakeApiServicer()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/this-is-unknown/12345/save", nil)
			router := makeTestUiRouter(t, entities, fakeApiService, &fakeTemplateServicer{})
			router.ServeHTTP(w, req)
			if http.StatusBadRequest != w.Code {
				t.Error("bad status: expected", http.StatusBadRequest, "got", w.Code, "body:", w.Body.String())
			}
		})

		t.Run("returns 400 when unknown entity provided without recordID", func(t *testing.T) {
			fakeApiService := NewFakeApiServicer()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/this-is-unknown/save", nil)
			router := makeTestUiRouter(t, entities, fakeApiService, &fakeTemplateServicer{})
			router.ServeHTTP(w, req)
			if http.StatusBadRequest != w.Code {
				t.Error("bad status: expected", http.StatusBadRequest, "got", w.Code, "body:", w.Body.String())
			}
		})

		t.Run("saves form data", func(t *testing.T) {
			fakeApiService := NewFakeApiServicer()
			fakeTemplateService := &fakeTemplateServicer{}
			uiRouteHandler := NewUiRoute(entities, fakeApiService, fakeTemplateService, testLog, testStatsD)
			router := chi.NewRouter()
			router.Route("/", uiRouteHandler)
			w := httptest.NewRecorder()
			form := url.Values{}
			form.Add("id", "12345")
			form.Add("name", "Bruce Lee")
			req, _ := http.NewRequest(http.MethodPost, "/users/12345/save", strings.NewReader(form.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			router.ServeHTTP(w, req)
			if http.StatusOK != w.Code {
				t.Error("bad status: expected", http.StatusOK, "got", w.Code, "body:", w.Body.String())
			}
			//expectedContext := map[string]interface{}{
			//	"create":   false,
			//	"entity":   userEntity,
			//	"recordID": "12345",
			//	"elementValues": []ElementValue{
			//		{
			//			Element: idElement,
			//			Value: "12345",
			//		},
			//		{
			//			Element: nameElement,
			//			Value: "Bruce Lee",
			//		},
			//	},
			//}

			expectedSavedClientRecord := &ClientRecord{
				KeyValues: KeyValues{
					{
						Key:   "id",
						Value: "12345",
					},
					{
						Key:   "name",
						Value: "Bruce Lee",
					},
				},
			}
			assert.True(t, fakeApiService.saveCalled, "Should have requested save()")
			assert.Equal(t, expectedSavedClientRecord, fakeApiService.saveClientRecord, "Should have been sent correct clientRecord")
			assert.Equal(t, ACTION_PUT, fakeApiService.saveAction, "Should have been PUT")
			assert.Equal(t, "12345", fakeApiService.saveRecordID, "Should have correct recordID")
			//assert.Equal(t, "form", fakeTemplateService.execTmplName, "Incorrect template name")
			//assert.Equal(t, expectedContext, fakeTemplateService.execContext, "Incorrect context")
		})
	})
}

func makeTestUiRouter(t *testing.T, entities Entities, apiService apiServicer, templateService templateServicer) http.Handler {
	testLog := NewTestLog(t)
	testStatsD := NewTestStatsD(t)

	uiRouteHandler := NewUiRoute(entities, apiService, templateService, testLog, testStatsD)
	router := chi.NewRouter()
	router.Route("/", uiRouteHandler)

	return router
}
