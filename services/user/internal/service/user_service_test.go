package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/repository"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// MockPasswordHasher is a mock implementation of PasswordHasher
type MockPasswordHasher struct {
	mock.Mock
}

// createTestLogger creates a test logger
func createTestLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func (m *MockPasswordHasher) Hash(password domain.Password) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Compare(hashedPassword domain.Password, plainPassword string) bool {
	args := m.Called(hashedPassword, plainPassword)
	return args.Bool(0)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	// 1. モックリポジトリを作成
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)

	// 2. サービスを作成（モックを注入）
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	// 3. モックの期待値を設定
	// パスワードハッシュ化
	mockHasher.On("Hash", domain.Password("testPass123")).
		Return("hashed_password_123", nil).Once()

	// GetByEmail が nil を返す（重複なし）
	mockRepo.On("GetByEmail",
		mock.Anything,
		domain.Email("test@example.com"),
	).Return(nil, errors.New("not found")).Once()

	// Create メソッドが呼ばれたら nil を返す（成功）
	mockRepo.On("Create",
		mock.Anything,                       // context
		mock.AnythingOfType("*domain.User"), // user引数
	).Return(nil).Once() // 1回だけ呼ばれることを期待

	// 4. テスト実行
	ctx := context.Background()
	req := service.CreateUserRequest{
		Email:    domain.Email("test@example.com"),
		Name:     domain.Name("Test User"),
		Password: domain.Password("testPass123"),
	}
	user, err := svc.CreateUser(ctx, req)

	// 5. 結果の検証
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, domain.Email("test@example.com"), user.Email)
	assert.Equal(t, domain.Name("Test User"), user.Name)
	assert.NotEmpty(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)

	// 6. モックが期待通りに呼ばれたか検証
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestUserService_CreateUser_RepositoryError(t *testing.T) {
	// 1. モックリポジトリを作成
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)

	// 2. サービスを作成
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	// 3. パスワードハッシュ化
	mockHasher.On("Hash", domain.Password("testPass123")).
		Return("hashed_password_123", nil).Once()

	// GetByEmail が nil を返す（重複なし）
	mockRepo.On("GetByEmail",
		mock.Anything,
		domain.Email("test@example.com"),
	).Return(nil, errors.New("not found")).Once()

	// モックがエラーを返すように設定
	expectedErr := errors.New("database error")
	mockRepo.On("Create",
		mock.Anything,
		mock.AnythingOfType("*domain.User"),
	).Return(expectedErr).Once()

	// 4. テスト実行
	ctx := context.Background()
	req := service.CreateUserRequest{
		Email:    domain.Email("test@example.com"),
		Name:     domain.Name("Test User"),
		Password: domain.Password("testPass123"),
	}
	user, err := svc.CreateUser(ctx, req)

	// 5. エラーが返されることを確認
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, expectedErr, err)

	// 6. モックの検証
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	// 1. モックリポジトリを作成
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)

	// 2. サービスを作成
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	// 3. 期待する返り値を準備
	expectedUser := &domain.User{
		ID:        uuid.New(),
		Email:     domain.Email("test@example.com"),
		Password:  domain.Password("testPass123"),
		Name:      domain.Name("Test User"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 4. モックの設定
	mockRepo.On("GetByID",
		mock.Anything,   // context
		expectedUser.ID, // 特定のIDで呼ばれることを期待
	).Return(expectedUser, nil).Once()

	// 5. テスト実行
	ctx := context.Background()
	user, err := svc.GetUserByID(ctx, expectedUser.ID)

	// 6. 結果の検証
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.CreatedAt, user.CreatedAt)
	assert.Equal(t, expectedUser.UpdatedAt, user.UpdatedAt)

	// 7. モックの検証
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	// 1. モックリポジトリを作成
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)

	// 2. サービスを作成
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	// 3. 存在しないユーザーID
	notFoundID := uuid.New()
	notFoundErr := errors.New("user not found")

	// 4. モックの設定（nilとエラーを返す）
	mockRepo.On("GetByID",
		mock.Anything,
		notFoundID,
	).Return(nil, notFoundErr).Once()

	// 5. テスト実行
	ctx := context.Background()
	user, err := svc.GetUserByID(ctx, notFoundID)

	// 6. 結果の検証
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, notFoundErr, err)

	// 7. モックの検証
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_InvalidID(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	ctx := context.Background()
	user, err := svc.GetUserByID(ctx, uuid.Nil)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, domain.ErrInvalidID, err)

	// GetByIDが呼ばれていないことを確認
	mockRepo.AssertNotCalled(t, "GetByID")
}

// ========== テーブルドリブンテストの例 ==========

