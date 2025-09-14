package services

import (
	"context"
	"fmt"
)

// ExampleService defines the business logic interface
type ExampleService interface {
	Echo(ctx context.Context, message string) (string, error)
	GetStatus(ctx context.Context) (string, error)
}

// exampleService implements ExampleService
type exampleService struct {
	// Add dependencies here (e.g., repositories, external services)
}

// NewExampleService creates a new instance of ExampleService
func NewExampleService() ExampleService {
	return &exampleService{}
}

// Echo processes and returns the message
func (s *exampleService) Echo(ctx context.Context, message string) (string, error) {
	// Business logic can be added here
	// For example: validation, transformation, enrichment

	// Check context cancellation
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("echo operation cancelled: %w", ctx.Err())
	default:
	}

	return message, nil
}

// GetStatus returns the current service status
func (s *exampleService) GetStatus(ctx context.Context) (string, error) {
	// Could check dependencies, database connections, etc.
	return "ok", nil
}