package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	authMiddleware "github.com/sepulsa/teleco/api/extl/v1/routes/middleware"
	authService "github.com/sepulsa/teleco/business/auth/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error body")
}

func TestAuthentication(t *testing.T) {
	endpoint := `/api/v1/order`
	key := "token"

	e := echo.New()
	authService := authService.New()
	middleware := authMiddleware.NewAuth(authService)

	authService.On("VerifyPartnerSignature", mock.Anything).Return(true, nil).Once()

	// valid
	purchaseData := `{"order_id": "1"}`
	req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(purchaseData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(authMiddleware.HeaderPartnerCodeKeyName, "p1")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	valid, err := middleware.PartnerSignatureValidator(key, c)
	if assert.Nil(t, err) {
		assert.Equal(t, true, valid)
	}

	// invalid json
	purchaseData = `{"order_id": test}`
	req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(purchaseData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(authMiddleware.HeaderPartnerCodeKeyName, "p1")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	_, err = middleware.PartnerSignatureValidator(key, c)
	if assert.NotNil(t, err) {
		assert.Equal(t, authMiddleware.ErrInvalidBody, err.Error())
	}

	// read body fail
	req = httptest.NewRequest(http.MethodPost, endpoint, errReader(0))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(authMiddleware.HeaderPartnerCodeKeyName, "p1")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	_, err = middleware.PartnerSignatureValidator(key, c)
	if assert.NotNil(t, err) {
		assert.Equal(t, authMiddleware.ErrReadBody, err.Error())
	}
}
