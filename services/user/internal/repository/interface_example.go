package repository

import (
	"context"
	"fmt"
)

// ============================================
// 例1: Javaスタイル（Goでは書けない！）
// ============================================
// class MockUserRepository implements UserRepository { // ❌ Goにはimplementsがない！
//     ...
// }

// ============================================
// 例2: Goの実際の仕組み
// ============================================

// シンプルな例のインターフェース
type SimpleRepository interface {
	Save(ctx context.Context, data string) error
	Find(ctx context.Context, id int) (string, error)
}

// 実装その1：メモリ実装
type MemoryRepository struct {
	storage map[int]string
}

// MemoryRepositoryのメソッド定義
// 注目：どこにも「implements SimpleRepository」とは書いていない！
func (m *MemoryRepository) Save(ctx context.Context, data string) error {
	// メソッドのシグネチャがSimpleRepositoryのSaveと完全一致
	m.storage[len(m.storage)] = data
	return nil
}

func (m *MemoryRepository) Find(ctx context.Context, id int) (string, error) {
	// メソッドのシグネチャがSimpleRepositoryのFindと完全一致
	if data, ok := m.storage[id]; ok {
		return data, nil
	}
	return "", fmt.Errorf("not found")
}

// 実装その2：ファイル実装
type FileRepository struct {
	filePath string
}

func (f *FileRepository) Save(ctx context.Context, data string) error {
	// 同じメソッド名、同じ引数、同じ戻り値
	// これだけでSimpleRepositoryインターフェースを満たす！
	return nil
}

func (f *FileRepository) Find(ctx context.Context, id int) (string, error) {
	// 同じメソッド名、同じ引数、同じ戻り値
	return "file data", nil
}

// ============================================
// Goがインターフェースを判定するタイミング
// ============================================

func UseRepository() {
	// 1. インターフェース型の変数を宣言
	var repo SimpleRepository

	// 2. 具体的な型を代入しようとする
	repo = &MemoryRepository{storage: make(map[int]string)}
	// ↑ ここでGoコンパイラがチェック！
	// MemoryRepositoryはSaveとFindメソッドを持っている？ → YES
	// メソッドのシグネチャは一致している？ → YES
	// じゃあSimpleRepositoryとして使える！ → OK

	// 3. インターフェースとして使用
	_ = repo.Save(context.Background(), "test")

	// 4. 別の実装に切り替えることも可能
	repo = &FileRepository{filePath: "/tmp/data"}
	// ↑ FileRepositoryも必要なメソッドを持っているのでOK
}

// ============================================
// MockUserRepositoryの場合
// ============================================

// MockUserRepositoryがUserRepositoryを満たす仕組み：
//
// UserRepositoryインターフェースが要求するメソッド：
// - Create(ctx context.Context, user *domain.User) error
// - GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
// - GetByEmail(ctx context.Context, email domain.Email) (*domain.User, error)
// - Update(ctx context.Context, user *domain.User) error
// - Delete(ctx context.Context, id uuid.UUID) error
// - ListUsers(ctx context.Context, limit int32, offset int32) ([]*domain.User, error)
//
// MockUserRepositoryが持っているメソッド：
// - Create(ctx context.Context, user *domain.User) error ✅
// - GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) ✅
// - GetByEmail(ctx context.Context, email domain.Email) (*domain.User, error) ✅
// - Update(ctx context.Context, user *domain.User) error ✅
// - Delete(ctx context.Context, id uuid.UUID) error ✅
// - ListUsers(ctx context.Context, limit int32, offset int32) ([]*domain.User, error) ✅
//
// 全て一致！だから MockUserRepository は UserRepository として使える！

// ============================================
// コンパイラのチェックポイント
// ============================================

func ExplainInterfaceCheck() {
	// Goコンパイラは以下をチェック：
	// 1. メソッド名が同じか？
	// 2. レシーバー型が適切か？（ポインタ or 値）
	// 3. 引数の数が同じか？
	// 4. 引数の型が同じか？（順序も重要）
	// 5. 戻り値の数が同じか？
	// 6. 戻り値の型が同じか？

	// 一つでも違うと、インターフェースを満たさない！
}

// ============================================
// なぜこの仕組みが便利か？
// ============================================

// 1. 疎結合：実装が後から追加できる
// 2. テスタビリティ：モックを簡単に作れる
// 3. 柔軟性：既存のコードを変更せずに新しい実装を追加できる
// 4. シンプル：継承階層などの複雑な仕組みが不要