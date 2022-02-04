package middleware

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/labstack/echo/v4"
	authPort "github.com/sepulsa/teleco/business/auth/port"
	"github.com/sepulsa/teleco/utils/config"
	"github.com/sepulsa/teleco/utils/minifier"
)

type Authentication struct {
	authService authPort.Service
}

func NewAuth(authService authPort.Service) *Authentication {
	return &Authentication{
		authService,
	}
}

var (
	ErrReadBody         = "failed get request body"
	ErrInvalidBody      = "invalid body"
	ErrPartnerCodeEmpty = "header partner-code is required"

	HeaderPartnerCodeKeyName = "partner-code"
)

func (handler *Authentication) PartnerSignatureValidator(key string, c echo.Context) (bool, error) {
	partnerCode := c.Request().Header.Get(HeaderPartnerCodeKeyName)
	if strings.TrimSpace(partnerCode) == "" {
		return false, errors.New(ErrPartnerCodeEmpty)
	}
	c.Set("partnercode", partnerCode)

	payload, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return false, errors.New(ErrReadBody)
	}

	if strings.TrimSpace(string(payload)) != "" {
		var err error
		payload, err = minifier.JSON(payload)
		if err != nil {
			return false, errors.New(ErrInvalidBody)
		}
	}

	// put it back
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	authData := authPort.Signature{
		Token:       key,
		PartnerCode: partnerCode,
		TimeLimit:   config.Signature.TimeLimit,
		Payload:     payload,
	}

	return handler.authService.VerifyPartnerSignature(authData)
}
