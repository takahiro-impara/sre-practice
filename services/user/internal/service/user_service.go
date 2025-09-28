package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/repository"
	"go.uber.org/zap"
)

// Ensure userService implements UserService interface
var _ UserService = (*userService)(nil)

// userService provides business logic for user operations
type userService struct {
	repo   repository.UserRepository
	hasher PasswordHasher
	logger *zap.Logger
}

// NewUserService creates a new UserService instance
func NewUserService(repo repository.UserRepository, hasher PasswordHasher, logger *zap.Logger) UserService {
	return &userService{
		repo:   repo,
		hasher: hasher,
		logger: logger,
	}
}

type UserResponse struct {
	ID        uuid.UUID    `json:"id"`
	Email     domain.Email `json:"email"`
	Name      domain.Name  `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type CreateUserRequest struct {
	Email    domain.Email    `json:"email"`
	Name     domain.Name     `json:"name"`
	Password domain.Password `json:"password"`
}

// CreateUser creates a new user with the given email and name
func (s *userService) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	if req.Email == "" {
		return nil, domain.ErrInvalidEmail
	}
	if req.Name == "" {
		return nil, domain.ErrInvalidName
	}
	if req.Password == "" {
		return nil, domain.ErrInvalidPassword
	}

	// まず元のパスワードでバリデーション
	s.logger.Info("Validating password", zap.String("password_length", fmt.Sprintf("%d", len(req.Password))))
	if err := domain.ValidatePassword(req.Password); err != nil {
		s.logger.Error("Password validation failed", zap.Error(err))
		return nil, fmt.Errorf("password validation failed: %w", err)
	}

	// パスワードのハッシュ化
	s.logger.Info("Hashing password")
	hashedPassword, err := s.hasher.Hash(req.Password)
	if err != nil {
		s.logger.Error("Password hashing failed", zap.Error(err))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// ドメインモデルの作成
	s.logger.Info("Creating domain user",
		zap.String("email", string(req.Email)),
		zap.String("name", string(req.Name)),
		zap.String("hashed_password_length", fmt.Sprintf("%d", len(hashedPassword))))
	user := domain.NewUser(req.Email, domain.Password(hashedPassword), req.Name)

	// ドメインモデルのバリデーション（パスワード以外）
	s.logger.Info("Validating email")
	if err := domain.ValidateEmail(user.Email); err != nil {
		s.logger.Error("Email validation failed", zap.Error(err))
		return nil, fmt.Errorf("email validation failed: %w", err)
	}
	s.logger.Info("Validating name")
	if err := domain.ValidateName(user.Name); err != nil {
		s.logger.Error("Name validation failed", zap.Error(err))
		return nil, fmt.Errorf("name validation failed: %w", err)
	}

	// 重複チェック
	if existing, _ := s.repo.GetByEmail(ctx, req.Email); existing != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// リポジトリで保存
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponse, error) {
	if id == uuid.Nil {
		return nil, domain.ErrInvalidID
	}
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email domain.Email) (*UserResponse, error) {
	if email == "" {
		return nil, domain.ErrInvalidEmail
	}
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

type UpdateUserRequest struct {
	ID    uuid.UUID    `json:"id"`
	Email domain.Email `json:"email"`
	Name  domain.Name  `json:"name"`
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, req UpdateUserRequest) error {
	// IDのバリデーション
	if req.ID == uuid.Nil {
		return domain.ErrInvalidID
	}

	if req.Email == "" && req.Name == "" {
		return domain.ErrInvalidUpdateInput
	}

	// 既存のユーザーを取得
	user, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if req.Email != "" {
		if err := user.UpdateEmail(req.Email); err != nil {
			return err
		}
	}
	if req.Name != "" {
		if err := user.UpdateName(req.Name); err != nil {
			return err
		}
	}

	if err := user.Validate(); err != nil {
		return err
	}

	// リポジトリで保存
	return s.repo.Update(ctx, user)
}

type DeleteUserRequest struct {
	ID uuid.UUID `json:"id"`
}

// DeleteUser deletes a user by ID
func (s *userService) DeleteUser(ctx context.Context, req DeleteUserRequest) error {
	if req.ID == uuid.Nil {
		return domain.ErrInvalidID
	}

	if err := s.repo.Delete(ctx, req.ID); err != nil {
		return err
	}

	return nil
}

type ListUsersRequest struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

// ListUsers retrieves a paginated list of users
func (s *userService) ListUsers(ctx context.Context, req ListUsersRequest) ([]*UserResponse, error) {
	if req.Limit <= 0 {
		return nil, domain.ErrInvalidLimit
	}
	if req.Offset < 0 {
		return nil, domain.ErrInvalidOffset
	}

	users, err := s.repo.ListUsers(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	return toUserResponses(users), nil
}

func toUserResponses(users []*domain.User) []*UserResponse {
	responses := make([]*UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, &UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
	return responses
}

type AuthenticateUserRequest struct {
	Email    domain.Email    `json:"email"`
	Password domain.Password `json:"password"`
}

func (s *userService) AuthenticateUser(ctx context.Context, req AuthenticateUserRequest) error {
	if req.Email == "" {
		return domain.ErrInvalidEmail
	}
	if req.Password == "" {
		return domain.ErrInvalidPassword
	}

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if !s.hasher.Compare(user.Password, string(req.Password)) {
		return fmt.Errorf("invalid credentials")
	}

	return nil
}
