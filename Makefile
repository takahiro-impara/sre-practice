# Database management
.PHONY: run-db stop-db restart-db
run-db:
	docker compose up -d

stop-db:
	docker compose down

restart-db:
	docker compose down
	docker compose up -d

# Code quality and linting
.PHONY: fmt vet lint lint-fix check test

# Go formatting
fmt:
	@echo "Running gofmt..."
	gofmt -s -w .
	@echo "Running goimports..."
	goimports -w .

# Go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Go linting with golangci-lint
lint:
	@echo "Running golangci-lint..."
	golangci-lint run

# Go linting with auto-fix
lint-fix:
	@echo "Running golangci-lint with auto-fix..."
	golangci-lint run --fix

# Shell script linting
shell-lint:
	@echo "Running shellcheck..."
	@if command -v shellcheck >/dev/null 2>&1; then \
		find . -type f -name "*.sh" -not -path "./vendor/*" -not -path "./.git/*" | xargs shellcheck; \
	else \
		echo "shellcheck not found. Install with: brew install shellcheck (macOS) or apt-get install shellcheck (Ubuntu)"; \
	fi

# YAML linting
yaml-lint:
	@echo "Running yamllint..."
	@if command -v yamllint >/dev/null 2>&1; then \
		yamllint -c .yamllint.yml .; \
	else \
		echo "yamllint not found. Install with: pip install yamllint"; \
	fi

# Docker Compose file validation
docker-lint:
	@echo "Validating docker-compose files..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose config --quiet; \
	else \
		docker compose config --quiet; \
	fi

# Comprehensive check (format, vet, lint, test)
check: fmt vet lint shell-lint yaml-lint docker-lint test
	@echo "All checks passed!"

# Run tests
.PHONY: test test-short test-race test-user test-user-service test-integration
test:
	@echo "Running all tests..."
	go test -v ./...

# Run short tests only (skip integration tests)
test-short:
	@echo "Running short tests..."
	go test -short -v ./...

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	go test -race -v ./...

test-user-all:
	@echo "Running user all tests..."
	go test -v ./services/user/internal/...

# Run specific domain tests
test-user:
	@echo "Running user domain tests..."
	go test -v ./services/user/internal/domain/...

# Run user service tests
test-user-service:
	@echo "Running user service tests..."
	cd services/user && go test -v ./...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v -tags=integration ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test -v -cover ./...
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests and generate coverage report for specific package
test-cover-user:
	@echo "Running user domain tests with coverage..."
	go test -v -cover -coverprofile=coverage-user.out ./services/user/internal/domain/...
	go tool cover -html=coverage-user.out -o coverage-user.html
	@echo "Coverage report generated: coverage-user.html"

# Run benchmarks
test-bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Code generation
.PHONY: generate generate-sqlc
generate: generate-sqlc
	@echo "All code generation complete!"

generate-sqlc:
	@echo "Generating SQLC code..."
	@if command -v sqlc >/dev/null 2>&1; then \
		cd services/user/db/sqlc && sqlc generate; \
		echo "SQLC code generated successfully"; \
	else \
		echo "sqlc not found. Install with: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest"; \
		exit 1; \
	fi

# Install development tools
.PHONY: install-tools install-linters
install-tools: install-linters
	@echo "Installing development tools..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "All tools installed!"

# Install linting tools
install-linters:
	@echo "Installing linting tools..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Install yamllint with: pip install yamllint"
	@echo "Install shellcheck with: brew install shellcheck (macOS) or apt-get install shellcheck (Ubuntu)"

# Migration shortcuts
.PHONY: migrate-up migrate-down migrate-status migrate-reset migrate-create
migrate-up:
	./scripts/migrate.sh up

migrate-down:
	./scripts/migrate.sh down

migrate-status:
	./scripts/migrate.sh status

migrate-reset:
	./scripts/migrate.sh reset

migrate-create:
	@read -p "Migration name: " name; \
	./scripts/migrate.sh create $$name
