# AGENTS.md

Guidance for AI coding agents working in this repository. Read this file before making changes.

## Project

**brightloop-subscription-manager** — a subscription management platform for Brightloop Software. It covers accounts and roles, products and plans, quotations, subscriptions, recurring invoices, taxes, discounts, payments and revenue reporting. The source requirements are in `Subscription_Management_System.pdf` (not committed).

## Stack

- **Backend:** Go 1.23+, standard-library `net/http` routing (Go 1.22+ method patterns)
- **Database:** PostgreSQL 16+, accessed with `database/sql` + `pgx/v5/stdlib`
- **Auth:** JWT bearer tokens (`golang-jwt/jwt/v5`) and bcrypt password hashes
- **Frontend:** not in scope for this repository yet

All application code lives in `backend/`. Do not add frontend code.

## Repository layout

```
.
├── AGENTS.md
├── LICENSE                 # GNU GPL v3.0
├── README.md
└── backend/
    ├── cmd/api/            # HTTP server entrypoint
    ├── cmd/migrate/        # SQL migration runner
    ├── internal/
    │   ├── api/            # router, handlers, middleware, response helpers
    │   ├── auth/           # JWT + bcrypt
    │   ├── config/         # environment config
    │   ├── database/       # PostgreSQL pool
    │   ├── domain/         # entities, enums, sentinel errors
    │   ├── repository/     # interfaces + postgres/ implementations
    │   ├── service/        # business rules
    │   └── validation/     # input validation
    ├── migrations/         # numbered .up.sql / .down.sql files
    ├── Dockerfile
    ├── docker-compose.yml
    ├── Makefile
    └── go.mod
```

Dependencies point inward: `api` → `service` → `repository` → `domain`. `domain` must not import other internal packages.

## Commands

Run from `backend/`:

| Task           | Command                |
| -------------- | ---------------------- |
| Run API        | `go run ./cmd/api`     |
| Run migrations | `go run ./cmd/migrate` |
| Build          | `go build ./...`       |
| Test           | `go test ./...`        |
| Format         | `gofmt -w .`           |
| Vet            | `go vet ./...`         |
| Tidy modules   | `go mod tidy`          |
| Docker         | `docker compose up --build` |

Configuration is read from environment variables; see `backend/.env.example`. `DATABASE_URL` and a `JWT_SECRET` of at least 32 characters are required for the API.

## Conventions

- **Money:** store integer minor units (cents) in `*Cents` fields and `BIGINT` columns. Never use floats for money.
- **Percentages:** store basis points in `*BasisPoints` fields (2000 = 20.00%).
- **Dates:** `DATE` for billing dates, `TIMESTAMPTZ` for audit timestamps.
- **Errors:** repositories return sentinel errors from `internal/domain` (`ErrNotFound`, `ErrConflict`, `ErrUnauthorized`, `ErrForbidden`, `ErrValidation`); `response.Error` maps them to HTTP status codes.
- **SQL:** only in `internal/repository/postgres`. Services depend on interfaces.
- **Validation:** business validation belongs in `internal/service` or `internal/validation`, not handlers.
- **Handlers:** decode with `response.DecodeJSON`, write with `response.JSON`, and return `response.Error` on failure.
- **Migrations:** additive, numbered (`0002_...up.sql` / `.down.sql`). Never edit an applied migration.
- **Tests:** place `*_test.go` next to the code under test.

## Adding a module

Use the auth slice as the reference implementation:

1. Add entities/enums in `internal/domain`.
2. Add a repository interface in `internal/repository` and a PostgreSQL implementation in `internal/repository/postgres`.
3. Add business rules in `internal/service`.
4. Add handlers in `internal/api/handler` and register routes in `internal/api/router.go`.
5. Add a numbered migration in `migrations/`.
6. Add tests.

## Workflow

1. Read `README.md`, `backend/README.md` and this file first.
2. Check `git status` before editing to avoid overwriting uncommitted work.
3. Make minimal, focused changes that match existing patterns.
4. Run `gofmt`, `go vet ./...` and `go test ./...` before finishing.
5. Update documentation when adding routes, commands or architectural decisions.

## Do not

- Commit `.env` files, secrets, generated binaries or dependency caches.
- Put SQL in handlers or services.
- Use floating-point types for money.
- Rewrite git history or force-push shared branches.
- Add frontend code or non-Go application code.
- Modify `LICENSE` unless explicitly asked.
