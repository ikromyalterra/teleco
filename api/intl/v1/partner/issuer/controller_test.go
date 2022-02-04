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
	partnerIssuerController "github.com/sepulsa/teleco/api/intl/v1/partner/issuer"
	partnerIssuerService "github.com/sepulsa/teleco/business/partner/issuer/mock"
	partnerIssuerPort "github.com/sepulsa/teleco/business/partner/issuer/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	TestID        = "6138813fb95630b0b528b160"
	TestPartnerID = "6138813fb95630b0b528b162"
	TestIssuerID  = "6138813fb95630b0b528b163"
	TestConfig    = "config"

	ErrRequiredCode = "partner_id is required"
)

func TestCreateData(t *testing.T) {
	e := echo.New()

	service := partnerIssuerService.New()
	partnerIssuer := partnerIssuerController.New(service)
	endpoint := `/api/v1/partnerIssuer`

	// 201
	reqData := fmt.Sprintf(`{"partner_id":"%s","issuer_id":"%s","config":"%s"}`, TestPartnerID, TestIssuerID, TestConfig)
	req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(reqData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	service.On("CreateData", mock.Anything).Return(nil).Once()
	if assert.NoError(t, partnerIssuer.CreateData(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	// 400 bind
	reqData = fmt.Sprintf(`{"issuer_id":"%s","config":"%s"}`, TestIssuerID, TestConfig)
	req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(reqData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, partnerIssuer.CreateData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	// 400 validate
	reqData = `{"partner_id":""}`
	req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(reqData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, partnerIssuer.CreateData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrRequiredCode)
	}
}

func TestReadData(t *testing.T) {
	e := echo.New()

	service := partnerIssuerService.New()
	partnerIssuer := partnerIssuerController.New(service)
	endpoint := `/api/v1/partnerIssuer`

	// 200
	dataService := partnerIssuerPort.PartnerIssuerService{
		ID:        TestID,
		PartnerId: TestPartnerID,
		IssuerId:  TestIssuerID,
		Config:    TestConfig,
	}

	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("ReadData", mock.Anything).Return(dataService, nil).Once()
	if assert.NoError(t, partnerIssuer.ReadData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response partnerIssuerController.ResponsePartnerIssuer
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
			assert.Equal(t, TestID, response.ID)
			assert.Equal(t, TestPartnerID, response.PartnerId)
			assert.Equal(t, TestIssuerID, response.IssuerId)
			assert.Equal(t, TestConfig, response.Config)
		}
	}

	// 400 empty ID
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("  ")
	if assert.NoError(t, partnerIssuer.ReadData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), partnerIssuerController.ErrRequiredID)
	}

	// 404 partnerIssuer not found
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("ReadData", mock.Anything).Return(partnerIssuerPort.PartnerIssuerService{}, errors.New(partnerIssuerController.ErrPartnerIssuerNotFound)).Once()
	if assert.NoError(t, partnerIssuer.ReadData(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), partnerIssuerController.ErrPartnerIssuerNotFound)
	}

	// 422 err service
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("ReadData", mock.Anything).Return(partnerIssuerPort.PartnerIssuerService{}, errors.New("")).Once()
	if assert.NoError(t, partnerIssuer.ReadData(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestUpdateData(t *testing.T) {
	e := echo.New()

	service := partnerIssuerService.New()
	partnerIssuer := partnerIssuerController.New(service)
	endpoint := `/api/v1/partnerIssuer/`

	// 200
	updDataSuccess := fmt.Sprintf(`{"partner_id":"%s","issuer_id":"%s","config":"%s"}`, TestPartnerID, TestIssuerID, TestConfig)
	req := httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updDataSuccess))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("UpdateData", mock.Anything).Return(nil).Once()
	if assert.NoError(t, partnerIssuer.UpdateData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 400 empty ID
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updDataSuccess))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("")
	if assert.NoError(t, partnerIssuer.UpdateData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), partnerIssuerController.ErrRequiredID)
	}

	// 400 bind
	updData := fmt.Sprintf(`{"issuer_id":"%s","config":"%s"}`, TestIssuerID, TestConfig)
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	if assert.NoError(t, partnerIssuer.UpdateData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	// 400 - validate
	updData = `{"code":"*"}`
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	if assert.NoError(t, partnerIssuer.UpdateData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	// 404 - partnerIssuer not found
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updDataSuccess))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("UpdateData", mock.Anything).Return(errors.New(partnerIssuerController.ErrPartnerIssuerNotFound)).Once()
	if assert.NoError(t, partnerIssuer.UpdateData(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), partnerIssuerController.ErrPartnerIssuerNotFound)
	}

	// 422 - err service
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updDataSuccess))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("UpdateData", mock.Anything).Return(errors.New("")).Once()
	if assert.NoError(t, partnerIssuer.UpdateData(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestDeleteData(t *testing.T) {
	e := echo.New()

	service := partnerIssuerService.New()
	partnerIssuer := partnerIssuerController.New(service)
	endpoint := `/api/v1/partnerIssuer/`

	// 200
	req := httptest.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("DeleteData", mock.Anything).Return(nil).Once()
	if assert.NoError(t, partnerIssuer.DeleteData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 400 empty ID
	req = httptest.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("  ")
	if assert.NoError(t, partnerIssuer.DeleteData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), partnerIssuerController.ErrRequiredID)
	}

	// 404 partnerIssuer not found
	req = httptest.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("DeleteData", mock.Anything).Return(errors.New(partnerIssuerController.ErrPartnerIssuerNotFound)).Once()
	if assert.NoError(t, partnerIssuer.DeleteData(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), partnerIssuerController.ErrPartnerIssuerNotFound)
	}

	// 422 error service
	req = httptest.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestID)
	service.On("DeleteData", mock.Anything).Return(errors.New("")).Once()
	if assert.NoError(t, partnerIssuer.DeleteData(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestListData(t *testing.T) {
	e := echo.New()

	service := partnerIssuerService.New()
	partnerIssuer := partnerIssuerController.New(service)
	endpoint := `/api/v1/partnerIssuer`

	// 200
	partnerIssuers := []partnerIssuerPort.PartnerIssuerService{
		{
			ID:        TestID,
			PartnerId: TestPartnerID,
			IssuerId:  TestIssuerID,
			Config:    TestConfig,
		},
	}
	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	service.On("ListData").Return(partnerIssuers, nil).Once()
	if assert.NoError(t, partnerIssuer.ListData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string][]partnerIssuerController.ResponsePartnerIssuer
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
			assert.Equal(t, TestID, response["data"][0].ID)
			assert.Equal(t, TestPartnerID, response["data"][0].PartnerId)
			assert.Equal(t, TestIssuerID, response["data"][0].IssuerId)
			assert.Equal(t, TestPartnerID, response["data"][0].PartnerId)
		}
	}

	// 200 empty data
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	service.On("ListData").Return([]partnerIssuerPort.PartnerIssuerService{}, nil).Once()
	if assert.NoError(t, partnerIssuer.ListData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string][]partnerIssuerController.ResponsePartnerIssuer
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
			assert.Equal(t, 0, len(response["data"]))
		}
	}

	// 422 err service
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	service.On("ListData").Return([]partnerIssuerPort.PartnerIssuerService{}, errors.New("")).Once()
	if assert.NoError(t, partnerIssuer.ListData(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}
