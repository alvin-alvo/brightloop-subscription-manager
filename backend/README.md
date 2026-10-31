# Brightloop Backend — Microservices (Draft)

Draft microservice backend for the Brightloop subscription management platform.
Each service is an independent Go module with a basic Gin server. HAProxy is
the single public entry point.

> This is a scaffold. Every business endpoint returns `501 Not Implemented` so
> the routing and deployment can be exercised before the logic is written.

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

| Service              | Folder                  | Internal port | Public prefix           |
| -------------------- | ----------------------- | ------------- | ----------------------- |
| Auth                 | `services/auth`         | 8081          | `/api/v1/auth`          |
| User                 | `services/user`         | 8082          | `/api/v1/users`         |
| Subscription         | `services/subscription` | 8083          | `/api/v1/subscriptions` |
| Notification         | `services/notification` | 8084          | `/api/v1/notifications` |

All services share the PostgreSQL 16 instance exposed on `postgres:5432`; the
schema is defined in `db/schema.sql`.

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

# Public auth stub
curl -X POST http://localhost:8080/api/v1/auth/login

# Protected stub (draft filter only checks for the header)
curl -H "Authorization: Bearer test" http://localhost:8080/api/v1/users/me

# HAProxy stats
# Visit http://localhost:8404/stats in a browser
```

Only the gateway and PostgreSQL ports are published. Services are reachable on
the compose network at `auth:8081`, `user:8082`, `subscription:8083` and
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

The port defaults to the service port; override it with `PORT`:

```bash
# macOS / Linux
PORT=8091 go run .

# PowerShell
$env:PORT = "8091"; go run .
```

## Endpoints (draft)

All service endpoints return `501 Not Implemented` with the feature name.

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

| Method | Path             | Responsibility       |
| ------ | ---------------- | -------------------- |
| GET    | `""`             | List subscriptions   |
| POST   | ``               | Create subscription  |
| GET    | `/analytics`     | Analytics            |
| GET    | `/reports`       | Reports              |
| GET    | `/:id`           | Get subscription     |
| PUT    | `/:id`           | Update subscription  |
| DELETE | `/:id`           | Cancel subscription  |
| GET    | `/:id/billing`   | Billing cycles       |
| POST   | `/:id/renew`     | Renew subscription   |

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
2. Add per-service migrations on top of the shared `db/schema.sql`.
3. Replace the draft gateway auth filter with real JWT validation.
4. Add service-to-service calls (subscription → notification, user → auth).
5. Add request validation, OpenAPI specs and tests per service.
