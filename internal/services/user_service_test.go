package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewUserService(t *testing.T) {
	service := NewUserService()
	if service == nil {
		t.Fatal("NewUserService should not return nil")
	}

	// Check if test data is initialized
	users, err := service.GetAllUsers(context.Background())
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 initial users, got %d", len(users))
	}
}

func TestGetUserByID(t *testing.T) {
	service := NewUserService()
	ctx := context.Background()

	tests := []struct {
		name    string
		userID  string
		wantErr error
		wantNil bool
	}{
		{
			name:    "existing user",
			userID:  "usr_001",
			wantErr: nil,
			wantNil: false,
		},
		{
			name:    "non-existing user",
			userID:  "usr_999",
			wantErr: ErrUserNotFound,
			wantNil: true,
		},
		{
			name:    "empty user ID",
			userID:  "",
			wantErr: ErrInvalidUserID,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.GetUserByID(ctx, tt.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNil && user != nil {
				t.Errorf("Expected nil user, got %+v", user)
			}
			if !tt.wantNil && user == nil {
				t.Error("Expected non-nil user, got nil")
			}
			if user != nil && user.ID != tt.userID {
				t.Errorf("Expected user ID %s, got %s", tt.userID, user.ID)
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	service := NewUserService()
	ctx := context.Background()

	users, err := service.GetAllUsers(ctx)
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	// Verify users are copies (not references)
	if len(users) > 0 {
		originalEmail := users[0].Email
		users[0].Email = "modified@example.com"

		// Fetch again to verify original wasn't modified
		users2, _ := service.GetAllUsers(ctx)
		if users2[0].Email != originalEmail {
			t.Error("GetAllUsers should return copies, not references")
		}
	}
}

func TestCreateUser(t *testing.T) {
	service := NewUserService()
	ctx := context.Background()

	tests := []struct {
		name     string
		email    string
		userName string
		wantErr  error
	}{
		{
			name:     "valid user",
			email:    "new@example.com",
			userName: "New User",
			wantErr:  nil,
		},
		{
			name:     "duplicate email",
			email:    "john.doe@example.com",
			userName: "Another John",
			wantErr:  ErrEmailAlreadyExists,
		},
		{
			name:     "empty email",
			email:    "",
			userName: "No Email",
			wantErr:  ErrInvalidEmail,
		},
		{
			name:     "empty name",
			email:    "valid@example.com",
			userName: "",
			wantErr:  errors.New("name is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.CreateUser(ctx, tt.email, tt.userName)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.wantErr)
				} else if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				if user != nil {
					t.Errorf("Expected nil user on error, got %+v", user)
				}
			} else {
				if err != nil {
					t.Errorf("CreateUser() unexpected error: %v", err)
				}
				if user == nil {
					t.Fatal("Expected non-nil user, got nil")
				}
				if user.Email != tt.email {
					t.Errorf("Expected email %s, got %s", tt.email, user.Email)
				}
				if user.Name != tt.userName {
					t.Errorf("Expected name %s, got %s", tt.userName, user.Name)
				}
				if user.Role != "user" {
					t.Errorf("Expected role 'user', got %s", user.Role)
				}
				if user.ID == "" {
					t.Error("User ID should not be empty")
				}
				if user.CreatedAt.IsZero() {
					t.Error("CreatedAt should not be zero")
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	service := NewUserService()
	ctx := context.Background()

	// Create a test user first
	testUser, _ := service.CreateUser(ctx, "test@example.com", "Test User")

	tests := []struct {
		name    string
		userID  string
		updates map[string]interface{}
		wantErr error
		verify  func(*User)
	}{
		{
			name:   "update name",
			userID: testUser.ID,
			updates: map[string]interface{}{
				"name": "Updated Name",
			},
			wantErr: nil,
			verify: func(u *User) {
				if u.Name != "Updated Name" {
					t.Errorf("Expected name 'Updated Name', got %s", u.Name)
				}
			},
		},
		{
			name:   "update email",
			userID: testUser.ID,
			updates: map[string]interface{}{
				"email": "updated@example.com",
			},
			wantErr: nil,
			verify: func(u *User) {
				if u.Email != "updated@example.com" {
					t.Errorf("Expected email 'updated@example.com', got %s", u.Email)
				}
			},
		},
		{
			name:   "update role",
			userID: testUser.ID,
			updates: map[string]interface{}{
				"role": "admin",
			},
			wantErr: nil,
			verify: func(u *User) {
				if u.Role != "admin" {
					t.Errorf("Expected role 'admin', got %s", u.Role)
				}
			},
		},
		{
			name:   "update to existing email",
			userID: testUser.ID,
			updates: map[string]interface{}{
				"email": "john.doe@example.com",
			},
			wantErr: ErrEmailAlreadyExists,
			verify:  nil,
		},
		{
			name:   "update non-existent user",
			userID: "usr_999",
			updates: map[string]interface{}{
				"name": "Ghost",
			},
			wantErr: ErrUserNotFound,
			verify:  nil,
		},
		{
			name:   "empty user ID",
			userID: "",
			updates: map[string]interface{}{
				"name": "No ID",
			},
			wantErr: ErrInvalidUserID,
			verify:  nil,
		},
		{
			name:   "multiple updates",
			userID: testUser.ID,
			updates: map[string]interface{}{
				"name":  "Multi Update",
				"email": "multi@example.com",
				"role":  "moderator",
			},
			wantErr: nil,
			verify: func(u *User) {
				if u.Name != "Multi Update" {
					t.Errorf("Expected name 'Multi Update', got %s", u.Name)
				}
				if u.Email != "multi@example.com" {
					t.Errorf("Expected email 'multi@example.com', got %s", u.Email)
				}
				if u.Role != "moderator" {
					t.Errorf("Expected role 'moderator', got %s", u.Role)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.UpdateUser(ctx, tt.userID, tt.updates)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				if user != nil {
					t.Errorf("Expected nil user on error, got %+v", user)
				}
			} else {
				if user == nil {
					t.Fatal("Expected non-nil user, got nil")
				}
				if tt.verify != nil {
					tt.verify(user)
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	service := NewUserService()
	ctx := context.Background()

	// Create a test user
	testUser, _ := service.CreateUser(ctx, "delete@example.com", "Delete Me")

	tests := []struct {
		name    string
		userID  string
		wantErr error
	}{
		{
			name:    "delete existing user",
			userID:  testUser.ID,
			wantErr: nil,
		},
		{
			name:    "delete non-existent user",
			userID:  "usr_999",
			wantErr: ErrUserNotFound,
		},
		{
			name:    "empty user ID",
			userID:  "",
			wantErr: ErrInvalidUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteUser(ctx, tt.userID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify user is actually deleted
			if err == nil {
				_, getErr := service.GetUserByID(ctx, tt.userID)
				if !errors.Is(getErr, ErrUserNotFound) {
					t.Error("User should not exist after deletion")
				}
			}
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	service := NewUserService()
	ctx := context.Background()

	// Test concurrent reads and writes
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				switch j % 4 {
				case 0:
					// Create
					email := fmt.Sprintf("concurrent%d_%d@example.com", id, j)
					_, _ = service.CreateUser(ctx, email, "Concurrent User")
				case 1:
					// Read all
					_, _ = service.GetAllUsers(ctx)
				case 2:
					// Read one
					_, _ = service.GetUserByID(ctx, "usr_001")
				case 3:
					// Update
					_, _ = service.UpdateUser(ctx, "usr_001", map[string]interface{}{
						"name": fmt.Sprintf("Updated %d", time.Now().UnixNano()),
					})
				}
			}
		}(i)
	}

	wg.Wait()

	// If we get here without deadlock or panic, the concurrent access is safe
}

func TestUserDataIsolation(t *testing.T) {
	service := NewUserService()
	ctx := context.Background()

	// Get a user
	user1, err := service.GetUserByID(ctx, "usr_001")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	originalName := user1.Name

	// Modify the returned user
	user1.Name = "Modified Outside"

	// Get the user again
	user2, err := service.GetUserByID(ctx, "usr_001")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// The modification should not affect the stored user
	if user2.Name != originalName {
		t.Errorf("User data not properly isolated. Expected %s, got %s", originalName, user2.Name)
	}
}
