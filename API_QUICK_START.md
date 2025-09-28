# API Quick Start Guide

## cURLコマンドでのテスト

### 1. サーバー起動確認

```bash
# ヘルスチェック
curl -i http://localhost:8080/healthz

# レディネスチェック
curl -i http://localhost:8080/readyz
```

### 2. ユーザー作成

```bash
# 新規ユーザー作成
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "name": "Alice Smith",
    "password": "Password123!"
  }'

# レスポンス例:
# {
#   "id": "550e8400-e29b-41d4-a716-446655440000",
#   "email": "alice@example.com",
#   "name": "Alice Smith",
#   "created_at": "2024-01-01T12:00:00Z",
#   "updated_at": "2024-01-01T12:00:00Z"
# }
```

### 3. ユーザー取得

```bash
# IDを指定してユーザー取得（上記で取得したIDを使用）
USER_ID="550e8400-e29b-41d4-a716-446655440000"
curl -i http://localhost:8080/api/v1/users/$USER_ID
```

### 4. ユーザー更新

```bash
# 名前のみ更新
curl -X PUT http://localhost:8080/api/v1/users/$USER_ID \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Johnson"
  }'

# メールアドレスのみ更新
curl -X PUT http://localhost:8080/api/v1/users/$USER_ID \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice.johnson@example.com"
  }'

# 両方更新
curl -X PUT http://localhost:8080/api/v1/users/$USER_ID \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice.new@example.com",
    "name": "Alice New Name"
  }'
```

### 5. ユーザー一覧取得

```bash
# デフォルト（limit=10, offset=0）
curl -i http://localhost:8080/api/v1/users

# ページネーション指定
curl -i "http://localhost:8080/api/v1/users?limit=5&offset=0"
curl -i "http://localhost:8080/api/v1/users?limit=5&offset=5"
```

### 6. ユーザー認証

```bash
# 正しいパスワードで認証
curl -X POST http://localhost:8080/api/v1/users/authenticate \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "Password123!"
  }'

# 間違ったパスワードで認証（エラーになる）
curl -X POST http://localhost:8080/api/v1/users/authenticate \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "WrongPassword"
  }'
```

### 7. ユーザー削除

```bash
# ユーザー削除
curl -X DELETE http://localhost:8080/api/v1/users/$USER_ID
```

## HTTPieを使用した場合

HTTPieをインストール:
```bash
brew install httpie  # macOS
# または
pip install httpie   # Python pip
```

### HTTPieコマンド例

```bash
# ヘルスチェック
http GET localhost:8080/healthz

# ユーザー作成
http POST localhost:8080/api/v1/users \
  email="bob@example.com" \
  name="Bob Smith" \
  password="Password123!"

# ユーザー取得
http GET localhost:8080/api/v1/users/$USER_ID

# ユーザー更新
http PUT localhost:8080/api/v1/users/$USER_ID \
  name="Bob Johnson"

# ユーザー一覧
http GET localhost:8080/api/v1/users limit==5 offset==0

# 認証
http POST localhost:8080/api/v1/users/authenticate \
  email="bob@example.com" \
  password="Password123!"

# ユーザー削除
http DELETE localhost:8080/api/v1/users/$USER_ID
```

## 一括テストスクリプト

`test-api.sh`として保存して実行:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "1. Health Check"
curl -s "$BASE_URL/healthz"
echo -e "\n"

echo "2. Create User"
RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/users" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "Password123!"
  }')

echo $RESPONSE | jq '.'
USER_ID=$(echo $RESPONSE | jq -r '.id')
echo "Created user with ID: $USER_ID"
echo -e "\n"

echo "3. Get User"
curl -s "$BASE_URL/api/v1/users/$USER_ID" | jq '.'
echo -e "\n"

echo "4. Update User"
curl -s -X PUT "$BASE_URL/api/v1/users/$USER_ID" \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Name"}' | jq '.'
echo -e "\n"

echo "5. List Users"
curl -s "$BASE_URL/api/v1/users?limit=5" | jq '.'
echo -e "\n"

echo "6. Authenticate User"
curl -s -X POST "$BASE_URL/api/v1/users/authenticate" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!"
  }' | jq '.'
echo -e "\n"

echo "7. Delete User"
curl -s -X DELETE "$BASE_URL/api/v1/users/$USER_ID"
echo "User deleted"
```

実行方法:
```bash
chmod +x test-api.sh
./test-api.sh
```

## Docker Composeでの起動

```bash
# データベース起動
docker-compose up -d

# データベース接続確認
docker-compose exec db psql -U app -d appdb -c "SELECT 1"

# サーバー起動（別ターミナル）
cd services/user
go run cmd/server/main.go

# または環境変数を指定して起動
DATABASE_URL="postgres://app:app@localhost:5432/appdb?sslmode=disable" \
  go run cmd/server/main.go
```

## トラブルシューティング

### Connection refused エラー
```bash
# サーバーが起動しているか確認
lsof -i :8080

# プロセスを確認
ps aux | grep "go run"
```

### Database connection エラー
```bash
# PostgreSQLが起動しているか確認
docker-compose ps

# ログを確認
docker-compose logs db

# 接続テスト
psql -h localhost -U app -d appdb -p 5432
```

### Invalid password エラー
パスワードは以下の条件を満たす必要があります：
- 最低8文字
- 大文字を含む
- 小文字を含む
- 数字を含む

例: `Password123!`, `SecurePass456`