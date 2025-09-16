package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mikko-kohtala/go-api/internal/metrics"
)

// Custom error types for better error handling
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidEmail       = errors.New("invalid email address")
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, email, name string) (*User, error)
	UpdateUser(ctx context.Context, id string, updates map[string]interface{}) (*User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	mu    sync.RWMutex // Protects concurrent access to the users map
	users map[string]*User
}

func NewUserService() UserService {
	// Initialize with some test data
	service := &userService{
		users: map[string]*User{
			"usr_001": {
				ID:        "usr_001",
				Email:     "john.doe@example.com",
				Name:      "John Doe",
				Role:      "admin",
				CreatedAt: time.Now().Add(-24 * time.Hour),
			},
			"usr_002": {
				ID:        "usr_002",
				Email:     "jane.smith@example.com",
				Name:      "Jane Smith",
				Role:      "user",
				CreatedAt: time.Now().Add(-48 * time.Hour),
			},
		},
	}

	// Update metrics
	metrics.UsersTotal.Set(float64(len(service.users)))

	return service
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*User, error) {
	if id == "" {
		metrics.UserOperations.WithLabelValues("get", "error").Inc()
		return nil, ErrInvalidUserID
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if user, exists := s.users[id]; exists {
		// Return a copy to prevent external modifications
		userCopy := *user
		metrics.UserOperations.WithLabelValues("get", "success").Inc()
		return &userCopy, nil
	}
	metrics.UserOperations.WithLabelValues("get", "not_found").Inc()
	return nil, ErrUserNotFound
}

func (s *userService) GetAllUsers(ctx context.Context) ([]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		// Return copies to prevent external modifications
		users = append(users, *user)
	}
	return users, nil
}

func (s *userService) CreateUser(ctx context.Context, email, name string) (*User, error) {
	// Basic validation
	if email == "" {
		metrics.UserOperations.WithLabelValues("create", "invalid_email").Inc()
		return nil, ErrInvalidEmail
	}
	if name == "" {
		metrics.UserOperations.WithLabelValues("create", "invalid_name").Inc()
		return nil, errors.New("name is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if email already exists
	for _, user := range s.users {
		if user.Email == email {
			metrics.UserOperations.WithLabelValues("create", "duplicate").Inc()
			return nil, ErrEmailAlreadyExists
		}
	}

	// Generate a simple ID (in production, use UUID)
	id := fmt.Sprintf("usr_%03d", len(s.users)+1)

	user := &User{
		ID:        id,
		Email:     email,
		Name:      name,
		Role:      "user",
		CreatedAt: time.Now(),
	}

	s.users[id] = user

	// Update metrics
	metrics.UsersTotal.Inc()
	metrics.UserOperations.WithLabelValues("create", "success").Inc()

	// Return a copy
	userCopy := *user
	return &userCopy, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) (*User, error) {
	if id == "" {
		return nil, ErrInvalidUserID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	// Apply updates with validation
	if name, ok := updates["name"].(string); ok && name != "" {
		user.Name = name
	}
	if email, ok := updates["email"].(string); ok && email != "" {
		// Check if new email already exists (except for current user)
		for uid, u := range s.users {
			if uid != id && u.Email == email {
				return nil, ErrEmailAlreadyExists
			}
		}
		user.Email = email
	}
	if role, ok := updates["role"].(string); ok && role != "" {
		user.Role = role
	}

	// Return a copy
	userCopy := *user
	return &userCopy, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		metrics.UserOperations.WithLabelValues("delete", "invalid_id").Inc()
		return ErrInvalidUserID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		metrics.UserOperations.WithLabelValues("delete", "not_found").Inc()
		return ErrUserNotFound
	}
	delete(s.users, id)

	// Update metrics
	metrics.UsersTotal.Dec()
	metrics.UserOperations.WithLabelValues("delete", "success").Inc()

	return nil
}
