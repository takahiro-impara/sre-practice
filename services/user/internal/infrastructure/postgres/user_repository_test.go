package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	postgresDriver "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/infrastructure/postgres"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

// テストスイート
type UserRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	repo repository.UserRepository
}

// テストスイートのセットアップ
func (suite *UserRepositoryTestSuite) SetupSuite() {
	// テスト用データベース接続
	// 環境変数からDSNを取得、デフォルトはテスト用DB
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		// デフォルトはテスト用のデータベース
		dsn = "postgres://app:app@localhost:5432/testdb?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	suite.db = db
	suite.repo = postgres.NewUserRepository(db)

	// マイグレーションを実行
	suite.setupTestTable()
}

// 各テストの前処理
func (suite *UserRepositoryTestSuite) SetupTest() {
	// テーブルをクリーンアップ
	suite.cleanupTestData()
}

// 各テストの後処理
func (suite *UserRepositoryTestSuite) TearDownTest() {
	suite.cleanupTestData()
}

// テストスイートの終了処理
func (suite *UserRepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
}

// テスト用テーブルのセットアップ（マイグレーションファイルを使用）
func (suite *UserRepositoryTestSuite) setupTestTable() {
	// プロジェクトルートからマイグレーションディレクトリへのパスを構築
	wd, err := os.Getwd()
	require.NoError(suite.T(), err)

	// プロジェクトルートまで移動（infrastructure/postgres/testからプロジェクトルートまで）
	projectRoot := filepath.Join(wd, "../../../../..")
	migrationPath := filepath.Join(projectRoot, "services/user/db/migrations")

	// file://プロトコルでマイグレーションディレクトリを指定
	sourceURL := fmt.Sprintf("file://%s", migrationPath)

	// マイグレーション用のドライバーを作成
	driver, err := postgresDriver.WithInstance(suite.db, &postgresDriver.Config{})
	require.NoError(suite.T(), err)

	// マイグレーションインスタンスを作成
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres", driver)
	if err != nil {
		// パスが見つからない場合はログを出力
		log.Printf("Failed to create migration instance. Migration path: %s, Error: %v", migrationPath, err)
		suite.T().Fatalf("Failed to create migration instance: %v", err)
	}

	// すべてのマイグレーションを適用
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		suite.T().Fatalf("Failed to run migrations: %v", err)
	}
}

// テストデータのクリーンアップ
func (suite *UserRepositoryTestSuite) cleanupTestData() {
	_, err := suite.db.Exec("DELETE FROM users")
	require.NoError(suite.T(), err)
}

// テスト: Create
func (suite *UserRepositoryTestSuite) TestCreate() {
	testCases := []struct {
		name    string
		user    *domain.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常系:ユーザー作成成功",
			user: &domain.User{
				Email:    domain.Email("test@example.com"),
				Password: domain.Password("testPass123"),
				Name:     domain.Name("Test User"),
			},
			wantErr: false,
		},
		{
			name: "正常系:日本語名のユーザー",
			user: &domain.User{
				Email:    domain.Email("yamada@example.jp"),
				Password: domain.Password("yamadaPass456"),
				Name:     domain.Name("山田太郎"),
			},
			wantErr: false,
		},
		{
			name: "異常系:重複したEmail",
			user: &domain.User{
				Email:    domain.Email("duplicate@example.com"),
				Password: domain.Password("dupPass789"),
				Name:     domain.Name("Duplicate User"),
			},
			wantErr: true,
			errMsg:  "email already exists",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.name == "異常系:重複したEmail" {
				// 事前に同じEmailのユーザーを作成
				firstUser := &domain.User{
					Email:    tc.user.Email,
					Password: domain.Password("firstPass123"),
					Name:     domain.Name("First User"),
				}
				err := suite.repo.Create(context.Background(), firstUser)
				require.NoError(suite.T(), err)
			}

			// テスト実行
			err := suite.repo.Create(context.Background(), tc.user)

			if tc.wantErr {
				assert.Error(suite.T(), err)
				if tc.errMsg != "" {
					assert.Contains(suite.T(), err.Error(), tc.errMsg)
				}
			} else {
				assert.NoError(suite.T(), err)
				assert.NotEmpty(suite.T(), tc.user.ID)
				assert.NotZero(suite.T(), tc.user.CreatedAt)
				assert.NotZero(suite.T(), tc.user.UpdatedAt)
			}
		})
	}
}

// テスト: GetByID
func (suite *UserRepositoryTestSuite) TestGetByID() {
	// テストデータを事前作成
	testUser := &domain.User{
		Email:    domain.Email("get@example.com"),
		Password: domain.Password("getPass123"),
		Name:     domain.Name("Get Test User"),
	}
	err := suite.repo.Create(context.Background(), testUser)
	require.NoError(suite.T(), err)

	testCases := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
		errMsg  string
	}{
		{
			name:    "正常系:存在するユーザーを取得",
			id:      testUser.ID,
			wantErr: false,
		},
		{
			name:    "異常系:存在しないユーザー",
			id:      uuid.New(),
			wantErr: true,
			errMsg:  "not found",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// テスト実行
			user, err := suite.repo.GetByID(context.Background(), tc.id)

			if tc.wantErr {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), user)
				if tc.errMsg != "" {
					assert.Contains(suite.T(), err.Error(), tc.errMsg)
				}
			} else {
				assert.NoError(suite.T(), err)
				require.NotNil(suite.T(), user)
				assert.Equal(suite.T(), testUser.ID, user.ID)
				assert.Equal(suite.T(), testUser.Email, user.Email)
				assert.Equal(suite.T(), testUser.Name, user.Name)
			}
		})
	}
}

