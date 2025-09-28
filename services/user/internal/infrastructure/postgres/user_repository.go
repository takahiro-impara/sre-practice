package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/lot-koichi/sre-skill-up-project/services/user/db/sqlc/generated"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/repository"
)

type postgresUserRepository struct {
	db      *sql.DB
	queries *db.Queries
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(database *sql.DB) repository.UserRepository {
	queries := db.New(database)
	return &postgresUserRepository{
		db:      database,
		queries: queries,
	}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	// Use converter function for parameters
	params := toCreateUserParams(user)

	createdUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return handlePostgresError(err)
	}

	// Update timestamps using converter function
	updateTimestamps(user, createdUser)
	return nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, handlePostgresError(err)
	}

	// Use converter function
	return toDomainUser(user), nil
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, string(email))
	if err != nil {
		return nil, handlePostgresError(err)
	}
	return toDomainUser(user), nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	// Use converter function for parameters
	params := toUpdateUserParams(user)

	updatedUser, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return handlePostgresError(err)
	}

	// Update only UpdatedAt using converter function
	updateUpdatedAt(user, updatedUser)
	return nil
}

func (r *postgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return handlePostgresError(err)
	}
	return nil
}

func (r *postgresUserRepository) ListUsers(ctx context.Context, limit int32, offset int32) ([]*domain.User, error) {
	// Use converter function for parameters
	params := toListUsersParams(limit, offset)

	users, err := r.queries.ListUsers(ctx, params)
	if err != nil {
		return nil, handlePostgresError(err)
	}

	// Use converter function for batch conversion
	return toDomainUsers(users), nil
}
