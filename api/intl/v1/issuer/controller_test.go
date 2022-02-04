package issuer_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	issuerController "github.com/sepulsa/teleco/api/intl/v1/issuer"
	issuerService "github.com/sepulsa/teleco/business/issuer/mock"
	issuerPort "github.com/sepulsa/teleco/business/issuer/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	TestID          = "6138813fb95630b0b528b160"
	TestCode        = "issuer"
	TestLabel       = "issuer testing"
	TestInvalidCode = 1001

	ErrRequiredCode = "code is required"
	ErrAlphanumCode = "code value must alphanumeric"
)

func TestCreateData(t *testing.T) {
	e := echo.New()

	service := issuerService.New()
	issuer := issuerController.New(service)
	endpoint := `/api/v1/issuer`

	// 201
	reqData := fmt.Sprintf(`{"code":"%s","label":"%s"}`, TestCode, TestLabel)
	req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(reqData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	service.On("CreateData", mock.Anything).Return(nil).Once()
	if assert.NoError(t, issuer.CreateData(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	// 400 bind
	reqData = fmt.Sprintf(`{"code":%d,"label":"%s"}`, TestInvalidCode, TestLabel)
	req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(reqData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, issuer.CreateData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	// 400 validate
	reqData = `{"code":""}`
	req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(reqData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, issuer.CreateData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrRequiredCode)
	}

	// 422 error service
	reqData = fmt.Sprintf(`{"code":"%s"}`, TestCode)
	req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(reqData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	service.On("CreateData", mock.Anything).Return(errors.New(issuerController.ErrDuplicateCode)).Once()
	if assert.NoError(t, issuer.CreateData(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Contains(t, rec.Body.String(), issuerController.ErrDuplicateCode)
	}
}

func TestReadData(t *testing.T) {
	e := echo.New()

	service := issuerService.New()
	issuer := issuerController.New(service)
	endpoint := `/api/v1/issuer`

	// 200
	dataService := issuerPort.IssuerService{
		ID:    TestID,
		Code:  TestCode,
		Label: TestLabel,
	}

	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("ReadData", mock.Anything).Return(dataService, nil).Once()
	if assert.NoError(t, issuer.ReadData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response issuerController.ResponseIssuer
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
			assert.Equal(t, TestID, response.ID)
			assert.Equal(t, TestCode, response.Code)
			assert.Equal(t, TestLabel, response.Label)
		}
	}

	// 400 empty ID
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("  ")
	if assert.NoError(t, issuer.ReadData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), issuerController.ErrRequiredID)
	}

	// 404 issuer not found
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("ReadData", mock.Anything).Return(issuerPort.IssuerService{}, errors.New(issuerController.ErrIssuerNotFound)).Once()
	if assert.NoError(t, issuer.ReadData(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), issuerController.ErrIssuerNotFound)
	}

	// 422 err service
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("ReadData", mock.Anything).Return(issuerPort.IssuerService{}, errors.New("")).Once()
	if assert.NoError(t, issuer.ReadData(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestUpdateData(t *testing.T) {
	e := echo.New()

	service := issuerService.New()
	issuer := issuerController.New(service)
	endpoint := `/api/v1/issuer/`

	// 200
	updDataSuccess := fmt.Sprintf(`{"code":"%s","label":"%s"}`, TestCode, TestLabel)
	req := httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updDataSuccess))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("UpdateData", mock.Anything).Return(nil).Once()
	if assert.NoError(t, issuer.UpdateData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 400 empty ID
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updDataSuccess))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("")
	if assert.NoError(t, issuer.UpdateData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), issuerController.ErrRequiredID)
	}

	// 400 bind
	updData := fmt.Sprintf(`{"code":%d,"label":"%s"}`, TestInvalidCode, TestLabel)
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	if assert.NoError(t, issuer.UpdateData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	// 404 - issuer not found
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updDataSuccess))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("UpdateData", mock.Anything).Return(errors.New(issuerController.ErrIssuerNotFound)).Once()
	if assert.NoError(t, issuer.UpdateData(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), issuerController.ErrIssuerNotFound)
	}

	// 422 - err service
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updDataSuccess))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("UpdateData", mock.Anything).Return(errors.New("")).Once()
	if assert.NoError(t, issuer.UpdateData(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestDeleteData(t *testing.T) {
	e := echo.New()

	service := issuerService.New()
	issuer := issuerController.New(service)
	endpoint := `/api/v1/issuer/`

	// 200
	req := httptest.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("DeleteData", mock.Anything).Return(nil).Once()
	if assert.NoError(t, issuer.DeleteData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 400 empty ID
	req = httptest.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("  ")
	if assert.NoError(t, issuer.DeleteData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), issuerController.ErrRequiredID)
	}

	// 404 issuer not found
	req = httptest.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("DeleteData", mock.Anything).Return(errors.New(issuerController.ErrIssuerNotFound)).Once()
	if assert.NoError(t, issuer.DeleteData(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), issuerController.ErrIssuerNotFound)
	}

	// 422 error service
	req = httptest.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("DeleteData", mock.Anything).Return(errors.New("")).Once()
	if assert.NoError(t, issuer.DeleteData(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestListData(t *testing.T) {
	e := echo.New()

	service := issuerService.New()
	issuer := issuerController.New(service)
	endpoint := `/api/v1/issuer`

	// 200
	issuers := []issuerPort.IssuerService{
		{
			ID:    TestID,
			Code:  TestCode,
			Label: TestLabel,
		},
	}
	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	service.On("ListData").Return(issuers, nil).Once()
	if assert.NoError(t, issuer.ListData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string][]issuerController.ResponseIssuer
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
			assert.Equal(t, TestID, response["data"][0].ID)
			assert.Equal(t, TestCode, response["data"][0].Code)
			assert.Equal(t, TestLabel, response["data"][0].Label)
		}
	}

	// 200 empty data
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	service.On("ListData").Return([]issuerPort.IssuerService{}, nil).Once()
	if assert.NoError(t, issuer.ListData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string][]issuerController.ResponseIssuer
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
			assert.Equal(t, 0, len(response["data"]))
		}
	}

	// 422 err service
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	service.On("ListData").Return([]issuerPort.IssuerService{}, errors.New("")).Once()
	if assert.NoError(t, issuer.ListData(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}