// テスト: ListUsers
func (suite *UserRepositoryTestSuite) TestListUsers() {
	// テストデータを複数作成
	users := []*domain.User{
		{Email: domain.Email("user1@example.com"), Password: domain.Password("pass1"), Name: domain.Name("User 1")},
		{Email: domain.Email("user2@example.com"), Password: domain.Password("pass2"), Name: domain.Name("User 2")},
		{Email: domain.Email("user3@example.com"), Password: domain.Password("pass3"), Name: domain.Name("User 3")},
		{Email: domain.Email("user4@example.com"), Password: domain.Password("pass4"), Name: domain.Name("User 4")},
		{Email: domain.Email("user5@example.com"), Password: domain.Password("pass5"), Name: domain.Name("User 5")},
	}

	for _, u := range users {
		err := suite.repo.Create(context.Background(), u)
		require.NoError(suite.T(), err)
		time.Sleep(10 * time.Millisecond) // 順序を保証するため
	}

	testCases := []struct {
		name       string
		limit      int32
		offset     int32
		wantCount  int
		checkOrder bool
	}{
		{
			name:       "正常系:全件取得(limit=10)",
			limit:      10,
			offset:     0,
			wantCount:  5,
			checkOrder: true,
		},
		{
			name:      "正常系:ページネーション(limit=2, offset=0)",
			limit:     2,
			offset:    0,
			wantCount: 2,
		},
		{
			name:      "正常系:ページネーション(limit=2, offset=2)",
			limit:     2,
			offset:    2,
			wantCount: 2,
		},
		{
			name:      "正常系:最終ページ(limit=2, offset=4)",
			limit:     2,
			offset:    4,
			wantCount: 1,
		},
		{
			name:      "境界値:offset超過",
			limit:     10,
			offset:    100,
			wantCount: 0,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// テスト実行
			result, err := suite.repo.ListUsers(context.Background(), tc.limit, tc.offset)

			assert.NoError(suite.T(), err)
			assert.Len(suite.T(), result, tc.wantCount)

			if tc.checkOrder && len(result) > 1 {
				// created_at DESCの順序確認
				for i := 0; i < len(result)-1; i++ {
					assert.True(suite.T(), result[i].CreatedAt.After(result[i+1].CreatedAt) ||
						result[i].CreatedAt.Equal(result[i+1].CreatedAt))
				}
			}
		})
	}
}

// テスト: Update
func (suite *UserRepositoryTestSuite) TestUpdate() {
	// テストデータを事前作成
	originalUser := &domain.User{
		Email:    domain.Email("original@example.com"),
		Password: domain.Password("originalPass"),
		Name:     domain.Name("Original Name"),
	}
	err := suite.repo.Create(context.Background(), originalUser)
	require.NoError(suite.T(), err)

	testCases := []struct {
		name    string
		update  func() *domain.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常系:名前の更新",
			update: func() *domain.User {
				user := *originalUser
				user.Name = domain.Name("Updated Name")
				return &user
			},
			wantErr: false,
		},
		{
			name: "正常系:メールアドレスの更新",
			update: func() *domain.User {
				user := *originalUser
				user.Email = domain.Email("updated@example.com")
				return &user
			},
			wantErr: false,
		},
		{
			name: "異常系:存在しないユーザーの更新",
			update: func() *domain.User {
				return &domain.User{
					ID:       uuid.New(),
					Email:    domain.Email("nonexistent@example.com"),
					Password: domain.Password("nonPass"),
					Name:     domain.Name("Non Existent"),
				}
			},
			wantErr: true,
			errMsg:  "not found",
		},
		// このテストケースは最初に移動するか、特別な処理が必要
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			updateUser := tc.update()
			originalUpdatedAt := updateUser.UpdatedAt

			// テスト実行
			err := suite.repo.Update(context.Background(), updateUser)

			if tc.wantErr {
				assert.Error(suite.T(), err)
				if tc.errMsg != "" {
					assert.Contains(suite.T(), err.Error(), tc.errMsg)
				}
			} else {
				assert.NoError(suite.T(), err)

				// 更新されたデータを確認
				updated, err := suite.repo.GetByID(context.Background(), updateUser.ID)
				require.NoError(suite.T(), err)
				assert.Equal(suite.T(), updateUser.Email, updated.Email)
				assert.Equal(suite.T(), updateUser.Name, updated.Name)
				assert.True(suite.T(), updated.UpdatedAt.After(originalUpdatedAt))
			}

			// クリーンアップ
			suite.cleanupTestData()
			// 次のテストケースのために再作成（存在しないユーザー以外）
			if tc.name != "異常系:存在しないユーザーの更新" {
				// originalUserを初期状態に戻して再作成
				originalUser.Email = domain.Email("original@example.com")
				originalUser.Password = domain.Password("originalPass")
				originalUser.Name = domain.Name("Original Name")
				err = suite.repo.Create(context.Background(), originalUser)
				require.NoError(suite.T(), err)
			}
		})
	}
}

