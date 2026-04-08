package router

import (
	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/interfaces/http/handler"
	"github.com/cureerel/cserver/pkg/logger"
)

// Deps aggregates HTTP dependencies for modular route registration.
type Deps struct {
	AuthService       *service.AuthService
	UserHandler       *handler.UserHandler
	BlogHandler       *handler.BlogHandler
	AuthHandler       *handler.AuthHandler
	ServiceHandler    *handler.ServiceHandler
	OrderHandler      *handler.OrderHandler
	PaymentHandler    *handler.PaymentHandler
	CouponHandler     *handler.CouponHandler
	PayoutHandler     *handler.PayoutHandler
	TicketHandler     *handler.TicketHandler
	DashboardHandler  *handler.DashboardHandler
	SuperAdminHandler *handler.SuperAdminHandler
	UploadHandler     *handler.UploadHandler
	MembershipHandler *handler.MembershipHandler
	PGHandler         *handler.PaymentGatewayHandler
	CoinHandler       *handler.CoinHandler
	Log               logger.Logger
	AllowedOrigins    []string
}
