package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	svc    service.UserService
	logger *zap.Logger
}

func NewUserHandler(svc service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		svc:    svc,
		logger: logger,
	}
}

func toUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     string(user.Email),
		Name:      string(user.Name),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	// デバッグ用ログ
	h.logger.Info("CreateUser request",
		zap.String("email", string(req.Email)),
		zap.String("name", string(req.Name)),
		zap.String("password_length", fmt.Sprintf("%d", len(req.Password))))

	svcReq := service.CreateUserRequest{
		Email:    domain.Email(req.Email),
		Name:     domain.Name(req.Name),
		Password: domain.Password(req.Password),
	}

	user, err := h.svc.CreateUser(ctx, svcReq)
	if err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	user, err := h.svc.GetUserByID(ctx, userID)
	if err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	email := chi.URLParam(r, "email")
	err := domain.ValidateEmail(domain.Email(email))
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid email format")
		return
	}

	user, err := h.svc.GetUserByEmail(ctx, domain.Email(email))
	if err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	render.JSON(w, r, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// URLパラメータからユーザーIDを取得
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	svcReq := service.UpdateUserRequest{
		ID: userID,
	}

	if req.Email != nil {
		svcReq.Email = domain.Email(*req.Email)
	}
	if req.Name != nil {
		svcReq.Name = domain.Name(*req.Name)
	}

	err = h.svc.UpdateUser(ctx, svcReq)
	if err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "User updated successfully"})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	svcReq := service.DeleteUserRequest{ID: userID}
	err = h.svc.DeleteUser(ctx, svcReq)
	if err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toUserResponses(users []*service.UserResponse) []*UserResponse {
	return make([]*UserResponse, 0, len(users))
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	req := service.ListUsersRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	users, err := h.svc.ListUsers(ctx, req)
	if err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	resp := ListUsersResponse{
		Users:      h.toUserResponses(users),
		TotalCount: len(users),
		Limit:      limit,
		Offset:     offset,
	}

	render.JSON(w, r, resp)
}

// AuthenticateUser handles user authentication
func (h *UserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req AuthenticateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		h.renderError(w, r, http.StatusBadRequest, "Email and password are required")
		return
	}

	authReq := service.AuthenticateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.svc.AuthenticateUser(ctx, authReq); err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	render.JSON(w, r, map[string]string{"message": "Authentication successful"})
}

// HealthCheck handles health check
func (h *UserHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// ReadinessCheck handles readiness check
func (h *UserHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	// TODO: Add actual readiness checks (e.g., database connectivity)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}

// Helper methods

func (h *UserHandler) toUserResponse(user *service.UserResponse) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     string(user.Email),
		Name:      string(user.Name),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *UserHandler) toUserResponses(users []*service.UserResponse) []*UserResponse {
	responses := make([]*UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, h.toUserResponse(user))
	}
	return responses
}
