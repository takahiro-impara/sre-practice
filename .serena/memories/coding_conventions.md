# コーディング規約とスタイル

## Go言語基本ポリシー
- Go公式のEffective Go、Code Review Comments、Golang Standardsに準拠
- gofmt、go vet、golint、golangci-lintの規則に従う
- パッケージ構成、命名規則、エラーハンドリング等はGoの慣習に従う
- interface設計、並行処理、メモリ管理はGoのベストプラクティスを適用

## 命名規則

### インターフェースと実装
- **インターフェース優先パターン**を採用
- インターフェースに主要な名前を与える（例：`UserService`）
- 実装は`impl`サフィックスか小文字始まり（例：`userServiceImpl`）
- 外部パッケージからはインターフェースのみを公開

### ファイル構成
```
services/user/
├── internal/
│   ├── domain/          # ドメインモデルとビジネスロジック
│   ├── handler/         # HTTP/gRPCハンドラ
│   ├── service/         # ビジネスロジックサービス
│   │   ├── interface.go # インターフェース定義
│   │   └── user_service.go # 実装
│   └── repository/      # データアクセス層
│       ├── interface.go # インターフェース定義
│       └── postgres.go  # PostgreSQL実装
```

### インターフェース設計
- 単一メソッド: `-er`形式（例：`Reader`, `Writer`, `Validator`）
- 複数メソッド: 役割名（例：`UserService`, `UserRepository`, `Manager`）
- インターフェース分離の原則（ISP）に従い、小さなインターフェースに分割

### 実装のチェック
- `var _ Interface = (*implementation)(nil)` でコンパイル時チェック
- ファクトリー関数はインターフェース型を返す
- 依存性はインターフェース型で受け取る

## golangci-lint設定
主要な有効化リンター:
- gofmt, goimports, govet
- errcheck, staticcheck, ineffassign
- gosimple, gocritic, revive
- misspell, bodyclose, stylecheck
- gosec, wrapcheck, errorlint

設定の詳細:
- タイムアウト: 5分
- テストファイル含む
- vendor, .gitディレクトリはスキップ
- Protocol Buffer生成ファイル（.pb.go）はスキップ

## 重要な規則
- ファイル末尾は必ず1行改行
- コメントは必要な場合のみ（過度なコメントは避ける）
- エラーハンドリングは確実に行う
- Context propagationを適切に実装
- セキュリティベストプラクティスに従う（シークレットをコードに含めない）