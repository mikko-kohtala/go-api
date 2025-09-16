package handlers

import (
	"log/slog"
	"net/http"

	"github.com/mikko-kohtala/go-api/internal/response"
	"github.com/mikko-kohtala/go-api/internal/services"
)

type StatsHandler struct {
	statsService services.StatsService
	logger       *slog.Logger
}

func NewStatsHandler(statsService services.StatsService, logger *slog.Logger) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
		logger:       logger,
	}
}

// GetSystemStats godoc
// @Summary      Get system statistics
// @Description  Returns current system statistics including memory usage, goroutines, etc.
// @Tags         stats
// @Produce      json
// @Success      200 {object} services.SystemStats
// @Failure      500 {object} map[string]interface{}
// @Router       /api/v1/stats/system [get]
func (h *StatsHandler) GetSystemStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.statsService.GetSystemStats(r.Context())
	if err != nil {
		h.logger.Error("failed to get system stats", slog.String("error", err.Error()))
		response.Error(w, r, http.StatusInternalServerError, "internal_error", "Failed to retrieve system stats", nil)
		return
	}

	response.JSON(w, r, http.StatusOK, stats)
}

// GetAPIStats godoc
// @Summary      Get API statistics
// @Description  Returns API usage statistics including request counts and latencies
// @Tags         stats
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/v1/stats/api [get]
func (h *StatsHandler) GetAPIStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.statsService.GetAPIStats(r.Context())
	if err != nil {
		h.logger.Error("failed to get API stats", slog.String("error", err.Error()))
		response.Error(w, r, http.StatusInternalServerError, "internal_error", "Failed to retrieve API stats", nil)
		return
	}

	response.JSON(w, r, http.StatusOK, stats)
}
