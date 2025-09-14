package services

import (
	"context"
	"fmt"
	"time"
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
	// In a real application, this would have a database connection
	// For now, we'll use in-memory storage as an example
	users map[string]*User
}

func NewUserService() UserService {
	// Initialize with some test data
	return &userService{
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
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*User, error) {
	if user, exists := s.users[id]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found: %s", id)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]User, error) {
	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, *user)
	}
	return users, nil
}

func (s *userService) CreateUser(ctx context.Context, email, name string) (*User, error) {
	// Generate a simple ID (in production, use UUID)
	id := fmt.Sprintf("usr_%03d", len(s.users)+1)

	// Check if email already exists
	for _, user := range s.users {
		if user.Email == email {
			return nil, fmt.Errorf("email already exists: %s", email)
		}
	}

	user := &User{
		ID:        id,
		Email:     email,
		Name:      name,
		Role:      "user",
		CreatedAt: time.Now(),
	}

	s.users[id] = user
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) (*User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", id)
	}

	// Apply updates (simplified version)
	if name, ok := updates["name"].(string); ok {
		user.Name = name
	}
	if email, ok := updates["email"].(string); ok {
		user.Email = email
	}
	if role, ok := updates["role"].(string); ok {
		user.Role = role
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	if _, exists := s.users[id]; !exists {
		return fmt.Errorf("user not found: %s", id)
	}
	delete(s.users, id)
	return nil
}