// テスト: Delete
func (suite *UserRepositoryTestSuite) TestDelete() {
	// テストデータを事前作成
	testUser := &domain.User{
		Email:    domain.Email("delete@example.com"),
		Password: domain.Password("deletePass"),
		Name:     domain.Name("Delete Test User"),
	}
	err := suite.repo.Create(context.Background(), testUser)
	require.NoError(suite.T(), err)

	testCases := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
		errMsg  string
		setup   func()
	}{
		{
			name:    "正常系:ユーザーの削除",
			id:      testUser.ID,
			wantErr: false,
		},
		{
			name:    "正常系:存在しないユーザーの削除(エラーなし)",
			id:      uuid.New(),
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.setup != nil {
				tc.setup()
			}

			// テスト実行
			err := suite.repo.Delete(context.Background(), tc.id)

			if tc.wantErr {
				assert.Error(suite.T(), err)
				if tc.errMsg != "" {
					assert.Contains(suite.T(), err.Error(), tc.errMsg)
				}
			} else {
				assert.NoError(suite.T(), err)

				// 削除確認
				if tc.name == "正常系:ユーザーの削除" {
					_, err := suite.repo.GetByID(context.Background(), tc.id)
					assert.Error(suite.T(), err)
					assert.Contains(suite.T(), err.Error(), "not found")
				}
			}
		})
	}
}

// トランザクションのテスト
func (suite *UserRepositoryTestSuite) TestTransaction() {
	suite.Run("トランザクション内での複数操作", func() {
		tx, err := suite.db.Begin()
		require.NoError(suite.T(), err)
		defer tx.Rollback()

		// トランザクション用のリポジトリ(実装が必要)
		// txRepo := repository.NewUserRepositoryWithTx(tx)

		// ここでトランザクション内での操作をテスト
		// 例: 複数ユーザーの作成、更新、削除

		err = tx.Commit()
		assert.NoError(suite.T(), err)
	})
}

// ベンチマークテスト
func BenchmarkUserRepository_Create(b *testing.B) {
	// セットアップ
	dsn := "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		b.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewUserRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &domain.User{
			Email:    domain.Email(fmt.Sprintf("bench%d@example.com", i)),
			Password: domain.Password("benchPass"),
			Name:     domain.Name(fmt.Sprintf("Bench User %d", i)),
		}
		_ = repo.Create(context.Background(), user)
	}
}

func BenchmarkUserRepository_GetByID(b *testing.B) {
	// セットアップ
	dsn := "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		b.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewUserRepository(db)

	// テストユーザー作成
	user := &domain.User{
		Email:    domain.Email("bench@example.com"),
		Password: domain.Password("benchPass"),
		Name:     domain.Name("Bench User"),
	}
	err = repo.Create(context.Background(), user)
	if err != nil {
		b.Fatalf("Failed to create test user: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetByID(context.Background(), user.ID)
	}
}

// テストスイートの実行
func TestUserRepositoryTestSuite(t *testing.T) {
	// 統合テストのスキップフラグ
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	suite.Run(t, new(UserRepositoryTestSuite))
}

// モックリポジトリのテスト
type MockUserRepository struct {
	CreateFunc    func(*domain.User) error
	GetByIDFunc   func(uuid.UUID) (*domain.User, error)
	ListUsersFunc func(int32, int32) ([]*domain.User, error)
	UpdateFunc    func(*domain.User) error
	DeleteFunc    func(uuid.UUID) error
}

func (m *MockUserRepository) Create(user *domain.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockUserRepository) ListUsers(limit, offset int32) ([]*domain.User, error) {
	if m.ListUsersFunc != nil {
		return m.ListUsersFunc(limit, offset)
	}
	return nil, nil
}

func (m *MockUserRepository) Update(user *domain.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return nil
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// モックを使ったユニットテストの例
func TestMockRepository(t *testing.T) {
	t.Run("モック:Create成功", func(t *testing.T) {
		mock := &MockUserRepository{
			CreateFunc: func(user *domain.User) error {
				user.ID = uuid.New()
				user.CreatedAt = time.Now()
				user.UpdatedAt = time.Now()
				return nil
			},
		}

		user := &domain.User{
			Email:    domain.Email("mock@example.com"),
			Password: domain.Password("mockPass"),
			Name:     domain.Name("Mock User"),
		}

		err := mock.Create(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID)
	})

	t.Run("モック:GetByID NotFound", func(t *testing.T) {
		mock := &MockUserRepository{
			GetByIDFunc: func(id uuid.UUID) (*domain.User, error) {
				return nil, fmt.Errorf("not found")
			},
		}

		_, err := mock.GetByID(uuid.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
