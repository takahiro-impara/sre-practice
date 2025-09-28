package domain_test

import (
	"testing"
	"time"

	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	testCases := []struct {
		name         string
		email        domain.Email
		userName     domain.Name
		password     domain.Password
		wantEmail    domain.Email
		wantName     domain.Name
		wantPassword domain.Password
	}{
		{
			name:         "正常系：一般的なユーザー",
			email:        domain.Email("test@example.com"),
			userName:     domain.Name("Test User"),
			password:     domain.Password("password123"),
			wantEmail:    domain.Email("test@example.com"),
			wantName:     domain.Name("Test User"),
			wantPassword: domain.Password("password123"),
		},
		{
			name:         "正常系：日本語名",
			email:        domain.Email("yamada@example.co.jp"),
			userName:     domain.Name("山田太郎"),
			password:     domain.Password("securePass456"),
			wantEmail:    domain.Email("yamada@example.co.jp"),
			wantName:     domain.Name("山田太郎"),
			wantPassword: domain.Password("securePass456"),
		},
		{
			name:         "正常系：長いメールアドレス",
			email:        domain.Email("very.long.email.address@subdomain.example.com"),
			userName:     domain.Name("Long Email User"),
			password:     domain.Password("LongEmailPass789"),
			wantEmail:    domain.Email("very.long.email.address@subdomain.example.com"),
			wantName:     domain.Name("Long Email User"),
			wantPassword: domain.Password("LongEmailPass789"),
		},
		{
			name:         "境界値：短い名前",
			email:        domain.Email("a@b.c"),
			userName:     domain.Name("A"),
			password:     domain.Password("shortPass1"),
			wantEmail:    domain.Email("a@b.c"),
			wantName:     domain.Name("A"),
			wantPassword: domain.Password("shortPass1"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 実行
			user := domain.NewUser(tc.email, tc.password, tc.userName)

			// 検証
			require.NotNil(t, user)
			assert.NotEmpty(t, user.ID)
			assert.Equal(t, tc.wantEmail, user.Email)
			assert.Equal(t, tc.wantName, user.Name)
			assert.Equal(t, tc.wantPassword, user.Password)
			assert.NotZero(t, user.CreatedAt)
			assert.NotZero(t, user.UpdatedAt)
			// CreatedAtとUpdatedAtは同じ時刻になるはずだが、ナノ秒単位の差は許容
			assert.WithinDuration(t, user.CreatedAt, user.UpdatedAt, time.Microsecond)
		})
	}
}

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		user    *domain.User
		wantErr bool
		errMsg  string
	}{
		{
			name:    "正常系：有効なユーザー",
			user:    domain.NewUser(domain.Email("valid@example.com"), domain.Password("validPass123"), domain.Name("Valid User")),
			wantErr: false,
		},
		{
			name:    "異常系：無効なメールアドレス（@なし）",
			user:    domain.NewUser(domain.Email("invalid-email"), domain.Password("password123"), domain.Name("User")),
			wantErr: true,
			errMsg:  "invalid email",
		},
		{
			name:    "異常系：空のメールアドレス",
			user:    domain.NewUser(domain.Email(""), domain.Password("password123"), domain.Name("User")),
			wantErr: true,
			errMsg:  "invalid email",
		},
		{
			name:    "異常系：空の名前",
			user:    domain.NewUser(domain.Email("test@example.com"), domain.Password("password123"), domain.Name("")),
			wantErr: true,
			errMsg:  "invalid name",
		},
		{
			name:    "異常系：メールアドレスと名前両方が無効",
			user:    domain.NewUser(domain.Email(""), domain.Password("password123"), domain.Name("")),
			wantErr: true,
		},
		{
			name:    "異常系：空のパスワード",
			user:    domain.NewUser(domain.Email("test@example.com"), domain.Password(""), domain.Name("User")),
			wantErr: true,
			errMsg:  "invalid password",
		},
		{
			name: "境界値：非常に長い名前（256文字）",
			user: func() *domain.User {
				longName := ""
				for i := 0; i < 256; i++ {
					longName += "a"
				}
				return domain.NewUser(domain.Email("test@example.com"), domain.Password("password123"), domain.Name(longName))
			}(),
			wantErr: true,
			errMsg:  "invalid name",
		},
		{
			name:    "境界値：短すぎる名前（2文字）",
			user:    domain.NewUser(domain.Email("test@example.com"), domain.Password("password123"), domain.Name("ab")),
			wantErr: true,
			errMsg:  "invalid name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 実行
			err := tc.user.Validate()

			// 検証
			if tc.wantErr {
				assert.Error(t, err)
				if tc.errMsg != "" {
					assert.Contains(t, err.Error(), tc.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_UpdateName(t *testing.T) {
	testCases := []struct {
		name        string
		initialName domain.Name
		newName     domain.Name
		wantName    domain.Name
		expectError bool
	}{
		{
			name:        "正常系：名前の更新",
			initialName: domain.Name("Old Name"),
			newName:     domain.Name("New Name"),
			wantName:    domain.Name("New Name"),
			expectError: false,
		},
		{
			name:        "正常系：日本語から英語",
			initialName: domain.Name("山田太郎"),
			newName:     domain.Name("Taro Yamada"),
			wantName:    domain.Name("Taro Yamada"),
			expectError: false,
		},
		{
			name:        "異常系：空文字への更新",
			initialName: domain.Name("Name"),
			newName:     domain.Name(""),
			wantName:    domain.Name("Name"), // エラー時は元の値が保持される
			expectError: true,
		},
		{
			name:        "境界値：長い名前への更新",
			initialName: domain.Name("Short"),
			newName:     domain.Name("Very Long Name With Many Characters"),
			wantName:    domain.Name("Very Long Name With Many Characters"),
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 準備
			user := domain.NewUser(domain.Email("test@example.com"), domain.Password("password123"), tc.initialName)
			originalUpdatedAt := user.UpdatedAt
			time.Sleep(time.Millisecond) // UpdatedAtの変更を確認するため

			// 実行
			err := user.UpdateName(tc.newName)

			// 検証
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.wantName, user.Name) // エラー時は元の値が保持される
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantName, user.Name)
				assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}

func TestUser_UpdateEmail(t *testing.T) {
	testCases := []struct {
		name         string
		initialEmail domain.Email
		newEmail     domain.Email
		wantEmail    domain.Email
		expectError  bool
	}{
		{
			name:         "正常系：メールアドレスの更新",
			initialEmail: domain.Email("old@example.com"),
			newEmail:     domain.Email("new@example.com"),
			wantEmail:    domain.Email("new@example.com"),
			expectError:  false,
		},
		{
			name:         "正常系：ドメイン変更",
			initialEmail: domain.Email("user@old-domain.com"),
			newEmail:     domain.Email("user@new-domain.com"),
			wantEmail:    domain.Email("user@new-domain.com"),
			expectError:  false,
		},
		{
			name:         "異常系：短いメールアドレスへの更新",
			initialEmail: domain.Email("very.long.email@example.com"),
			newEmail:     domain.Email("a@b.c"),
			wantEmail:    domain.Email("very.long.email@example.com"), // エラー時は元の値が保持される
			expectError:  true,
		},
		{
			name:         "異常系：無効なメールアドレスへの更新",
			initialEmail: domain.Email("valid@example.com"),
			newEmail:     domain.Email("invalid-email"),
			wantEmail:    domain.Email("valid@example.com"), // エラー時は元の値が保持される
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 準備
			user := domain.NewUser(tc.initialEmail, domain.Password("password123"), domain.Name("Test User"))
			originalUpdatedAt := user.UpdatedAt
			time.Sleep(time.Millisecond) // UpdatedAtの変更を確認するため

			// 実行
			err := user.UpdateEmail(tc.newEmail)

			// 検証
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.wantEmail, user.Email) // エラー時は元の値が保持される
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantEmail, user.Email)
				assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}

func TestUser_ComplexScenarios(t *testing.T) {
	testCases := []struct {
		name     string
		scenario func(t *testing.T)
	}{
		{
			name: "複数回の更新",
			scenario: func(t *testing.T) {
				user := domain.NewUser(domain.Email("initial@example.com"), domain.Password("initialPass123"), domain.Name("Initial Name"))

				// 1回目の更新
				err := user.UpdateName(domain.Name("First Update"))
				require.NoError(t, err)
				assert.Equal(t, domain.Name("First Update"), user.Name)

				// 2回目の更新
				err = user.UpdateEmail(domain.Email("second@example.com"))
				require.NoError(t, err)
				assert.Equal(t, domain.Email("second@example.com"), user.Email)

				// 3回目の更新
				err = user.UpdateName(domain.Name("Final Name"))
				require.NoError(t, err)
				assert.Equal(t, domain.Name("Final Name"), user.Name)
			},
		},
		{
			name: "更新後のバリデーション",
			scenario: func(t *testing.T) {
				user := domain.NewUser(domain.Email("valid@example.com"), domain.Password("validPass123"), domain.Name("Valid Name"))

				// 有効な状態の確認
				err := user.Validate()
				require.NoError(t, err)

				// 無効なメールアドレスに更新（エラーが期待される）
				err = user.UpdateEmail(domain.Email("invalid"))
				require.Error(t, err) // 更新は失敗

				// 元の値が保持される
				err = user.Validate()
				assert.NoError(t, err)
			},
		},
		{
			name: "IDの不変性",
			scenario: func(t *testing.T) {
				user := domain.NewUser(domain.Email("test@example.com"), domain.Password("testPass123"), domain.Name("Test"))
				originalID := user.ID

				// 各種更新
				_ = user.UpdateName(domain.Name("New Name"))
				_ = user.UpdateEmail(domain.Email("new@example.com"))

				// IDは変わらない
				assert.Equal(t, originalID, user.ID)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.scenario(t)
		})
	}
}
