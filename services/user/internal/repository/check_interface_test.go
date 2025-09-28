package repository_test

import (
	"testing"

	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/repository"
)

func TestInterfaceImplementation(t *testing.T) {
	// ここでGoコンパイラがチェック！
	// MockUserRepositoryがUserRepositoryインターフェースを満たしているか

	// パターン1: 変数への代入でチェック
	var repo repository.UserRepository
	repo = &repository.MockUserRepository{} // ← コンパイラがここでチェック！

	// もしMockUserRepositoryに必要なメソッドがなければ
	// コンパイルエラー:
	// "cannot use &repository.MockUserRepository{} as repository.UserRepository"

	// パターン2: 関数の引数でチェック
	processRepository(repo)

	// パターン3: 型アサーションでチェック
	mock := &repository.MockUserRepository{}
	_, ok := interface{}(mock).(repository.UserRepository)
	if !ok {
		t.Fatal("MockUserRepository does not implement UserRepository")
	}

	t.Log("✅ MockUserRepositoryはUserRepositoryインターフェースを満たしています！")
}

func processRepository(repo repository.UserRepository) {
	// UserRepository型として使用
	// ここに到達できる = インターフェースを満たしている
}