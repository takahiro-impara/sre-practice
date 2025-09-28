# SRE Skill-up Microservices Training Project

## プロジェクトの目的
SREスキル向上のための1ヶ月ハンズオントレーニングプロジェクト。マイクロサービス実践を通じて、実際のSRE業務に即した設計/実装/デプロイを体験する。

## 主要機能
- ユーザー登録システム
- イベント駆動通知システム（Outboxパターン）
- 非同期メッセージング

## サービス構成
1. **User Service**
   - REST API (chi framework)
   - gRPC/Connect API
   - PostgreSQL (Cloud SQL) でユーザー情報保存
   - Outboxパターンでイベント発行

2. **Notification Service**
   - Pub/Sub購読
   - UserCreatedイベント処理
   - 冪等性保証（correlation_id）

3. **Infrastructure**
   - GKE (Google Kubernetes Engine)
   - Cloud SQL (PostgreSQL)
   - Cloud Pub/Sub
   - Terraform管理

## 学習ロードマップ
- Week 1: REST API + PostgreSQL (CRUD + sqlc + Index)
- Week 2: gRPC(Connect) + goroutine/channel
- Week 3: Pub/Sub 非同期処理 + サービス分割
- Week 4: Terraform + GKE/Cloud SQL/CI/CD + 可観測性