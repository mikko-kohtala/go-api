package handlers

import (
    "encoding/json"
    "log/slog"
    "net/http"
    "github.com/mikko-kohtala/go-api/internal/logging"
)

// Health godoc
// @Summary      Liveness probe
// @Description  Simple health check indicating the service is up.
// @Tags         health
// @Success      200 {object} map[string]string
// @Router       /healthz [get]
func Health(w http.ResponseWriter, r *http.Request) {
    respondJSON(r, w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready godoc
// @Summary      Readiness probe
// @Description  Indicates whether the service is ready to accept traffic.
// @Tags         health
// @Success      200 {object} map[string]string
// @Router       /readyz [get]
func Ready(w http.ResponseWriter, r *http.Request) {
    // In a real app, check dependencies (DB, cache, etc.)
    respondJSON(r, w, http.StatusOK, map[string]string{"ready": "true"})
}

func respondJSON(r *http.Request, w http.ResponseWriter, code int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    if err := json.NewEncoder(w).Encode(v); err != nil {
        if l := logging.FromContext(r.Context()); l != nil {
            l.Error("failed to encode JSON response", slog.String("error", err.Error()))
        }
        http.Error(w, "internal server error", http.StatusInternalServerError)
    }
}
