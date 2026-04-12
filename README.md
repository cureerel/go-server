# Cureerel Server API Endpoints

A modular Go backend for the Cureerel platform. All API routes are prefixed with `/api`.

## 🔐 Authentication (`/auth`)
| Method | Path | Description |
| :--- | :--- | :--- |
| POST | `/auth/register/init` | Start user registration |
| POST | `/auth/register/verify` | Verify registration OTP |
| POST | `/auth/password/reset/init` | Start password reset |
| POST | `/auth/password/reset/verify` | Verify reset OTP |
| POST | `/auth/signup` | Complete signup |
| POST | `/auth/login` | User login |
| POST | `/auth/refresh` | Refresh JWT tokens |
| POST | `/auth/logout` | Logout (Auth Required) |

## ✍️ Blogs (`/blog`, `/blogs`)
| Method | Path | Auth/Role | Description |
| :--- | :--- | :--- | :--- |
| GET | `/blog` | Public | List all blogs |
| GET | `/blog/:id` | Public | Get blog by ID |
| GET | `/blog/slug/:slug` | Public | Get blog by slug |
| POST | `/blogs/:id/unlock` | Auth | Unlock paid content |
| GET | `/blogs/mine` | Writer | List my blogs |
| POST | `/blogs` | Writer | Create a blog |
| PUT | `/blogs/:id` | Writer | Update blog |
| DELETE | `/blogs/:id` | Writer | Delete blog |
| GET | `/reviewer/blogs/pending` | Reviewer | List pending blogs |
| POST | `/reviewer/blogs/:id/approve` | Reviewer | Approve blog |

## 📦 Products (`/products`)
| Method | Path | Role | Description |
| :--- | :--- | :--- | :--- |
| GET | `/products` | Public | List all products |
| GET | `/products/:id` | Public | Get product details |
| POST | `/products` | Admin | Create new product |
| PUT | `/products/:id` | Admin | Update product |
| DELETE | `/products/:id` | Admin | Delete product |

## 🛠️ Services (`/services`)
| Method | Path | Role | Description |
| :--- | :--- | :--- | :--- |
| GET | `/services` | Public | List available services |
| POST | `/services` | Partner | Create a new service |
| POST | `/services/:id/approve` | Admin | Approve a service |

## 💳 Payments & Orders (`/payments`, `/orders`)
| Method | Path | Role | Description |
| :--- | :--- | :--- | :--- |
| POST | `/orders` | Auth | Create a new order |
| GET | `/orders/me` | Auth | List my orders |
| POST | `/payments/razorpay/create-order` | Auth | Init Razorpay order |
| POST | `/payments/stripe/create-session` | Auth | Init Stripe session |
| POST | `/payments/:id/complete` | Admin | Mark payment as completed |

## 🎟️ Support Tickets (`/tickets`)
| Method | Path | Role | Description |
| :--- | :--- | :--- | :--- |
| POST | `/tickets` | Auth | Create a support ticket |
| GET | `/tickets/me` | Auth | List my tickets |
| GET | `/tickets` | Worker | List all tickets |

## 🎣 Webhooks (`/webhooks`)
| Method | Path | Provider | Description |
| :--- | :--- | :--- | :--- |
| POST | `/webhooks/stripe` | Stripe | Specialized Stripe event handler |
| POST | `/webhooks/razorpay` | Razorpay | Specialized Razorpay event handler |
| POST | `/payments/stripe/webhook` | Stripe | Legacy/Direct Stripe handler |

## 📊 Dashboard & System
- **User**: `GET /api/users/me`
- **Dashboard**: `GET /api/dashboard/[user|writer|partner|worker|admin]`
- **Memberships**: `POST /api/memberships/activate`, `POST /api/memberships/upgrade`
- **Coins**: `GET /api/coins/balance`
- **SuperAdmin**: `PATCH /api/superadmin/users/:id/role`, `GET /api/superadmin/stats`
- **Uploads**: `POST /api/upload/image`

## 🚀 Running the Project
```bash
make run
```


```
# ayth return 
{
  "data": {
    "access_token": "...",
    "refresh_token": "...",
    "expires_at": "..."
  }
}
```