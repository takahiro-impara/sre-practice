package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponse, error)
	GetUserByEmail(ctx context.Context, email domain.Email) (*UserResponse, error)
	UpdateUser(ctx context.Context, req UpdateUserRequest) error
	DeleteUser(ctx context.Context, req DeleteUserRequest) error
	ListUsers(ctx context.Context, req ListUsersRequest) ([]*UserResponse, error)
	AuthenticateUser(ctx context.Context, req AuthenticateUserRequest) error
}
