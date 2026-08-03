package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"cmd/api/internal/models"

	"github.com/google/uuid"
)

func doJSONRequest(router http.Handler, method, path string, body any, token string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}

	req, _ := http.NewRequest(method, path, &buf)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func clearTables() {
	testDB.Exec("DELETE FROM ticket_assignment_history")
	testDB.Exec("DELETE FROM tickets")
	testDB.Exec("DELETE FROM users")
	testDB.Exec("DELETE FROM categories")
}

func createUser(t *testing.T, role string, email string) models.User {
	t.Helper()

	u := models.User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     email,
		Role:      models.Role(role),
	}
	if err := testDB.Create(&u).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	return u
}

func createTicket(t *testing.T, creatorID uuid.UUID, title string, status models.Status, assignedTo *uuid.UUID) models.Ticket {
	t.Helper()

	ticket := models.Ticket{
		Title:       title,
		Description: "desc",
		Priority:    models.PriorityHigh,
		Status:      status,
		CreatedBy:   &creatorID,
		AssignedTo:  assignedTo,
	}
	if err := testDB.Create(&ticket).Error; err != nil {
		t.Fatalf("create ticket: %v", err)
	}
	return ticket
}

func TestTicketsAPI(t *testing.T) {
	router := SetupRouter()

	t.Run("create ticket ok standard user", func(t *testing.T) {
		clearTables()
		u := createUser(t, string(models.RoleStandard), "std1@test.com")
		token := GenerateTestToken(u.ID.String(), string(models.RoleStandard))

		w := doJSONRequest(router, http.MethodPost, "/api/v1/tickets", map[string]any{
			"title":       "VPN down",
			"description": "Cannot connect",
			"priority":    "high",
		}, token)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected %d got %d body=%s", http.StatusCreated, w.Code, w.Body.String())
		}

		var ticket models.Ticket
		if err := json.Unmarshal(w.Body.Bytes(), &ticket); err != nil {
			t.Fatalf("unmarshal: %v body=%s", err, w.Body.String())
		}
		if ticket.Title != "VPN down" {
			t.Fatalf("expected title VPN down got %s", ticket.Title)
		}
		if ticket.Status != models.StatusOpen {
			t.Fatalf("expected status open got %s", ticket.Status)
		}
		if ticket.AssignedTo != nil {
			t.Fatalf("expected no assignment")
		}
	})

	t.Run("create ticket missing title", func(t *testing.T) {
		clearTables()
		u := createUser(t, string(models.RoleStandard), "std2@test.com")
		token := GenerateTestToken(u.ID.String(), string(models.RoleStandard))

		w := doJSONRequest(router, http.MethodPost, "/api/v1/tickets", map[string]any{
			"description": "No title",
			"priority":    "high",
		}, token)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected %d got %d body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	t.Run("create ticket invalid priority", func(t *testing.T) {
		clearTables()
		u := createUser(t, string(models.RoleStandard), "std3@test.com")
		token := GenerateTestToken(u.ID.String(), string(models.RoleStandard))

		w := doJSONRequest(router, http.MethodPost, "/api/v1/tickets", map[string]any{
			"title":       "Bad priority",
			"description": "x",
			"priority":    "super",
		}, token)

		if w.Code == http.StatusCreated {
			t.Fatalf("expected failure for invalid priority, got created body=%s", w.Body.String())
		}
	})

	t.Run("create ticket without auth", func(t *testing.T) {
		clearTables()
		w := doJSONRequest(router, http.MethodPost, "/api/v1/tickets", map[string]any{
			"title":       "No auth",
			"description": "x",
			"priority":    "high",
		}, "")
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected %d got %d body=%s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	t.Run("create with assignment by standard forbidden", func(t *testing.T) {
		clearTables()
		u := createUser(t, string(models.RoleStandard), "std4@test.com")
		assignee := createUser(t, string(models.RoleTech), "tech1@test.com")
		token := GenerateTestToken(u.ID.String(), string(models.RoleStandard))

		w := doJSONRequest(router, http.MethodPost, "/api/v1/tickets", map[string]any{
			"title":       "Assign forbidden",
			"description": "x",
			"priority":    "high",
			"assigned_to": assignee.ID.String(),
		}, token)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected %d got %d body=%s", http.StatusForbidden, w.Code, w.Body.String())
		}
	})

	t.Run("create with assignment by tech ok and history written", func(t *testing.T) {
		clearTables()
		tech := createUser(t, string(models.RoleTech), "tech2@test.com")
		assignee := createUser(t, string(models.RoleTech), "tech3@test.com")
		token := GenerateTestToken(tech.ID.String(), string(models.RoleTech))

		w := doJSONRequest(router, http.MethodPost, "/api/v1/tickets", map[string]any{
			"title":       "Assigned ticket",
			"description": "x",
			"priority":    "high",
			"assigned_to": assignee.ID.String(),
		}, token)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected %d got %d body=%s", http.StatusCreated, w.Code, w.Body.String())
		}

		var ticket models.Ticket
		if err := json.Unmarshal(w.Body.Bytes(), &ticket); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if ticket.AssignedTo == nil || *ticket.AssignedTo != assignee.ID {
			t.Fatalf("assignment not saved")
		}

		var count int64
		if err := testDB.Model(&models.TicketAssignmentHistory{}).Where("ticket_id = ?", ticket.ID).Count(&count).Error; err != nil {
			t.Fatalf("count history: %v", err)
		}
		if count != 1 {
			t.Fatalf("expected 1 history row got %d", count)
		}
	})

	t.Run("get ticket by id ok", func(t *testing.T) {
		clearTables()
		u := createUser(t, string(models.RoleStandard), "std5@test.com")
		token := GenerateTestToken(u.ID.String(), string(models.RoleStandard))
		ticket := createTicket(t, u.ID, "my ticket", models.StatusOpen, nil)

		path := fmt.Sprintf("/api/v1/tickets/%s", ticket.ID.String())
		w := doJSONRequest(router, http.MethodGet, path, nil, token)

		if w.Code != http.StatusOK {
			t.Fatalf("expected %d got %d body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("get ticket invalid id", func(t *testing.T) {
		clearTables()
		u := createUser(t, string(models.RoleStandard), "std6@test.com")
		token := GenerateTestToken(u.ID.String(), string(models.RoleStandard))

		w := doJSONRequest(router, http.MethodGet, "/api/v1/tickets/not-a-uuid", nil, token)
		if w.Code != http.StatusNotFound && w.Code != http.StatusBadRequest {
			t.Fatalf("expected bad request or not found got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("get ticket not found", func(t *testing.T) {
		clearTables()
		u := createUser(t, string(models.RoleStandard), "std7@test.com")
		token := GenerateTestToken(u.ID.String(), string(models.RoleStandard))

		w := doJSONRequest(router, http.MethodGet, "/api/v1/tickets/"+uuid.New().String(), nil, token)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected %d got %d body=%s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})

	t.Run("standard user cannot update others ticket", func(t *testing.T) {
		clearTables()
		creator := createUser(t, string(models.RoleStandard), "creator@test.com")
		other := createUser(t, string(models.RoleStandard), "other@test.com")
		token := GenerateTestToken(other.ID.String(), string(models.RoleStandard))
		ticket := createTicket(t, creator.ID, "t1", models.StatusOpen, nil)

		path := fmt.Sprintf("/api/v1/tickets/%s", ticket.ID.String())
		w := doJSONRequest(router, http.MethodPatch, path, map[string]any{
			"title": "hacked",
		}, token)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected %d got %d body=%s", http.StatusForbidden, w.Code, w.Body.String())
		}
	})

	t.Run("standard user can update own ticket title only", func(t *testing.T) {
		clearTables()
		creator := createUser(t, string(models.RoleStandard), "creator2@test.com")
		token := GenerateTestToken(creator.ID.String(), string(models.RoleStandard))
		ticket := createTicket(t, creator.ID, "old", models.StatusOpen, nil)

		path := fmt.Sprintf("/api/v1/tickets/%s", ticket.ID.String())
		w := doJSONRequest(router, http.MethodPatch, path, map[string]any{
			"title": "new title",
		}, token)

		if w.Code != http.StatusOK {
			t.Fatalf("expected %d got %d body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("standard user cannot change status", func(t *testing.T) {
		clearTables()
		creator := createUser(t, string(models.RoleStandard), "creator3@test.com")
		token := GenerateTestToken(creator.ID.String(), string(models.RoleStandard))
		ticket := createTicket(t, creator.ID, "old", models.StatusOpen, nil)

		path := fmt.Sprintf("/api/v1/tickets/%s", ticket.ID.String())
		w := doJSONRequest(router, http.MethodPatch, path, map[string]any{
			"status": "in_progress",
		}, token)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected %d got %d body=%s", http.StatusForbidden, w.Code, w.Body.String())
		}
	})

	t.Run("tech can assign and history grows", func(t *testing.T) {
		clearTables()
		tech := createUser(t, string(models.RoleTech), "tech4@test.com")
		a1 := createUser(t, string(models.RoleTech), "a1@test.com")
		a2 := createUser(t, string(models.RoleTech), "a2@test.com")
		token := GenerateTestToken(tech.ID.String(), string(models.RoleTech))
		ticket := createTicket(t, tech.ID, "assign", models.StatusOpen, nil)

		path := fmt.Sprintf("/api/v1/tickets/%s", ticket.ID.String())

		w1 := doJSONRequest(router, http.MethodPatch, path, map[string]any{
			"assigned_to": a1.ID.String(),
		}, token)
		if w1.Code != http.StatusOK {
			t.Fatalf("first assign expected ok got %d body=%s", w1.Code, w1.Body.String())
		}

		w2 := doJSONRequest(router, http.MethodPatch, path, map[string]any{
			"assigned_to": a2.ID.String(),
		}, token)
		if w2.Code != http.StatusOK {
			t.Fatalf("reassign expected ok got %d body=%s", w2.Code, w2.Body.String())
		}

		var count int64
		if err := testDB.Model(&models.TicketAssignmentHistory{}).Where("ticket_id = ?", ticket.ID).Count(&count).Error; err != nil {
			t.Fatalf("count history: %v", err)
		}
		if count != 2 {
			t.Fatalf("expected 2 history rows got %d", count)
		}
	})

	t.Run("patch invalid json", func(t *testing.T) {
		clearTables()
		u := createUser(t, string(models.RoleAdmin), "admin1@test.com")
		token := GenerateTestToken(u.ID.String(), string(models.RoleAdmin))
		ticket := createTicket(t, u.ID, "badjson", models.StatusOpen, nil)

		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/tickets/"+ticket.ID.String(), bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected %d got %d body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	t.Run("delete ticket admin ok", func(t *testing.T) {
		clearTables()
		admin := createUser(t, string(models.RoleAdmin), "admin2@test.com")
		token := GenerateTestToken(admin.ID.String(), string(models.RoleAdmin))
		ticket := createTicket(t, admin.ID, "delete me", models.StatusOpen, nil)

		path := fmt.Sprintf("/api/v1/tickets/%s", ticket.ID.String())
		w := doJSONRequest(router, http.MethodDelete, path, nil, token)
		if w.Code != http.StatusOK {
			t.Fatalf("expected %d got %d body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("delete ticket non admin forbidden", func(t *testing.T) {
		clearTables()
		tech := createUser(t, string(models.RoleTech), "tech5@test.com")
		token := GenerateTestToken(tech.ID.String(), string(models.RoleTech))
		ticket := createTicket(t, tech.ID, "delete no", models.StatusOpen, nil)

		path := fmt.Sprintf("/api/v1/tickets/%s", ticket.ID.String())
		w := doJSONRequest(router, http.MethodDelete, path, nil, token)
		if w.Code != http.StatusForbidden {
			t.Fatalf("expected %d got %d body=%s", http.StatusForbidden, w.Code, w.Body.String())
		}
	})

	t.Run("history endpoint ok", func(t *testing.T) {
		clearTables()
		admin := createUser(t, string(models.RoleAdmin), "admin3@test.com")
		tech := createUser(t, string(models.RoleTech), "tech6@test.com")
		token := GenerateTestToken(admin.ID.String(), string(models.RoleAdmin))
		ticket := createTicket(t, admin.ID, "history", models.StatusOpen, nil)

		path := fmt.Sprintf("/api/v1/tickets/%s", ticket.ID.String())
		w1 := doJSONRequest(router, http.MethodPatch, path, map[string]any{
			"assigned_to": tech.ID.String(),
		}, token)
		if w1.Code != http.StatusOK {
			t.Fatalf("assign failed got %d body=%s", w1.Code, w1.Body.String())
		}

		historyPath := fmt.Sprintf("/api/v1/tickets/%s/history", ticket.ID.String())
		w2 := doJSONRequest(router, http.MethodGet, historyPath, nil, token)
		if w2.Code != http.StatusOK {
			t.Fatalf("history failed got %d body=%s", w2.Code, w2.Body.String())
		}

		var history []models.TicketAssignmentHistory
		if err := json.Unmarshal(w2.Body.Bytes(), &history); err != nil {
			t.Fatalf("unmarshal history: %v body=%s", err, w2.Body.String())
		}
		if len(history) == 0 {
			t.Fatalf("expected history rows")
		}
	})
}
