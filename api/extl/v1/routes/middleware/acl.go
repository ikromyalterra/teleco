package middleware

import (
	"github.com/labstack/echo/v4"
)

//ACL is method for checking user permisson
func ACL(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
