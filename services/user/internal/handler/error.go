package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
	"go.uber.org/zap"
)

func (h *UserHandler) handleServiceError(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Error("Service error", zap.Error(err))

	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		h.renderError(w, r, http.StatusNotFound, "User not found")
	case errors.Is(err, domain.ErrUserAlreadyExists):
		h.renderError(w, r, http.StatusConflict, "User already exists")
	case errors.Is(err, domain.ErrInvalidEmail):
		h.renderError(w, r, http.StatusBadRequest, "Invalid email address")
	case errors.Is(err, domain.ErrInvalidName):
		h.renderError(w, r, http.StatusBadRequest, "Invalid name")
	case errors.Is(err, domain.ErrInvalidPassword):
		h.renderError(w, r, http.StatusBadRequest, "Invalid password")
	default:
		// 詳細なエラーメッセージを表示
		h.renderError(w, r, http.StatusInternalServerError, err.Error())
	}
}

func (h *UserHandler) renderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	h.renderErrorWithDetails(w, r, status, message, nil)
}

func (h *UserHandler) renderErrorWithDetails(w http.ResponseWriter, r *http.Request, status int, message string, details map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResp := ErrorResponse{
		Error:   message,
		Code:    http.StatusText(status),
		Details: details,
	}

	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		h.logger.Error("Failed to encode error response", zap.Error(err))
	}
}

func (h *UserHandler) renderValidationError(w http.ResponseWriter, r *http.Request, validationErrors map[string]string) {
	h.renderErrorWithDetails(w, r, http.StatusBadRequest, "Validation failed", validationErrors)
}
