// internal/interfaces/dto/superadmin_dto.go
package dto

type SetRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type RequestUpgradeRequest struct {
	ToRole string `json:"to_role" binding:"required"`
}

type ReviewUpgradeRequest struct {
	Approve bool `json:"approve"`
}

type UpgradeRequestResponse struct {
	ID         uint    `json:"id"`
	UserID     uint    `json:"user_id"`
	FromRole   string  `json:"from_role"`
	ToRole     string  `json:"to_role"`
	Status     string  `json:"status"`
	ReviewedBy *uint   `json:"reviewed_by,omitempty"`
	ReviewedAt *string `json:"reviewed_at,omitempty"`
	CreatedAt  string  `json:"created_at"`
}

type UpgradeRequestListResponse struct {
	Data  []UpgradeRequestResponse `json:"data"`
	Total int64                    `json:"total"`
	Page  int                      `json:"page"`
	Limit int                      `json:"limit"`
}