# Brightloop Backend — Microservices

Microservice backend for the Brightloop subscription management platform. Each
service is an independent Go module with a Gin server. HAProxy is the single
public entry point.

> The **subscription** service is fully implemented: subscription lifecycle,
> billing schedules, recurring invoices, renewals, transactions, dashboard
> analytics and revenue reports. Auth, user and notification are still
> scaffolds whose business endpoints return `501 Not Implemented`.

## Architecture

```
                          +------------------------+
                          |        HAProxy         |
                          |   API Gateway :8080    |
                          |   Stats       :8404    |
                          +-----------+------------+
                                      |
          +---------------+-----------+-----------+---------------+
          |               |                       |               |
     +----v----+     +----v----+             +----v-----+    +----v-----+
     |  auth   |     |  user   |             |subscript.|    |notification|
     |  :8081  |     |  :8082  |             |  :8083   |    |  :8084    |
     +---------+     +---------+             +----------+    +-----------+
```

| Service              | Folder                  | Internal port | Public prefix                                   |
| -------------------- | ----------------------- | ------------- | ----------------------------------------------- |
| Auth                 | `services/auth`         | 8081          | `/api/v1/auth`                                  |
| User                 | `services/user`         | 8082          | `/api/v1/users`                                 |
| Subscription         | `services/subscription` | 8083          | `/api/v1/subscriptions`, `/api/subscriptions`   |
| Notification         | `services/notification` | 8084          | `/api/v1/notifications`                         |

All services share the PostgreSQL 16 instance exposed on `postgres:5432`; the
schema is defined in `db/schema.sql`.

The frontend's `VITE_API_URL` defaults to `http://localhost:3000/api`, so the
gateway also publishes port `3000` and the subscription service registers the
`/api/subscriptions` alias next to the canonical `/api/v1/subscriptions`
prefix. No frontend changes are required.

## Gateway responsibilities

- **Routing** by path prefix to the correct service.
- **Load balancing** with round-robin (one replica per service for now).
- **Health checks** against each service's `/health`.
- **Rate limiting**: 100 requests per 10 seconds per source IP.
- **Authentication filter (draft)**: protected routes require an
  `Authorization: Bearer ...` header. This only checks that the header exists;
  replace it with real JWT validation (forward-auth/SPOE or an auth-service
  call) before production.
- **Request/response handling**: HTTP keep-alive, forwarding headers and
  structured request logs.
- **Stats dashboard**: `http://localhost:8404/stats`.

## Layout

```
backend/
├── db/
│   └── schema.sql           # PostgreSQL schema, auto-loaded on first start
├── gateway/
│   └── haproxy.cfg          # routing, rate limit, health checks
├── services/
│   ├── auth/                # register/login/logout, JWT, password reset, RBAC
│   │   ├── main.go
│   │   ├── go.mod
│   │   └── Dockerfile
│   ├── user/                # profile, preferences, account settings
│   ├── subscription/        # subscription CRUD, billing, renewals, reports
│   └── notification/        # email, in-app, push, templates
├── monolith/                # archived previous single-service draft
├── docker-compose.yml
├── Makefile
└── README.md
```

`monolith/` is kept only as a reference for the domain model and auth code. It
is **not** part of the compose stack and can be deleted once it is no longer
useful.

## Subscription service

The subscription service follows a strict inward dependency chain:

```
Routes/Controllers (api) → Services (business rules) → Repositories (SQL) → PostgreSQL
```

```
services/subscription/
├── main.go                  # config, database, repository, service and router wiring
└── internal/
    ├── api/                 # Gin router, handlers, middleware, DTOs, response helpers
    ├── config/              # environment configuration
    ├── database/            # pgx-backed database/sql pool
    ├── domain/              # entities, enums, filters, sentinel errors
    ├── repository/          # interfaces + postgres/ implementations
    └── service/             # subscription, billing and analytics business rules
```

Conventions:

- **Money** is stored and transported as integer minor units (`*Cents`) and also
  exposed as decimal major units (`price`, `amount`, `mrr`, …) for the frontend.
