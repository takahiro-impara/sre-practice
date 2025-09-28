# プロジェクト構造

## ルートディレクトリ構成
```
.
├── services/           # マイクロサービス
│   └── user/          # ユーザーサービス
├── platform/          # インフラストラクチャ
│   ├── terraform/     # IaC (GKE, Cloud SQL, Pub/Sub, IAM, VPC)
│   └── k8s/          # Kubernetes manifests or Helm charts
├── ops/              # 運用関連
│   ├── runbooks/     # 運用手順書
│   └── slo/          # SLO定義
├── docs/             # ドキュメント
│   └── week1/        # Week1関連ドキュメント
├── scripts/          # ユーティリティスクリプト
│   └── migrate.sh    # マイグレーションスクリプト
└── .github/workflows/ # CI/CDワークフロー
```

## User Serviceの詳細構造
```
services/user/
├── cmd/
│   └── server/
│       └── main.go              # メインエントリポイント
├── internal/
│   ├── domain/                  # ドメインモデル
│   │   ├── user.go             # Userエンティティ
│   │   ├── user_test.go        # ユニットテスト
│   │   ├── validate.go         # バリデーション
│   │   └── error.go            # ドメインエラー定義
│   ├── handler/                 # HTTPハンドラ
│   ├── service/                 # ビジネスロジック層
│   │   ├── interface.go        # サービスインターフェース
│   │   ├── user_service.go     # サービス実装
│   │   └── user_service_test.go
│   ├── repository/              # リポジトリ層
│   │   ├── interface.go        # リポジトリインターフェース
│   │   ├── interface_example.go
│   │   ├── mock_user_repository.go  # モック
│   │   └── check_interface_test.go
│   └── infrastructure/
│       └── postgres/            # PostgreSQL実装
│           ├── user_repository.go
│           ├── user_repository_test.go
│           ├── converter.go     # ドメインモデル変換
│           └── errors.go        # DBエラー処理
├── db/
│   ├── migrations/              # データベースマイグレーション
│   │   ├── 000001_create_users_table.up.sql
│   │   ├── 000001_create_users_table.down.sql
│   │   ├── 000002_add_users_password.up.sql
│   │   └── 000002_add_users_password.down.sql
│   └── sqlc/                    # SQLC設定
│       ├── sqlc.yaml           # SQLC設定ファイル
│       ├── queries.sql         # SQLクエリ定義
│       └── generated/          # 生成コード
│           ├── db.go
│           ├── models.go
│           ├── queries.sql.go
│           └── querier.go
├── proto/                       # Protocol Buffer定義（将来実装）
└── Dockerfile                   # コンテナイメージ定義
```

## 重要な設定ファイル
- **go.mod**: Goモジュール定義
- **go.sum**: 依存関係のチェックサム
- **Makefile**: ビルド・テスト・デプロイタスク
- **.golangci.yml**: Goリンター設定
- **.yamllint.yml**: YAMLリンター設定
- **docker-compose.yml**: ローカル開発環境
- **CLAUDE.md**: Claude Code用プロジェクト指示
- **.gitignore**: Git無視ファイル設定

## データモデル

### Users テーブル
- id (uuid): 主キー
- email (text): UNIQUE制約
- name (text)
- created_at (timestamptz)
- updated_at (timestamptz)

### Outbox Events テーブル（将来実装）
- seq_id (bigserial): 主キー
- event_id (uuid): UNIQUE
- event_type (text)
- payload (jsonb)
- status (text): pending|published|failed
- その他メタデータ

### Notification Jobs テーブル（将来実装）
- id (bigserial): 主キー
- correlation_id (uuid): 冪等キー
- user_id (uuid): 外部キー
- その他通知関連フィールド