package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cmd/api/internal/models"
)

func TestUserRoutes_Functional(t *testing.T) {
	router := SetupRouter()

	// Nettoyage de la BDD avant les tests
	testDB.Exec("DELETE FROM users")

	// ------------------------------------------------------------------
	// 1. TEST : Créer un utilisateur via POST /api/v1/users (Admin uniquement)
	// ------------------------------------------------------------------
	t.Run("POST /users - Should create user when caller is Admin", func(t *testing.T) {
		adminToken := GenerateTestToken("admin-uuid-123", "admin")

		body := map[string]interface{}{
			"firstname": "John",
			"lastname":  "Doe",
			"email":     "john.doe@example.com",
			"password":  "password123",
			"role":      "standard",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d. Body: %s", resp.Code, resp.Body.String())
		}
	})

	// ------------------------------------------------------------------
	// 2. TEST : Interdire la création si rôle Standard (403 Forbidden)
	// ------------------------------------------------------------------
	t.Run("POST /users - Should return 403 Forbidden when caller is Standard user", func(t *testing.T) {
		userToken := GenerateTestToken("user-uuid-456", "standard")

		body := map[string]interface{}{
			"firstname": "Jane",
			"lastname":  "Doe",
			"email":     "jane.doe@example.com",
			"password":  "password123",
			"role":      "standard",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.Code)
		}
	})

	// ------------------------------------------------------------------
	// 3. TEST : Lister les utilisateurs GET /api/v1/users
	// ------------------------------------------------------------------
	t.Run("GET /users - Should return list of users", func(t *testing.T) {
		token := GenerateTestToken("user-uuid-456", "standard")

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/?role=standard", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.Code)
		}

		var users []models.User
		_ = json.Unmarshal(resp.Body.Bytes(), &users)

		if len(users) == 0 {
			t.Error("Expected at least 1 user in response")
		}
	})

	// ------------------------------------------------------------------
	// 4. TEST : Modifier son propre profil via PATCH /api/v1/users/:id
	// ------------------------------------------------------------------
	t.Run("PATCH /users/:id - User can update their own profile", func(t *testing.T) {
		// Récupérer un utilisateur existant en BDD
		var user models.User
		testDB.First(&user)

		token := GenerateTestToken(user.ID.String(), string(user.Role))

		updateBody := map[string]interface{}{
			"firstname": "JohnUpdated",
		}
		jsonBody, _ := json.Marshal(updateBody)

		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/users/"+user.ID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
		}
	})

	// ------------------------------------------------------------------
	// 5. TEST : Refus d'accès sans token (401 Unauthorized)
	// ------------------------------------------------------------------
	t.Run("GET /users - Should return 401 when no token is provided", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/", nil)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.Code)
		}
	})
}
