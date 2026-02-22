package dto

type ActivateMembershipRequest struct {
    Plan string `json:"plan" binding:"required,oneof=free basic pro enterprise"`
}

type UpgradeMembershipRequest struct {
    Plan string `json:"plan" binding:"required,oneof=free basic pro enterprise"`
}

type MembershipResponse struct {
    ID        uint   `json:"id"`
    UserID    uint   `json:"user_id"`
    Plan      string `json:"plan"`
    Status    string `json:"status"`
    StartsAt  string `json:"starts_at"`
    ExpiresAt string `json:"expires_at"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}