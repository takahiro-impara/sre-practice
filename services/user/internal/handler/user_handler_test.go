package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, req service.CreateUserRequest) (*service.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id uuid.UUID) (*service.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email domain.Email) (*service.UserResponse, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.UserResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, req service.UpdateUserRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, req service.DeleteUserRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockUserService) ListUsers(ctx context.Context, req service.ListUsersRequest) ([]*service.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*service.UserResponse), args.Error(1)
}

func (m *MockUserService) AuthenticateUser(ctx context.Context, req service.AuthenticateUserRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func TestUserHandler_CreateUser(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "成功: 正常なユーザー作成",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"name":     "Test User",
				"password": "Password123",
			},
			mockSetup: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, mock.MatchedBy(func(req service.CreateUserRequest) bool {
					return req.Email == "test@example.com" &&
						req.Name == "Test User" &&
						req.Password == "Password123"
				})).Return(&service.UserResponse{
					ID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Email:     "test@example.com",
					Name:      "Test User",
					CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "失敗: 無効なJSON",
			requestBody:    "invalid json",
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "失敗: メールアドレス重複",
			requestBody: map[string]string{
				"email":    "existing@example.com",
				"name":     "Test User",
				"password": "Password123",
			},
			mockSetup: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, mock.Anything).
					Return(nil, domain.ErrUserAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "失敗: 無効なメールアドレス",
			requestBody: map[string]string{
				"email":    "invalid-email",
				"name":     "Test User",
				"password": "Password123",
			},
			mockSetup: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, mock.Anything).
					Return(nil, domain.ErrInvalidEmail)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "失敗: サーバーエラー",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"name":     "Test User",
				"password": "Password123",
			},
			mockSetup: func(m *MockUserService) {
				m.On("CreateUser", mock.Anything, mock.Anything).
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSvc := new(MockUserService)
			tt.mockSetup(mockSvc)

			handler := NewUserHandler(mockSvc, logger)

			// Create request
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			// Execute
			handler.CreateUser(rec, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetUserByID(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockUserService)
		expectedStatus int
	}{
		{
			name:   "成功: ユーザー取得",
			userID: userID.String(),
			mockSetup: func(m *MockUserService) {
				m.On("GetUserByID", mock.Anything, userID).
					Return(&service.UserResponse{
						ID:        userID,
						Email:     "test@example.com",
						Name:      "Test User",
						CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "失敗: 無効なUUID",
			userID:         "invalid-uuid",
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "失敗: ユーザーが見つからない",
			userID: userID.String(),
			mockSetup: func(m *MockUserService) {
				m.On("GetUserByID", mock.Anything, userID).
					Return(nil, domain.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSvc := new(MockUserService)
			tt.mockSetup(mockSvc)

			handler := NewUserHandler(mockSvc, logger)

			// Create request with chi context
			req := httptest.NewRequest("GET", "/api/v1/users/"+tt.userID, nil)
			rec := httptest.NewRecorder()

			// Set up chi context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Execute
			handler.GetUserByID(rec, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name           string
		userID         string
		requestBody    interface{}
		mockSetup      func(*MockUserService)
		expectedStatus int
	}{
		{
			name:   "成功: 名前のみ更新",
			userID: userID.String(),
			requestBody: map[string]interface{}{
				"name": "Updated Name",
			},
			mockSetup: func(m *MockUserService) {
				// UpdateUser is called
				m.On("UpdateUser", mock.Anything, mock.MatchedBy(func(req service.UpdateUserRequest) bool {
					return req.Name == "Updated Name"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "成功: メールアドレスのみ更新",
			userID: userID.String(),
			requestBody: map[string]interface{}{
				"email": "newemail@example.com",
			},
			mockSetup: func(m *MockUserService) {
				m.On("UpdateUser", mock.Anything, mock.MatchedBy(func(req service.UpdateUserRequest) bool {
					return req.Email == "newemail@example.com"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "失敗: 無効なJSON",
			userID:         userID.String(),
			requestBody:    "invalid json",
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "失敗: ユーザーが見つからない",
			userID: userID.String(),
			requestBody: map[string]interface{}{
				"name": "Updated Name",
			},
			mockSetup: func(m *MockUserService) {
				m.On("UpdateUser", mock.Anything, mock.Anything).
					Return(domain.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSvc := new(MockUserService)
			tt.mockSetup(mockSvc)

			handler := NewUserHandler(mockSvc, logger)

			// Create request
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("PUT", "/api/v1/users/"+tt.userID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Set up chi context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Execute
			handler.UpdateUser(rec, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockUserService)
		expectedStatus int
	}{
		{
			name:   "成功: ユーザー削除",
			userID: userID.String(),
			mockSetup: func(m *MockUserService) {
				m.On("DeleteUser", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "失敗: 無効なUUID",
			userID:         "invalid-uuid",
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "失敗: ユーザーが見つからない",
			userID: userID.String(),
			mockSetup: func(m *MockUserService) {
				m.On("DeleteUser", mock.Anything, mock.Anything).
					Return(domain.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSvc := new(MockUserService)
			tt.mockSetup(mockSvc)

			handler := NewUserHandler(mockSvc, logger)

			// Create request
			req := httptest.NewRequest("DELETE", "/api/v1/users/"+tt.userID, nil)
			rec := httptest.NewRecorder()

			// Set up chi context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Execute
			handler.DeleteUser(rec, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_ListUsers(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockUserService)
		expectedStatus int
		validateBody   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "成功: デフォルトパラメータ",
			queryParams: "",
			mockSetup: func(m *MockUserService) {
				m.On("ListUsers", mock.Anything, service.ListUsersRequest{
					Limit:  10,
					Offset: 0,
				}).Return([]*service.UserResponse{
					{
						ID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Email:     "user1@example.com",
						Name:      "User 1",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
						Email:     "user2@example.com",
						Name:      "User 2",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp ListUsersResponse
				err := json.NewDecoder(rec.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, 2, resp.TotalCount)
				assert.Equal(t, 10, resp.Limit)
				assert.Equal(t, 0, resp.Offset)
			},
		},
		{
			name:        "成功: カスタムパラメータ",
			queryParams: "?limit=20&offset=10",
			mockSetup: func(m *MockUserService) {
				m.On("ListUsers", mock.Anything, service.ListUsersRequest{
					Limit:  20,
					Offset: 10,
				}).Return([]*service.UserResponse{}, nil)
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp ListUsersResponse
				err := json.NewDecoder(rec.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, 0, resp.TotalCount)
				assert.Equal(t, 20, resp.Limit)
				assert.Equal(t, 10, resp.Offset)
			},
		},
		{
			name:        "失敗: 無効なlimit",
			queryParams: "?limit=invalid",
			mockSetup: func(m *MockUserService) {
				// ListUsers is still called with default values because we handle invalid params
				m.On("ListUsers", mock.Anything, service.ListUsersRequest{
					Limit:  10,
					Offset: 0,
				}).Return([]*service.UserResponse{}, nil)
			},
			expectedStatus: http.StatusOK, // Handler uses default values for invalid params
		},
		{
			name:        "失敗: サービスエラー",
			queryParams: "",
			mockSetup: func(m *MockUserService) {
				m.On("ListUsers", mock.Anything, mock.Anything).
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSvc := new(MockUserService)
			tt.mockSetup(mockSvc)

			handler := NewUserHandler(mockSvc, logger)

			// Create request
			req := httptest.NewRequest("GET", "/api/v1/users"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			// Execute
			handler.ListUsers(rec, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.validateBody != nil {
				tt.validateBody(t, rec)
			}
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_AuthenticateUser(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockUserService)
		expectedStatus int
	}{
		{
			name: "成功: 認証成功",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "Password123",
			},
			mockSetup: func(m *MockUserService) {
				m.On("AuthenticateUser", mock.Anything, service.AuthenticateUserRequest{
					Email:    "test@example.com",
					Password: "Password123",
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "失敗: 認証失敗",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "WrongPassword",
			},
			mockSetup: func(m *MockUserService) {
				m.On("AuthenticateUser", mock.Anything, mock.Anything).
					Return(errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "失敗: 無効なJSON",
			requestBody:    "invalid json",
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "失敗: メールアドレスがない",
			requestBody: map[string]string{
				"password": "Password123",
			},
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "失敗: パスワードがない",
			requestBody: map[string]string{
				"email": "test@example.com",
			},
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSvc := new(MockUserService)
			tt.mockSetup(mockSvc)

			handler := NewUserHandler(mockSvc, logger)

			// Create request
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/api/v1/users/authenticate", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Execute
			handler.AuthenticateUser(rec, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_HealthCheck(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockSvc := new(MockUserService)
	handler := NewUserHandler(mockSvc, logger)

	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.HealthCheck(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestUserHandler_ReadinessCheck(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockSvc := new(MockUserService)
	handler := NewUserHandler(mockSvc, logger)

	req := httptest.NewRequest("GET", "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.ReadinessCheck(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Ready", rec.Body.String())
}
