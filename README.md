# brightloop-subscription-manager

A subscription management platform for Brightloop Software: products and
plans, quotations, subscriptions, recurring invoices, taxes, discounts,
payments and revenue reporting.

The backend is a Go HTTP API. See [`backend/README.md`](backend/README.md) for
the layout, setup and commands.

## Repository layout

```
.
├── AGENTS.md        # guidance for AI coding agents
├── LICENSE          # GNU GPL v3.0
├── README.md
└── backend/         # Go API (the only application code for now)
```

## Quick start

```bash
cd backend
cp .env.example .env
go mod tidy
go run ./cmd/migrate
go run ./cmd/api
```

Requires Go 1.23+ and PostgreSQL 16+, or use `docker compose up --build` for
the database.
