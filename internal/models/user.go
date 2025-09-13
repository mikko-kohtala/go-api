package models

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" example:"John Doe"`
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" binding:"omitempty,email" example:"john@example.com"`
}

// Example represents an example resource
type Example struct {
	ID          int       `json:"id" example:"1"`
	Title       string    `json:"title" example:"Example Title"`
	Description string    `json:"description" example:"Example Description"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// CreateExampleRequest represents the request payload for creating an example
type CreateExampleRequest struct {
	Title       string `json:"title" binding:"required" example:"Example Title"`
	Description string `json:"description" example:"Example Description"`
}