func TestUserService_CreateUser_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		req       service.CreateUserRequest
		mockSetup func(*repository.MockUserRepository)
		wantErr   bool
		errEquals error
	}{
		{
			name: "正常系：ユーザー作成成功",
			req: service.CreateUserRequest{
				Email:    domain.Email("test@example.com"),
				Name:     domain.Name("Test User"),
				Password: domain.Password("testPass123"),
			},
			mockSetup: func(m *repository.MockUserRepository) {
				m.On("GetByEmail", mock.Anything, domain.Email("test@example.com")).
					Return(nil, errors.New("not found")).Once()
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
					Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "異常系：メールアドレス重複",
			req: service.CreateUserRequest{
				Email:    domain.Email("duplicate@example.com"),
				Name:     domain.Name("Duplicate User"),
				Password: domain.Password("dupPass456"),
			},
			mockSetup: func(m *repository.MockUserRepository) {
				existingUser := &domain.User{
					ID:    uuid.New(),
					Email: domain.Email("duplicate@example.com"),
				}
				m.On("GetByEmail", mock.Anything, domain.Email("duplicate@example.com")).
					Return(existingUser, nil).Once()
			},
			wantErr:   true,
			errEquals: domain.ErrUserAlreadyExists,
		},
		{
			name: "異常系：空のメールアドレス",
			req: service.CreateUserRequest{
				Email:    domain.Email(""),
				Name:     domain.Name("Test User"),
				Password: domain.Password("testPass123"),
			},
			mockSetup: func(m *repository.MockUserRepository) {},
			wantErr:   true,
			errEquals: domain.ErrInvalidEmail,
		},
		{
			name: "異常系：空の名前",
			req: service.CreateUserRequest{
				Email:    domain.Email("test@example.com"),
				Name:     domain.Name(""),
				Password: domain.Password("testPass123"),
			},
			mockSetup: func(m *repository.MockUserRepository) {},
			wantErr:   true,
			errEquals: domain.ErrInvalidName,
		},
		{
			name: "異常系：空のパスワード",
			req: service.CreateUserRequest{
				Email:    domain.Email("test@example.com"),
				Name:     domain.Name("Test User"),
				Password: domain.Password(""),
			},
			mockSetup: func(m *repository.MockUserRepository) {},
			wantErr:   true,
			errEquals: domain.ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックをセットアップ
			mockRepo := new(repository.MockUserRepository)
			mockHasher := new(MockPasswordHasher)
			tt.mockSetup(mockRepo)

			// パスワードハッシュ化のモック（空パスワード以外）
			if tt.req.Password != "" {
				mockHasher.On("Hash", tt.req.Password).
					Return("hashed_"+string(tt.req.Password), nil).Maybe()
			}

			// サービスを作成
			svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

			// テスト実行
			ctx := context.Background()
			user, err := svc.CreateUser(ctx, tt.req)

			// 検証
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errEquals != nil {
					assert.Equal(t, tt.errEquals, err)
				}
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}

			// モックの期待値を検証
			mockRepo.AssertExpectations(t)
			mockHasher.AssertExpectations(t)
		})
	}
}

// ========== UpdateUser テスト ==========

