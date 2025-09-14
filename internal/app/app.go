package app

import (
	"context"
	"log/slog"

	"github.com/mikko-kohtala/go-api/internal/config"
	"github.com/mikko-kohtala/go-api/internal/services"
	"github.com/mikko-kohtala/go-api/internal/telemetry"
)

// App holds all application dependencies
type App struct {
	Config         *config.Config
	Logger         *slog.Logger
	Metrics        *telemetry.Metrics
	ExampleService services.ExampleService
	// Add more services and repositories here as needed
}

// New creates a new application with all dependencies
func New(cfg *config.Config, logger *slog.Logger) *App {
	return &App{
		Config:         cfg,
		Logger:         logger,
		Metrics:        telemetry.NewMetrics(),
		ExampleService: services.NewExampleService(),
	}
}

// Shutdown gracefully shuts down application dependencies
func (a *App) Shutdown(ctx context.Context) error {
	// Clean up resources here (database connections, etc.)
	a.Logger.Info("application shutdown complete")
	return nil
}