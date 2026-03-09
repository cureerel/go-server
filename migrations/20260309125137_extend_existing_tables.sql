-- migrations/002_extend_existing_tables.sql
-- Extends users, blogs, orders, payments with new fields.
-- Safe to run on existing data — all new columns have defaults or are nullable.

-- ── Users: extend role check + add new fields ─────────────────
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
  CHECK (role IN ('user', 'writer', 'partner', 'worker', 'admin', 'superadmin'));

ALTER TABLE users
  ADD COLUMN IF NOT EXISTS is_verified          BOOLEAN     NOT NULL DEFAULT false,
  ADD COLUMN IF NOT EXISTS upgrade_requested_at TIMESTAMPTZ;

-- Existing users are considered verified (they were created before OTP flow)
UPDATE users SET is_verified = true WHERE is_verified = false;

-- ── Blogs: cover image + view counter ────────────────────────
ALTER TABLE blogs
  ADD COLUMN IF NOT EXISTS cover_image_url TEXT,
  ADD COLUMN IF NOT EXISTS cover_image_key TEXT,
  ADD COLUMN IF NOT EXISTS views_total     BIGINT NOT NULL DEFAULT 0;



-- ── Payments: provider IDs + refund tracking ─────────────────
ALTER TABLE payments
  ADD COLUMN IF NOT EXISTS provider_payment_id TEXT,
  ADD COLUMN IF NOT EXISTS refund_id           TEXT,
  ADD COLUMN IF NOT EXISTS refunded_at         TIMESTAMPTZ;