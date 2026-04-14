# Cureerel Server API Endpoints

A modular Go backend for the Cureerel platform. All API routes are prefixed with `/api`.


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
- **User**: 
- **Dashboard**:
- **Memberships**: 
- **Coins**: 
- **Admin**: 
- **Uploads**: 

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

```
 https://deferred-banking-glaring.ngrok-free.dev/api/payments/stripe/webhook
```