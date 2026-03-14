

ALTER TABLE blogs
  ADD COLUMN IF NOT EXISTS excerpt TEXT;

-- Index for potential full-text search on excerpt
CREATE INDEX IF NOT EXISTS idx_blogs_excerpt ON blogs USING gin(to_tsvector('english', coalesce(excerpt, '')));