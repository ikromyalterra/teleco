package routes

import (
	intlMiddleware "github.com/sepulsa/teleco/api/intl/v1/routes/middleware"

	partnerController "github.com/sepulsa/teleco/api/intl/v1/partner"
	partnerService "github.com/sepulsa/teleco/business/partner"
	partnerRepository "github.com/sepulsa/teleco/modules/repository/mongodb/partner"

	partnerIssuerController "github.com/sepulsa/teleco/api/intl/v1/partner/issuer"
	partnerIssuerService "github.com/sepulsa/teleco/business/partner/issuer"
	partnerIssuerRepository "github.com/sepulsa/teleco/modules/repository/mongodb/partner/issuer"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	authController "github.com/sepulsa/teleco/api/intl/v1/auth"
	authService "github.com/sepulsa/teleco/business/auth"
	userTokenRepository "github.com/sepulsa/teleco/modules/repository/mongodb/usertoken"

	issuerController "github.com/sepulsa/teleco/api/intl/v1/issuer"
	issuerService "github.com/sepulsa/teleco/business/issuer"
	issuerRepository "github.com/sepulsa/teleco/modules/repository/mongodb/issuer"

	userController "github.com/sepulsa/teleco/api/intl/v1/user"
	userService "github.com/sepulsa/teleco/business/user"
	userRepository "github.com/sepulsa/teleco/modules/repository/mongodb/user"

	"github.com/sepulsa/teleco/utils/config"
)

func API(e *echo.Echo) {
	db := config.Mgo

	// User && Auth
	userRepo := userRepository.New(db)
	userTokenRepo := userTokenRepository.New(db)
	userServ := userService.New(userRepo, userTokenRepo)
	authServ := authService.New(userRepo, userTokenRepo, nil)

	authMiddleware := intlMiddleware.NewAuth(authServ)
	JWTCustomConfig := middleware.JWTConfig{
		Skipper:        intlMiddleware.AuthAPISkipper,
		ParseTokenFunc: authMiddleware.CustomParseToken,
	}

	e.Use(middleware.JWTWithConfig(JWTCustomConfig))

	authHandler := authController.New(authServ)
	auth := e.Group("/api/v1/auth")
	auth.POST("/user", authHandler.UserRegister)
	auth.POST("/login", authHandler.UserLogin)
	auth.POST("/refresh", authHandler.UserRefreshToken)
	auth.POST("/logout", authHandler.UserLogout)

	userHandler := userController.New(userServ)
	user := e.Group("/api/v1/user")
	user.GET("", userHandler.ListData)
	user.POST("", userHandler.CreateData)
	user.GET("/:id", userHandler.ReadData)
	user.PUT("/:id", userHandler.UpdateData)
	user.DELETE("/:id", userHandler.DeleteData)

	// Issuer
	issuerRepo := issuerRepository.New(db)
	issuerServ := issuerService.New(issuerRepo)
	issuerHandler := issuerController.New(issuerServ)
	issuer := e.Group("/api/v1/issuer")
	issuer.POST("", issuerHandler.CreateData)
	issuer.GET("/:id", issuerHandler.ReadData)
	issuer.PUT("/:id", issuerHandler.UpdateData)
	issuer.DELETE("/:id", issuerHandler.DeleteData)
	issuer.GET("", issuerHandler.ListData)

	// Partner Mapping
	partnerRepository := partnerRepository.New(db)
	partnerService := partnerService.New(partnerRepository)
	partnerController := partnerController.New(partnerService)
	partner := e.Group("/api/v1/partner")
	partner.POST("", partnerController.CreateData)
	partner.GET("/:id", partnerController.ReadData)
	partner.PUT("/:id", partnerController.UpdateData)
	partner.DELETE("/:id", partnerController.DeleteData)
	partner.GET("", partnerController.ListData)

	// Partner Issuer Mapping
	partnerIssuerRepository := partnerIssuerRepository.New(db)
	partnerIssuerService := partnerIssuerService.New(partnerIssuerRepository)
	partnerIssuerController := partnerIssuerController.New(partnerIssuerService)
	partnerIssuer := e.Group("/api/v1/partner/issuer")
	partnerIssuer.POST("", partnerIssuerController.CreateData)
	partnerIssuer.GET("/:id", partnerIssuerController.ReadData)
	partnerIssuer.PUT("/:id", partnerIssuerController.UpdateData)
	partnerIssuer.DELETE("/:id", partnerIssuerController.DeleteData)
	partnerIssuer.GET("", partnerIssuerController.ListData)
}