- **Percentages** are stored as basis points (`2000` = `20.00%`) and returned as
  human-readable percentages where the UI needs them.
- **Dates** use `DATE` for billing periods and `TIMESTAMPTZ` for audit stamps.
- **Errors** are mapped to `{"success": false, "error": "<code>", "message":
  "..."}`. The frontend `apiClient` reads `message`.
- **Pagination** responses match `PaginatedResponse<T>`: `items`, `total`,
  `page`, `pageSize`, `totalPages`.

### Configuration

The service reads its own environment variables (also wired through
`docker-compose.yml`):

| Variable                | Default                                                              | Purpose                                                        |
| ----------------------- | -------------------------------------------------------------------- | -------------------------------------------------------------- |
| `PORT`                  | `8083`                                                               | Internal HTTP port                                             |
| `DATABASE_URL`          | `postgres://brightloop:brightloop@localhost:5432/brightloop?sslmode=disable` | PostgreSQL connection string                            |
| `JWT_SECRET`            | empty                                                                | When set, bearer tokens are signature-verified (HS256/384/512) |
| `AUTH_REQUIRED`         | `false`                                                              | Reject missing/invalid tokens (`true` in shared environments)  |
| `AUTH_DEFAULT_ROLE`     | `admin`                                                              | Role assigned in development when no token is present          |
| `REPORTING_CURRENCY`    | `USD`                                                                | Currency used to format dashboard metric values                |
| `CORS_ALLOWED_ORIGINS`  | `*`                                                                  | Comma-separated list of allowed browser origins                |

With `AUTH_REQUIRED=false` the service accepts the frontend demo token
(`mock_jwt_token_demo`) and falls back to `AUTH_DEFAULT_ROLE`. Tokens of the form
`mock_token_admin`, `mock_token_staff` and `mock_token_customer` select the
matching role for local RBAC testing. In production set `AUTH_REQUIRED=true`
and a 32+ character `JWT_SECRET`.

### Endpoints

All routes are served under both `/api/v1/subscriptions` and
`/api/subscriptions`. `GET /health` is also available at the service root.

| Method   | Path                  | Access   | Responsibility                                              |
| -------- | --------------------- | -------- | ----------------------------------------------------------- |
| `GET`    | `""`                  | any      | Paginated subscription list (customers only see their own)   |
| `POST`   | `""`                  | staff    | Create a subscription, optional trial/discount              |
| `GET`    | `/plans`              | any      | Active plans for assignment forms (`?activeOnly=false`)     |
| `GET`    | `/customers`          | staff    | Customer lookup for assignment (`?search=`)                 |
| `GET`    | `/dashboard`          | staff    | Frontend-ready dashboard payload                            |
| `GET`    | `/analytics`, `/metrics` | staff | Metric cards + raw analytics                                |
| `GET`    | `/reports`            | staff    | Revenue report, plan/customer breakdown, invoice aging      |
| `GET`    | `/transactions`       | staff    | Paginated settlement stream (`page`, `pageSize`, `status`, `customerId`) |
| `GET`    | `/billing/overview`   | staff    | Pending/overdue invoices and recent billing runs            |
| `GET`    | `/billing/runs`       | staff    | Billing batches derived from invoice issue dates            |
| `POST`   | `/billing/run`        | staff    | Run recurring billing (`{"asOf":"YYYY-MM-DD","dryRun":true}`) |
| `GET`    | `/:id`                | any      | Get one subscription                                         |
| `PUT`    | `/:id`                | staff    | Update plan, quantity, price, period, status, discount      |
| `DELETE` | `/:id`                | staff    | Cancel (`?immediate=true` for immediate cancellation)       |
| `POST`   | `/:id/cancel`         | staff    | Same as `DELETE`                                             |
| `POST`   | `/:id/pause`          | staff    | Pause billing                                                |
| `POST`   | `/:id/resume`         | staff    | Resume a paused subscription                                 |
| `POST`   | `/:id/renew`          | staff    | Advance one period and issue the renewal invoice            |
| `GET`    | `/:id/billing`        | any      | Past/current/upcoming billing cycles                         |

