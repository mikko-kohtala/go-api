package services

import (
	"context"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	svc := NewUserService()

	user, err := svc.CreateUser(context.Background(), "new.user@example.com", "New User")
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}
	if user.ID == "" {
		t.Fatalf("expected ID to be set")
	}
	if user.Email != "new.user@example.com" {
		t.Fatalf("expected email to match, got %s", user.Email)
	}

	if _, err := svc.CreateUser(context.Background(), "john.doe@example.com", "Dup"); err != ErrEmailAlreadyExists {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	svc := NewUserService()

	user, err := svc.CreateUser(context.Background(), "update@example.com", "Update User")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	updated, err := svc.UpdateUser(context.Background(), user.ID, map[string]interface{}{"name": "Renamed", "role": "admin"})
	if err != nil {
		t.Fatalf("UpdateUser returned error: %v", err)
	}
	if updated.Name != "Renamed" || updated.Role != "admin" {
		t.Fatalf("expected updated fields to persist")
	}

	// Attempt to reuse existing email
	if _, err := svc.UpdateUser(context.Background(), user.ID, map[string]interface{}{"email": "john.doe@example.com"}); err != ErrEmailAlreadyExists {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}

	if _, err := svc.UpdateUser(context.Background(), "missing", map[string]interface{}{}); err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	svc := NewUserService()

	if err := svc.DeleteUser(context.Background(), "usr_001"); err != nil {
		t.Fatalf("DeleteUser returned error: %v", err)
	}
	if err := svc.DeleteUser(context.Background(), "usr_001"); err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_GetUserByIDValidation(t *testing.T) {
	svc := NewUserService()

	if _, err := svc.GetUserByID(context.Background(), ""); err != ErrInvalidUserID {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
}
