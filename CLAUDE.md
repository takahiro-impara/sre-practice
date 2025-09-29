`# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Japanese SRE skill-up training project that demonstrates microservices architecture using Go, PostgreSQL, gRPC (Connect), Pub/Sub, Terraform, and GKE. The project implements a user registration system with event-driven notifications using the Outbox pattern.

## Development Commands

### Database Setup
```bash
# Start PostgreSQL database
docker compose up -d

# Generate SQLC code (when implemented)
(cd services/user/db/sqlc && sqlc generate)
```

### Service Development
```bash
# Run User Service
go run ./services/user/cmd/server

# Health check
curl -i http://localhost:8080/healthz
```

### Protocol Buffer Generation (when implemented)
```bash
# Generate Connect/gRPC code using buf
buf generate
```

### Testing
```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Architecture

### Services Structure
- **User Service**: REST (chi) and gRPC (Connect) APIs with PostgreSQL storage
- **Notification Service**: Pub/Sub subscriber for handling user creation events
- **Infrastructure**: Terraform-managed GKE, Cloud SQL, Pub/Sub

### Key Directories
```
services/user/
├── cmd/server/           # Main application entry point
├── internal/
│   ├── domain/          # Domain models and business logic
│   ├── handler/         # HTTP and gRPC handlers
│   ├── repository/      # Data access layer
│   └── service/         # Business logic services
├── db/
│   ├── migrations/      # Database schema migrations
│   └── sqlc/           # SQLC configuration and generated code
└── proto/              # Protocol buffer definitions

platform/
├── terraform/          # Infrastructure as Code
└── k8s/               # Kubernetes manifests

ops/
└── runbooks/          # Operational documentation
```

### Data Models
- **Users**: Core user entity (id, email, name, timestamps)
- **Outbox Events**: Event sourcing for reliable message publishing
- **Notification Jobs**: Asynchronous notification processing with idempotency

## Technology Stack

- **Language**: Go 1.24+
- **Web Framework**: net/http + chi
- **gRPC**: Connect (connectrpc.com) + buf for schema management
- **Database**: PostgreSQL + pgx + sqlc
- **Messaging**: Google Cloud Pub/Sub
- **Infrastructure**: Terraform + GKE + Cloud SQL
- **Observability**: OpenTelemetry + pprof + Cloud Monitoring

## API Endpoints

### REST API
- `POST /users` - Create new user
- `GET /users/{id}` - Get user by ID
- `GET /users?limit=&offset=` - List users with pagination
- `GET /healthz` - Health check

### gRPC/Connect (when implemented)
- `CreateUser` - Create new user
- `GetUser` - Get user by ID
- `ListUsers` - List users

## Event-Driven Architecture

### Pub/Sub Events
Topic: `user.events`

Event structure:
```json
{
  "event_id": "uuid",
  "event_type": "user.created",
  "version": 1,
  "occurred_at": "2025-09-27T12:00:00Z",
  "producer": "user-service",
  "payload": {
    "user_id": "uuid",
    "email": "alice@example.com",
    "name": "Alice"
  },
  "trace": {
    "trace_id": "xxx",
    "span_id": "yyy"
  }
}
```

## Development Workflow

This project follows a 4-week learning progression:

1. **Week 1**: REST API + PostgreSQL with CRUD operations and indexing
2. **Week 2**: gRPC (Connect) implementation with concurrent processing
3. **Week 3**: Pub/Sub async processing and service decomposition
4. **Week 4**: Terraform + GKE deployment with observability

## Monitoring & SLO

### Key Metrics
- HTTP/gRPC success rate and latency (p95)
- Outbox pending events and publish failures
- Notification backlog and retry counts
- Database connection pool utilization

### SLO Examples
- **User API success rate**: 99.5% over 30 days
- **User API p95 latency**: ≤300ms
- **Outbox publish lag p95**: ≤30 seconds

## Important Patterns

### Outbox Pattern
User creation atomically writes to both users table and outbox_events table. A background worker publishes events to Pub/Sub and updates the outbox status.

### Idempotency
Notification service uses correlation_id for idempotent message processing to handle duplicate events gracefully.

### Context Propagation
All services should properly propagate request context for tracing and timeout handling across service boundaries.