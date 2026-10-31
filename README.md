# Brightloop — B2B SaaS Subscription Management Platform

Brightloop is an enterprise-grade B2B subscription and revenue management platform designed for SaaS operators to automate recurring billing, coordinate sales quotations, manage tiered product catalogs, and provide self-service customer portals with strict Role-Based Access Control (RBAC).

> **TL;DR:** Brightloop pairs a Go microservice backend (HAProxy gateway, Gin, PostgreSQL 16) with a high-performance React + TypeScript frontend following the Stripe design system (`DESIGN.md`). Launch the full stack using `docker compose up --build` in `backend/` and `npm run dev` in `frontend/`.

---

## Architecture Overview

```mermaid
flowchart TD
    subgraph Client ["Client Layer"]
        SPA["React 19 + TypeScript SPA (:5173)"]
    end

    subgraph Ingress ["API Gateway Layer"]
        Gateway["HAProxy Gateway (:8080 / :3000)"]
    end

    subgraph Services ["Go Microservices Layer"]
        AuthSvc["Auth Service (:8081)"]
        UserSvc["User Service (:8082)"]
        SubSvc["Subscription Service (:8083)"]
        NotifSvc["Notification Service (:8084)"]
    end

    subgraph Data ["Data Storage Layer"]
        DB[("PostgreSQL 16 (:5432)")]
    end

    SPA -->|"HTTP REST Requests"| Gateway
    Gateway -->|"/api/v1/auth"| AuthSvc
    Gateway -->|"/api/v1/users"| UserSvc
    Gateway -->|"/api/v1/subscriptions"| SubSvc
    Gateway -->|"/api/v1/notifications"| NotifSvc

    AuthSvc --> DB
    UserSvc --> DB
    SubSvc --> DB
    NotifSvc --> DB
```

### Network & Port Directory

| Component | Technology | Host Port | Routing Prefix / Scope |
|---|---|---|---|
| **Frontend SPA** | React 19, Vite 8, Tailwind CSS | `5173` | Browser Client (`http://localhost:5173`) |
| **API Gateway** | HAProxy 3.0 Alpine | `8080`, `3000` | Reverse proxy, rate limiting, health checks |
| **HAProxy Stats** | Built-in Telemetry | `8404` | Metrics UI (`http://localhost:8404/stats`) |
| **Subscription Service** | Go 1.23, Gin, pgx/v5 | `8083` (internal) | `/api/v1/subscriptions`, `/api/subscriptions` |
| **Auth Service** | Go 1.23, Gin, JWT | `8081` (internal) | `/api/v1/auth` |
| **User Service** | Go 1.23, Gin | `8082` (internal) | `/api/v1/users` |
| **Notification Service** | Go 1.23, Gin | `8084` (internal) | `/api/v1/notifications` |
| **Database** | PostgreSQL 16 Alpine | `5432` | Shared relational database (`brightloop`) |

---

## Role-Based Capabilities (RBAC)

The application dynamically partitions operational surfaces and data contracts across three distinct user roles:

1. **`Admin` (Full System Control)**
   - **Finance Dashboard:** Real-time MRR, net retention rate, customer lifetime value, and settlement logs.
   - **Products & Plans:** Tiered catalog creation, billing interval configurations (Monthly/Annual), and entitlement feature matrices.
   - **Discounts:** Promotional coupon generators, percentage/fixed credits, and redemption capacity ceilings.
   - **Staff Management:** Team onboarding, RBAC permission auditing, and department tracking.

2. **`Staff` (Commercial & Billing Operations)**
   - **CRM / Customers:** Customer directory, contracted account standings, and MRR run-rates.
   - **Quotations:** Custom proposal authoring, enterprise term negotiation, and automated quote-to-subscription conversion.
   - **Subscriptions Lifecycle:** Renewal schedules, upgrade/downgrade workflows, and cancellation handling.
   - **Billing Runs:** Automated recurring batch invoices, dunning retries, and ledger settlements.

3. **`Customer` (Self-Service Portal)**
   - **Portal Overview:** Account standing, active license count, next renewal date, and total balance due.
   - **My Subscriptions:** Transparent license cards displaying subscription numbers, payment terms, and contract entitlements.
   - **Invoices & Receipts:** Tax invoices with itemized tax/discount breakdowns, PDF receipts, and live **"Pay Now"** instant settlement.

---

## Quickstart

Follow these steps to run the complete environment locally:

### 1. Prerequisites
- **Node.js:** >= 20.11.0 (tested on v24.12.0) with npm >= 10.0.0
- **Docker & Docker Compose:** Docker Desktop installed and running
- **Go (optional for local non-container development):** >= 1.23

### 2. Start Backend Microservices
Run the containerized gateway, database, and Go microservices:

```powershell
cd backend
docker compose up --build
```
*Expected result:* PostgreSQL initializes schema from `db/schema.sql`, microservices start on internal ports, and HAProxy becomes available at `http://localhost:8080` (and `http://localhost:3000`).

### 3. Launch Frontend Application
In a separate terminal, install dependencies and start the Vite dev server:

