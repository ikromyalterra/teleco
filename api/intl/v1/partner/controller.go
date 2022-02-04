package partner

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sepulsa/teleco/utils/net/httperror"
	"github.com/sepulsa/teleco/utils/validator"

	partnerPort "github.com/sepulsa/teleco/business/partner/port"
)

var (
	ErrRequiredID      = "ID can't be empty"
	ErrPartnerNotFound = "Partner not found"
)

type Controller struct {
	partnerService partnerPort.Service
}

// partnerService *portPartner.Service
func New(partnerService partnerPort.Service) *Controller {
	return &Controller{
		partnerService,
	}
}

// CreateData godoc
// @Summary Add an partner
// @Description add an partner
// @Tags Partner
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param body body RequestPartner true "please refer to partner.RequestPartner models below"
// @Success 200
// @Failure 400
// @Failure 422
// @Router /partner [post]
func (controller *Controller) CreateData(c echo.Context) error {
	reqData := new(RequestPartner)
	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}
	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	data := partnerPort.PartnerService{
		Code:        reqData.Code,
		Name:        reqData.Name,
		Pic:         reqData.Pic,
		Address:     reqData.Address,
		CallbackUrl: reqData.CallbackUrl,
		IpWhitelist: reqData.IpWhitelist,
		Status:      reqData.Status,
		SecretKey:   reqData.SecretKey,
	}

	if err := controller.partnerService.CreateData(data); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, "")
}

// CreateData godoc
// @Summary Get an partner
// @Description get an partner
// @Tags Partner
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param id path string true "id"
// @Success 200
// @Failure 400
// @Failure 422
// @Router /partner/{id} [get]
func (controller *Controller) ReadData(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, ErrRequiredID)
	}

	data, err := controller.partnerService.ReadData(id)
	if err != nil {
		if err.Error() == ErrPartnerNotFound {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: ErrPartnerNotFound})
		}
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}
	partner := ResponsePartner{
		ID:          data.ID,
		Code:        data.Code,
		Name:        data.Name,
		Pic:         data.Pic,
		Address:     data.Address,
		CallbackUrl: data.CallbackUrl,
		IpWhitelist: data.IpWhitelist,
		Status:      data.Status,
		SecretKey:   data.SecretKey,
	}

	return c.JSON(http.StatusOK, partner)
}

// CreateData godoc
// @Summary Edit an partner
// @Description edit an partner
// @Tags Partner
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param id path string true "id"
// @Param body body RequestPartner true "please refer to partner.RequestPartner models below"
// @Success 200
// @Failure 400
// @Failure 422
// @Router /partner/{id} [put]
func (controller *Controller) UpdateData(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: ErrRequiredID})
	}

	reqData := new(RequestPartner)
	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}
	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	data := partnerPort.PartnerService{
		ID:          id,
		Code:        reqData.Code,
		Name:        reqData.Name,
		Pic:         reqData.Pic,
		Address:     reqData.Address,
		CallbackUrl: reqData.CallbackUrl,
		IpWhitelist: reqData.IpWhitelist,
		Status:      reqData.Status,
		SecretKey:   reqData.SecretKey,
	}
	if err := controller.partnerService.UpdateData(data); err != nil {
		if err.Error() == ErrPartnerNotFound {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: ErrPartnerNotFound})
		}
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "")
}

// CreateData godoc
// @Summary Delete an partner
// @Description delete an partner
// @Tags Partner
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param id path string true "id"
// @Success 200
// @Failure 400
// @Failure 422
// @Router /partner/{id} [delete]
func (controller *Controller) DeleteData(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: ErrRequiredID})
	}

	if err := controller.partnerService.DeleteData(id); err != nil {
		if err.Error() == ErrPartnerNotFound {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: ErrPartnerNotFound})
		}
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, "")
}

// CreateData godoc
// @Summary Get List an partner
// @Description get list an partner
// @Tags Partner
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Success 200
// @Failure 400
// @Failure 422
// @Router /partner [get]
func (controller *Controller) ListData(c echo.Context) error {
	datas, err := controller.partnerService.ListData()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	partners := make([]ResponsePartner, 0)
	if len(datas) > 0 {
		d, _ := json.Marshal(datas)
		json.Unmarshal(d, &partners)
	}

	return c.JSON(http.StatusOK, map[string][]ResponsePartner{"data": partners})
}
