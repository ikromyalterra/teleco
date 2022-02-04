package issuer

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sepulsa/teleco/utils/net/httperror"
	"github.com/sepulsa/teleco/utils/validator"

	issuerPort "github.com/sepulsa/teleco/business/issuer/port"
)

var (
	ErrRequiredID     = "ID can't be empty"
	ErrDuplicateCode  = "Code already in use"
	ErrIssuerNotFound = "Issuer not found"
)

type Controller struct {
	issuerService issuerPort.Service
}

func New(issuerService issuerPort.Service) *Controller {
	return &Controller{
		issuerService,
	}
}

// CreateData godoc
// @Summary Add an issuer
// @Description add an issuer
// @Tags Issuer
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param body body RequestIssuer true "please refer to issuer.RequestIssuer models below"
// @Success 201
// @Failure 400
// @Failure 422
// @Router /issuer [post]
func (controller *Controller) CreateData(c echo.Context) error {
	reqData := new(RequestIssuer)
	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}
	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	data := issuerPort.IssuerService{
		Code:             reqData.Code,
		Label:            reqData.Label,
		Config:           reqData.Config,
		ThreadNum:        reqData.ThreadNum,
		ThreadTimeout:    reqData.ThreadTimeout,
		QueueWorkerLimit: reqData.QueueWorkerLimit,
	}
	if err := controller.issuerService.CreateData(data); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, "")
}

// ReadData godoc
// @Summary Get detail an issuer
// @Description get detail an issuer
// @Tags Issuer
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param id path string true "Issuer ID"
// @Success 200 {object} ResponseIssuer
// @Failure 400
// @Failure 404
// @Failure 422
// @Router /issuer/{id} [get]
func (controller *Controller) ReadData(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, ErrRequiredID)
	}

	data, err := controller.issuerService.ReadData(id)
	if err != nil {
		if err.Error() == ErrIssuerNotFound {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: ErrIssuerNotFound})
		}
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}
	issuer := ResponseIssuer{
		ID:               data.ID,
		Code:             data.Code,
		Label:            data.Label,
		Config:           data.Config,
		ThreadNum:        data.ThreadNum,
		ThreadTimeout:    data.ThreadTimeout,
		QueueWorkerLimit: data.QueueWorkerLimit,
	}

	return c.JSON(http.StatusOK, issuer)
}

// UpdateData godoc
// @Summary Update an issuer
// @Description update an issuer
// @Tags Issuer
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param id path string true "Issuer ID"
// @Param body body RequestIssuer true "please refer to issuer.RequestIssuer models below"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 422
// @Router /issuer/{id} [put]
func (controller *Controller) UpdateData(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: ErrRequiredID})
	}

	reqData := new(RequestIssuer)
	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}
	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	data := issuerPort.IssuerService{
		ID:               id,
		Code:             reqData.Code,
		Label:            reqData.Label,
		Config:           reqData.Config,
		ThreadNum:        reqData.ThreadNum,
		ThreadTimeout:    reqData.ThreadTimeout,
		QueueWorkerLimit: reqData.QueueWorkerLimit,
	}
	if err := controller.issuerService.UpdateData(data); err != nil {
		if err.Error() == ErrIssuerNotFound {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: ErrIssuerNotFound})
		}
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "")
}

// DeleteData godoc
// @Summary Remove an issuer
// @Description remove an issuer
// @Tags Issuer
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param id path string true "Issuer ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 422
// @Router /issuer/{id} [delete]
func (controller *Controller) DeleteData(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: ErrRequiredID})
	}

	if err := controller.issuerService.DeleteData(id); err != nil {
		if err.Error() == ErrIssuerNotFound {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: ErrIssuerNotFound})
		}
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, "")
}

// ListData godoc
// @Summary List issuers
// @Description list issuers
// @Tags Issuer
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Success 200
// @Failure 422
// @Router /issuer [get]
func (controller *Controller) ListData(c echo.Context) error {
	datas, err := controller.issuerService.ListData()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	issuers := make([]ResponseIssuer, 0)
	if len(datas) > 0 {
		d, _ := json.Marshal(datas)
		json.Unmarshal(d, &issuers)
	}

	return c.JSON(http.StatusOK, map[string][]ResponseIssuer{"data": issuers})
}
