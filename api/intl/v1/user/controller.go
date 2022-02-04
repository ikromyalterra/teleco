package user

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	userPort "github.com/sepulsa/teleco/business/user/port"
	"github.com/sepulsa/teleco/utils/net/httperror"
	"github.com/sepulsa/teleco/utils/validator"
)

type Controller struct {
	userService userPort.Service
}

func New(userService userPort.Service) *Controller {
	return &Controller{
		userService,
	}
}

var (
	ErrRequiredID   = "user id is required"
	ErrUserNotFound = "user not found"
)

// CreateData godoc
// @Summary Add an user
// @Description add an user
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param body body RequestUser true "please refer to user.RequestUser models below"
// @Success 201
// @Failure 400
// @Failure 422
// @Router /user [post]
func (controller *Controller) CreateData(c echo.Context) error {
	reqData := new(RequestUser)
	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}
	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	var data userPort.UserService

	data.Email = reqData.Email
	data.Fullname = reqData.Fullname
	data.Password = reqData.Password

	if err := controller.userService.CreateData(data); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, "")
}

// ReadData godoc
// @Summary Get detail an user
// @Description get detail an user
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param id path string true "User ID"
// @Success 200 {object} ResponseUser
// @Failure 400
// @Failure 404
// @Failure 422
// @Router /user/{id} [get]
func (controller *Controller) ReadData(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, ErrRequiredID)
	}

	data, err := controller.userService.ReadData(id)
	if err != nil {
		if err.Error() == ErrUserNotFound {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: ErrUserNotFound})
		}
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	var user ResponseUser
	user.ID = data.ID
	user.Email = data.Email
	user.Fullname = data.Fullname

	return c.JSON(http.StatusOK, user)
}

// UpdateData godoc
// @Summary Update an user
// @Description update an user
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param id path string true "User ID"
// @Param body body RequestUser true "please refer to user.RequestUser models below"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 422
// @Router /user/{id} [put]
func (controller *Controller) UpdateData(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: ErrRequiredID})
	}

	reqData := new(RequestUser)
	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}
	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	var data userPort.UserService
	data.ID = id
	data.Email = reqData.Email
	data.Password = reqData.Password
	data.Fullname = reqData.Fullname

	if err := controller.userService.UpdateData(data); err != nil {
		if err.Error() == ErrUserNotFound {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: ErrUserNotFound})
		}
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "")
}

// DeleteData godoc
// @Summary Remove an user
// @Description remove an user
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Param id path string true "User ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 422
// @Router /user/{id} [delete]
func (controller *Controller) DeleteData(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: ErrRequiredID})
	}

	if err := controller.userService.DeleteData(id); err != nil {
		if err.Error() == ErrUserNotFound {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: ErrUserNotFound})
		}
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, "")
}

// ListData godoc
// @Summary List users
// @Description list users
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Success 200
// @Failure 422
// @Router /user [get]
func (controller *Controller) ListData(c echo.Context) error {
	datas, err := controller.userService.ListData()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	var user ResponseUser
	users := make([]ResponseUser, 0, len(datas))
	for i := range datas {
		user.ID = datas[i].ID
		user.Email = datas[i].Email
		user.Fullname = datas[i].Fullname

		users = append(users, user)
	}

	return c.JSON(http.StatusOK, map[string][]ResponseUser{"data": users})
}
