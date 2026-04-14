# Cureerel Server API Endpoints

A modular Go backend for the Cureerel platform. All API routes are prefixed with `/api`.


## 🎣 Webhooks (`/webhooks`)
| Method | Path | Provider | Description |
| :--- | :--- | :--- | :--- |
| POST | `/webhooks/stripe` | Stripe | Specialized Stripe event handler |
| POST | `/webhooks/razorpay` | Razorpay | Specialized Razorpay event handler |
| POST | `/payments/stripe/webhook` | Stripe | Legacy/Direct Stripe handler |

## Dashboard 
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