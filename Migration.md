# Migration Guide

## Workflow

```bash
# 1. Create new migration file
make migrate-create NAME=your_migration_name

# 2. Hash the migrations dir after editing
atlas migrate hash --dir "file://migrations"

# 3. Apply migrations (local)
make migrate-up
```

## Latest Migration
`20260414080000_refactor_single_person.sql`

Changes:
- `users.name` → `users.username` (required)
- Added `first_name`, `last_name`, `country`, `phone_number`, `address` to users
- Role constraint: `user | partner | admin` only (removed reviewer/writer/superadmin)
- Blog: replaced `author_id`, `tags`, `cover_image_url`, `views_total`, review fields
  with `keyword`, `tag`, `thumbnail`, `views`
- Blog status: `draft | published | archived` only
- Dropped `blog_authors` join table (single-author now)
- Order status: `in_cart | paid | refunded`
- Added `delivery_status`: `created | in_progress | pending | completed | review`
- Added `payment_id` to orders, dropped `payment_provider`, `service_id`
- Added `product_id` to `order_items`
- Membership plans: `free | basic | pro` (removed enterprise)
- Created `user_wallets`, `coin_ledger`, `blog_unlocks`, `blog_views` (idempotent)


