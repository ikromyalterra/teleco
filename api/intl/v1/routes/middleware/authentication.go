package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	authPort "github.com/sepulsa/teleco/business/auth/port"
)

type Authentication struct {
	authService authPort.Service
}

func NewAuth(authService authPort.Service) *Authentication {
	return &Authentication{
		authService,
	}
}

func (handler *Authentication) CustomParseToken(tokenString string, c echo.Context) (interface{}, error) {
	return handler.authService.VerifyUserToken(tokenString)
}

func AuthAPISkipper(c echo.Context) bool {
	paths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/refresh",
		"/api/v1/auth/user",
		"/api/v1/swagger",
		"/health",
	}

	requestURI := c.Request().RequestURI
	for i := range paths {
		if strings.Contains(requestURI, paths[i]) {
			return true
		}
	}
	indexAlwaysAllowed := requestURI == "/" || requestURI == ""

	return indexAlwaysAllowed
}
