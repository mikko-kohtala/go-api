package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mikko-kohtala/go-api/internal/services"
)

func testUserHandler() (*UserHandler, services.UserService) {
	svc := services.NewUserService()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewUserHandler(svc, logger), svc
}

func TestUserHandler_CreateUserSuccess(t *testing.T) {
	handler, _ := testUserHandler()
	rr := httptest.NewRecorder()
	body, _ := json.Marshal(map[string]string{"email": "success@example.com", "name": "Example"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	handler.CreateUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rr.Code)
	}
}

func TestUserHandler_CreateUserValidation(t *testing.T) {
	handler, _ := testUserHandler()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{"email": "invalid"}`))
	req.Header.Set("Content-Type", "application/json")

	handler.CreateUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestUserHandler_GetUserByID(t *testing.T) {
	handler, svc := testUserHandler()
	user, _ := svc.CreateUser(context.Background(), "byid@example.com", "By ID")

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+user.ID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", user.ID)
	req = req.WithContext(contextWithRoute(req.Context(), rctx))

	handler.GetUserByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestUserHandler_GetUserByID_NotFound(t *testing.T) {
	handler, _ := testUserHandler()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/unknown", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "unknown")
	req = req.WithContext(contextWithRoute(req.Context(), rctx))

	handler.GetUserByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	handler, svc := testUserHandler()
	user, _ := svc.CreateUser(context.Background(), "update.handler@example.com", "Handler")

	rr := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"name":"Updated"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/"+user.ID, body)
	req.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", user.ID)
	req = req.WithContext(contextWithRoute(req.Context(), rctx))

	handler.UpdateUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	handler, svc := testUserHandler()
	user, _ := svc.CreateUser(context.Background(), "delete.handler@example.com", "Delete")

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+user.ID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", user.ID)
	req = req.WithContext(contextWithRoute(req.Context(), rctx))

	handler.DeleteUser(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func contextWithRoute(ctx context.Context, routeCtx *chi.Context) context.Context {
	return context.WithValue(ctx, chi.RouteCtxKey, routeCtx)
}