func TestUserService_UpdateUser_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	existingUser := &domain.User{
		ID:        uuid.New(),
		Email:     domain.Email("old@example.com"),
		Name:      domain.Name("Old Name"),
		Password:  domain.Password("hashedPassword"),
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	mockRepo.On("GetByID",
		mock.Anything,
		existingUser.ID,
	).Return(existingUser, nil).Once()

	mockRepo.On("Update",
		mock.Anything,
		mock.AnythingOfType("*domain.User"),
	).Return(nil).Once()

	ctx := context.Background()
	req := service.UpdateUserRequest{
		ID:    existingUser.ID,
		Email: domain.Email("new@example.com"),
		Name:  domain.Name("New Name"),
	}
	err := svc.UpdateUser(ctx, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_UserNotFound(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	userID := uuid.New()
	mockRepo.On("GetByID",
		mock.Anything,
		userID,
	).Return(nil, errors.New("user not found")).Once()

	ctx := context.Background()
	req := service.UpdateUserRequest{
		ID:    userID,
		Email: domain.Email("new@example.com"),
		Name:  domain.Name("New Name"),
	}
	err := svc.UpdateUser(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_InvalidInput(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	ctx := context.Background()
	req := service.UpdateUserRequest{
		ID:    uuid.New(),
		Email: domain.Email(""),
		Name:  domain.Name(""),
	}
	err := svc.UpdateUser(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidUpdateInput, err)
	mockRepo.AssertNotCalled(t, "GetByEmail")
}

// ========== DeleteUser テスト ==========

func TestUserService_DeleteUser_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	userID := uuid.New()
	mockRepo.On("Delete",
		mock.Anything,
		userID,
	).Return(nil).Once()

	ctx := context.Background()
	req := service.DeleteUserRequest{ID: userID}
	err := svc.DeleteUser(ctx, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_InvalidID(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	ctx := context.Background()
	req := service.DeleteUserRequest{ID: uuid.Nil}
	err := svc.DeleteUser(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidID, err)
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestUserService_DeleteUser_RepositoryError(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	userID := uuid.New()
	expectedErr := errors.New("database error")
	mockRepo.On("Delete",
		mock.Anything,
		userID,
	).Return(expectedErr).Once()

	ctx := context.Background()
	req := service.DeleteUserRequest{ID: userID}
	err := svc.DeleteUser(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

// ========== ListUsers テスト ==========

func TestUserService_ListUsers_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	mockUsers := []*domain.User{
		{
			ID:        uuid.New(),
			Email:     domain.Email("user1@example.com"),
			Name:      domain.Name("User 1"),
			Password:  domain.Password("hash1"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Email:     domain.Email("user2@example.com"),
			Name:      domain.Name("User 2"),
			Password:  domain.Password("hash2"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockRepo.On("ListUsers",
		mock.Anything,
		int32(10),
		int32(0),
	).Return(mockUsers, nil).Once()

	ctx := context.Background()
	req := service.ListUsersRequest{
		Limit:  10,
		Offset: 0,
	}
	users, err := svc.ListUsers(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, mockUsers[0].Email, users[0].Email)
	assert.Equal(t, mockUsers[1].Email, users[1].Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_ListUsers_InvalidLimit(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	tests := []struct {
		name    string
		limit   int32
		offset  int32
		wantErr error
	}{
		{
			name:    "リミットが0",
			limit:   0,
			offset:  0,
			wantErr: domain.ErrInvalidLimit,
		},
		{
			name:    "リミットが負の値",
			limit:   -1,
			offset:  0,
			wantErr: domain.ErrInvalidLimit,
		},
		{
			name:    "オフセットが負の値",
			limit:   10,
			offset:  -1,
			wantErr: domain.ErrInvalidOffset,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			req := service.ListUsersRequest{
				Limit:  tt.limit,
				Offset: tt.offset,
			}
			users, err := svc.ListUsers(ctx, req)

			assert.Error(t, err)
			assert.Equal(t, tt.wantErr, err)
			assert.Nil(t, users)
			mockRepo.AssertNotCalled(t, "ListUsers")
		})
	}
}

func TestUserService_ListUsers_EmptyResult(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	mockRepo.On("ListUsers",
		mock.Anything,
		int32(10),
		int32(100),
	).Return([]*domain.User{}, nil).Once()

	ctx := context.Background()
	req := service.ListUsersRequest{
		Limit:  10,
		Offset: 100,
	}
	users, err := svc.ListUsers(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 0)
	mockRepo.AssertExpectations(t)
}

// ========== Password Hashing Tests ==========

func TestPasswordHasher_Hash(t *testing.T) {
	hasher := service.NewPasswordHasher(bcrypt.MinCost) // 高速なテスト用

	tests := []struct {
		name     string
		password domain.Password
		wantErr  bool
	}{
		{
			name:     "正常系：通常のパスワード",
			password: domain.Password("testPassword123"),
			wantErr:  false,
		},
		{
			name:     "正常系：長いパスワード",
			password: domain.Password("ThisIsAVeryLongPasswordWith123Numbers!@#"),
			wantErr:  false,
		},
		{
			name:     "正常系：特殊文字を含むパスワード",
			password: domain.Password("P@ssw0rd!#$%^&*()"),
			wantErr:  false,
		},
		{
			name:     "境界値：72バイト（bcryptの最大）",
			password: domain.Password(string(make([]byte, 72, 72))),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := hasher.Hash(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hashed)
				assert.NotEqual(t, string(tt.password), hashed)
				// bcryptハッシュの形式を確認
				assert.True(t, len(hashed) > 50)
			}
		})
	}
}

func TestPasswordHasher_Compare(t *testing.T) {
	hasher := service.NewPasswordHasher(bcrypt.MinCost)

	// 事前にハッシュを生成
	password := "testPassword123"
	hashedPassword, err := hasher.Hash(domain.Password(password))
	assert.NoError(t, err)

	tests := []struct {
		name           string
		hashedPassword domain.Password
		plainPassword  string
		want           bool
	}{
		{
			name:           "正常系：正しいパスワード",
			hashedPassword: domain.Password(hashedPassword),
			plainPassword:  password,
			want:           true,
		},
		{
			name:           "異常系：間違ったパスワード",
			hashedPassword: domain.Password(hashedPassword),
			plainPassword:  "wrongPassword",
			want:           false,
		},
		{
			name:           "異常系：空のパスワード",
			hashedPassword: domain.Password(hashedPassword),
			plainPassword:  "",
			want:           false,
		},
		{
			name:           "異常系：無効なハッシュ",
			hashedPassword: domain.Password("invalid_hash"),
			plainPassword:  password,
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasher.Compare(tt.hashedPassword, tt.plainPassword)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestUserService_CreateUser_WithPasswordHashing(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	// パスワードハッシュ化の期待値設定
	plainPassword := "securePassword123"
	hashedPassword := "$2a$10$hashedPasswordExample"

	mockHasher.On("Hash", domain.Password(plainPassword)).
		Return(hashedPassword, nil).Once()

	mockRepo.On("GetByEmail",
		mock.Anything,
		domain.Email("test@example.com"),
	).Return(nil, errors.New("not found")).Once()

	// ハッシュ化されたパスワードで作成されることを確認
	mockRepo.On("Create",
		mock.Anything,
		mock.MatchedBy(func(user *domain.User) bool {
			return user.Password == domain.Password(hashedPassword)
		}),
	).Return(nil).Once()

	ctx := context.Background()
	req := service.CreateUserRequest{
		Email:    domain.Email("test@example.com"),
		Name:     domain.Name("Test User"),
		Password: domain.Password(plainPassword),
	}

	user, err := svc.CreateUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestUserService_CreateUser_HashingError(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	// ハッシュ化でエラーを返す
	mockHasher.On("Hash", domain.Password("testPass123")).
		Return("", errors.New("hashing failed")).Once()

	ctx := context.Background()
	req := service.CreateUserRequest{
		Email:    domain.Email("test@example.com"),
		Name:     domain.Name("Test User"),
		Password: domain.Password("testPass123"),
	}

	user, err := svc.CreateUser(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to hash password")
	assert.Nil(t, user)
	mockHasher.AssertExpectations(t)
	// GetByEmailが呼ばれていないことを確認
	mockRepo.AssertNotCalled(t, "GetByEmail")
}

// ========== Authentication Tests ==========

func TestUserService_AuthenticateUser_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	hashedPassword := "$2a$10$hashedPasswordExample"
	existingUser := &domain.User{
		ID:       uuid.New(),
		Email:    domain.Email("test@example.com"),
		Password: domain.Password(hashedPassword),
		Name:     domain.Name("Test User"),
	}

	mockRepo.On("GetByEmail",
		mock.Anything,
		domain.Email("test@example.com"),
	).Return(existingUser, nil).Once()

	mockHasher.On("Compare",
		domain.Password(hashedPassword),
		"correctPassword",
	).Return(true).Once()

	ctx := context.Background()
	req := service.AuthenticateUserRequest{
		Email:    domain.Email("test@example.com"),
		Password: domain.Password("correctPassword"),
	}

	err := svc.AuthenticateUser(ctx, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestUserService_AuthenticateUser_InvalidPassword(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	hashedPassword := "$2a$10$hashedPasswordExample"
	existingUser := &domain.User{
		ID:       uuid.New(),
		Email:    domain.Email("test@example.com"),
		Password: domain.Password(hashedPassword),
		Name:     domain.Name("Test User"),
	}

	mockRepo.On("GetByEmail",
		mock.Anything,
		domain.Email("test@example.com"),
	).Return(existingUser, nil).Once()

	mockHasher.On("Compare",
		domain.Password(hashedPassword),
		"wrongPassword",
	).Return(false).Once()

	ctx := context.Background()
	req := service.AuthenticateUserRequest{
		Email:    domain.Email("test@example.com"),
		Password: domain.Password("wrongPassword"),
	}

	err := svc.AuthenticateUser(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestUserService_AuthenticateUser_UserNotFound(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := service.NewUserService(mockRepo, mockHasher, createTestLogger())

	mockRepo.On("GetByEmail",
		mock.Anything,
		domain.Email("nonexistent@example.com"),
	).Return(nil, errors.New("user not found")).Once()

	ctx := context.Background()
	req := service.AuthenticateUserRequest{
		Email:    domain.Email("nonexistent@example.com"),
		Password: domain.Password("password"),
	}

	err := svc.AuthenticateUser(ctx, req)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
	// Compareが呼ばれていないことを確認
	mockHasher.AssertNotCalled(t, "Compare")
}

func TestPasswordHasher_UniquenessOfHashes(t *testing.T) {
	hasher := service.NewPasswordHasher(bcrypt.MinCost)
	password := domain.Password("samePassword123")

	// 同じパスワードでも異なるハッシュが生成されることを確認
	hash1, err1 := hasher.Hash(password)
	hash2, err2 := hasher.Hash(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2) // ソルトにより異なるハッシュ

	// 両方のハッシュが元のパスワードと一致することを確認
	assert.True(t, hasher.Compare(domain.Password(hash1), string(password)))
	assert.True(t, hasher.Compare(domain.Password(hash2), string(password)))
}
