-- Brightloop subscription management schema.
--
-- This file is mounted into the PostgreSQL container at
-- /docker-entrypoint-initdb.d/01-schema.sql and is executed automatically the
-- first time the database volume is created.
--
-- Conventions:
--   * money            -> BIGINT minor units (cents)
--   * percentages      -> INTEGER basis points (2000 = 20.00%)
--   * billing dates    -> DATE
--   * audit timestamps -> TIMESTAMPTZ

BEGIN;

-- ---------------------------------------------------------------------------
-- Enumerated types
-- ---------------------------------------------------------------------------

CREATE TYPE user_role AS ENUM ('admin', 'staff', 'customer');
CREATE TYPE user_status AS ENUM ('active', 'invited', 'inactive', 'suspended');
CREATE TYPE customer_status AS ENUM ('active', 'inactive', 'suspended');
CREATE TYPE billing_interval AS ENUM ('monthly', 'quarterly', 'yearly');
CREATE TYPE subscription_status AS ENUM (
    'active', 'trialing', 'past_due', 'canceled', 'paused', 'unpaid'
);
CREATE TYPE invoice_status AS ENUM ('draft', 'open', 'paid', 'uncollectible', 'void');
CREATE TYPE quotation_status AS ENUM ('draft', 'sent', 'accepted', 'declined', 'expired');
CREATE TYPE discount_kind AS ENUM ('percentage', 'fixed_amount');
CREATE TYPE discount_duration AS ENUM ('once', 'repeating', 'forever');
CREATE TYPE payment_method AS ENUM ('card', 'bank_transfer', 'cash', 'other');
CREATE TYPE payment_status AS ENUM ('pending', 'succeeded', 'failed', 'refunded');
CREATE TYPE notification_channel AS ENUM ('email', 'push', 'in_app');

-- ---------------------------------------------------------------------------
-- Shared trigger used to keep updated_at current
-- ---------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ---------------------------------------------------------------------------
-- Accounts and authentication
-- ---------------------------------------------------------------------------

