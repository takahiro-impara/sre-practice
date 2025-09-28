package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	ListUsers(ctx context.Context, limit int32, offset int32) ([]*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email domain.Email) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
