

-- ── Users ─────────────────────────────────────────────────────
CREATE TABLE users (
    id            BIGSERIAL    PRIMARY KEY,
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255),
    role          VARCHAR(20)  NOT NULL DEFAULT 'user',
    is_active     BOOLEAN      NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_users_email     ON users(email);
CREATE INDEX        idx_users_deleted_at ON users(deleted_at);

-- ── Refresh Tokens ────────────────────────────────────────────
CREATE TABLE refresh_tokens (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked    BOOLEAN     NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id    ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);

-- ── Blacklisted Tokens ────────────────────────────────────────
CREATE TABLE blacklisted_tokens (
    id         BIGSERIAL    PRIMARY KEY,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_blacklisted_tokens_hash ON blacklisted_tokens(token_hash);

-- ── Blogs ─────────────────────────────────────────────────────
CREATE TABLE blogs (
    id         BIGSERIAL    PRIMARY KEY,
    title      VARCHAR(200) NOT NULL,
    slug       VARCHAR(200) NOT NULL,
    content    TEXT,
    author_id  BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status     VARCHAR(20)  NOT NULL DEFAULT 'draft',
    tags       VARCHAR(500),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_blogs_slug       ON blogs(slug);
CREATE INDEX        idx_blogs_author_id  ON blogs(author_id);
CREATE INDEX        idx_blogs_status     ON blogs(status);
CREATE INDEX        idx_blogs_deleted_at ON blogs(deleted_at);

-- ── Products ──────────────────────────────────────────────────
CREATE TABLE products (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(200) NOT NULL,
    description TEXT,
    type        VARCHAR(20)  NOT NULL CHECK (type IN ('physical', 'digital')),
    price       BIGINT       NOT NULL CHECK (price > 0),
    currency    VARCHAR(10)  NOT NULL,
    is_active   BOOLEAN      NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_type      ON products(type);
CREATE INDEX idx_products_is_active ON products(is_active);

-- ── Orders ────────────────────────────────────────────────────
CREATE TABLE orders (
    id           BIGSERIAL   PRIMARY KEY,
    user_id      BIGINT      NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status       VARCHAR(20) NOT NULL DEFAULT 'pending'
                     CHECK (status IN ('pending','confirmed','dispatched','delivered','completed','cancelled')),
    total_amount BIGINT      NOT NULL CHECK (total_amount >= 0),
    currency     VARCHAR(10) NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status  ON orders(status);

-- ── Order Items ───────────────────────────────────────────────
CREATE TABLE order_items (
    id         BIGSERIAL   PRIMARY KEY,
    order_id   BIGINT      NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT      NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    type       VARCHAR(20) NOT NULL CHECK (type IN ('physical', 'digital')),
    quantity   INT         NOT NULL CHECK (quantity > 0),
    unit_price BIGINT      NOT NULL CHECK (unit_price > 0)
);

CREATE INDEX idx_order_items_order_id   ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);

-- ── Payments ──────────────────────────────────────────────────
CREATE TABLE payments (
    id               VARCHAR(100) PRIMARY KEY,
    user_id          BIGINT       NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    order_id         VARCHAR(100) NOT NULL,
    amount           BIGINT       NOT NULL CHECK (amount > 0),
    currency         VARCHAR(10)  NOT NULL,
    status           VARCHAR(20)  NOT NULL DEFAULT 'pending'
                         CHECK (status IN ('pending','completed','failed','refunded')),
    provider         VARCHAR(50)  NOT NULL CHECK (provider IN ('stripe','razorpay')),
    provider_txn_id  VARCHAR(255),
    customer_email   VARCHAR(100) NOT NULL,
    description      TEXT,
    refunded_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_user_id  ON payments(user_id);
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_status   ON payments(status);
CREATE INDEX idx_payments_provider ON payments(provider);

-- ── Webhook Events ────────────────────────────────────────────
CREATE TABLE webhook_events (
    id         VARCHAR(100) PRIMARY KEY,
    provider   VARCHAR(50)  NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload    BYTEA,
    signature  VARCHAR(255),
    processed  BOOLEAN      NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_webhook_events_provider   ON webhook_events(provider);
CREATE INDEX idx_webhook_events_processed  ON webhook_events(processed);

-- ── Memberships ───────────────────────────────────────────────
CREATE TABLE memberships (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan       VARCHAR(20) NOT NULL CHECK (plan IN ('free','basic','pro','enterprise')),
    status     VARCHAR(20) NOT NULL DEFAULT 'active'
                   CHECK (status IN ('active','cancelled','expired')),
    starts_at  TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_memberships_user_id ON memberships(user_id);
CREATE INDEX        idx_memberships_status  ON memberships(status);