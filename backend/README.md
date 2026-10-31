# Brightloop Subscription Manager — Backend

Go HTTP API for the subscription management platform described in
`Subscription_Management_System.pdf`. It covers accounts and roles, products,
plans, quotations, subscriptions, invoices, taxes, discounts, payments and
reporting.

## Layout

```
backend/
├── cmd/
│   ├── api/                 # HTTP server entrypoint
│   └── migrate/             # SQL migration runner
├── internal/
│   ├── api/                 # HTTP layer
│   │   ├── handler/         # request/response handlers
│   │   ├── middleware/      # auth, logging, panic recovery
│   │   ├── response/        # JSON helpers and error mapping
│   │   └── router.go        # route registration
│   ├── auth/                # JWT issuing/parsing and bcrypt helpers
│   ├── config/              # environment configuration
│   ├── database/            # PostgreSQL connection pool
│   ├── domain/              # entities, enums and sentinel errors
│   ├── repository/          # persistence interfaces
│   │   └── postgres/        # PostgreSQL implementations
│   ├── service/             # business rules
│   └── validation/          # reusable input validation
├── migrations/              # numbered up/down SQL files
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

Dependencies point inward: `api` → `service` → `repository` → `domain`.
`domain` must not import any other internal package.

## Running locally

Requirements: Go 1.23+, PostgreSQL 16+ (or Docker).

```bash
cd backend
cp .env.example .env      # then edit values
```

Load the environment variables into your shell (or use your tool of choice),
then:

```bash
go mod tidy
go run ./cmd/migrate      # apply migrations
go run ./cmd/api          # start the API on :8080
```

With Docker:

```bash
docker compose up --build
```

The API container does not run migrations automatically; run
`go run ./cmd/migrate` against the exposed PostgreSQL port.

## Commands

| Task            | Command                 |
| --------------- | ----------------------- |
| Run API         | `go run ./cmd/api`      |
| Run migrations  | `go run ./cmd/migrate`  |
| Build           | `go build ./...`        |
| Test            | `go test ./...`         |
| Format          | `gofmt -w .`            |
| Vet             | `go vet ./...`          |
| Tidy modules    | `go mod tidy`           |

## Endpoints

| Method | Path                    | Access        |
| ------ | ----------------------- | ------------- |
| GET    | `/health/live`          | public        |
| GET    | `/health/ready`         | public        |
| POST   | `/api/v1/auth/register` | public        |
| POST   | `/api/v1/auth/login`    | public        |
| GET    | `/api/v1/auth/me`       | authenticated |
| GET    | `/api/v1/admin/ping`    | admin only    |

`POST /api/v1/auth/register` accepts:

```json
{
  "email": "owner@example.com",
  "password": "Str0ng!Pass",
  "full_name": "Owner Name"
}
```

`POST /api/v1/auth/login` accepts `email` and `password`. Both return a
`token` to send as `Authorization: Bearer <token>`.

## Conventions

- **Money is stored in minor units.** Every `*Cents` field is an integer
  amount in cents. Never use floats for money.
- **Percentages are basis points.** A 20% tax rate is `2000`; a 7.5% discount
  is `750`.
- **Dates.** Billing dates use `DATE`; audit timestamps use `TIMESTAMPTZ`.
- **Errors.** Repositories return the sentinel errors in `internal/domain`
  (`ErrNotFound`, `ErrConflict`, ...); `response.Error` maps them to HTTP.
- **Validation.** Handlers decode and delegate; business validation belongs in
  `internal/service` or `internal/validation`.
- **No SQL outside the repository layer.** Services depend on interfaces, not
  `*sql.DB`.

## Adding a module

Use the auth slice as the reference:

1. Add or extend entities in `internal/domain`.
2. Add a repository interface in `internal/repository` and a PostgreSQL
   implementation in `internal/repository/postgres`.
3. Put business rules in `internal/service`.
4. Add HTTP handlers in `internal/api/handler` and register routes in
   `internal/api/router.go`.
5. Add a numbered migration under `migrations/`.
6. Add tests next to the code (`*_test.go`).

## Suggested next modules

- Products and plans (admin CRUD, variants, availability rules)
- Subscriptions and quotations (state machine: draft → quotation → confirmed
  → active → closed)
- Billing run (idempotent invoice generation per subscription and period)
- Payments and outstanding/overdue invoice queries
- Taxes and discounts (automatic application and usage limits)
- Reporting (active subscriptions, revenue, outstanding balance, overdue
  invoices) with pagination and indexed filters
