package handlers

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mikko-kohtala/go-api/internal/models"
)

func TestHealthEndpoint(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	h := New(logger)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	h.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got %v", response["status"])
	}
}

func TestCreateUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	h := New(logger)

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "valid user",
			payload: models.CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing email",
			payload: models.CreateUserRequest{
				Name: "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation failed",
		},
		{
			name: "missing name",
			payload: models.CreateUserRequest{
				Email: "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation failed",
		},
		{
			name: "invalid email",
			payload: models.CreateUserRequest{
				Email: "invalid-email",
				Name:  "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.CreateUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if tt.expectedError != "" {
				if response["error"] != tt.expectedError {
					t.Errorf("expected error '%s', got %v", tt.expectedError, response["error"])
				}
			} else {
				if response["id"] == "" {
					t.Error("expected user ID in response")
				}
			}
		})
	}
}

func TestEmailUniqueness(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	h := New(logger)

	// Create first user
	user1 := models.CreateUserRequest{
		Email: "unique@example.com",
		Name:  "User One",
	}
	body1, _ := json.Marshal(user1)
	req1 := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	h.CreateUser(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Fatalf("failed to create first user: %d", w1.Code)
	}

	// Try to create second user with same email
	user2 := models.CreateUserRequest{
		Email: "unique@example.com",
		Name:  "User Two",
	}
	body2, _ := json.Marshal(user2)
	req2 := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	h.CreateUser(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("expected status %d for duplicate email, got %d", http.StatusConflict, w2.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w2.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["error"] != "email already exists" {
		t.Errorf("expected error 'email already exists', got %v", response["error"])
	}
}

func TestGetUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	h := New(logger)

	// Create a user first
	user := models.CreateUserRequest{
		Email: "get@example.com",
		Name:  "Get User",
	}
	body, _ := json.Marshal(user)
	createReq := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	h.CreateUser(createW, createReq)

	var createResp map[string]interface{}
	json.NewDecoder(createW.Body).Decode(&createResp)
	userID := createResp["id"].(string)

	// Test getting the user
	getReq := httptest.NewRequest("GET", "/api/v1/users/"+userID, nil)
	getReq.SetPathValue("id", userID)
	getW := httptest.NewRecorder()
	h.GetUser(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, getW.Code)
	}

	var getResp map[string]interface{}
	if err := json.NewDecoder(getW.Body).Decode(&getResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if getResp["email"] != "get@example.com" {
		t.Errorf("expected email 'get@example.com', got %v", getResp["email"])
	}

	// Test getting non-existent user
	notFoundReq := httptest.NewRequest("GET", "/api/v1/users/00000000-0000-0000-0000-000000000000", nil)
	notFoundReq.SetPathValue("id", "00000000-0000-0000-0000-000000000000")
	notFoundW := httptest.NewRecorder()
	h.GetUser(notFoundW, notFoundReq)

	if notFoundW.Code != http.StatusNotFound {
		t.Errorf("expected status %d for non-existent user, got %d", http.StatusNotFound, notFoundW.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	h := New(logger)

	// Create a user first
	user := models.CreateUserRequest{
		Email: "delete@example.com",
		Name:  "Delete User",
	}
	body, _ := json.Marshal(user)
	createReq := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	h.CreateUser(createW, createReq)

	var createResp map[string]interface{}
	json.NewDecoder(createW.Body).Decode(&createResp)
	userID := createResp["id"].(string)

	// Delete the user
	deleteReq := httptest.NewRequest("DELETE", "/api/v1/users/"+userID, nil)
	deleteReq.SetPathValue("id", userID)
	deleteW := httptest.NewRecorder()
	h.DeleteUser(deleteW, deleteReq)

	if deleteW.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, deleteW.Code)
	}

	// Verify user is deleted
	getReq := httptest.NewRequest("GET", "/api/v1/users/"+userID, nil)
	getReq.SetPathValue("id", userID)
	getW := httptest.NewRecorder()
	h.GetUser(getW, getReq)

	if getW.Code != http.StatusNotFound {
		t.Errorf("expected status %d for deleted user, got %d", http.StatusNotFound, getW.Code)
	}
}