List filters: `status`, `customerId`, `planId`, `search`, `sortBy`
(`created_at`, `updated_at`, `current_period_end`, `unit_price_cents`,
`customer_name`, `status`), `sortOrder` (`asc`/`desc`), `page`, `pageSize`.

Create request body:

```json
{
  "customerId": "uuid",
  "planId": "uuid",
  "quantity": 2,
  "status": "active",
  "currentPeriodStart": "2026-11-01",
  "currentPeriodEnd": "2026-12-01",
  "discountCode": "LAUNCH20",
  "trialDays": 14
}
```

### Dashboard payload

`GET /dashboard` returns everything `FinanceDashboardPage` needs in a single
response:

```json
{
  "generatedAt": "2026-10-31T09:00:00Z",
  "currency": "USD",
  "metrics": [
    {
      "key": "mrr",
      "title": "Monthly Recurring Revenue",
      "value": "$48,250.00",
      "change": "+14.2%",
      "trend": "up",
      "subtitle": "vs. $42,240.00 last month",
      "icon": "DollarSign"
    }
  ],
  "raw": { "mrr": 48250, "mrrCents": 4825000, "activeSubscriptions": 1284 },
  "recentTransactions": [
    {
      "id": "tx_98124",
      "customer": "Acme Enterprises Inc.",
      "plan": "Enterprise Annual",
      "amount": "$14,400.00",
      "amountCents": 1440000,
      "status": "success",
      "date": "Today, 10:14 AM"
    }
  ],
  "revenueTrend": [{ "month": "2026-05", "label": "May 2026", "revenue": 42000 }],
  "statusBreakdown": [{ "status": "active", "count": 1284, "percent": 78.4 }],
  "upcomingRenewals": [
    { "subscriptionId": "uuid", "customerName": "Acme", "date": "2026-11-04", "amount": 499 }
  ]
}
```

The metric card titles are always `Monthly Recurring Revenue`, `Active
Subscriptions`, `Net Retention Rate` and `Total Customers`. `recentTransactions`
also carries the structured fields (`invoiceId`, `customerId`, `planId`,
`paymentStatus`, …) used by other views.

## Run with Docker

```bash
cd backend
cp .env.example .env
docker compose up --build
```

Then:

```bash
# Gateway health
curl http://localhost:8080/health

# Subscription health (through the gateway)
curl http://localhost:8080/api/subscriptions/health

# Frontend-compatible base URL (same gateway, port 3000)
curl http://localhost:3000/api/subscriptions/dashboard \
  -H "Authorization: Bearer mock_jwt_token_demo"

# Create a subscription (staff/admin)
curl -X POST http://localhost:3000/api/subscriptions \
  -H "Authorization: Bearer mock_token_admin" \
  -H "Content-Type: application/json" \
  -d '{"customerId":"<uuid>","planId":"<uuid>"}'

# Public auth stub
curl -X POST http://localhost:8080/api/v1/auth/login

# Protected stub (draft filter only checks for the header)
curl -H "Authorization: Bearer test" http://localhost:8080/api/v1/users/me

# HAProxy stats
# Visit http://localhost:8404/stats in a browser
```

The gateway publishes `8080` (canonical), `3000` (frontend default API URL) and
`8404` (stats); PostgreSQL publishes `5432`. Services are reachable on the
compose network at `auth:8081`, `user:8082`, `subscription:8083` and
`notification:8084`.

## Database

The stack runs `postgres:16-alpine` with a named `postgres-data` volume.
`db/schema.sql` is mounted into `/docker-entrypoint-initdb.d/`, so PostgreSQL
executes it automatically the first time the volume is created. The credentials
below are development defaults; change them before deploying anywhere shared.

```bash
# Connect with psql from the host
psql "postgres://brightloop:brightloop@localhost:5432/brightloop"

# Or through the container
docker compose exec postgres psql -U brightloop -d brightloop
```

Credentials and the connection string come from `.env` (see `.env.example`):

