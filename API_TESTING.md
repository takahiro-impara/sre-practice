# User Service API Testing Guide

## Postmanコレクションの使用方法

### 1. インポート手順

1. Postmanを開く
2. 左サイドバーの「Collections」タブをクリック
3. 「Import」ボタンをクリック
4. 以下のファイルをドラッグ&ドロップまたは選択：
   - `user-service-api.postman_collection.json` - APIコレクション
   - `user-service.postman_environment.json` - 環境設定

### 2. 環境設定

1. 右上の環境選択ドロップダウンから「User Service Environment」を選択
2. 必要に応じて変数を編集：
   - `base_url`: APIサーバーのURL（デフォルト: http://localhost:8080）
   - `test_email`: テスト用メールアドレス
   - `test_password`: テスト用パスワード

### 3. サーバーの起動

```bash
# データベースを起動
docker-compose up -d

# サーバーを起動
cd services/user
go run cmd/server/main.go
```

### 4. APIのテスト実行

#### 個別テスト
1. コレクション内の任意のリクエストを選択
2. 「Send」ボタンをクリック

#### 自動テストの実行順序
1. **Health Check** - サーバーの稼働確認
2. **Create User** - 新規ユーザー作成（user_idを環境変数に保存）
3. **Get User by ID** - 作成したユーザーの取得
4. **Update User** - ユーザー情報の更新
5. **List Users** - ユーザー一覧の取得
6. **Authenticate User** - ユーザー認証
7. **Delete User** - ユーザーの削除

#### コレクション全体の実行
1. コレクション名の右にある「...」をクリック
2. 「Run collection」を選択
3. 実行順序を確認して「Run User Service API」をクリック

## APIエンドポイント一覧

### ヘルスチェック
- `GET /healthz` - ヘルスチェック
- `GET /readyz` - レディネスチェック

### ユーザー管理
- `POST /api/v1/users` - ユーザー作成
- `GET /api/v1/users/{id}` - ユーザー取得
- `PUT /api/v1/users/{id}` - ユーザー更新
- `DELETE /api/v1/users/{id}` - ユーザー削除
- `GET /api/v1/users` - ユーザー一覧（ページネーション対応）

### 認証
- `POST /api/v1/users/authenticate` - ユーザー認証

## リクエスト/レスポンス例

### ユーザー作成
**リクエスト:**
```json
POST /api/v1/users
{
    "email": "john.doe@example.com",
    "name": "John Doe",
    "password": "Password123!"
}
```

**レスポンス (201 Created):**
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "john.doe@example.com",
    "name": "John Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
}
```

### ユーザー更新
**リクエスト:**
```json
PUT /api/v1/users/{id}
{
    "name": "Updated Name",
    "email": "newemail@example.com"
}
```

**レスポンス (200 OK):**
```json
{
    "message": "User updated successfully"
}
```

### ユーザー一覧
**リクエスト:**
```
GET /api/v1/users?limit=10&offset=0
```

**レスポンス (200 OK):**
```json
{
    "users": [
        {
            "id": "123e4567-e89b-12d3-a456-426614174000",
            "email": "john.doe@example.com",
            "name": "John Doe",
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z"
        }
    ],
    "total_count": 1,
    "limit": 10,
    "offset": 0
}
```

## エラーレスポンス

### 400 Bad Request
```json
{
    "error": "Invalid email address",
    "code": "Bad Request"
}
```

### 404 Not Found
```json
{
    "error": "User not found",
    "code": "Not Found"
}
```

### 409 Conflict
```json
{
    "error": "User already exists",
    "code": "Conflict"
}
```

### 500 Internal Server Error
```json
{
    "error": "Internal server error",
    "code": "Internal Server Error"
}
```

## テストのアサーション

各リクエストには自動テストが含まれています：
- ステータスコードの検証
- レスポンス構造の検証
- 必須フィールドの存在確認
- 環境変数への値の保存

## トラブルシューティング

### サーバーに接続できない
1. サーバーが起動しているか確認
2. `base_url`環境変数が正しいか確認
3. ポート8080が使用可能か確認

### 認証エラー
1. ユーザーが作成されているか確認
2. パスワードが正しいか確認（最低8文字、大文字・小文字・数字を含む）

### データベースエラー
1. PostgreSQLが起動しているか確認
2. データベース接続設定を確認
3. マイグレーションが実行されているか確認

## 開発用コマンド

```bash
# データベースの初期化
docker-compose down -v
docker-compose up -d

# ログの確認
docker-compose logs -f db

# サーバーのビルドと実行
cd services/user
go build -o bin/server cmd/server/main.go
./bin/server
```