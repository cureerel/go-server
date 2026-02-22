# gotemplate — Production Go REST API

A production-grade REST API built with Go, following **Clean Architecture** principles. Designed to be scalable, maintainable, and infrastructure-agnostic.

---

## Architecture

This project follows Clean Architecture with four distinct layers:

```
domain → application → infrastructure → interfaces
```

Each layer depends only on the layer inside it. The domain layer has zero external dependencies.

### Layer Breakdown

**`domain/`** — Pure business logic. No frameworks, no GORM, no HTTP.
- `entity/` — Business structs (User, Blog, Order, Product, Membership, Payment, Webhook)
- `repository/` — Interfaces only. Defines what data operations exist, not how.

**`application/service/`** — Use cases. Orchestrates domain logic.
- Calls repository interfaces
- Enforces business rules (e.g. valid order transitions, membership upgrades)
- Never imports GORM or Gin

**`infrastructure/`** — All external concerns.
- `postgres/models/` — GORM structs with `ToDomain()` / `FromDomain()` mappers
- `postgres/repositories/` — Implements domain repository interfaces using GORM
- `redis/` — Redis client
- `factory.go` — Database and Redis constructors

**`interfaces/http/`** — HTTP delivery layer.
- `handler/` — Gin handlers. Parse request → call service → return response
- `middleware/` — Auth, RBAC, logger, request ID, error handler
- `router/` — Route registration with role guards
- `dto/` — Request/response structs

---

## Project Structure

```
.
├── cmd/server/main.go          # Entrypoint — wires everything together
├── configs/config.yaml         # All configuration (loaded at startup)
├── migrations/                 # Atlas SQL migration files
├── atlas.hcl                   # Atlas migration config
├── internal/
│   ├── domain/
│   │   ├── entity/             # Business entities (no GORM tags)
│   │   └── repository/         # Repository interfaces
│   ├── application/service/    # Business logic / use cases
│   ├── infrastructure/
│   │   ├── postgres/
│   │   │   ├── models/         # GORM models + domain mappers
│   │   │   └── repositories/   # Repository implementations
│   │   ├── redis/              # Redis client
│   │   ├── dbtypes/            # DB interface types
│   │   └── factory.go          # DB + Redis constructors
│   └── interfaces/http/
│       ├── handler/            # HTTP handlers
│       ├── middleware/         # Gin middleware
│       ├── router/             # Route setup
│       └── dto/                # Request/response DTOs
└── pkg/
    ├── logger/                 # Structured logger
    ├── apperror/               # App error types
    └── utils/                  # Slug generator etc.
```

---

## Tech Stack

| Concern | Library |
|---|---|
| HTTP | Gin |
| ORM | GORM |
| Database | PostgreSQL |
| Cache | Redis |
| Migrations | Atlas |
| Auth | JWT (golang-jwt/jwt) |
| Password | bcrypt |
| Config | YAML |
| Logging | Custom structured logger |

---

## Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL
- Redis
- Atlas CLI (`brew install ariga/tap/atlas`)

### Setup

```bash
# Clone and install dependencies
git clone https://github.com/cureerel/gotemplate
cd gotemplate
make deps

# Configure environment
cp .env.example .env
# Edit .env with your DB credentials

# Run migrations
make migrate-up

# Start server
make run
```

### Environment Variables (`.env`)

```env
DATABASE_URL=postgres://postgres:secret@localhost:5432/mydb?sslmode=disable
DEV_DATABASE_URL=postgres://postgres:secret@localhost:5432/mydb_dev?sslmode=disable
```

### Config (`configs/config.yaml`)

```yaml
server:
  port: "8080"
  env: "development"

database:
  driver: "postgres"
  dsn: "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable"

jwt:
  access_secret: "your-access-secret"
  refresh_secret: "your-refresh-secret"

webhook:
  stripe_secret: "your-stripe-secret"
  razorpay_secret: "your-razorpay-secret"

redis:
  addr: "localhost:6379"
  password: ""
  db: 0
```

---

## Make Commands

| Command | Description |
|---|---|
| `make run` | Run in development |
| `make build` | Build binary |
| `make prod-build` | Build Linux/amd64 binary for production |
| `make deps` | Download and tidy dependencies |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Revert last migration |
| `make migrate-status` | Show migration status |
| `make migrate-create NAME=x` | Create new migration file |
| `make test` | Run tests |
| `make lint` | Run linter |
| `make fmt` | Format code |

---

## Key Design Decisions

**Why separate `entity` from `models`?**
Domain entities have no GORM tags. This means you can swap Postgres for MongoDB without touching business logic. Models handle the DB mapping; services never know GORM exists.

**Why context everywhere?**
All repository and service methods accept `context.Context` as the first argument. This allows proper request cancellation, timeouts, and tracing propagation.

**Why typed constants over raw strings?**
`entity.OrderStatus`, `entity.MembershipPlan`, `entity.Currency` etc. are typed string aliases. This prevents invalid values being passed around silently and makes code self-documenting.

**Why is Payment in WebhookRepository?**
Payment records are created as a direct result of webhook events for compliance — every payment must trace back to a webhook event. A dedicated `PaymentRepository` should be split out when payment history queries (by user, date range, etc.) are needed.