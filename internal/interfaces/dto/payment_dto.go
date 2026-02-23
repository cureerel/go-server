package dto

type PaymentResponse struct {
    ID            string `json:"id"`
    UserID        uint   `json:"user_id"`
    OrderID       string `json:"order_id"`
    Amount        int64  `json:"amount"`
    Currency      string `json:"currency"`
    Status        string `json:"status"`
    Provider      string `json:"provider"`
    ProviderTxnID string `json:"provider_txn_id"`
    CustomerEmail string `json:"customer_email"`
    Description   string `json:"description"`
    CreatedAt     string `json:"created_at"`
    UpdatedAt     string `json:"updated_at"`
}