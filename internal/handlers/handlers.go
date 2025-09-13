package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/user/go-api-template/internal/models"
	"github.com/user/go-api-template/pkg/logger"
	"github.com/user/go-api-template/pkg/response"
	"github.com/user/go-api-template/pkg/validation"
)

type Handler struct {
	logger     *slog.Logger
	users      map[string]*models.User
	emailIndex map[string]string // email -> userID
	mu         sync.RWMutex
}

func New(log *slog.Logger) *Handler {
	return &Handler{
		logger:     log,
		users:      make(map[string]*models.User),
		emailIndex: make(map[string]string),
		mu:         sync.RWMutex{},
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	h.mu.RLock()
	users := make([]*models.UserResponse, 0, len(h.users))
	for _, user := range h.users {
		users = append(users, &models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
	h.mu.RUnlock()

	log.Info("fetched users", "count", len(users))

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"users": users,
		"total": len(users),
	})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())
	id := r.PathValue("id")

	if _, err := uuid.Parse(id); err != nil {
		log.Warn("invalid user ID format", "id", id)
		response.Error(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	h.mu.RLock()
	user, exists := h.users[id]
	h.mu.RUnlock()

	if !exists {
		log.Info("user not found", "user_id", id)
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	log.Info("fetched user", "user_id", user.ID)

	response.JSON(w, http.StatusOK, &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("failed to decode request body", "error", err)
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validation.Validate(&req); err != nil {
		validationErrors := validation.FormatValidationErrors(err)
		log.Info("validation failed", "errors", validationErrors)
		response.ErrorWithDetails(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	h.mu.RLock()
	if _, exists := h.emailIndex[req.Email]; exists {
		h.mu.RUnlock()
		log.Info("email already exists", "email", req.Email)
		response.Error(w, http.StatusConflict, "email already exists")
		return
	}
	h.mu.RUnlock()

	user := &models.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	h.mu.Lock()
	h.users[user.ID] = user
	h.emailIndex[user.Email] = user.ID
	h.mu.Unlock()

	log.Info("user created",
		"user_id", user.ID,
		"email", user.Email,
	)

	response.JSON(w, http.StatusCreated, &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())
	id := r.PathValue("id")

	if _, err := uuid.Parse(id); err != nil {
		log.Warn("invalid user ID format", "id", id)
		response.Error(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	h.mu.RLock()
	user, exists := h.users[id]
	h.mu.RUnlock()

	if !exists {
		log.Info("user not found", "user_id", id)
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("failed to decode request body", "error", err)
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validation.Validate(&req); err != nil {
		validationErrors := validation.FormatValidationErrors(err)
		log.Info("validation failed", "errors", validationErrors)
		response.ErrorWithDetails(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	if req.Email != "" && req.Email != user.Email {
		h.mu.RLock()
		if existingID, exists := h.emailIndex[req.Email]; exists && existingID != id {
			h.mu.RUnlock()
			log.Info("email already exists", "email", req.Email)
			response.Error(w, http.StatusConflict, "email already exists")
			return
		}
		h.mu.RUnlock()
	}

	h.mu.Lock()
	if req.Email != "" && req.Email != user.Email {
		delete(h.emailIndex, user.Email)
		user.Email = req.Email
		h.emailIndex[req.Email] = user.ID
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	user.UpdatedAt = time.Now()
	h.mu.Unlock()

	log.Info("user updated",
		"user_id", user.ID,
		"email", user.Email,
	)

	response.JSON(w, http.StatusOK, &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())
	id := r.PathValue("id")

	if _, err := uuid.Parse(id); err != nil {
		log.Warn("invalid user ID format", "id", id)
		response.Error(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	h.mu.Lock()
	user, exists := h.users[id]
	if !exists {
		h.mu.Unlock()
		log.Info("user not found", "user_id", id)
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	delete(h.emailIndex, user.Email)
	delete(h.users, id)
	h.mu.Unlock()

	log.Info("user deleted", "user_id", id)

	w.WriteHeader(http.StatusNoContent)
}

