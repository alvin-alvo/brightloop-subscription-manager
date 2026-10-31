# AGENTS.md

Guidance and strict constraints for AI coding agents working in this repository. Read this file before making changes.

---

## 1. Project Overview

**brightloop-subscription-manager** is a B2B SaaS subscription and revenue management platform for Brightloop Software. It covers accounts and RBAC, products and tiered plans, quotations, recurring subscriptions, invoices, discounts, payments, and financial reporting.

---

## 2. Architecture & Tech Stack

### Backend (`backend/`)
- **Language & Runtime:** Go 1.23+, standard-library `net/http` routing (Go 1.22+ method patterns).
- **Database:** PostgreSQL 16+, accessed via `database/sql` + `pgx/v5/stdlib`.
- **Auth:** JWT bearer tokens (`golang-jwt/jwt/v5`) and bcrypt password hashes.
- **Dependency Flow:** Strictly inward: `api` → `service` → `repository` → `domain`. `domain` must not import other internal packages.

### Frontend (`frontend/`)
- **Language & Framework:** React, TypeScript, Vite.
- **Styling & Design System:** Tailwind CSS exclusively adhering to the Stripe design guidelines in `frontend/DESIGN.md` (clean surfaces `#ffffff`/`#f6f9fc`, hairline borders `#e3e8ee`, electric indigo `#533afd`, high-contrast navy ink `#0d253d`, subtle shadows, tabular figures `tnum`, and `ss01` font features). No custom CSS files.
- **State & Data Fetching:** TanStack Query hitting the local backend (`http://localhost:3000/api/...`). **No static JSON files for data fetching.** If a backend endpoint is not yet available, mock the API layer using simulated async delays (`setTimeout`), structured identically to real API responses.
- **Typing:** Strict TypeScript without `any`. Standard SaaS data contracts defined for `User`, `Product`, `Plan`, `Subscription`, and `Invoice`.
- **Modularity & Performance:** Micro-components (`StatCard`, `StatusBadge`, `DataTable`). Code-split route views wrapped in React `Suspense` with `DashboardSkeleton` to guarantee perceived loads under 2 seconds.

---

## 3. Business Logic, Validation & RBAC

### Authentication & Passwords
- Strict registration validation enforced via Zod (`frontend/src/utils/validations.ts`).
- Passwords **must be > 8 characters** and contain at least one uppercase letter, one lowercase letter, and one special character.

### Role-Based Access Control (RBAC)
The platform isolates capabilities across three distinct roles:
1. **`Admin`**: Full system control — financial dashboards, staff management, product/tier configuration, and system-wide discounts.
2. **`Staff`**: Operational workflows — customer CRM, quotation proposals, active subscription lifecycle, and recurring billing runs.
3. **`Customer`**: Self-service portal — viewing contracted subscriptions, entitlements, and downloading tax invoices.

Protected routes in `react-router-dom` dynamically verify role permissions before rendering layout and views.

---

## 4. Repository Layout

```text
.
├── AGENTS.md               # Unified instructions and constraints for AI agents
├── LICENSE                 # GNU General Public License v3.0
├── README.md               # Project documentation and quickstart
├── backend/                # Go backend microservice
│   ├── cmd/api/            # HTTP server entrypoint
│   ├── cmd/migrate/        # SQL migration runner
│   ├── internal/
│   │   ├── api/            # Router, handlers, middleware, response helpers
│   │   ├── auth/           # JWT and bcrypt helpers
│   │   ├── config/         # Environment configuration
│   │   ├── database/       # PostgreSQL connection pool
│   │   ├── domain/         # Entities, enums, sentinel errors
│   │   ├── repository/     # Repository interfaces + postgres/ implementations
│   │   ├── service/        # Business logic rules
│   │   └── validation/     # Input validation
│   ├── migrations/         # Numbered .up.sql / .down.sql files
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── go.mod
└── frontend/               # React + TypeScript + Vite frontend
    ├── index.html          # Application entry HTML
    ├── package.json        # Dependencies and scripts
    ├── tailwind.config.js  # Stripe design system tokens
    ├── tsconfig.json       # TypeScript configuration
    ├── vite.config.ts      # Vite bundler config (uses import.meta.dirname)
    ├── DESIGN.md           # Visual design specifications
    └── src/
        ├── components/     # Layout, micro-components, auth guards
        ├── context/        # Auth and Toast notification providers
        ├── hooks/          # React hooks (useAuth, etc.)
        ├── pages/          # Admin, Staff, and Customer view pages
        ├── routes/         # Dynamic TypeScript routing definitions
        ├── services/       # Local API client
        ├── types/          # Domain TypeScript interfaces
        └── utils/          # Zod validation schemas
```

---

## 5. Development Commands

### Backend (`backend/`)
| Task | Command |
|---|---|
| Run API | `go run ./cmd/api` |
| Run migrations | `go run ./cmd/migrate` |
| Build | `go build ./...` |
| Test | `go test ./...` |
| Format & Vet | `gofmt -w . && go vet ./...` |
| Docker Compose | `docker compose up --build` |

### Frontend (`frontend/`)
| Task | Command |
|---|---|
| Install | `npm install` |
| Dev server | `npm run dev` |
| Type check | `npx tsc -b` |
| Production build | `npm run build` |
| Preview build | `npm run preview` |

---

## 6. Coding Conventions

- **Money:** Store integer minor units (cents) in `*Cents` fields / `BIGINT` columns. Never use floating-point types for currency.
- **Percentages:** Store basis points in `*BasisPoints` fields (e.g., `2000` = 20.00%).
- **Dates & Times:** Use `DATE` for billing dates and `TIMESTAMPTZ` for audit timestamps.
- **Backend Errors:** Repositories return domain sentinel errors (`ErrNotFound`, `ErrConflict`, `ErrUnauthorized`, `ErrForbidden`, `ErrValidation`); `response.Error` maps them to standard HTTP status codes.
- **Database Access:** SQL queries live exclusively in `internal/repository/postgres`. Handlers and services must not write SQL.
- **Self-Contained Architecture:** Do not hallucinate external third-party payment gateways (like external Stripe APIs). The system is fully self-contained using the local backend.

---

## 7. Strict Prohibitions ("Do Not")

- Do not commit `.env` files, secrets, generated binaries, build folders (`dist/`), or dependency caches (`node_modules/`).
- Do not use floating-point types for financial calculations.
- Do not write SQL inside handlers or business services.
- Do not use static JSON files for data fetching — use TanStack Query with real or simulated local API endpoints.
- Do not use TypeScript `any`.
- Do not rewrite git history or force-push shared branches.
- Do not modify `LICENSE` unless explicitly requested.
