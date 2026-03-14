// internal/interfaces/dto/ticket_dto.go
package dto

type CreateTicketRequest struct {
	Subject     string `json:"subject"     binding:"required,min=5,max=200"`
	Description string `json:"description" binding:"required,min=10"`
	Priority    string `json:"priority"    binding:"omitempty,oneof=low medium high urgent"`
}

type SendMessageRequest struct {
	Message string `json:"message" binding:"required,min=1"`
}

type AssignTicketRequest struct {
	WorkerID uint `json:"worker_id" binding:"required,min=1"`
}

type TicketMessageResponse struct {
	ID        uint   `json:"id"`
	TicketID  uint   `json:"ticket_id"`
	SenderID  uint   `json:"sender_id"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

type TicketResponse struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"user_id"`
	Subject     string  `json:"subject"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Priority    string  `json:"priority"`
	AssignedTo  *uint   `json:"assigned_to,omitempty"`
	ClosedAt    *string `json:"closed_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type TicketListResponse struct {
	Data  []TicketResponse `json:"data"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}