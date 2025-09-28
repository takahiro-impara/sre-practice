# 技術スタック

## 言語とフレームワーク
- **言語**: Go 1.24.6
- **Webフレームワーク**: net/http + chi
- **gRPC**: Connect (connectrpc.com) + buf
- **テストフレームワーク**: testify

## データベース
- **PostgreSQL 16**
- **pgx**: PostgreSQLドライバ
- **sqlc**: SQLからGoコード生成
- **golang-migrate**: マイグレーション管理

## メッセージング
- Google Cloud Pub/Sub

## インフラストラクチャ
- **コンテナ**: Docker
- **オーケストレーション**: Kubernetes (GKE)
- **IaC**: Terraform
- **CI/CD**: GitHub Actions

## 開発ツール
- **リンター**: golangci-lint
- **フォーマッター**: gofmt, goimports
- **YAML検証**: yamllint
- **シェルスクリプト検証**: shellcheck

## 可観測性
- OpenTelemetry
- pprof
- Cloud Monitoring

## システム環境
- OS: Darwin (macOS)
- Git使用
- Docker Compose for local development