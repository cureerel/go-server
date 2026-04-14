-- Refactor: single-person platform (user/partner/admin only)
-- Drops old many-to-many blog_authors, reviewer/writer roles, old order statuses
-- Adds: username, user profile fields, new blog fields, new order statuses + delivery_status


-- USERS: rename name→username, add profile fields

ALTER TABLE users RENAME COLUMN name TO username;

ALTER TABLE users
  ADD COLUMN IF NOT EXISTS first_name   VARCHAR(100),
  ADD COLUMN IF NOT EXISTS last_name    VARCHAR(100),
  ADD COLUMN IF NOT EXISTS country      VARCHAR(100),
  ADD COLUMN IF NOT EXISTS phone_number VARCHAR(50),
  ADD COLUMN IF NOT EXISTS address      TEXT;

-- Role constraint: only 3 roles
-- First remap any legacy roles to the nearest valid one
UPDATE users SET role = 'admin'   WHERE role IN ('superadmin', 'reviewer');
UPDATE users SET role = 'partner' WHERE role IN ('writer', 'worker');
-- Now safe to add constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
  CHECK (role IN ('user', 'partner', 'admin'));


-- BLOGS: replace old fields with new schema

ALTER TABLE blogs
  DROP COLUMN IF EXISTS author_id,
  DROP COLUMN IF EXISTS tags,
  DROP COLUMN IF EXISTS cover_image_url,
  DROP COLUMN IF EXISTS cover_image_key,
  DROP COLUMN IF EXISTS views_total,
  DROP COLUMN IF EXISTS submitted_for_review_at,
  DROP COLUMN IF EXISTS reviewed_by_id,
  DROP COLUMN IF EXISTS review_note;

ALTER TABLE blogs
  ADD COLUMN IF NOT EXISTS keyword      VARCHAR(500),
  ADD COLUMN IF NOT EXISTS tag          VARCHAR(500),
  ADD COLUMN IF NOT EXISTS thumbnail    TEXT,
  ADD COLUMN IF NOT EXISTS thumbnail_key TEXT,
  ADD COLUMN IF NOT EXISTS views        BIGINT NOT NULL DEFAULT 0;

-- Blog status: only draft/published/archived
ALTER TABLE blogs DROP CONSTRAINT IF EXISTS blogs_status_check;
ALTER TABLE blogs ADD CONSTRAINT blogs_status_check
  CHECK (status IN ('draft', 'published', 'archived'));

-- Drop many-to-many co-authors table
DROP TABLE IF EXISTS blog_authors;


-- ORDERS: new statuses + delivery_status

ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_status_check;
ALTER TABLE orders ADD CONSTRAINT orders_status_check
  CHECK (status IN ('in_cart', 'paid', 'refunded'));

ALTER TABLE orders
  DROP COLUMN IF EXISTS service_id,
  DROP COLUMN IF EXISTS payment_provider;

ALTER TABLE orders
  ADD COLUMN IF NOT EXISTS delivery_status VARCHAR(20) NOT NULL DEFAULT 'created',
  ADD COLUMN IF NOT EXISTS payment_id      VARCHAR(100);

ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_delivery_status_check;
ALTER TABLE orders ADD CONSTRAINT orders_delivery_status_check
  CHECK (delivery_status IN ('created', 'in_progress', 'pending', 'completed', 'review'));

-- ORDER ITEMS: add product_id
ALTER TABLE order_items
  ADD COLUMN IF NOT EXISTS product_id BIGINT REFERENCES products(id) ON DELETE SET NULL;


-- PRODUCTS: add SKU, type, stock, image_url

ALTER TABLE products
  ADD COLUMN IF NOT EXISTS sku       VARCHAR(30)  UNIQUE,
  ADD COLUMN IF NOT EXISTS type      VARCHAR(20)  NOT NULL DEFAULT 'digital',
  ADD COLUMN IF NOT EXISTS stock     INTEGER      NOT NULL DEFAULT -1,
  ADD COLUMN IF NOT EXISTS image_url TEXT;

ALTER TABLE products DROP CONSTRAINT IF EXISTS products_type_check;
ALTER TABLE products ADD CONSTRAINT products_type_check
  CHECK (type IN ('physical', 'digital'));


-- MEMBERSHIPS: remove enterprise plan

ALTER TABLE memberships DROP CONSTRAINT IF EXISTS memberships_plan_check;
ALTER TABLE memberships ADD CONSTRAINT memberships_plan_check
  CHECK (plan IN ('free', 'basic', 'pro'));


-- COIN WALLET & LEDGER (idempotent)

CREATE TABLE IF NOT EXISTS user_wallets (
    user_id    BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    balance    BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS coin_ledger (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    delta         BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    reason        VARCHAR(50) NOT NULL,
    ref_type      VARCHAR(50),
    ref_id        BIGINT,
    meta          JSONB,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_coin_ledger_user_id ON coin_ledger(user_id);

CREATE TABLE IF NOT EXISTS blog_unlocks (
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blog_id     BIGINT NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    coins_spent BIGINT NOT NULL CHECK (coins_spent > 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, blog_id)
);


-- BLOG VIEWS dedup table (idempotent)

CREATE TABLE IF NOT EXISTS blog_views (
    blog_id      BIGINT NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    visitor_hash VARCHAR(64) NOT NULL,
    viewed_date  DATE NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (blog_id, visitor_hash, viewed_date)
);
