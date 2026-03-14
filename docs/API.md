# cureerel API — Testing Guide

Base URL: `http://localhost:8080`

---

## Table of Contents

1. [Setup](#setup)
2. [Auth — Register & Login](#1-auth--register--login)
3. [Auth — Password Reset](#2-auth--password-reset)
4. [Blog — Public & Writer](#3-blog)
5. [Services — Public & Partner](#4-services)
6. [Upload Images](#5-upload-images)
7. [Orders & Payments](#6-orders--payments)
8. [Coupons](#7-coupons)
9. [Payouts](#8-payouts)
10. [Tickets & Messages](#9-tickets--messages)
11. [Dashboards](#10-dashboards)
12. [Superadmin Panel](#11-superadmin-panel)
13. [Role Hierarchy](#role-hierarchy)
14. [Postman Setup](#postman-setup)

---

## Setup

### Curl — set your token once per session

```bash
TOKEN="paste_access_token_here"
```

### Postman

1. Create a collection `cureerel`
2. Add a collection variable `baseUrl` = `http://localhost:8080`
3. Add a collection variable `token` (fill after login)
4. In collection **Authorization** tab → Bearer Token → `{{token}}`

---

## 1. Auth — Register & Login

### 1.1 Send OTP (register)

```bash
curl -X POST http://localhost:8080/api/auth/register/init \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@example.com"}'
```

**Expected:** `200 OK` — OTP sent to email (check server logs if using noop email)

---

### 1.2 Verify OTP & Create Account

```bash
curl -X POST http://localhost:8080/api/auth/register/verify \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice",
    "email": "alice@example.com",
    "password": "secret123",
    "code": "123456"
  }'
```

**Expected:** `201 Created`
```json
{
  "user": { "id": 1, "name": "Alice", "email": "alice@example.com", "role": "user" },
  "access_token": "eyJ...",
  "refresh_token": "eyJ..."
}
```

> Copy `access_token` → set `TOKEN=...` in terminal or `{{token}}` in Postman

---

### 1.3 Login (existing user)

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@example.com", "password": "secret123"}'
```

---

### 1.4 Refresh Token

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "eyJ..."}'
```

---

### 1.5 Logout

```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer $TOKEN"
```

---

### 1.6 Get My Profile

```bash
curl http://localhost:8080/api/users/me \
  -H "Authorization: Bearer $TOKEN"
```

---

## 2. Auth — Password Reset

### 2.1 Send Reset OTP

```bash
curl -X POST http://localhost:8080/api/auth/password/reset/init \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@example.com"}'
```

---

### 2.2 Verify OTP & Set New Password

```bash
curl -X POST http://localhost:8080/api/auth/password/reset/verify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "code": "123456",
    "new_password": "newpass456"
  }'
```

---

## 3. Blog

### Public — no auth required

#### Get all posts

```bash
curl "http://localhost:8080/api/blog?page=1&limit=10"
```

#### Get by slug

```bash
curl http://localhost:8080/api/blog/slug/my-first-post
```

#### Get by ID

```bash
curl http://localhost:8080/api/blog/1
```

---

### Writer+ — requires writer role

#### Create post

```bash
curl -X POST http://localhost:8080/api/blogs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Post",
    "content": "Hello world content here.",
    "published": true
  }'
```

#### Update post

```bash
curl -X PUT http://localhost:8080/api/blogs/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Title", "content": "Updated content."}'
```

#### Patch post

```bash
curl -X PATCH http://localhost:8080/api/blogs/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"published": false}'
```

#### Delete post

```bash
curl -X DELETE http://localhost:8080/api/blogs/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

## 4. Services

### Public — no auth required

#### Get all services

```bash
curl "http://localhost:8080/api/services?page=1&limit=10"
```

#### Get service by ID

```bash
curl http://localhost:8080/api/services/1
```

---

### Partner+ — requires partner role

#### Get my services

```bash
curl http://localhost:8080/api/services/mine \
  -H "Authorization: Bearer $TOKEN"
```

#### Create service

```bash
curl -X POST http://localhost:8080/api/services \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Logo Design",
    "description": "Professional logo design service.",
    "price_usd_cents": 4999
  }'
```

#### Update service

```bash
curl -X PUT http://localhost:8080/api/services/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Logo Design Pro", "price_usd_cents": 5999}'
```

#### Set service live

```bash
curl -X POST http://localhost:8080/api/services/1/live \
  -H "Authorization: Bearer $TOKEN"
```

#### Pause service

```bash
curl -X POST http://localhost:8080/api/services/1/pause \
  -H "Authorization: Bearer $TOKEN"
```

#### Delete service

```bash
curl -X DELETE http://localhost:8080/api/services/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

### Admin+ — approve / reject services

#### Approve

```bash
curl -X POST http://localhost:8080/api/services/1/approve \
  -H "Authorization: Bearer $TOKEN"
```

#### Reject

```bash
curl -X POST http://localhost:8080/api/services/1/reject \
  -H "Authorization: Bearer $TOKEN"
```

---

## 5. Upload Images

### Upload

```bash
curl -X POST http://localhost:8080/api/upload/image \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/image.png"
```

**Expected:**
```json
{ "url": "https://res.cloudinary.com/...", "key": "cureerel/abc123" }
```

### Delete

```bash
curl -X DELETE http://localhost:8080/api/upload/image \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"key": "cureerel/abc123"}'
```

---

## 6. Orders & Payments

### Create order (also creates pending payment)

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "service_id": 1,
    "provider": "stripe",
    "coupon_code": ""
  }'
```

**Expected:** `201 Created`
```json
{
  "order": { "id": 1, "status": "pending", "total_cents": 4999 },
  "payment": { "id": "stripe-1-1234567890", "status": "pending" }
}
```

> To apply a coupon, set `"coupon_code": "SAVE10"`

---

### Get my orders

```bash
curl "http://localhost:8080/api/orders/me?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

### Get order by ID

```bash
curl http://localhost:8080/api/orders/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Get payment by ID

```bash
curl http://localhost:8080/api/payments/stripe-1-1234567890 \
  -H "Authorization: Bearer $TOKEN"
```

---

### Admin — manage orders & payments

#### Get all orders

```bash
curl "http://localhost:8080/api/orders?page=1&limit=10&status=pending" \
  -H "Authorization: Bearer $TOKEN"
```

#### Update order status

```bash
curl -X PATCH http://localhost:8080/api/orders/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "cancelled"}'
```

#### Get all payments

```bash
curl "http://localhost:8080/api/payments?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

#### Mark payment complete (confirms order)

```bash
curl -X POST http://localhost:8080/api/payments/stripe-1-1234567890/complete \
  -H "Authorization: Bearer $TOKEN"
```

#### Mark payment failed (cancels order)

```bash
curl -X POST http://localhost:8080/api/payments/stripe-1-1234567890/fail \
  -H "Authorization: Bearer $TOKEN"
```

#### Refund payment (cancels order)

```bash
curl -X POST http://localhost:8080/api/payments/stripe-1-1234567890/refund \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"refund_id": "re_abc123"}'
```

---

## 7. Coupons

### Validate coupon — public, no auth

```bash
curl "http://localhost:8080/api/coupons/validate?code=SAVE10"
```

**Expected:**
```json
{
  "valid": true,
  "code": "SAVE10",
  "type": "discount",
  "discount_usd_cents": 1000,
  "commission_pct": 0
}
```

---

### Create coupon — partner+

```bash
curl -X POST http://localhost:8080/api/coupons \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "SAVE10",
    "type": "discount",
    "discount_usd_cents": 1000,
    "max_discount_cents": 1000,
    "commission_pct": 0,
    "usage_limit": 100
  }'
```

> `type` must be one of: `discount`, `affiliate`, `both`

### Create affiliate coupon

```bash
curl -X POST http://localhost:8080/api/coupons \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "PARTNER20",
    "type": "affiliate",
    "discount_usd_cents": 500,
    "max_discount_cents": 500,
    "commission_pct": 20
  }'
```

### Get coupon by ID — partner+ (own) or admin

```bash
curl http://localhost:8080/api/coupons/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

### Admin — manage coupons

#### Get all coupons

```bash
curl "http://localhost:8080/api/coupons?page=1&limit=10&status=pending" \
  -H "Authorization: Bearer $TOKEN"
```

#### Approve coupon

```bash
curl -X POST http://localhost:8080/api/coupons/1/approve \
  -H "Authorization: Bearer $TOKEN"
```

#### Reject coupon

```bash
curl -X POST http://localhost:8080/api/coupons/1/reject \
  -H "Authorization: Bearer $TOKEN"
```

---

## 8. Payouts

### Get my payouts — any auth user

```bash
curl "http://localhost:8080/api/payouts/me?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

### Admin — get all payouts

```bash
curl "http://localhost:8080/api/payouts?page=1&limit=10&status=pending" \
  -H "Authorization: Bearer $TOKEN"
```

### Admin — mark payout paid

```bash
curl -X POST http://localhost:8080/api/payouts/1/pay \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"reference": "bank-txn-abc123"}'
```

---

## 9. Tickets & Messages

### Create ticket — any auth user

```bash
curl -X POST http://localhost:8080/api/tickets \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "Cannot access my order",
    "description": "I placed an order but cannot see it in my dashboard.",
    "priority": "medium"
  }'
```

> `priority`: `low`, `medium`, `high`, `urgent`

### Get my tickets

```bash
curl "http://localhost:8080/api/tickets/me?page=1&limit=10&status=open" \
  -H "Authorization: Bearer $TOKEN"
```

### Get ticket by ID — owner or worker+

```bash
curl http://localhost:8080/api/tickets/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Send message on ticket — owner or worker+

```bash
curl -X POST http://localhost:8080/api/tickets/1/messages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"message": "I have checked and the order is under ID #42."}'
```

### Get messages on ticket — owner or worker+

```bash
curl http://localhost:8080/api/tickets/1/messages \
  -H "Authorization: Bearer $TOKEN"
```

### Close ticket — owner or admin

```bash
curl -X POST http://localhost:8080/api/tickets/1/close \
  -H "Authorization: Bearer $TOKEN"
```

---

### Worker+ — manage tickets

#### Get all tickets

```bash
curl "http://localhost:8080/api/tickets?page=1&limit=10&status=open&priority=urgent" \
  -H "Authorization: Bearer $TOKEN"
```

#### Assign ticket to worker

```bash
curl -X POST http://localhost:8080/api/tickets/1/assign \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"worker_id": 5}'
```

#### Resolve ticket

```bash
curl -X POST http://localhost:8080/api/tickets/1/resolve \
  -H "Authorization: Bearer $TOKEN"
```

---

## 10. Dashboards

### Auto-dispatch — returns data based on your role

```bash
curl http://localhost:8080/api/dashboard \
  -H "Authorization: Bearer $TOKEN"
```

### Explicit per-role endpoints

```bash
# Any auth user
curl http://localhost:8080/api/dashboard/user \
  -H "Authorization: Bearer $TOKEN"

# writer+
curl http://localhost:8080/api/dashboard/writer \
  -H "Authorization: Bearer $TOKEN"

# partner+
curl http://localhost:8080/api/dashboard/partner \
  -H "Authorization: Bearer $TOKEN"

# worker+
curl http://localhost:8080/api/dashboard/worker \
  -H "Authorization: Bearer $TOKEN"

# admin+
curl http://localhost:8080/api/dashboard/admin \
  -H "Authorization: Bearer $TOKEN"
```

---

## 11. Superadmin Panel

### Request role upgrade — any auth user

```bash
curl -X POST http://localhost:8080/api/upgrade-requests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"to_role": "partner"}'
```

> `to_role` options: `writer`, `partner`, `worker`, `admin`

### Check my upgrade request

```bash
curl http://localhost:8080/api/upgrade-requests/me \
  -H "Authorization: Bearer $TOKEN"
```

---

### Superadmin only

#### Platform stats

```bash
curl http://localhost:8080/api/superadmin/stats \
  -H "Authorization: Bearer $TOKEN"
```

**Expected:**
```json
{
  "total_users": 42,
  "by_role": { "user": 30, "partner": 8, "admin": 4 },
  "total_orders": 120,
  "revenue_cents": 598800,
  "open_tickets": 5,
  "pending_upgrades": 3
}
```

#### Get pending upgrade requests

```bash
curl "http://localhost:8080/api/superadmin/upgrades?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

#### Approve upgrade request

```bash
curl -X POST http://localhost:8080/api/superadmin/upgrades/1/review \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"approve": true}'
```

#### Reject upgrade request

```bash
curl -X POST http://localhost:8080/api/superadmin/upgrades/1/review \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"approve": false}'
```

#### Directly set a user's role

```bash
curl -X PATCH http://localhost:8080/api/superadmin/users/2/role \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "admin"}'
```

---

## Role Hierarchy

| Rank | Role | Can access |
|------|------|-----------|
| 1 | `user` | Own orders, tickets, dashboard/user |
| 2 | `writer` | + Blog write endpoints, dashboard/writer |
| 3 | `partner` | + Services, coupons, dashboard/partner |
| 4 | `worker` | + All tickets management, dashboard/worker |
| 5 | `admin` | + Users, orders, payments, payouts, dashboard/admin |
| 6 | `superadmin` | + Role management, platform stats, upgrade reviews |

---

## Postman Setup

### Environment variables

| Variable | Value |
|----------|-------|
| `baseUrl` | `http://localhost:8080` |
| `token` | _(fill after login)_ |
| `refreshToken` | _(fill after login)_ |

### Auto-save token on login

In the **Tests** tab of your login request:

```javascript
const res = pm.response.json();
if (res.access_token) {
    pm.collectionVariables.set("token", res.access_token);
    pm.collectionVariables.set("refreshToken", res.refresh_token);
}
```

### Authorization header (collection level)

Set **Auth Type** → `Bearer Token` → value: `{{token}}`

All requests inherit this automatically. Override per-request for public endpoints by setting auth to `No Auth`.

### Recommended test flow

```
1. POST /api/auth/register/init       → get OTP (check server log)
2. POST /api/auth/register/verify     → create account, save token
3. POST /api/superadmin/users/1/role  → (as superadmin) promote user
4. POST /api/services                 → create service (partner)
5. POST /api/services/1/live          → set live
6. POST /api/coupons                  → create coupon
7. POST /api/coupons/1/approve        → approve it (admin)
8. GET  /api/coupons/validate?code=X  → verify it works
9. POST /api/orders                   → place order with coupon
10. POST /api/payments/:id/complete   → confirm payment (admin)
11. GET  /api/dashboard               → check your dashboard
12. GET  /api/superadmin/stats        → platform overview
```