| Variable            | Default                                                              |
| ------------------- | -------------------------------------------------------------------- |
| `POSTGRES_DB`       | `brightloop`                                                         |
| `POSTGRES_USER`     | `brightloop`                                                         |
| `POSTGRES_PASSWORD` | `brightloop`                                                         |
| `POSTGRES_PORT`     | `5432`                                                               |
| `DATABASE_URL`      | `postgres://brightloop:brightloop@postgres:5432/brightloop?sslmode=disable` |

Services receive `DATABASE_URL` from compose and wait for the database
healthcheck before starting.

The schema covers users and auth sessions, customers, products, plans, tax
rates, discounts, quotations, subscriptions, invoices, payments, notifications
and an audit log. Money is stored as `BIGINT` minor units, percentages as
basis points and billing dates as `DATE`.

### Demo seed data

`db/seed.sql` (mounted as `02-seed.sql`) inserts a small demo dataset on a fresh
volume: one product, three plan tiers, two customers, two active subscriptions
with upcoming renewals, a paid invoice and a settled payment. This makes the
dashboard, analytics, reports and transactions endpoints return non-zero values
on first boot. Omit the mount in production to skip seeding.

### Reset the database

Init scripts only run on an empty data directory. To re-apply the schema from
scratch:

```bash
docker compose down -v
docker compose up -d postgres

# or
make db-reset
```

## Run a service locally

Each service is its own Go module:

```bash
cd services/auth
go mod tidy
go run .
```

For the subscription service, point `DATABASE_URL` at a running PostgreSQL and
run the same commands from `services/subscription`:

```bash
cd services/subscription
go run .
```

The port defaults to the service port; override it with `PORT`:

```bash
# macOS / Linux
PORT=8091 go run .

# PowerShell
$env:PORT = "8091"; go run .
```

Checks for the subscription module:

```bash
cd services/subscription
gofmt -l .
go vet ./...
go test ./...
go build ./...
```

## Endpoints

The subscription service is implemented; see
[Subscription service](#subscription-service) for its full route table and
payloads. Auth, user and notification remain drafts and return
`501 Not Implemented`.

### Auth (`/api/v1/auth`)

| Method | Path             | Responsibility              |
| ------ | ---------------- | --------------------------- |
| POST   | `/register`      | Create an account           |
| POST   | `/login`         | Authenticate and issue JWT  |
| POST   | `/logout`        | Revoke the session          |
| POST   | `/password-reset`| Send a reset link           |
| GET    | `/validate`      | Validate a token            |

### User (`/api/v1/users`)

| Method | Path               | Responsibility          |
| ------ | ------------------ | ----------------------- |
| GET    | `/me`              | Get profile             |
| PUT    | `/me`              | Update profile          |
| GET    | `/me/preferences`  | Get preferences         |
| PUT    | `/me/preferences`  | Update preferences      |
| DELETE | `/me`              | Delete account          |
| GET    | `""`               | List users (admin)      |
| GET    | `/:id`             | Get user (admin)        |

### Subscription (`/api/v1/subscriptions`)

Implemented. The table in [Subscription service](#subscription-service) is the
authoritative reference; `/api/subscriptions/...` serves the same routes.

### Notification (`/api/v1/notifications`)

| Method | Path          | Responsibility               |
| ------ | ------------- | ---------------------------- |
| GET    | `""`          | List in-app notifications    |
| POST   | `/email`      | Send email notification      |
| POST   | `/push`       | Send push notification       |
| GET    | `/templates`  | List templates               |
| POST   | `/templates`  | Create template              |

## Next steps

1. Implement auth: bcrypt hashes, JWT issuing/validation, refresh/revocation,
   password reset tokens and RBAC claims.
2. Implement the user and notification services.
3. Replace the draft gateway bearer-presence filter with real JWT validation
   (forward-auth/SPOE or an auth-service call).
4. Add service-to-service calls (subscription → notification, user → auth).
5. Add per-service migrations on top of the shared `db/schema.sql` and
   PostgreSQL-backed integration tests for the subscription repositories.
6. Publish OpenAPI specs per service.
