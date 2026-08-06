package tests

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"cmd/api/internal/models"
)

func TestCategoriesApi(t *testing.T) {
	router := SetupRouter()

	t.Run("Create category", func(t *testing.T) {
		clearTables()
		u := createUser(t, string(models.RoleAdmin), "adminCat@test.com")
		token := GenerateTestToken(u.ID.String(), string(models.RoleAdmin))

		w := doJSONRequest(router, http.MethodPost, "/api/v1/categories", map[string]any{
			"name": "Hardware",
		}, token)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected %d got %d body=%s", http.StatusCreated, w.Code, w.Body.String())
		}

		var category models.Category

		if err := json.Unmarshal(w.Body.Bytes(), &category); err != nil {
			t.Fatalf("unmarshal: %v body=%s", err, w.Body.String())
		}
		if category.Name != "Hardware" {
			t.Fatalf("expected name Hardware got %s", category.Name)
		}
	})

	t.Run("Create category with name missing", func(t *testing.T) {
		clearTables()
		u := createUser(t, string(models.RoleAdmin), "adminCat@test.com")
		token := GenerateTestToken(u.ID.String(), string(models.RoleAdmin))

		w := doJSONRequest(router, http.MethodPost, "/api/v1/categories", map[string]any{}, token)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected %d got %d body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	t.Run("Create category while being a standard user", func(t *testing.T) {
		clearTables()
		u := createUser(t, string(models.RoleStandard), "adminCat@test.com")
		token := GenerateTestToken(u.ID.String(), string(models.RoleStandard))

		w := doJSONRequest(router, http.MethodPost, "/api/v1/categories", map[string]any{
			"name": "IT",
		}, token)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected %d got %d body=%s", http.StatusForbidden, w.Code, w.Body.String())
		}
	})

	t.Run("Get category by ID", func(t *testing.T) {
		clearTables()
		u := createUser(t, string(models.RoleAdmin), "adminGet@test.com")
		token := GenerateTestToken(u.ID.String(), string(models.RoleAdmin))

		w := doJSONRequest(router, http.MethodPost, "/api/v1/categories", map[string]any{
			"name": "IT",
		}, token)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected %d got %d body=%s", http.StatusCreated, w.Code, w.Body.String())
		}

		var category models.Category
		if err := json.Unmarshal(w.Body.Bytes(), &category); err != nil {
			t.Fatalf("unmarshal: %v body=%s", err, w.Body.String())
		}

		path := "/api/v1/categories/" + strconv.FormatUint(uint64(category.ID), 10)
		w = doJSONRequest(router, http.MethodGet, path, nil, token)

		if w.Code != http.StatusOK {
			t.Fatalf("expected %d got %d body=%s", http.StatusOK, w.Code, w.Body.String())
		}

		var fetchedCategory models.Category
		if err := json.Unmarshal(w.Body.Bytes(), &fetchedCategory); err != nil {
			t.Fatalf("unmarshal: %v body=%s", err, w.Body.String())
		}

		if fetchedCategory.ID != category.ID || fetchedCategory.Name != category.Name {
			t.Fatalf("expected category %+v got %+v", category, fetchedCategory)
		}
	})
}
