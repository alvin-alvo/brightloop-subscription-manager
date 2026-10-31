-- Users and roles
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    full_name     TEXT NOT NULL DEFAULT '',
    role          TEXT NOT NULL DEFAULT 'customer'
                  CHECK (role IN ('admin', 'staff', 'customer')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Products and variants
CREATE TABLE IF NOT EXISTS products (
    id                UUID PRIMARY KEY,
    name              TEXT NOT NULL,
    type              TEXT NOT NULL DEFAULT 'service',
    sales_price_cents BIGINT NOT NULL DEFAULT 0 CHECK (sales_price_cents >= 0),
    cost_price_cents  BIGINT NOT NULL DEFAULT 0 CHECK (cost_price_cents >= 0),
    notes             TEXT NOT NULL DEFAULT '',
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS product_variants (
    id                UUID PRIMARY KEY,
    product_id        UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name              TEXT NOT NULL,
    price_delta_cents BIGINT NOT NULL DEFAULT 0,
    UNIQUE (product_id, name)
);

-- Plans
CREATE TABLE IF NOT EXISTS plans (
    id                UUID PRIMARY KEY,
    product_id        UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    name              TEXT NOT NULL,
    price_cents       BIGINT NOT NULL CHECK (price_cents >= 0),
    billing_interval  TEXT NOT NULL
                      CHECK (billing_interval IN ('daily', 'weekly', 'monthly', 'yearly')),
    min_quantity      INTEGER NOT NULL DEFAULT 1 CHECK (min_quantity > 0),
    available_from    DATE,
    available_to      DATE,
    auto_close        BOOLEAN NOT NULL DEFAULT TRUE,
    allow_early_close BOOLEAN NOT NULL DEFAULT FALSE,
    allow_pause       BOOLEAN NOT NULL DEFAULT FALSE,
    allow_renewal     BOOLEAN NOT NULL DEFAULT FALSE,
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Taxes
CREATE TABLE IF NOT EXISTS tax_rates (
    id                UUID PRIMARY KEY,
    name              TEXT NOT NULL,
    rate_basis_points INTEGER NOT NULL CHECK (rate_basis_points >= 0),
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Discounts
CREATE TABLE IF NOT EXISTS discounts (
    id                      UUID PRIMARY KEY,
    name                    TEXT NOT NULL,
    kind                    TEXT NOT NULL CHECK (kind IN ('fixed', 'percentage')),
    fixed_amount_cents      BIGINT NOT NULL DEFAULT 0 CHECK (fixed_amount_cents >= 0),
    percentage_basis_points INTEGER NOT NULL DEFAULT 0 CHECK (percentage_basis_points >= 0),
    min_purchase_cents      BIGINT NOT NULL DEFAULT 0 CHECK (min_purchase_cents >= 0),
    min_quantity            INTEGER NOT NULL DEFAULT 0 CHECK (min_quantity >= 0),
    starts_at               TIMESTAMPTZ,
    ends_at                 TIMESTAMPTZ,
    max_uses                INTEGER CHECK (max_uses IS NULL OR max_uses >= 0),
    used_count              INTEGER NOT NULL DEFAULT 0 CHECK (used_count >= 0),
    applies_to              TEXT NOT NULL DEFAULT 'subscription'
                            CHECK (applies_to IN ('product', 'subscription')),
    product_id              UUID REFERENCES products(id) ON DELETE CASCADE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Quotation templates
CREATE TABLE IF NOT EXISTS quotation_templates (
    id         UUID PRIMARY KEY,
    name       TEXT NOT NULL,
    plan_id    UUID NOT NULL REFERENCES plans(id) ON DELETE RESTRICT,
    valid_days INTEGER NOT NULL DEFAULT 30 CHECK (valid_days > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS quotation_template_lines (
    id               UUID PRIMARY KEY,
    template_id      UUID NOT NULL REFERENCES quotation_templates(id) ON DELETE CASCADE,
    product_id       UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity         INTEGER NOT NULL CHECK (quantity > 0),
    unit_price_cents BIGINT NOT NULL CHECK (unit_price_cents >= 0)
);

-- Subscriptions
CREATE TABLE IF NOT EXISTS subscriptions (
    id            UUID PRIMARY KEY,
    number        TEXT NOT NULL UNIQUE,
    customer_id   UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    plan_id       UUID NOT NULL REFERENCES plans(id) ON DELETE RESTRICT,
    status        TEXT NOT NULL DEFAULT 'draft'
                  CHECK (status IN ('draft', 'quotation', 'confirmed', 'active', 'closed')),
    start_date    DATE NOT NULL,
    expiry_date   DATE,
    payment_terms TEXT NOT NULL DEFAULT 'immediate',
    notes         TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS subscription_lines (
    id               UUID PRIMARY KEY,
    subscription_id  UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    product_id       UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity         INTEGER NOT NULL CHECK (quantity > 0),
    unit_price_cents BIGINT NOT NULL CHECK (unit_price_cents >= 0),
    tax_rate_id      UUID REFERENCES tax_rates(id) ON DELETE SET NULL,
    amount_cents     BIGINT NOT NULL CHECK (amount_cents >= 0)
);

-- Invoices
CREATE TABLE IF NOT EXISTS invoices (
    id                   UUID PRIMARY KEY,
    number               TEXT NOT NULL UNIQUE,
    subscription_id      UUID REFERENCES subscriptions(id) ON DELETE SET NULL,
    customer_id          UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status               TEXT NOT NULL DEFAULT 'draft'
                         CHECK (status IN ('draft', 'confirmed', 'paid')),
    issue_date           DATE NOT NULL,
    due_date             DATE NOT NULL,
    subtotal_cents       BIGINT NOT NULL DEFAULT 0 CHECK (subtotal_cents >= 0),
    tax_total_cents      BIGINT NOT NULL DEFAULT 0 CHECK (tax_total_cents >= 0),
    discount_total_cents BIGINT NOT NULL DEFAULT 0 CHECK (discount_total_cents >= 0),
    total_cents          BIGINT NOT NULL DEFAULT 0 CHECK (total_cents >= 0),
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS invoice_lines (
    id                    UUID PRIMARY KEY,
    invoice_id            UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    product_id            UUID REFERENCES products(id) ON DELETE SET NULL,
    description           TEXT NOT NULL,
    quantity              INTEGER NOT NULL CHECK (quantity > 0),
    unit_price_cents      BIGINT NOT NULL CHECK (unit_price_cents >= 0),
    tax_rate_basis_points INTEGER NOT NULL DEFAULT 0 CHECK (tax_rate_basis_points >= 0),
    amount_cents          BIGINT NOT NULL CHECK (amount_cents >= 0)
);

-- Payments
CREATE TABLE IF NOT EXISTS payments (
    id           UUID PRIMARY KEY,
    invoice_id   UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    method       TEXT NOT NULL,
    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    paid_at      TIMESTAMPTZ NOT NULL,
    note         TEXT NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes for common lookups and reports
CREATE INDEX IF NOT EXISTS idx_subscriptions_customer ON subscriptions (customer_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions (status);
CREATE INDEX IF NOT EXISTS idx_subscriptions_expiry ON subscriptions (expiry_date);
CREATE INDEX IF NOT EXISTS idx_invoices_customer ON invoices (customer_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices (status);
CREATE INDEX IF NOT EXISTS idx_invoices_due_date ON invoices (due_date);
CREATE INDEX IF NOT EXISTS idx_payments_invoice ON payments (invoice_id);
CREATE INDEX IF NOT EXISTS idx_payments_paid_at ON payments (paid_at);
CREATE INDEX IF NOT EXISTS idx_discounts_window ON discounts (starts_at, ends_at);
CREATE INDEX IF NOT EXISTS idx_plans_product ON plans (product_id);
