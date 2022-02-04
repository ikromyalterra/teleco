package issuer

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sepulsa/teleco/utils/net/httperror"
	"github.com/sepulsa/teleco/utils/validator"

	partnerIssuerPort "github.com/sepulsa/teleco/business/partner/issuer/port"
)

var (
	ErrRequiredID            = "ID can't be empty"
	ErrPartnerIssuerNotFound = "PartnerIssuer not found"
)

type Controller struct {
	partnerIssuerService partnerIssuerPort.Service
}

func New(partnerIssuerService partnerIssuerPort.Service) *Controller {
	return &Controller{
		partnerIssuerService,
	}
}

// CreateData godoc
// @Summary Add an partner issuer mapping
// @Description add an partner issuer mapping
// @Tags PartnerIssuerMapping
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param body body RequestPartnerIssuer true "please refer to partnerIssuer.RequestPartnerIssuer models below"
// @Success 201
// @Failure 400
// @Failure 422
// @Router /partner/issuer [post]
func (controller *Controller) CreateData(c echo.Context) error {
	reqData := new(RequestPartnerIssuer)
	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}
	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	data := partnerIssuerPort.PartnerIssuerService{
		PartnerId: reqData.PartnerId,
		IssuerId:  reqData.IssuerId,
		Config:    reqData.Config,
	}
	if err := controller.partnerIssuerService.CreateData(data); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, "")
}

// ReadData godoc
// @Summary Get detail an partner issuer mapping
// @Description get detail an partner issuer mapping
// @Tags PartnerIssuerMapping
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param id path string true "Partner Issuer Mapping ID"
// @Success 200 {object} ResponsePartnerIssuer
// @Failure 400
// @Failure 404
// @Failure 422
// @Router /partner/issuer/{id} [get]
func (controller *Controller) ReadData(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, ErrRequiredID)
	}

	data, err := controller.partnerIssuerService.ReadData(id)
	if err != nil {
		if err.Error() == ErrPartnerIssuerNotFound {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: ErrPartnerIssuerNotFound})
		}
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}
	partnerIssuer := ResponsePartnerIssuer{
		ID:        data.ID,
		PartnerId: data.PartnerId,
		IssuerId:  data.IssuerId,
		Config:    data.Config,
	}

	return c.JSON(http.StatusOK, partnerIssuer)
}

// UpdateData godoc
// @Summary Update an partner issuer mapping
// @Description update an partner issuer mapping
// @Tags PartnerIssuerMapping
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param id path string true "Partner Issuer Mapping ID"
// @Param body body RequestPartnerIssuer true "please refer to partnerIssuer.RequestPartnerIssuer models below"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 422
// @Router /partner/issuer/{id} [put]
func (controller *Controller) UpdateData(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: ErrRequiredID})
	}

	reqData := new(RequestPartnerIssuer)
	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}
	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	data := partnerIssuerPort.PartnerIssuerService{
		ID:        id,
		PartnerId: reqData.PartnerId,
		IssuerId:  reqData.IssuerId,
		Config:    reqData.Config,
	}
	if err := controller.partnerIssuerService.UpdateData(data); err != nil {
		if err.Error() == ErrPartnerIssuerNotFound {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: ErrPartnerIssuerNotFound})
		}
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "")
}

// DeleteData godoc
// @Summary Remove an partner issuer mapping
// @Description remove detail an partner issuer mapping
// @Tags PartnerIssuerMapping
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param id path string true "Partner Issuer Mapping ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 422
// @Router /partner/issuer/{id} [delete]
func (controller *Controller) DeleteData(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: ErrRequiredID})
	}

	if err := controller.partnerIssuerService.DeleteData(id); err != nil {
		if err.Error() == ErrPartnerIssuerNotFound {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: ErrPartnerIssuerNotFound})
		}
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, "")
}

// ListData godoc
// @Summary List partner issuer mapping
// @Description list partner issuer mapping
// @Tags PartnerIssuerMapping
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Success 200
// @Failure 422
// @Router /partner/issuer [get]
func (controller *Controller) ListData(c echo.Context) error {
	datas, err := controller.partnerIssuerService.ListData()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	partnerIssuers := make([]ResponsePartnerIssuer, 0)
	if len(datas) > 0 {
		d, _ := json.Marshal(datas)
		json.Unmarshal(d, &partnerIssuers)
	}

	return c.JSON(http.StatusOK, map[string][]ResponsePartnerIssuer{"data": partnerIssuers})
}
