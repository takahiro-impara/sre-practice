package handler

import (
	"github.com/google/uuid"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
)

type CreateUserRequest struct {
	Email    domain.Email    `json:"email" validate:"required,email"`
	Name     domain.Name     `json:"name" validate:"required,min=1,max=255"`
	Password domain.Password `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
	Email *domain.Email `json:"email,omitempty" validate:"omitempty,email"`
	Name  *domain.Name  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
}

type DeleteUserRequest struct {
	ID uuid.UUID `json:"id"`
}

type AuthenticateUserRequest struct {
	Email    domain.Email    `json:"email" validate:"required,email"`
	Password domain.Password `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type ListUsersResponse struct {
	Users      []*UserResponse `json:"users"`
	TotalCount int             `json:"total_count"`
	Limit      int             `json:"limit"`
	Offset     int             `json:"offset"`
}

type ErrorResponse struct {
	Error   string            `json:"error"`
	Code    string            `json:"code,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}
