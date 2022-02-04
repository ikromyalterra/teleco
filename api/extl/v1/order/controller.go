package order

import (
	"net/http"

	"github.com/sepulsa/teleco/utils/net/httperror"
	"github.com/sepulsa/teleco/utils/validator"

	"github.com/labstack/echo/v4"

	orderPort "github.com/sepulsa/teleco/business/order/port"
)

type Controller struct {
	OrderService orderPort.Service
}

func New(OrderService orderPort.Service) *Controller {
	return &Controller{OrderService}
}

var (
	// refer to middleware
	PartnerCodeContextKey = "partnercode"
)

// Purchase godoc
// @Summary Purchase
// @Description Add an order
// @Tags Order
// @Accept  json
// @Param partner-code header string true "fill with partner code value" default(partner001)
// @Param Authorization header string true "Authentication Bearer Token, token format ===> b64(unixTime:hmacSHA256(unixTime:JSONminify(body)))" default(Bearer token)
// @Produce  json
// @Param body body PurchaseRequestOrder true "please refer to order.PurchaseRequestOrder models below"
// @Success 200 {object} ResponseOrder
// @Failure 400
// @Failure 401
// @Router /order/purchase [post]
func (controller *Controller) Purchase(c echo.Context) error {
	reqData := new(PurchaseRequestOrder)

	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	// get Partner code form header (signature authentication)
	partnerCode := c.Get(PartnerCodeContextKey).(string)

	data := orderPort.OrderService{
		TransactionId:   reqData.TransactionId,
		IssuerProductId: reqData.IssuerProductId,
		CustomerNumber:  reqData.CustomerNumber,
		PartnerCode:     partnerCode,
		IssuerCode:      reqData.IssuerCode,
	}
	result, err := controller.OrderService.Purchase(data)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ResponseOrder{
		IssuerTransactionId: result.IssuerTransactionId,
		SerialNumber:        result.SerialNumber,
		IssuerRescode:       result.IssuerRescode,
		Message:             result.Message,
		RawData:             result.RawData,
	})
}

// Advice godoc
// @Summary Advice
// @Description Check order status
// @Tags Order
// @Accept  json
// @Param partner-code header string true "fill with partner code value" default(partner001)
// @Param Authorization header string true "Authentication Bearer Token, token format ===> b64(unixTime:hmacSHA256(unixTime:JSONminify(body)))" default(Bearer token)
// @Produce  json
// @Param body body AdviseRequestOrder true "please refer to order.AdviseRequestOrder models below"
// @Success 200 {object} ResponseOrder
// @Failure 400
// @Failure 401
// @Router /order/advise [post]
func (controller *Controller) Advise(c echo.Context) error {
	reqData := new(AdviseRequestOrder)

	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	partnerCode := c.Get(PartnerCodeContextKey).(string)

	data := orderPort.OrderService{
		IssuerTransactionId: reqData.IssuerTransactionId,
		PartnerCode:         partnerCode,
		IssuerCode:          reqData.IssuerCode,
	}
	result, err := controller.OrderService.Advise(data)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ResponseOrder{
		IssuerTransactionId: result.IssuerTransactionId,
		SerialNumber:        result.SerialNumber,
		IssuerRescode:       result.IssuerRescode,
		Message:             result.Message,
		RawData:             result.RawData,
	})
}

// Reversal godoc
// @Summary Reversal
// @Description Cancel an order
// @Tags Order
// @Accept  json
// @Param partner-code header string true "fill with partner code value" default(partner001)
// @Param Authorization header string true "Authentication Bearer Token, token format ===> b64(unixTime:hmacSHA256(unixTime:JSONminify(body)))" default(Bearer token)
// @Produce  json
// @Param body body ReversalRequestOrder true "please refer to order.ReversalRequestOrder models below"
// @Success 200 {object} ResponseOrder
// @Failure 400
// @Failure 401
// @Router /order/reversal [post]
func (controller *Controller) Reversal(c echo.Context) error {
	reqData := new(ReversalRequestOrder)

	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	partnerCode := c.Get(PartnerCodeContextKey).(string)

	data := orderPort.OrderService{
		IssuerTransactionId: reqData.IssuerTransactionId,
		PartnerCode:         partnerCode,
		IssuerCode:          reqData.IssuerCode,
	}
	result, err := controller.OrderService.Reversal(data)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ResponseOrder{
		IssuerTransactionId: result.IssuerTransactionId,
		SerialNumber:        result.SerialNumber,
		IssuerRescode:       result.IssuerRescode,
		Message:             result.Message,
		RawData:             result.RawData,
	})
}