CREATE TABLE users (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name              TEXT NOT NULL,
    email             TEXT NOT NULL,
    password_hash     TEXT NOT NULL,
    role              user_role NOT NULL DEFAULT 'customer',
    status            user_status NOT NULL DEFAULT 'active',
    company_name      TEXT,
    department        TEXT NOT NULL DEFAULT '',
    phone             TEXT,
    avatar_url        TEXT,
    preferences       JSONB NOT NULL DEFAULT '{}'::jsonb,
    email_verified_at TIMESTAMPTZ,
    last_login_at     TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX users_email_key ON users (lower(email));
CREATE INDEX users_role_idx ON users (role);
CREATE INDEX users_status_idx ON users (status);

CREATE TABLE auth_sessions (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id            UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    refresh_token_hash TEXT NOT NULL UNIQUE,
    user_agent         TEXT,
    ip_address         INET,
    expires_at         TIMESTAMPTZ NOT NULL,
    revoked_at         TIMESTAMPTZ,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX auth_sessions_user_id_idx ON auth_sessions (user_id);
CREATE INDEX auth_sessions_expires_at_idx ON auth_sessions (expires_at);

CREATE TABLE password_reset_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX password_reset_tokens_user_id_idx ON password_reset_tokens (user_id);

-- ---------------------------------------------------------------------------
-- Customers
-- ---------------------------------------------------------------------------

CREATE TABLE customers (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                UUID UNIQUE REFERENCES users (id) ON DELETE SET NULL,
    name                   TEXT NOT NULL,
    email                  TEXT NOT NULL,
    company_name           TEXT,
    phone                  TEXT,
    billing_address_line1  TEXT,
    billing_address_line2  TEXT,
    billing_city           TEXT,
    billing_state          TEXT,
    billing_postal_code    TEXT,
    billing_country        CHAR(2),
    currency               CHAR(3) NOT NULL DEFAULT 'USD',
    status                 customer_status NOT NULL DEFAULT 'active',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX customers_email_idx ON customers (lower(email));
CREATE INDEX customers_status_idx ON customers (status);

-- ---------------------------------------------------------------------------
-- Product catalogue
-- ---------------------------------------------------------------------------

CREATE TABLE products (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    sku         TEXT,
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    features    TEXT[] NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX products_name_key ON products (lower(name));
CREATE UNIQUE INDEX products_sku_key ON products (sku) WHERE sku IS NOT NULL;

CREATE TABLE plans (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id        UUID NOT NULL REFERENCES products (id) ON DELETE RESTRICT,
    name              TEXT NOT NULL,
    description       TEXT NOT NULL DEFAULT '',
    billing_interval  billing_interval NOT NULL DEFAULT 'monthly',
    price_cents       BIGINT NOT NULL DEFAULT 0 CHECK (price_cents >= 0),
    currency          CHAR(3) NOT NULL DEFAULT 'USD',
    features          TEXT[] NOT NULL DEFAULT '{}',
    trial_period_days INTEGER NOT NULL DEFAULT 0 CHECK (trial_period_days >= 0),
    active            BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT plans_product_name_key UNIQUE (product_id, name)
);

CREATE INDEX plans_product_id_idx ON plans (product_id);
CREATE INDEX plans_active_idx ON plans (active) WHERE active;

CREATE TABLE tax_rates (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name               TEXT NOT NULL,
    country_code       CHAR(2) NOT NULL,
    region             TEXT,
    rate_basis_points  INTEGER NOT NULL CHECK (rate_basis_points BETWEEN 0 AND 10000),
    inclusive          BOOLEAN NOT NULL DEFAULT FALSE,
    active             BOOLEAN NOT NULL DEFAULT TRUE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE discounts (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code                TEXT NOT NULL,
    description         TEXT NOT NULL DEFAULT '',
    kind                discount_kind NOT NULL,
    percent_basis_points INTEGER CHECK (percent_basis_points BETWEEN 0 AND 10000),
    amount_cents        BIGINT CHECK (amount_cents >= 0),
    currency            CHAR(3),
    duration            discount_duration NOT NULL DEFAULT 'once',
    duration_months     INTEGER CHECK (duration_months IS NULL OR duration_months > 0),
    max_redemptions     INTEGER CHECK (max_redemptions IS NULL OR max_redemptions > 0),
    times_redeemed      INTEGER NOT NULL DEFAULT 0 CHECK (times_redeemed >= 0),
    starts_at           TIMESTAMPTZ,
    ends_at             TIMESTAMPTZ,
    active              BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT discounts_value_check CHECK (
        (kind = 'percentage' AND percent_basis_points IS NOT NULL AND amount_cents IS NULL)
        OR (kind = 'fixed_amount' AND amount_cents IS NOT NULL AND percent_basis_points IS NULL)
    ),
    CONSTRAINT discounts_window_check CHECK (
        ends_at IS NULL OR starts_at IS NULL OR ends_at > starts_at
    )
);

CREATE UNIQUE INDEX discounts_code_key ON discounts (lower(code));

-- ---------------------------------------------------------------------------
-- Quotations
-- ---------------------------------------------------------------------------

CREATE TABLE quotations (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quote_number   TEXT NOT NULL UNIQUE,
    customer_id    UUID NOT NULL REFERENCES customers (id) ON DELETE RESTRICT,
    status         quotation_status NOT NULL DEFAULT 'draft',
    currency       CHAR(3) NOT NULL DEFAULT 'USD',
    subtotal_cents BIGINT NOT NULL DEFAULT 0 CHECK (subtotal_cents >= 0),
    discount_cents BIGINT NOT NULL DEFAULT 0 CHECK (discount_cents >= 0),
    tax_cents      BIGINT NOT NULL DEFAULT 0 CHECK (tax_cents >= 0),
    total_cents    BIGINT NOT NULL DEFAULT 0 CHECK (total_cents >= 0),
    valid_until    DATE,
    notes          TEXT,
    created_by     UUID REFERENCES users (id) ON DELETE SET NULL,
    accepted_at    TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX quotations_customer_id_idx ON quotations (customer_id);
CREATE INDEX quotations_status_idx ON quotations (status);

CREATE TABLE quotation_items (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quotation_id           UUID NOT NULL REFERENCES quotations (id) ON DELETE CASCADE,
    plan_id                UUID REFERENCES plans (id) ON DELETE SET NULL,
    description            TEXT NOT NULL,
    quantity               INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price_cents       BIGINT NOT NULL CHECK (unit_price_cents >= 0),
    discount_basis_points  INTEGER NOT NULL DEFAULT 0
                               CHECK (discount_basis_points BETWEEN 0 AND 10000),
    tax_basis_points       INTEGER NOT NULL DEFAULT 0
                               CHECK (tax_basis_points BETWEEN 0 AND 10000),
    line_total_cents       BIGINT NOT NULL CHECK (line_total_cents >= 0),
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX quotation_items_quotation_id_idx ON quotation_items (quotation_id);

-- ---------------------------------------------------------------------------
-- Subscriptions
-- ---------------------------------------------------------------------------

CREATE TABLE subscriptions (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id           UUID NOT NULL REFERENCES customers (id) ON DELETE RESTRICT,
    plan_id               UUID NOT NULL REFERENCES plans (id) ON DELETE RESTRICT,
    status                subscription_status NOT NULL DEFAULT 'trialing',
    quantity              INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price_cents      BIGINT NOT NULL CHECK (unit_price_cents >= 0),
    discount_cents        BIGINT NOT NULL DEFAULT 0 CHECK (discount_cents >= 0),
    currency              CHAR(3) NOT NULL DEFAULT 'USD',
    current_period_start  DATE NOT NULL,
    current_period_end    DATE NOT NULL,
    cancel_at_period_end  BOOLEAN NOT NULL DEFAULT FALSE,
    trial_ends_at         TIMESTAMPTZ,
    started_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    canceled_at           TIMESTAMPTZ,
    ended_at              TIMESTAMPTZ,
    created_by            UUID REFERENCES users (id) ON DELETE SET NULL,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT subscriptions_period_check CHECK (current_period_end >= current_period_start)
);

CREATE INDEX subscriptions_customer_id_idx ON subscriptions (customer_id);
CREATE INDEX subscriptions_plan_id_idx ON subscriptions (plan_id);
CREATE INDEX subscriptions_status_idx ON subscriptions (status);
CREATE INDEX subscriptions_period_end_idx ON subscriptions (current_period_end);

CREATE TABLE subscription_discounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID NOT NULL REFERENCES subscriptions (id) ON DELETE CASCADE,
    discount_id     UUID NOT NULL REFERENCES discounts (id) ON DELETE RESTRICT,
    applied_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    ends_at         TIMESTAMPTZ,
    CONSTRAINT subscription_discounts_key UNIQUE (subscription_id, discount_id)
);

CREATE INDEX subscription_discounts_discount_id_idx ON subscription_discounts (discount_id);

-- ---------------------------------------------------------------------------
-- Invoices and payments
-- ---------------------------------------------------------------------------

CREATE TABLE invoices (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number    TEXT NOT NULL UNIQUE,
    customer_id       UUID NOT NULL REFERENCES customers (id) ON DELETE RESTRICT,
    subscription_id   UUID REFERENCES subscriptions (id) ON DELETE SET NULL,
    status            invoice_status NOT NULL DEFAULT 'draft',
    currency          CHAR(3) NOT NULL DEFAULT 'USD',
    subtotal_cents    BIGINT NOT NULL DEFAULT 0 CHECK (subtotal_cents >= 0),
    discount_cents    BIGINT NOT NULL DEFAULT 0 CHECK (discount_cents >= 0),
    tax_cents         BIGINT NOT NULL DEFAULT 0 CHECK (tax_cents >= 0),
    amount_due_cents  BIGINT NOT NULL DEFAULT 0 CHECK (amount_due_cents >= 0),
    amount_paid_cents BIGINT NOT NULL DEFAULT 0 CHECK (amount_paid_cents >= 0),
    issue_date        DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date          DATE NOT NULL,
    paid_at           TIMESTAMPTZ,
    notes             TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT invoices_due_date_check CHECK (due_date >= issue_date)
);

CREATE INDEX invoices_customer_id_idx ON invoices (customer_id);
CREATE INDEX invoices_subscription_id_idx ON invoices (subscription_id);
CREATE INDEX invoices_status_idx ON invoices (status);
CREATE INDEX invoices_due_date_idx ON invoices (due_date);

CREATE TABLE invoice_line_items (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id       UUID NOT NULL REFERENCES invoices (id) ON DELETE CASCADE,
    plan_id          UUID REFERENCES plans (id) ON DELETE SET NULL,
    description      TEXT NOT NULL,
    quantity         INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price_cents BIGINT NOT NULL CHECK (unit_price_cents >= 0),
    amount_cents     BIGINT NOT NULL CHECK (amount_cents >= 0),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX invoice_line_items_invoice_id_idx ON invoice_line_items (invoice_id);

CREATE TABLE payments (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id   UUID NOT NULL REFERENCES invoices (id) ON DELETE RESTRICT,
    customer_id  UUID NOT NULL REFERENCES customers (id) ON DELETE RESTRICT,
    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    currency     CHAR(3) NOT NULL DEFAULT 'USD',
    method       payment_method NOT NULL,
    status       payment_status NOT NULL DEFAULT 'pending',
    reference    TEXT,
    paid_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX payments_invoice_id_idx ON payments (invoice_id);
CREATE INDEX payments_customer_id_idx ON payments (customer_id);
CREATE INDEX payments_status_idx ON payments (status);

-- ---------------------------------------------------------------------------
-- Notifications
-- ---------------------------------------------------------------------------

CREATE TABLE notifications (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    channel    notification_channel NOT NULL DEFAULT 'in_app',
    title      TEXT NOT NULL,
    body       TEXT NOT NULL,
    read_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX notifications_user_id_idx ON notifications (user_id, created_at DESC);
CREATE INDEX notifications_unread_idx ON notifications (user_id) WHERE read_at IS NULL;

CREATE TABLE notification_templates (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code       TEXT NOT NULL UNIQUE,
    channel    notification_channel NOT NULL,
    subject    TEXT,
    body       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------------
-- Audit trail
-- ---------------------------------------------------------------------------

CREATE TABLE audit_logs (
    id          BIGSERIAL PRIMARY KEY,
    actor_id    UUID REFERENCES users (id) ON DELETE SET NULL,
    action      TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id   UUID,
    metadata    JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX audit_logs_actor_id_idx ON audit_logs (actor_id);
CREATE INDEX audit_logs_entity_idx ON audit_logs (entity_type, entity_id);
CREATE INDEX audit_logs_created_at_idx ON audit_logs (created_at DESC);

-- ---------------------------------------------------------------------------
-- updated_at triggers
-- ---------------------------------------------------------------------------

CREATE TRIGGER users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER customers_set_updated_at
    BEFORE UPDATE ON customers
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER products_set_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER plans_set_updated_at
    BEFORE UPDATE ON plans
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER tax_rates_set_updated_at
    BEFORE UPDATE ON tax_rates
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER discounts_set_updated_at
    BEFORE UPDATE ON discounts
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER quotations_set_updated_at
    BEFORE UPDATE ON quotations
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER subscriptions_set_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER invoices_set_updated_at
    BEFORE UPDATE ON invoices
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER payments_set_updated_at
    BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER notification_templates_set_updated_at
    BEFORE UPDATE ON notification_templates
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMIT;
