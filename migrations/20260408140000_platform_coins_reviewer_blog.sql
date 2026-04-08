-- Platform: reviewer role, coin wallet, blog workflow, post access, payment providers
-- Safe to apply on existing DBs (IF NOT EXISTS / DROP CONSTRAINT patterns).

-- ── Users: reviewer role ────────────────────────────────────────
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
  CHECK (role IN ('user', 'writer', 'reviewer', 'partner', 'worker', 'admin', 'superadmin'));

-- ── Payments: extra providers (dodpayments, coins internal) ───────
ALTER TABLE payments DROP CONSTRAINT IF EXISTS payments_provider_check;
ALTER TABLE payments ADD CONSTRAINT payments_provider_check
  CHECK (provider IN ('stripe', 'razorpay', 'dodpayments', 'coins'));

-- ── Orders: coin settlements ────────────────────────────────────
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_payment_provider_check;
ALTER TABLE orders ADD CONSTRAINT orders_payment_provider_check
  CHECK (payment_provider IS NULL OR payment_provider IN ('stripe', 'razorpay', 'dodpayments', 'coins'));

-- ── Coin wallet & ledger ────────────────────────────────────────
CREATE TABLE IF NOT EXISTS user_wallets (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS coin_ledger (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    delta       BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    reason      VARCHAR(50) NOT NULL,
    ref_type    VARCHAR(50),
    ref_id      BIGINT,
    meta        JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_coin_ledger_user_id ON coin_ledger(user_id);
CREATE INDEX IF NOT EXISTS idx_coin_ledger_created ON coin_ledger(created_at DESC);

-- ── Blog: access, review workflow, co-authors ───────────────────
ALTER TABLE blogs
  ADD COLUMN IF NOT EXISTS access_type VARCHAR(20) NOT NULL DEFAULT 'free',
  ADD COLUMN IF NOT EXISTS coin_price BIGINT NOT NULL DEFAULT 0 CHECK (coin_price >= 0),
  ADD COLUMN IF NOT EXISTS submitted_for_review_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS reviewed_by_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS review_note TEXT;

ALTER TABLE blogs DROP CONSTRAINT IF EXISTS blogs_access_type_check;
ALTER TABLE blogs ADD CONSTRAINT blogs_access_type_check
  CHECK (access_type IN ('free', 'member', 'paid_coins'));

CREATE TABLE IF NOT EXISTS blog_authors (
    blog_id BIGINT NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role    VARCHAR(20) NOT NULL DEFAULT 'co_author' CHECK (role IN ('co_author')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (blog_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_blog_authors_user_id ON blog_authors(user_id);

CREATE TABLE IF NOT EXISTS blog_unlocks (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blog_id BIGINT NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    coins_spent BIGINT NOT NULL CHECK (coins_spent > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, blog_id)
);

CREATE INDEX IF NOT EXISTS idx_blog_unlocks_blog_id ON blog_unlocks(blog_id);
