package routes

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikko-kohtala/go-api/internal/handlers"
	"github.com/mikko-kohtala/go-api/internal/services"
)

type Routes struct {
	logger       *slog.Logger
	userService  services.UserService
	statsService services.StatsService
	userHandler  *handlers.UserHandler
	statsHandler *handlers.StatsHandler
	env          string
}

func NewRoutes(
	logger *slog.Logger,
	userService services.UserService,
	statsService services.StatsService,
) *Routes {
	return &Routes{
		logger:       logger,
		userService:  userService,
		statsService: statsService,
		userHandler:  handlers.NewUserHandler(userService, logger),
		statsHandler: handlers.NewStatsHandler(statsService, logger),
		env:          "", // Will be set by SetEnvironment
	}
}

// SetEnvironment sets the environment for conditional route registration
func (rt *Routes) SetEnvironment(env string) {
	rt.env = env
}

// SetupHealthRoutes configures health check endpoints
func (rt *Routes) SetupHealthRoutes(r chi.Router) {
	r.Get("/healthz", handlers.Health)
	r.Get("/readyz", handlers.Ready)
}

// SetupAPIV1Routes configures API v1 endpoints
func (rt *Routes) SetupAPIV1Routes(r chi.Router) {
	// Example endpoints (existing)
	r.Get("/ping", handlers.Ping)
	r.Post("/echo", handlers.Echo)

	// User endpoints (new)
	r.Route("/users", func(r chi.Router) {
		r.Get("/", rt.userHandler.GetAllUsers)
		r.Post("/", rt.userHandler.CreateUser)
		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", rt.userHandler.GetUserByID)
			r.Put("/", rt.userHandler.UpdateUser)
			r.Delete("/", rt.userHandler.DeleteUser)
		})
	})

	// Stats endpoints (new)
	r.Route("/stats", func(r chi.Router) {
		r.Get("/system", rt.statsHandler.GetSystemStats)
		r.Get("/api", rt.statsHandler.GetAPIStats)
	})
}

// SetupRootRoute configures the root endpoint
func (rt *Routes) SetupRootRoute(r chi.Router) {
	r.Get("/", handlers.Root)
}

// SetupTestRoutes configures test/debug endpoints
func (rt *Routes) SetupTestRoutes(r chi.Router) {
	// Only enable test routes in development/test environments
	if rt.env == "production" || rt.env == "prod" {
		rt.logger.Info("Test routes disabled in production environment")
		return
	}

	r.Get("/logs", handlers.TestLogs)
}

// SetupSwaggerRoutes configures Swagger documentation routes
func (rt *Routes) SetupSwaggerRoutes(r chi.Router, swaggerHandler http.HandlerFunc) {
	r.Get("/swagger/*", swaggerHandler)

	// Alias the Swagger UI under /api-docs as well
	r.Get("/api-docs", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/api-docs/index.html", http.StatusTemporaryRedirect)
	})
	r.Get("/api-docs/*", swaggerHandler)
}
