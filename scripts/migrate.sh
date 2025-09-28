#!/bin/bash

# マイグレーション実行スクリプト
set -euo pipefail

# スクリプトのディレクトリを取得
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# 色付きの出力
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 環境変数ファイル読み込み（シンプルなsource方式）
load_env() {
    local env_file="${PROJECT_ROOT}/.env.local"

    if [ -f "$env_file" ]; then
        log_info "Loading environment from $env_file"
        set -a  # 自動export
        # shellcheck disable=SC1090
        source "$env_file"
        set +a
    else
        log_warn "Environment file $env_file not found, using defaults"
    fi
}

# デフォルト設定
set_defaults() {
    DATABASE_URL=${DATABASE_URL:-"postgres://app:app@localhost:5432/appdb?sslmode=disable"}
    MIGRATIONS_PATH=${MIGRATIONS_PATH:-"services/user/db/migrations"}

    # 絶対パス化
    if [[ ! "$MIGRATIONS_PATH" = /* ]]; then
        MIGRATIONS_PATH="${PROJECT_ROOT}/${MIGRATIONS_PATH}"
    fi
}

# golang-migrate の存在確認
check_migrate_tool() {
    if ! command -v migrate &> /dev/null; then
        log_error "golang-migrate not found!"
        echo "Install: brew install golang-migrate"
        exit 1
    fi
}

# マイグレーション実行
migrate_up() {
    log_info "Running migrations..."
    log_info "Path: $MIGRATIONS_PATH"

    if [[ -n "${1:-}" ]]; then
        migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" up "$1"
        log_success "Applied $1 migration(s)"
    else
        migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" up
        log_success "All pending migrations applied"
    fi
}

# マイグレーションロールバック
migrate_down() {
    local steps=${1:-1}

    log_warn "Rolling back $steps migration(s)..."

    # 本番環境チェック
    if [[ "${ENV:-}" == "production" ]]; then
        echo -n "⚠️  Production environment! Continue? [y/N] "
        read -r confirm
        [[ "$confirm" != "y" ]] && exit 0
    fi

    migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" down "$steps"
    log_success "Rolled back $steps migration(s)"
}

# マイグレーション状態確認
migrate_status() {
    log_info "Current migration status:"
    migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" version || {
        if [[ $? -eq 1 ]]; then
            log_info "No migrations applied yet"
        fi
    }
}

# 特定のバージョンに移動
migrate_goto() {
    local version=$1

    if [[ -z "$version" ]]; then
        log_error "Version number required"
        exit 1
    fi

    log_info "Migrating to version $version..."
    migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" goto "$version"
    log_success "Migrated to version $version"
}

# データベースリセット
migrate_reset() {
    if [[ "${ENV:-}" == "production" ]]; then
        log_error "Reset not allowed in production!"
        exit 1
    fi

    log_warn "Resetting database..."
    echo -n "This will destroy all data! Continue? [y/N] "
    read -r confirm

    if [[ "$confirm" == "y" ]]; then
        migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" drop -f
        migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" up
        log_success "Database reset completed"
    fi
}

# 新しいマイグレーション作成
create_migration() {
    local name=${1:-}

    if [[ -z "$name" ]]; then
        echo -n "Migration name: "
        read -r name
    fi

    if [[ -z "$name" ]]; then
        log_error "Migration name required"
        exit 1
    fi

    log_info "Creating migration: $name"
    migrate create -ext sql -dir "$MIGRATIONS_PATH" -seq "$name"
    log_success "Migration files created"
}

# ヘルプ表示
show_help() {
    cat << EOF
Migration management script

USAGE:
    $0 <command> [arguments]

COMMANDS:
    up [steps]      Run migrations
    down [steps]    Rollback migrations (default: 1)
    status          Show current status
    goto <version>  Migrate to version
    reset           Reset database (dev only)
    create [name]   Create new migration
    help            Show this help

EXAMPLES:
    $0 up           # Run all pending migrations
    $0 down         # Rollback last migration
    $0 status       # Check current version
    $0 create       # Create new migration
EOF
}

# メイン処理
main() {
    local command=${1:-help}

    # 環境変数読み込み
    load_env
    set_defaults

    case "$command" in
        up)
            check_migrate_tool
            migrate_up "${2:-}"
            ;;
        down)
            check_migrate_tool
            migrate_down "${2:-1}"
            ;;
        status)
            check_migrate_tool
            migrate_status
            ;;
        goto)
            check_migrate_tool
            migrate_goto "${2:-}"
            ;;
        reset)
            check_migrate_tool
            migrate_reset
            ;;
        create)
            check_migrate_tool
            create_migration "${2:-}"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# スクリプト実行
main "$@"