```powershell
cd frontend
npm install
npm run dev
```
*Expected result:* Vite serves the React application at [http://localhost:5173](http://localhost:5173).

### 4. Verify System
- Visit [http://localhost:5173](http://localhost:5173) to access the Brightloop Dashboard.
- Use the **Quick Role Switcher** in the top header to seamlessly test `Admin`, `Staff`, and `Customer` views.
- Visit [http://localhost:8404/stats](http://localhost:8404/stats) to inspect gateway traffic and health check probes.

---

## Repository Layout

```text
brightloop-subscription-manager/
├── AGENTS.md                 # Unified constraints and guidelines for AI agents
├── LICENSE                   # GNU General Public License v3.0
├── README.md                 # Project documentation and quickstart
├── backend/                  # Go microservices backend
│   ├── db/
│   │   └── schema.sql        # Database schema (PostgreSQL 16)
│   ├── gateway/
│   │   └── haproxy.cfg       # HAProxy configuration, load balancing, rate limits
│   ├── services/
│   │   ├── auth/             # JWT authentication service
│   │   ├── notification/     # In-app and webhook alerts
│   │   ├── subscription/     # Recurring billing, analytics, and subscriptions
│   │   └── user/             # Account profile management
│   ├── docker-compose.yml    # Complete backend multi-container stack
│   ├── Makefile              # Local automation tasks
│   └── README.md             # In-depth backend architecture guide
└── frontend/                 # React + TypeScript + Vite frontend
    ├── DESIGN.md             # Stripe-inspired design system tokens
    ├── index.html            # Single-page application HTML entrypoint
    ├── package.json          # Dependencies and scripts
    ├── tailwind.config.js    # Custom Tailwind styling tokens
    ├── tsconfig.json         # TypeScript configuration
    ├── vite.config.ts        # Vite bundler configuration
    └── src/
        ├── App.tsx           # Route declarations and root providers
        ├── components/
        │   ├── auth/         # ProtectedRoute guards and role checks
        │   ├── layout/       # DashboardLayout, Sidebar, Header, DashboardSkeleton
        │   └── StatusBadge.tsx # Modular badge indicators
        ├── context/
        │   ├── AuthContext.tsx  # Session persistence and role switching
        │   └── ToastContext.tsx # Real-time notification and webhook simulator
        ├── pages/
        │   ├── admin/        # Finance, Products, Discounts, Staff views
        │   ├── customer/     # Customer Dashboard, Subscriptions, Invoices
        │   └── staff/        # CRM, Quotations, Subscriptions, Billing views
        ├── routes/
        │   └── routes.config.tsx # Dynamic typed route definitions
        ├── services/
        │   ├── api.ts              # Fetch client hitting http://localhost:3000/api
        │   └── mockCustomerData.ts # Relative Date.now() dynamic data generator
        ├── types/
        │   └── index.ts      # Core domain contracts (Subscription, Invoice, etc.)
        └── utils/
            └── validations.ts # Zod registration and authentication schemas
```

---

## Configuration & Environment Variables

### Frontend (`frontend/.env`)
| Variable | Description | Default |
|---|---|---|
| `VITE_API_URL` | Base URL for API Gateway endpoints | `http://localhost:3000/api` |

### Backend (`backend/.env` / `backend/.env.example`)
| Variable | Description | Default |
|---|---|---|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://brightloop:brightloop@postgres:5432/brightloop?sslmode=disable` |
| `JWT_SECRET` | Secret key for signing bearer tokens | Minimum 32-character secure secret |
| `GATEWAY_PORT` | HAProxy public ingress port | `8080` |
| `GATEWAY_COMPAT_PORT` | HAProxy frontend compatibility port | `3000` |

---

## Available Scripts

### Frontend (`frontend/`)
| Command | Action |
|---|---|
| `npm run dev` | Launches local Vite development server on port `5173` |
| `npm run build` | Compiles TypeScript (`tsc -b`) and generates production bundle in `dist/` |
| `npm run preview` | Previews the production build locally |

### Backend (`backend/`)
| Command | Action |
|---|---|
| `docker compose up --build` | Builds and launches all microservices, database, and gateway |
| `docker compose down -v` | Stops containers and cleans up database volumes |
| `make test` | Runs unit and integration test suites across all Go modules |

---

## Design & Engineering Conventions

- **Financial Integrity:** Monetary values are modeled as integer minor units (`*Cents` / `BIGINT`) in backend storage to avoid floating-point imprecision.
- **Stripe Design System:** Frontend UI strictly implements the token specs in `frontend/DESIGN.md`:
  - Surfaces: `#ffffff` canvas, `#f6f9fc` soft canvas, `#e3e8ee` hairline borders.
  - Palette: `#533afd` electric indigo, `#0d253d` high-contrast ink text.
  - Typography: Inter with global `ss01` font features and tabular figures (`tnum`) on financial figures.
- **Perceived Responsiveness:** Heavy dashboard and table views use React `Suspense` paired with `DashboardSkeleton` to guarantee page transitions under 2 seconds.
- **Dynamic Simulation:** No static JSON files are used; customer entities are seeded dynamically relative to `Date.now()` with interactive payment mutations.

---

## License

This project is licensed under the terms of the [GNU General Public License v3.0](LICENSE).
