-- migrations/005_orders_payments.sql
-- Orders and Payments tables for Phase 6.
-- All money stored in integer cents — no floats.

-- ── Orders ────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS orders (
  id               BIGSERIAL    PRIMARY KEY,
  user_id          BIGINT       NOT NULL REFERENCES users(id),
  service_id       BIGINT       REFERENCES services(id) ON DELETE SET NULL,
  status           VARCHAR(20)  NOT NULL DEFAULT 'pending',
  total_cents      BIGINT       NOT NULL CHECK (total_cents >= 0),
  currency         VARCHAR(10)  NOT NULL DEFAULT 'USD',
  payment_provider VARCHAR(20),
  coupon_id        BIGINT       REFERENCES coupons(id) ON DELETE SET NULL,
  affiliate_id     BIGINT       REFERENCES users(id)   ON DELETE SET NULL,
  created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_orders_user_id    ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status     ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_service_id ON orders(service_id);

-- ── Order items ───────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS order_items (
  id         BIGSERIAL    PRIMARY KEY,
  order_id   BIGINT       NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  service_id BIGINT       REFERENCES services(id) ON DELETE SET NULL,
  title      VARCHAR(200) NOT NULL,
  quantity   INT          NOT NULL DEFAULT 1 CHECK (quantity > 0),
  unit_price BIGINT       NOT NULL CHECK (unit_price >= 0),
  created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);

-- ── Payments ──────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS payments (
  id                  VARCHAR(64)  PRIMARY KEY,
  order_id            BIGINT       NOT NULL REFERENCES orders(id),
  user_id             BIGINT       NOT NULL REFERENCES users(id),
  amount_cents        BIGINT       NOT NULL CHECK (amount_cents >= 0),
  currency            VARCHAR(10)  NOT NULL DEFAULT 'USD',
  status              VARCHAR(20)  NOT NULL DEFAULT 'pending',
  provider            VARCHAR(20)  NOT NULL,
  provider_txn_id     VARCHAR(128),
  provider_payment_id VARCHAR(128),
  customer_email      VARCHAR(100),
  description         TEXT,
  refund_id           VARCHAR(128),
  refunded_at         TIMESTAMPTZ,
  created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_user_id  ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_status   ON payments(status);