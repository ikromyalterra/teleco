package routes

import (
	orderController "github.com/sepulsa/teleco/api/extl/v1/order"
	extlMiddleware "github.com/sepulsa/teleco/api/extl/v1/routes/middleware"
	authService "github.com/sepulsa/teleco/business/auth"
	orderService "github.com/sepulsa/teleco/business/order"
	issuerApi "github.com/sepulsa/teleco/modules/issuerapi"
	issuerRepository "github.com/sepulsa/teleco/modules/repository/mongodb/issuer"
	orderRepository "github.com/sepulsa/teleco/modules/repository/mongodb/order"
	partnerRepository "github.com/sepulsa/teleco/modules/repository/mongodb/partner"
	partnerIssuerRepository "github.com/sepulsa/teleco/modules/repository/mongodb/partner/issuer"
	"github.com/sepulsa/teleco/utils/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func API(e *echo.Echo) {
	db := config.Mgo

	issuerRepo := issuerRepository.New(db)
	partnerRepo := partnerRepository.New(db)
	partnerIssuerRepo := partnerIssuerRepository.New(db)
	orderRepo := orderRepository.New(db)
	issuerApi := issuerApi.New()
	orderServiceHandler := orderService.New(issuerRepo, partnerRepo, partnerIssuerRepo, orderRepo, issuerApi)
	orderHandler := orderController.New(orderServiceHandler)
	authService := authService.New(nil, nil, partnerRepo)
	authMiddleware := extlMiddleware.NewAuth(authService)
	order := e.Group("/api/v1/order", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Validator: authMiddleware.PartnerSignatureValidator,
	}))
	order.POST("/purchase", orderHandler.Purchase)
	order.POST("/advise", orderHandler.Advise)
	order.POST("/reversal", orderHandler.Reversal)
}
