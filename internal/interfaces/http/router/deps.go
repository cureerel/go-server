package router

import (
	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/interfaces/http/handler"
	"github.com/cureerel/cserver/pkg/logger"
)


type Deps struct {
	AuthService       *service.AuthService
	UserHandler       *handler.UserHandler
	BlogHandler       *handler.BlogHandler
	AuthHandler       *handler.AuthHandler
	ServiceHandler    *handler.ServiceHandler
	OrderHandler      *handler.OrderHandler
	PaymentHandler    *handler.PaymentHandler
	CouponHandler     *handler.CouponHandler

	TicketHandler     *handler.TicketHandler
	DashboardHandler  *handler.DashboardHandler
	SuperAdminHandler *handler.AdminHandler
	UploadHandler     *handler.UploadHandler
	MembershipHandler *handler.MembershipHandler
	PGHandler         *handler.PaymentGatewayHandler
	CoinHandler       *handler.CoinHandler
	ProductHandler    *handler.ProductHandler
	WebhookHandler    *handler.WebhookHandler
	Log               logger.Logger
	AllowedOrigins    []string
}
