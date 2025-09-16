package handlers

import (
	"net/http"

	"github.com/mikko-kohtala/go-api/internal/response"
	pkglogger "github.com/mikko-kohtala/go-api/pkg/logger"
)

type RootResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Docs    string `json:"docs"`
	Status  string `json:"status"`
}

// Root godoc
// @Summary      API root endpoint
// @Description  Returns basic API information
// @Tags         root
// @Produce      json
// @Success      200 {object} RootResponse
// @Router       / [get]
func Root(w http.ResponseWriter, r *http.Request) {
	// Get logger from context
	if l := pkglogger.FromContext(r.Context()); l != nil {
		l.Info("Root endpoint accessed")
	}

	// Use proper JSON marshaling instead of manual string building
	resp := RootResponse{
		Name:    "go-api",
		Version: "1.0.0",
		Docs:    "/swagger/index.html",
		Status:  "healthy",
	}

	response.JSON(w, r, http.StatusOK, resp)
}
