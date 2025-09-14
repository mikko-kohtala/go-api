package handlers

import (
    "net/http"
    "github.com/mikko-kohtala/go-api/internal/response"
    "github.com/mikko-kohtala/go-api/internal/services"
)

// Health godoc
// @Summary      Liveness probe
// @Description  Simple health check indicating the service is up.
// @Tags         health
// @Success      200 {object} map[string]string
// @Router       /healthz [get]
func Health(w http.ResponseWriter, r *http.Request) {
    response.JSON(w, r, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready godoc
// @Summary      Readiness probe
// @Description  Indicates whether the service is ready to accept traffic.
// @Tags         health
// @Success      200 {object} map[string]string
// @Router       /readyz [get]
func Ready(w http.ResponseWriter, r *http.Request) {
    // In a real app, check dependencies (DB, cache, etc.)
    response.JSON(w, r, http.StatusOK, map[string]string{"ready": "true"})
}

// HealthHandler holds dependencies for health checks
type HealthHandler struct {
    services *services.ServiceContainer
}

// NewHealthHandler creates a new health handler with dependencies
func NewHealthHandler(services *services.ServiceContainer) *HealthHandler {
    return &HealthHandler{
        services: services,
    }
}

// ReadyWithDependencies checks readiness including service dependencies
func (h *HealthHandler) ReadyWithDependencies(w http.ResponseWriter, r *http.Request) {
    // Check service health
    if err := h.services.Echo.Health(r.Context()); err != nil {
        response.JSON(w, r, http.StatusServiceUnavailable, map[string]string{
            "ready": "false",
            "error": err.Error(),
        })
        return
    }
    
    response.JSON(w, r, http.StatusOK, map[string]string{"ready": "true"})
}
