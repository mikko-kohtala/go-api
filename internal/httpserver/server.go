package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	docs "github.com/mikko-kohtala/go-api/internal/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/mikko-kohtala/go-api/internal/config"
	"github.com/mikko-kohtala/go-api/internal/routes"
	"github.com/mikko-kohtala/go-api/internal/services"
)

// NewRouter assembles the chi router with middleware and routes.
// This function only builds the server structure - all handlers are defined in the handlers package.
func NewRouter(cfg *config.Config, appLogger *slog.Logger) http.Handler {
	// Initialize services
	userService := services.NewUserService()
	statsService := services.NewStatsService()

	// Initialize routes with services
	routesHandler := routes.NewRoutes(appLogger, userService, statsService)

	r := chi.NewRouter()

	// Setup middleware
	setupMiddleware(r, cfg, appLogger)

	// Setup rate limiting
	apiRate := setupRateLimiting(cfg, appLogger)

	// Setup all routes
	setupRoutes(r, routesHandler, apiRate)

	// Setup Swagger documentation
	setupSwagger(r, routesHandler)

	return r
}

// setupMiddleware configures all middleware for the router
func setupMiddleware(r chi.Router, cfg *config.Config, appLogger *slog.Logger) {
	// Core middleware (place timeout early to bound all work)
	r.Use(middleware.Timeout(cfg.RequestTimeout))
	r.Use(BodyLimit(cfg.BodyLimitBytes))
	r.Use(RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Compress(cfg.CompressionLevel))
	r.Use(LoggingMiddleware(appLogger))
	r.Use(middleware.Recoverer)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedMethods:   cfg.CORSAllowedMethods,
		AllowedHeaders:   cfg.CORSAllowedHeaders,
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Warn if permissive CORS in production
	if cfg.Env == "production" || cfg.Env == "prod" {
		for _, o := range cfg.CORSAllowedOrigins {
			if o == "*" {
				appLogger.Warn("CORS allows all origins in production; consider restricting AllowedOrigins")
				break
			}
		}
	}
}

// setupRateLimiting configures rate limiting middleware
func setupRateLimiting(cfg *config.Config, appLogger *slog.Logger) func(http.Handler) http.Handler {
	if !cfg.RateLimitEnabled {
		return func(h http.Handler) http.Handler { return h }
	}

	period, err := time.ParseDuration(cfg.RateLimitPeriod)
	if err != nil || period <= 0 {
		appLogger.Error("invalid rate limit period; disabling rate limit",
			slog.String("period", cfg.RateLimitPeriod),
			slog.Any("error", err))
		return func(h http.Handler) http.Handler { return h }
	}

	return httprate.LimitByIP(cfg.RateLimit, period)
}

// setupRoutes configures all application routes
func setupRoutes(r chi.Router, routesHandler *routes.Routes, apiRate func(http.Handler) http.Handler) {
	// Health endpoints (no rate limiting)
	r.Group(func(r chi.Router) {
		routesHandler.SetupHealthRoutes(r)
	})

	// API v1 routes (with rate limiting)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(apiRate)
		routesHandler.SetupAPIV1Routes(r)
	})

	// Test routes
	r.Route("/test", func(r chi.Router) {
		routesHandler.SetupTestRoutes(r)
	})

	// Root route
	routesHandler.SetupRootRoute(r)
}

// setupSwagger configures Swagger documentation endpoints
func setupSwagger(r chi.Router, routesHandler *routes.Routes) {
	// Configure Swagger info
	docs.SwaggerInfo.Title = "Init Codex API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"

	// Create Swagger handler
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)

	// Setup Swagger routes
	routesHandler.SetupSwaggerRoutes(r, swaggerHandler)
}