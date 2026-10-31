-- Brightloop demo seed data.
--
-- This file runs after db/schema.sql on a fresh data volume
-- (/docker-entrypoint-initdb.d/02-seed.sql). It inserts a small, coherent
-- dataset so the subscription service dashboard, analytics and transactions
-- endpoints return meaningful values on first boot. Safe to skip in production.
--
-- Idempotent: guarded by ON CONFLICT DO NOTHING where a natural key exists.

BEGIN;

-- Reference data: one product, three tiers, two customers.
INSERT INTO products (id, name, description, active, features)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Brightloop Platform',
    'Subscription management SaaS platform',
    TRUE,
    ARRAY['Recurring billing', 'Customer portal', 'Revenue reports']
) ON CONFLICT DO NOTHING;

INSERT INTO plans (id, product_id, name, description, billing_interval,
                   price_cents, currency, features, trial_period_days, active)
VALUES
    ('00000000-0000-0000-0000-000000000011', '00000000-0000-0000-0000-000000000001',
     'Growth Tier', 'Monthly growth plan', 'monthly', 49900, 'USD',
     '{}', 14, TRUE),
    ('00000000-0000-0000-0000-000000000012', '00000000-0000-0000-0000-000000000001',
     'Scale Tier', 'Quarterly scale plan', 'quarterly', 149000, 'USD',
     '{}', 0, TRUE),
    ('00000000-0000-0000-0000-000000000013', '00000000-0000-0000-0000-000000000001',
     'Enterprise Annual', 'Annual enterprise plan', 'yearly', 1440000, 'USD',
     '{}', 0, TRUE)
ON CONFLICT (product_id, name) DO NOTHING;

INSERT INTO customers (id, user_id, name, email, company_name, currency, status)
VALUES
    ('00000000-0000-0000-0000-000000000021', NULL,
     'Acme Enterprises Inc.', 'billing@acme-enterprises.com',
     'Acme Enterprises Inc.', 'USD', 'active'),
    ('00000000-0000-0000-0000-000000000022', NULL,
     'Vertex Technologies', 'finance@vertex-tech.io',
     'Vertex Technologies', 'USD', 'active')
ON CONFLICT (id) DO NOTHING;

-- Two subscriptions with periods ending soon so renewals appear on the
-- dashboard and the recurring billing run has work to do.
INSERT INTO subscriptions (
    id, customer_id, plan_id, status, quantity, unit_price_cents,
    discount_cents, currency, current_period_start, current_period_end,
    cancel_at_period_end, started_at
)
VALUES
    ('00000000-0000-0000-0000-000000000031', '00000000-0000-0000-0000-000000000021',
     '00000000-0000-0000-0000-000000000013', 'active', 1, 1440000,
     0, 'USD', CURRENT_DATE - 350, CURRENT_DATE + 15, FALSE,
     now() - interval '350 days'),
    ('00000000-0000-0000-0000-000000000032', '00000000-0000-0000-0000-000000000022',
     '00000000-0000-0000-0000-000000000011', 'active', 1, 49900,
     0, 'USD', CURRENT_DATE - 10, CURRENT_DATE + 20, FALSE,
     now() - interval '10 days')
ON CONFLICT (id) DO NOTHING;

-- A past, paid invoice + matching payment so the transaction ledger and
-- revenue trend have at least one settled row.
INSERT INTO invoices (
    id, invoice_number, customer_id, subscription_id, status, currency,
    subtotal_cents, discount_cents, tax_cents, amount_due_cents,
    amount_paid_cents, issue_date, due_date, paid_at
)
VALUES (
    '00000000-0000-0000-0000-000000000041', 'INV-DEMO-0001',
    '00000000-0000-0000-0000-000000000021',
    '00000000-0000-0000-0000-000000000031',
    'paid', 'USD', 1440000, 0, 0, 1440000, 1440000,
    CURRENT_DATE - interval '15 days', CURRENT_DATE - interval '5 days',
    now() - interval '5 days'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO invoice_line_items (
    invoice_id, plan_id, description, quantity, unit_price_cents, amount_cents
)
VALUES (
    '00000000-0000-0000-0000-000000000041',
    '00000000-0000-0000-0000-000000000013',
    'Enterprise Annual', 1, 1440000, 1440000
)
ON CONFLICT DO NOTHING;

INSERT INTO payments (
    id, invoice_id, customer_id, amount_cents, currency, method, status,
    reference, paid_at
)
VALUES (
    '00000000-0000-0000-0000-000000000051', '00000000-0000-0000-0000-000000000041',
    '00000000-0000-0000-0000-000000000021', 1440000, 'USD', 'card', 'succeeded',
    'pay_demo_0001', now() - interval '5 days'
)
ON CONFLICT (id) DO NOTHING;

COMMIT;