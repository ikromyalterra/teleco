package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sepulsa/teleco/api/extl/v1/routes"
	"github.com/sepulsa/teleco/utils/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	_ "github.com/sepulsa/teleco/app/api/extl/docs/v1"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Teleco External API
// @version 1.0
// @description REST API for Partner or Biller
// @termsOfService http://swagger.io/terms/

// @contact.name Teleco
// @contact.url https://teleco-dev.sumpahpalapa.com/
// @contact.email ikromy@alterra.id

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes https
// @host teleco-dev.sumpahpalapa.com
// @BasePath /api/v1
func main() {

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Output: logger.MiddlewareLog,
		}),
		middleware.Recover())
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		logger.Error().Msgf("[ERROR] -- %s, Request Path %s", err.Error(), c.Request().URL.Path)
		e.DefaultHTTPErrorHandler(err, c)
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	// Handler for hooking any request in routers registered and log it
	e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Handler: logger.APILogHandler,
		Skipper: logger.APILogSkipper,
	}))
	// Handler for putting app request and response timestamp. This used for get elapsed time
	e.Use(ServiceRequestTime)

	routes.API(e)

	e.GET("/", func(c echo.Context) error {
		message := `Aku adalah ...

	-- External, 2021`
		return c.String(http.StatusOK, message)
	})

	//health check
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(200)
	})

	// swagger
	e.GET("/api/v1/swagger/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) { c.URL = "./doc.json" }))

	// Start server
	go func() {
		if err := e.Start(":" + viper.GetString("port")); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// ServiceRequestTime middleware adds a `Server` header to the response.
func ServiceRequestTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set("X-App-RequestTime", time.Now().Format(time.RFC3339))
		return next(c)
	}
}
