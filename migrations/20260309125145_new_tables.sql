
CREATE TABLE IF NOT EXISTS otps (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      REFERENCES users(id) ON DELETE CASCADE,
    email      VARCHAR(100) NOT NULL,
    code       VARCHAR(6)  NOT NULL,
    type       VARCHAR(20) NOT NULL CHECK (type IN ('register', 'reset')),
    used       BOOLEAN     NOT NULL DEFAULT false,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_otps_email      ON otps(email);
CREATE INDEX idx_otps_type       ON otps(type);
CREATE INDEX idx_otps_expires_at ON otps(expires_at);

-- ── Services (partner listings) ───────────────────────────────
CREATE TABLE IF NOT EXISTS services (
    id               BIGSERIAL      PRIMARY KEY,
    owner_id         BIGINT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title            VARCHAR(200)   NOT NULL,
    description      TEXT,
    price_usd_cents  BIGINT         NOT NULL CHECK (price_usd_cents > 0), -- stored in cents
    status           VARCHAR(20)    NOT NULL DEFAULT 'pending'
                         CHECK (status IN ('pending', 'approved', 'rejected', 'live', 'paused')),
    cover_image_url  TEXT,
    cover_image_key  TEXT,
    views_total      BIGINT         NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_services_owner_id  ON services(owner_id);
CREATE INDEX idx_services_status    ON services(status);
CREATE INDEX idx_services_deleted_at ON services(deleted_at);

-- ── Coupons ───────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS coupons (
    id                  BIGSERIAL      PRIMARY KEY,
    creator_id          BIGINT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code                VARCHAR(32)    NOT NULL,
    type                VARCHAR(20)    NOT NULL CHECK (type IN ('discount', 'affiliate', 'both')),
    discount_usd_cents  BIGINT         NOT NULL DEFAULT 0,   -- flat $ off in cents
    max_discount_cents  BIGINT         NOT NULL DEFAULT 10000, -- cap in cents (default $100)
    commission_pct      DECIMAL(5,2)   NOT NULL DEFAULT 0,   -- affiliate % of sale
    status              VARCHAR(20)    NOT NULL DEFAULT 'pending'
                            CHECK (status IN ('pending', 'approved', 'rejected')),
    usage_limit         INT,           -- NULL = unlimited
    used_count          INT            NOT NULL DEFAULT 0,
    expires_at          TIMESTAMPTZ,   -- NULL = never expires
    approved_by         BIGINT         REFERENCES users(id) ON DELETE SET NULL,
    approved_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_coupons_code      ON coupons(code);
CREATE INDEX        idx_coupons_creator   ON coupons(creator_id);
CREATE INDEX        idx_coupons_status    ON coupons(status);

-- ── Coupon usages ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS coupon_usages (
    id                    BIGSERIAL   PRIMARY KEY,
    coupon_id             BIGINT      NOT NULL REFERENCES coupons(id) ON DELETE CASCADE,
    order_id              BIGINT      NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    user_id               BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    discount_applied_cents BIGINT     NOT NULL DEFAULT 0,
    commission_usd_cents  BIGINT      NOT NULL DEFAULT 0,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_coupon_usages_coupon_id ON coupon_usages(coupon_id);
CREATE INDEX idx_coupon_usages_order_id  ON coupon_usages(order_id);
CREATE INDEX idx_coupon_usages_user_id   ON coupon_usages(user_id);


ALTER TABLE orders
  ADD COLUMN IF NOT EXISTS service_id       BIGINT      REFERENCES services(id) ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS amount_usd       BIGINT      NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS payment_provider VARCHAR(20) CHECK (payment_provider IN ('stripe', 'razorpay')),
  ADD COLUMN IF NOT EXISTS coupon_id        BIGINT      REFERENCES coupons(id) ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS affiliate_id     BIGINT      REFERENCES users(id) ON DELETE SET NULL;

-- ── Payouts ───────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS payouts (
    id             BIGSERIAL      PRIMARY KEY,
    recipient_id   BIGINT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type           VARCHAR(30)    NOT NULL CHECK (type IN ('partner_sale', 'affiliate_commission')),
    amount_cents   BIGINT         NOT NULL CHECK (amount_cents > 0),
    status         VARCHAR(20)    NOT NULL DEFAULT 'pending'
                       CHECK (status IN ('pending', 'paid')),
    order_id       BIGINT         REFERENCES orders(id) ON DELETE SET NULL,
    reference      TEXT,          -- Stripe transfer ID or manual note
    paid_by        BIGINT         REFERENCES users(id) ON DELETE SET NULL,
    paid_at        TIMESTAMPTZ,
    created_at     TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payouts_recipient ON payouts(recipient_id);
CREATE INDEX idx_payouts_status    ON payouts(status);
CREATE INDEX idx_payouts_type      ON payouts(type);

-- ── Ticket messages ───────────────────────────────────────────
-- (tickets table already exists from previous migration)
CREATE TABLE IF NOT EXISTS ticket_messages (
    id         BIGSERIAL   PRIMARY KEY,
    ticket_id  BIGINT      NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    sender_id  BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ticket_messages_ticket_id ON ticket_messages(ticket_id);
CREATE INDEX idx_ticket_messages_sender_id ON ticket_messages(sender_id);

-- Extend tickets: add assigned_to + closed_at
ALTER TABLE tickets
  ADD COLUMN IF NOT EXISTS assigned_to BIGINT      REFERENCES users(id) ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS closed_at   TIMESTAMPTZ;

-- ── Blog analytics ────────────────────────────────────────────
-- visitor_hash = SHA256(ip + date + user_agent) — privacy safe, no PII stored
CREATE TABLE IF NOT EXISTS blog_views (
    id            BIGSERIAL   PRIMARY KEY,
    blog_id       BIGINT      NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    visitor_hash  VARCHAR(64) NOT NULL,
    viewed_date   DATE        NOT NULL DEFAULT CURRENT_DATE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_blog_views_unique
  ON blog_views(blog_id, visitor_hash, viewed_date);
CREATE INDEX idx_blog_views_blog_id    ON blog_views(blog_id);
CREATE INDEX idx_blog_views_date       ON blog_views(viewed_date);

-- ── Service analytics ─────────────────────────────────────────
CREATE TABLE IF NOT EXISTS service_views (
    id            BIGSERIAL   PRIMARY KEY,
    service_id    BIGINT      NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    visitor_hash  VARCHAR(64) NOT NULL,
    viewed_date   DATE        NOT NULL DEFAULT CURRENT_DATE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_service_views_unique
  ON service_views(service_id, visitor_hash, viewed_date);
CREATE INDEX idx_service_views_service_id ON service_views(service_id);
CREATE INDEX idx_service_views_date       ON service_views(viewed_date);

-- ── Partner upgrade requests ──────────────────────────────────
CREATE TABLE IF NOT EXISTS upgrade_requests (
    id             BIGSERIAL   PRIMARY KEY,
    user_id        BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    from_role   VARCHAR(20) NOT NULL,
    to_role VARCHAR(20) NOT NULL DEFAULT 'partner',
    status         VARCHAR(20) NOT NULL DEFAULT 'pending'
                       CHECK (status IN ('pending', 'approved', 'rejected')),
    reviewed_by    BIGINT      REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at    TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_upgrade_requests_user_id ON upgrade_requests(user_id);
CREATE INDEX idx_upgrade_requests_status  ON upgrade_requests(status);


-- ── Saved blogs (already exists, kept for reference) ──────────
-- saved_blogs table was created in previous migration, no changes needed.