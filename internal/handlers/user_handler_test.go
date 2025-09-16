package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mikko-kohtala/go-api/internal/services"
)

// Mock user service for testing
type mockUserService struct {
	getUserByIDFunc func(ctx context.Context, id string) (*services.User, error)
	getAllUsersFunc func(ctx context.Context) ([]services.User, error)
	createUserFunc  func(ctx context.Context, email, name string) (*services.User, error)
	updateUserFunc  func(ctx context.Context, id string, updates map[string]interface{}) (*services.User, error)
	deleteUserFunc  func(ctx context.Context, id string) error
}

func (m *mockUserService) GetUserByID(ctx context.Context, id string) (*services.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return nil, services.ErrUserNotFound
}

func (m *mockUserService) GetAllUsers(ctx context.Context) ([]services.User, error) {
	if m.getAllUsersFunc != nil {
		return m.getAllUsersFunc(ctx)
	}
	return []services.User{}, nil
}

func (m *mockUserService) CreateUser(ctx context.Context, email, name string) (*services.User, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, email, name)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) (*services.User, error) {
	if m.updateUserFunc != nil {
		return m.updateUserFunc(ctx, id, updates)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) DeleteUser(ctx context.Context, id string) error {
	if m.deleteUserFunc != nil {
		return m.deleteUserFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(testDiscard{}, nil))
}

type testDiscard struct{}

func (testDiscard) Write(p []byte) (int, error) { return len(p), nil }

func TestGetAllUsers(t *testing.T) {
	tests := []struct {
		name           string
		mockService    *mockUserService
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success with users",
			mockService: &mockUserService{
				getAllUsersFunc: func(ctx context.Context) ([]services.User, error) {
					return []services.User{
						{ID: "1", Email: "test1@example.com", Name: "Test User 1"},
						{ID: "2", Email: "test2@example.com", Name: "Test User 2"},
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "empty users list",
			mockService: &mockUserService{
				getAllUsersFunc: func(ctx context.Context) ([]services.User, error) {
					return []services.User{}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "service error",
			mockService: &mockUserService{
				getAllUsersFunc: func(ctx context.Context) ([]services.User, error) {
					return nil, errors.New("database error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewUserHandler(tt.mockService, testLogger())

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			rec := httptest.NewRecorder()

			handler.GetAllUsers(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if _, ok := response["users"]; !ok {
					t.Error("Response should contain 'users' field")
				}
				if _, ok := response["count"]; !ok {
					t.Error("Response should contain 'count' field")
				}
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockService    *mockUserService
		expectedStatus int
	}{
		{
			name:   "success",
			userID: "123",
			mockService: &mockUserService{
				getUserByIDFunc: func(ctx context.Context, id string) (*services.User, error) {
					return &services.User{
						ID:    id,
						Email: "test@example.com",
						Name:  "Test User",
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "user not found",
			userID: "999",
			mockService: &mockUserService{
				getUserByIDFunc: func(ctx context.Context, id string) (*services.User, error) {
					return nil, services.ErrUserNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "service error",
			userID: "123",
			mockService: &mockUserService{
				getUserByIDFunc: func(ctx context.Context, id string) (*services.User, error) {
					return nil, errors.New("database error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewUserHandler(tt.mockService, testLogger())

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+tt.userID, nil)
			rec := httptest.NewRecorder()

			// Set up chi context with URL param
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handler.GetUserByID(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		mockService    *mockUserService
		expectedStatus int
	}{
		{
			name: "success",
			body: CreateUserRequest{
				Email: "new@example.com",
				Name:  "New User",
			},
			mockService: &mockUserService{
				createUserFunc: func(ctx context.Context, email, name string) (*services.User, error) {
					return &services.User{
						ID:    "new123",
						Email: email,
						Name:  name,
						Role:  "user",
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid email",
			body: CreateUserRequest{
				Email: "invalid-email",
				Name:  "User",
			},
			mockService:    &mockUserService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing name",
			body: CreateUserRequest{
				Email: "test@example.com",
				Name:  "",
			},
			mockService:    &mockUserService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "email already exists",
			body: CreateUserRequest{
				Email: "existing@example.com",
				Name:  "User",
			},
			mockService: &mockUserService{
				createUserFunc: func(ctx context.Context, email, name string) (*services.User, error) {
					return nil, services.ErrEmailAlreadyExists
				},
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "invalid json",
			body:           "invalid json",
			mockService:    &mockUserService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			mockService: &mockUserService{
				createUserFunc: func(ctx context.Context, email, name string) (*services.User, error) {
					return nil, errors.New("database error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewUserHandler(tt.mockService, testLogger())

			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.CreateUser(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		body           interface{}
		mockService    *mockUserService
		expectedStatus int
	}{
		{
			name:   "success - update name",
			userID: "123",
			body: UpdateUserRequest{
				Name: "Updated Name",
			},
			mockService: &mockUserService{
				updateUserFunc: func(ctx context.Context, id string, updates map[string]interface{}) (*services.User, error) {
					return &services.User{
						ID:    id,
						Email: "test@example.com",
						Name:  "Updated Name",
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "user not found",
			userID: "999",
			body: UpdateUserRequest{
				Name: "Updated Name",
			},
			mockService: &mockUserService{
				updateUserFunc: func(ctx context.Context, id string, updates map[string]interface{}) (*services.User, error) {
					return nil, services.ErrUserNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "email already exists",
			userID: "123",
			body: UpdateUserRequest{
				Email: "existing@example.com",
			},
			mockService: &mockUserService{
				updateUserFunc: func(ctx context.Context, id string, updates map[string]interface{}) (*services.User, error) {
					return nil, services.ErrEmailAlreadyExists
				},
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:   "invalid role",
			userID: "123",
			body: UpdateUserRequest{
				Role: "superadmin",
			},
			mockService:    &mockUserService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid json",
			userID:         "123",
			body:           "invalid json",
			mockService:    &mockUserService{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewUserHandler(tt.mockService, testLogger())

			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPut, "/api/v1/users/"+tt.userID, bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Set up chi context with URL param
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handler.UpdateUser(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockService    *mockUserService
		expectedStatus int
	}{
		{
			name:   "success",
			userID: "123",
			mockService: &mockUserService{
				deleteUserFunc: func(ctx context.Context, id string) error {
					return nil
				},
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "user not found",
			userID: "999",
			mockService: &mockUserService{
				deleteUserFunc: func(ctx context.Context, id string) error {
					return services.ErrUserNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "service error",
			userID: "123",
			mockService: &mockUserService{
				deleteUserFunc: func(ctx context.Context, id string) error {
					return errors.New("database error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewUserHandler(tt.mockService, testLogger())

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+tt.userID, nil)
			rec := httptest.NewRecorder()

			// Set up chi context with URL param
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handler.DeleteUser(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
