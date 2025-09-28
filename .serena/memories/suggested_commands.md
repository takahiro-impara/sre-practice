# 開発コマンド一覧

## データベース管理
```bash
# PostgreSQL起動
docker compose up -d
make run-db

# PostgreSQL停止
docker compose down
make stop-db

# PostgreSQL再起動
make restart-db

# マイグレーション
make migrate-up      # マイグレーション適用
make migrate-down    # マイグレーション取り消し
make migrate-status  # マイグレーション状態確認
make migrate-reset   # リセット
make migrate-create  # 新規マイグレーション作成
```

## コード生成
```bash
# SQLC生成（SQLからGoコード生成）
make generate-sqlc
(cd services/user/db/sqlc && sqlc generate)

# すべてのコード生成
make generate
```

## サービス起動
```bash
# User Service起動
go run ./services/user/cmd/server

# ヘルスチェック
curl -i http://localhost:8080/healthz
```

## テスト実行
```bash
# すべてのテスト
make test
go test ./...

# ショートテストのみ
make test-short

# Race detector付きテスト
make test-race

# カバレッジ測定
make test-cover
go test -cover ./...

# 特定ドメインのテスト
make test-user

# 統合テスト
make test-integration

# ベンチマーク
make test-bench
```

## コード品質チェック
```bash
# フォーマット
make fmt

# すべてのチェック（fmt, vet, lint, test含む）
make check

# Go vet
make vet

# リント
make lint

# リント（自動修正付き）
make lint-fix

# Shellスクリプトのリント
make shell-lint

# YAMLファイルのリント
make yaml-lint

# Docker Composeファイルの検証
make docker-lint
```

## 開発ツールインストール
```bash
# すべての開発ツールインストール
make install-tools

# リンターのみインストール
make install-linters
```

## Git操作
```bash
# 基本的なGitコマンド
git status
git diff
git log
git add .
git commit -m "message"
git push
```

## システムコマンド（macOS/Darwin）
```bash
# ディレクトリ操作
ls -la         # ファイル一覧（隠しファイル含む）
cd [dir]       # ディレクトリ移動
pwd            # 現在のディレクトリ表示

# ファイル操作
cat [file]     # ファイル内容表示
grep [pattern] # パターン検索
find . -name   # ファイル検索

# プロセス管理
ps aux         # プロセス一覧
kill -9 [pid]  # プロセス終了
```