package httperror

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// New example
func NewValidationError(c echo.Context, status int, err error) error {
	var errMsg string
	if castedObject, ok := err.(validator.ValidationErrors); ok {
		for _, err := range castedObject {
			switch err.Tag() {
			case "required":
				errMsg = fmt.Sprintf("%s is required", err.Field())
			case "email":
				errMsg = fmt.Sprintf("invalid %s format", err.Field())
			case "alphanum":
				errMsg = fmt.Sprintf("%s value must alphanumeric", err.Field())
			case "min":
				errMsg = fmt.Sprintf("%s length too short (min=%s)", err.Field(), err.Param())
			case "max":
				errMsg = fmt.Sprintf("%s length too long (max=%s)", err.Field(), err.Param())
			case "gte":
				errMsg = fmt.Sprintf("%s value must be greater than %s", err.Field(), err.Param())
			case "lte":
				errMsg = fmt.Sprintf("%s value must be lower than %s", err.Field(), err.Param())
			}

			if errMsg != "" {
				break
			}
		}
	}

	return c.JSON(status, echo.HTTPError{Message: errMsg})
}
