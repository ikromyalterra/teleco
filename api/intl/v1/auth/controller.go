package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	authPort "github.com/sepulsa/teleco/business/auth/port"
	"github.com/sepulsa/teleco/utils/auth"
	"github.com/sepulsa/teleco/utils/net/httperror"
	"github.com/sepulsa/teleco/utils/validator"
)

type Controller struct {
	authService authPort.Service
}

func New(authService authPort.Service) *Controller {
	return &Controller{
		authService,
	}
}

// UserRegister godoc
// @Summary Register new user
// @Description register new user
// @Tags UserAuthentication
// @Accept  json
// @Produce  json
// @Param body body UserRegisterRequest true "please refer to auth.UserRegisterRequest models below"
// @Success 200
// @Failure 422
// @Failure 400
// @Router /auth/user [post]
func (controller *Controller) UserRegister(c echo.Context) error {
	reqData := new(UserRegisterRequest)
	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}
	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	var data authPort.UserAuthService
	data.Email = reqData.Email
	data.Password = reqData.Password
	data.Fullname = reqData.Fullname

	if err := controller.authService.UserRegister(data); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, "")
}

// UserLogin godoc
// @Summary User login
// @Description user login for get JWT
// @Tags UserAuthentication
// @Accept  json
// @Produce  json
// @Param body body UserRequestLogin true "please refer to auth.UserRequestLogin models below"
// @Success 200 {object} UserResponsetLogin
// @Failure 400
// @Router /auth/login [post]
func (controller *Controller) UserLogin(c echo.Context) error {
	reqData := new(UserRequestLogin)
	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}
	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	userAuth := new(authPort.UserAuthService)
	userAuth.Email = reqData.Email
	userAuth.Password = reqData.Password
	userAuth.KeepLogin = reqData.KeepLogin
	if err := controller.authService.UserLogin(userAuth); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	var response UserResponsetLogin
	response.Token = userAuth.Token
	response.RefreshToken = userAuth.RefreshToken

	return c.JSON(http.StatusOK, response)
}

// UserLogout godoc
// @Summary User logout
// @Description user logout will revoke token
// @Tags UserAuthentication
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authentication Bearer Token (JWT)" default(Bearer token)
// @Success 200
// @Failure 400
// @Failure 401
// @Router /auth/logout [post]
func (controller *Controller) UserLogout(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	tokenClaims := user.Claims.(*auth.JWTClaims)

	if err := controller.authService.UserLogout(tokenClaims.ID); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "")
}

// UserRefreshToken godoc
// @Summary Refresh user token
// @Description renewal user token
// @Tags UserAuthentication
// @Accept  json
// @Produce  json
// @Param body body UserRefreshToken true "please refer to auth.UserRefreshToken models below"
// @Success 200 {object} UserResponsetLogin
// @Failure 400
// @Router /auth/refresh [post]
func (controller *Controller) UserRefreshToken(c echo.Context) error {
	reqData := new(UserRefreshToken)
	if err := c.Bind(reqData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}
	if err := validator.GetValidator().Struct(reqData); err != nil {
		return httperror.NewValidationError(c, http.StatusBadRequest, err)
	}

	tokens, err := controller.authService.UserRefreshToken(reqData.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	var response UserResponsetLogin
	response.Token = tokens.Token
	response.RefreshToken = tokens.RefreshToken

	return c.JSON(http.StatusOK, response)
}
