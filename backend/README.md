# Brightloop Backend — Microservices

Microservice backend for the Brightloop subscription management platform. Each
service is an independent Go module with a Gin server. HAProxy is the single
public entry point.

> The **subscription**, **auth**, **user** and **notification** services are
> fully implemented.
> - **Auth** — register, login, logout, token validation and password reset
>   with bcrypt hashing and signed JWTs.
> - **User** — profile & preferences, account deletion, and staff management
>   (list, invite, role/status updates).
> - **Notification** — in-app notification list/badge, mark read / read-all,
>   send (in-app/email/push) and templates.
> - **Subscription** — lifecycle, recurring billing, invoices, renewals,
>   transactions, dashboard analytics and revenue reports.
>
> The frontend staff-management page and the header notification bell call the
> user and notification services (with local fallbacks when the backend is not
> running). Only the notification service's email/push delivery integrations
> and per-service migrations remain future work.

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

## Auth service

The auth service implements accounts, session authentication and password
flows with a strict inward dependency chain identical to the subscription
service: `api → service → repository → domain`, with JWT/bcrypt helpers in
`internal/auth`.

```
services/auth/
├── main.go                  # wiring: config, database, token manager, router
└── internal/
    ├── api/                 # Gin router, handlers, DTOs, CORS, response helpers
    ├── auth/                # JWT issuing/verification, bcrypt hashing, token helpers
    ├── config/              # environment configuration
    ├── database/            # pgx-backed database/sql pool
    ├── domain/              # user/session/reset entities, enums, sentinel errors
    ├── repository/          # interfaces + postgres/ implementations
    └── service/             # register, login, validate, logout, reset business rules
```

### Configuration

| Variable           | Default       | Purpose                                              |
| ------------------ | ------------- | ---------------------------------------------------- |
| `PORT`             | `8081`        | Internal HTTP port                                   |
| `DATABASE_URL`     | `postgres://brightloop:brightloop@localhost:5432/brightloop?sslmode=disable` | PostgreSQL connection string |
| `JWT_SECRET`       | empty         | Signing secret; falls back to a stable dev secret    |
| `JWT_ISSUER`       | `brightloop-auth` | Issuer claim on issued tokens                    |
| `ACCESS_TOKEN_TTL` | `168h`        | Access-token lifetime                                |
| `REFRESH_TOKEN_TTL`| `720h`        | Session/refresh lifetime                             |
| `PASSWORD_RESET_TTL` | `30m`       | Password-reset code lifetime                         |
| `CORS_ALLOWED_ORIGINS` | `*`       | Comma-separated allowed browser origins              |

The auth and subscription services share `JWT_SECRET`; when enabled, the
subscription middleware verifies tokens that auth issued, so RBAC claims carry
across services.

### Endpoints

| Method | Path                         | Responsibility                                      |
| ------ | ---------------------------- | --------------------------------------------------- |
| `POST` | `/api/v1/auth/register`      | Create a customer account, return `{ user, token }` |
| `POST` | `/api/v1/auth/login`         | Verify credentials, return `{ user, token }`        |
| `POST` | `/api/v1/auth/logout`        | Revoke the active session (`Authorization` bearer)  |
| `GET`  | `/api/v1/auth/validate`      | Verify a bearer token, return the active user       |
| `POST` | `/api/v1/auth/password-reset`| Issue a reset code; always returns generic success  |

Register and login return the `AuthSessionResponse` shape consumed by the
frontend:

```json
{
  "user": {
    "id": "usr_admin_01",
    "name": "Alex Vance",
    "email": "alex.admin@brightloop.io",
    "role": "admin",
    "status": "active",
    "companyName": "Brightloop HQ",
    "createdAt": "2024-01-15T08:00:00Z",
    "updatedAt": "2024-10-31T09:00:00Z"
  },
  "token": "<signed JWT>"
}
```

Behavior notes:

- **Passwords** are hashed with bcrypt; the same policy as the frontend is
  enforced server-side (>8 chars, uppercase, lowercase, special character).
- **Self-registration** is limited to `customer`; admin/staff roles are created
  by staff management, not public registration.
- **Login** returns a generic "Invalid email or password" message to avoid user
  enumeration; unknown emails and wrong passwords are indistinguishable.
- **Reset** always returns the same success message whether or not an account
  exists (anti-enumeration).
- **Logout/validate** check the session row, so a revoked token stops
  validating even before it expires.
- JWTs carry `sub`, `user_id`, `name`, `email`, `role` and `sid` (session id)
  claims, which the subscription service middleware already understands.

Seeded demo accounts (password `Brightloop2026!`) match the Quick Demo Logins
on the login page: admin, staff and customer.

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

## User service

The user service manages profiles, preferences and staff administration with
the same layered layout as the other services
(`api → service → repository → domain`).

### Endpoints

