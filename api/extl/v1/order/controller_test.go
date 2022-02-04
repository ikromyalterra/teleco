package order_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	orderController "github.com/sepulsa/teleco/api/extl/v1/order"
	orderService "github.com/sepulsa/teleco/business/order/mock"
	orderPort "github.com/sepulsa/teleco/business/order/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	TestPartnerCode = "partnercode"
)

func TestPurchase(t *testing.T) {
	e := echo.New()
	service := orderService.New()
	order := orderController.New(service)
	endpoint := `/api/v1/order/purchase`

	// 202
	purchaseData := `{"transaction_id": "1001", "issuer_product_id":"abcdef", "customer_number":"08123456789", "issuer_code":"issuer"}`
	req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(purchaseData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(orderController.PartnerCodeContextKey, TestPartnerCode)

	result := orderPort.OrderServiceResult{}

	service.On("Purchase", mock.Anything).Return(result, nil).Once()
	if assert.NoError(t, order.Purchase(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 400 bind
	purchaseData = `{"order_id": 1001}`
	req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(purchaseData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, order.Purchase(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	// 400 validate
	purchaseData = `{"order_id": ""}`
	req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(purchaseData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, order.Purchase(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestAdvice(t *testing.T) {
	e := echo.New()
	service := orderService.New()
	order := orderController.New(service)
	endpoint := `/api/v1/order/advise`

	// 202
	adviseData := `{"issuer_transaction_id": "1001", "issuer_code":"issuer"}`
	req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(adviseData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(orderController.PartnerCodeContextKey, TestPartnerCode)

	result := orderPort.OrderServiceResult{}

	service.On("Advise", mock.Anything).Return(result, nil).Once()
	if assert.NoError(t, order.Advise(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 400 bind
	adviseData = `{"order_id": 1001}`
	req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(adviseData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, order.Advise(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	// 400 validate
	adviseData = `{"order_id": ""}`
	req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(adviseData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, order.Advise(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestReversal(t *testing.T) {
	e := echo.New()
	service := orderService.New()
	order := orderController.New(service)
	endpoint := `/api/v1/order/reversal`

	// 202
	reversalData := `{"issuer_transaction_id": "1001", "issuer_code":"issuer"}`
	req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(reversalData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(orderController.PartnerCodeContextKey, TestPartnerCode)

	result := orderPort.OrderServiceResult{}

	service.On("Reversal", mock.Anything).Return(result, nil).Once()
	if assert.NoError(t, order.Reversal(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 400 bind
	reversalData = `{"order_id": 1001}`
	req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(reversalData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, order.Reversal(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	// 400 validate
	reversalData = `{"order_id": ""}`
	req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(reversalData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, order.Reversal(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}
