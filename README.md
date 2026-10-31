# Brightloop - B2B SaaS Subscription Management System

Brightloop is a multi-tenant B2B subscription and billing management portal for SaaS companies to automate recurring billing, manage product catalogs, and provide self-service customer portals with role-based access control.

> **TL;DR:** Brightloop provides role-based dashboards (`Admin`, `Staff`, `Customer`) for recurring billing, subscription lifecycle operations, and customer invoicing. Get started by running `cd frontend && npm install && npm run dev`.

---

## Why Brightloop Exists

Managing modern B2B SaaS revenue operations requires granular role separation, customizable billing tiers, and immediate visibility into customer accounts:
- **Admin**: Complete management of staff users, products, pricing plans, and system-wide discounts.
- **Staff**: CRM workflows, quotation generation, and real-time billing inspections.
- **Customer**: Dedicated self-service portal to review active subscriptions and download invoices.

---

## Quickstart

Follow these 4 steps to run the frontend application locally:

### 1. Prerequisites
Ensure you have Node.js and npm installed:
- Node.js >= 20.11.0 (tested on v24.12.0)
- npm >= 10.0.0

### 2. Install dependencies
```powershell
cd frontend
npm install
```
*Expected result:* All dependencies install cleanly without peer conflicts.

### 3. Start development server
```powershell
npm run dev
```
*Expected result:* Vite launches dev server on `http://localhost:5173`.

### 4. Open in browser
Visit [http://localhost:5173](http://localhost:5173) to verify the running application.

---

## Project Structure

```text
brightloop-subscription-manager/
├── AGENT.md                 # System constraints and hackathon requirements
├── LICENSE                  # GNU General Public License v3.0
├── README.md                # Project documentation and quickstart
└── frontend/                # Vite + React + TypeScript frontend
    ├── index.html           # Single-page application entry HTML
    ├── package.json         # Frontend dependencies and run scripts
    ├── postcss.config.js    # PostCSS configuration
    ├── tailwind.config.js   # Tailwind CSS design system configuration
    ├── tsconfig.json        # TypeScript configuration root
    ├── vite.config.ts       # Vite bundler configuration
    └── src/
        ├── main.tsx         # React root with TanStack Query provider
        ├── App.tsx          # Main application router container
        ├── index.css        # Tailwind directives and base layer styles
        ├── components/      # Reusable micro-components (e.g. StatusBadge)
        ├── hooks/           # Custom React hooks (e.g. useAuth)
        ├── pages/           # Application views and dashboard layouts
        ├── services/        # API clients and data layer fetching
        ├── types/           # Core SaaS TypeScript contracts and interfaces
        └── utils/           # Validation schemas (Zod) and utility helpers
```

---

## Configuration

Frontend environment variables can be customized in `frontend/.env`:

| Variable | Description | Default |
|---|---|---|
| `VITE_API_URL` | Base URL for the local backend API service | `http://localhost:3000/api` |

---

## Usage Examples

### Validating User Registration
Registration uses Zod schemas enforcing complex passwords (8+ chars, uppercase, lowercase, special character):

```typescript
import { registerSchema } from './src/utils/validations';

const result = registerSchema.safeParse({
  name: "Jane Doe",
  email: "jane@company.com",
  password: "SecurePassword123!",
  confirmPassword: "SecurePassword123!",
  role: "customer"
});

console.log(result.success); // true
```

### Type Contracts
Standard data models are strongly typed in `src/types/index.ts`:

```typescript
import type { Subscription, Invoice, User } from './src/types';
```

---

## Available Scripts

From the `frontend` directory:

| Command | Action |
|---|---|
| `npm run dev` | Starts Vite local development server on port `5173` |
| `npm run build` | Compiles TypeScript (`tsc -b`) and generates production bundle in `dist/` |
| `npm run preview` | Previews the production build locally |

---

## Contributing

1. Create a feature branch (`git checkout -b feature/auth-flow`).
2. Run TypeScript checks before committing:
   ```powershell
   cd frontend
   npm run build
   ```
3. Commit with concise messages and submit a pull request.

---

## License

This project is licensed under the terms of the [GNU General Public License v3.0](LICENSE).
