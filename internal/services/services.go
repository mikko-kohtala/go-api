package services

import (
	"context"
	"log/slog"
)

// Service defines the interface for business logic services
type Service interface {
	Health(ctx context.Context) error
}

// EchoService handles echo-related business logic
type EchoService struct {
	logger *slog.Logger
}

// NewEchoService creates a new EchoService
func NewEchoService(logger *slog.Logger) *EchoService {
	return &EchoService{
		logger: logger,
	}
}

// Echo processes an echo request
func (s *EchoService) Echo(ctx context.Context, message string) (string, error) {
	if s.logger != nil {
		s.logger.InfoContext(ctx, "processing echo request", slog.String("message", message))
	}
	
	// Add any business logic here
	// For example: validation, transformation, external API calls, etc.
	
	return message, nil
}

// Health checks the health of the echo service
func (s *EchoService) Health(ctx context.Context) error {
	// Add health checks for dependencies here
	// For example: database connectivity, external service availability, etc.
	return nil
}

// ServiceContainer holds all services
type ServiceContainer struct {
	Echo *EchoService
}

// NewServiceContainer creates a new service container with all services
func NewServiceContainer(logger *slog.Logger) *ServiceContainer {
	return &ServiceContainer{
		Echo: NewEchoService(logger),
	}
}