| Method | Path                           | Access   | Responsibility                                |
| ------ | ------------------------------ | -------- | --------------------------------------------- |
| `GET`  | `/api/v1/users/me`             | self     | Current user profile (`User` shape)           |
| `PUT`  | `/api/v1/users/me`             | self     | Update name, company, phone, avatar           |
| `DELETE` | `/api/v1/users/me`           | self     | Deactivate the account                        |
| `GET`  | `/api/v1/users/me/preferences` | self     | Read preferences (JSON object)                |
| `PUT`  | `/api/v1/users/me/preferences` | self     | Merge preferences                            |
| `GET`  | `/api/v1/users`                | staff    | List internal members (`AdminStaff` shape)    |
| `POST` | `/api/v1/users`                | staff    | Invite a staff/admin member                   |
| `GET`  | `/api/v1/users/:id`            | staff    | Get one member                                |
| `PUT`  | `/api/v1/users/:id`            | staff    | Update role / status / department / name      |

The staff list returns the `AdminStaff` shape the frontend settings page
renders (`id`, `name`, `email`, `role`, `department`, `status`, `lastActive`,
`createdAt`). `lastActive` is computed server-side as a relative string
("Just now", "2h ago", …) and invited users report "Just invited".

Inviting a member requires `name`, `email`, `role` (`admin`/`staff`) and an
optional `department`; the service generates the account with status `invited`
and a temporary password (the invitee uses password reset to set their own).

```
POST /api/v1/users
{ "name": "Taylor Evans", "email": "taylor@brightloop.io", "role": "staff", "department": "Ops" }
```

## Notification service

The notification service powers the in-app bell, sends and templates.

### Endpoints

| Method | Path                                     | Access   | Responsibility                  |
| ------ | ---------------------------------------- | -------- | ------------------------------- |
| `GET`  | `/api/v1/notifications`                  | self     | List my notifications + unread count (`?limit`, `?unreadOnly`) |
| `GET`  | `/api/v1/notifications/unread-count`     | self     | Unread badge count              |
| `POST` | `/api/v1/notifications/:id/read`         | self     | Mark one notification read      |
| `POST` | `/api/v1/notifications/read-all`         | self     | Mark all read                   |
| `POST` | `/api/v1/notifications`                  | staff    | Create an in-app notification   |
| `POST` | `/api/v1/notifications/email`            | staff    | Record an email notification    |
| `POST` | `/api/v1/notifications/push`             | staff    | Record a push notification      |
| `GET`  | `/api/v1/notifications/templates`        | staff    | List templates                  |
| `POST` | `/api/v1/notifications/templates`        | staff    | Create a template               |

The bell list returns `{ "items": [...], "unreadCount": n }`; each item is
`{ id, type, title, body, read, createdAt, displayTime }` where `displayTime`
is a relative string. In development both the list and the target user default
to the seeded admin account via `AUTH_DEFAULT_USER_ID` when no bearer identity
is present, so the bell is populated out of the box.

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

# Sign in with a seeded demo account
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alex.admin@brightloop.io","password":"Brightloop2026!"}'

# Register a new customer account
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@company.com","password":"StrongPass1!"}'

# Validate a token returned by login/register
curl http://localhost:8080/api/v1/auth/validate \
  -H "Authorization: Bearer <token>"

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

All four services — subscription, auth, user and notification — are
implemented; see their dedicated sections for the authoritative route tables
and payloads:

- [Auth service](#auth-service)
- [Subscription service](#subscription-service)
- [User service](#user-service)
- [Notification service](#notification-service)

### Auth (`/api/v1/auth`)

Implemented. Routes are `POST /register`, `POST /login`, `POST /logout`,
`GET /validate` and `POST /password-reset`; see
[Auth service](#auth-service) for the authoritative table.

### User (`/api/v1/users`)

Implemented. Profile and preferences (`/me`, `/me/preferences`, `DELETE /me`)
and staff administration (list, invite, `/:id`); see
[User service](#user-service) for the authoritative table.

### Subscription (`/api/v1/subscriptions`)

Implemented. The table in [Subscription service](#subscription-service) is the
authoritative reference; `/api/subscriptions/...` serves the same routes.

### Notification (`/api/v1/notifications`)

Implemented. List/badge, mark read / read-all, send (in-app/email/push) and
templates; see [Notification service](#notification-service) for the
authoritative table.

## Next steps

1. Lay on real email/push delivery for the notification service.
2. Replace the draft gateway bearer-presence filter with real JWT validation
   (forward-auth/SPOE or a call to the auth `/validate` endpoint).
3. Add service-to-service calls (subscription → notification, user → auth).
4. Add per-service migrations on top of the shared `db/schema.sql` and
   PostgreSQL-backed integration tests for all four service repositories.
5. Publish OpenAPI specs per service.
