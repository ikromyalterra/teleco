package partner_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	partnerController "github.com/sepulsa/teleco/api/intl/v1/partner"
	partnerService "github.com/sepulsa/teleco/business/partner/mock"
	partnerPort "github.com/sepulsa/teleco/business/partner/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	TestPartnerID1          = "613f12af2edd3a56323f0d1d"
	TestPartnerCode1        = "p1"
	TestPartnerName1        = "partner1"
	TestPartnerPic1         = "pic1"
	TestPartnerAddress1     = "address1"
	TestPartnerCallbackUrl1 = "callback_url1"
	TestPartnerIpwhitelist1 = []string{"192.168.1.1"}
	TestPartnerStatus1      = "active"
	TestPartnerSecretKey1   = "SECRETKEY"

	TestPartnerID2          = "613f131d2edd3a56323f0d46"
	TestPartnerCode2        = "p2"
	TestPartnerName2        = "partner2"
	TestPartnerPic2         = "pic2"
	TestPartnerAddress2     = "address2"
	TestPartnerCallbackUrl2 = "callback_url2"
	TestPartnerIpwhitelist2 = []string{"192.168.1.1"}
	TestPartnerStatus2      = "active"
	TestPartnerSecretKey2   = "SECRETKEY"

	ErrRequiredName = "code is required"
)

func TestCreate(t *testing.T) {
	t.Run("Expect partner created", func(t *testing.T) {
		e := echo.New()

		service := partnerService.New()
		partner := partnerController.New(service)
		endpoint := `/api/v1/partner`
		b, _ := json.Marshal(TestPartnerIpwhitelist1)

		// 201
		reqData := fmt.Sprintf(`{"code":"%s","name":"%s","pic":"%s","address":"%s","callback_url":"%s","ip_whitelist":%s,"status":"%s", "secret_key": "%s"}`, TestPartnerCode1, TestPartnerName1, TestPartnerPic1, TestPartnerAddress1, TestPartnerCallbackUrl1, string(b), TestPartnerStatus1, TestPartnerSecretKey1)
		req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(reqData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		service.On("CreateData", mock.Anything).Return(nil).Once()
		if assert.NoError(t, partner.CreateData(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}

		// 400 bind
		reqData = fmt.Sprintf(`{"partner_id":"%s","partner_name":"%s"}`, TestPartnerID1, TestPartnerName1)
		req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(reqData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		if assert.NoError(t, partner.CreateData(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}

		// 400 validate
		reqData = `{"name":""}`
		req = httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(reqData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		if assert.NoError(t, partner.CreateData(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), ErrRequiredName)
		}
	})
}

func TestRead(t *testing.T) {
	e := echo.New()

	service := partnerService.New()
	partner := partnerController.New(service)
	endpoint := `/api/v1/partner`

	// 200
	dataService := partnerPort.PartnerService{
		ID:          TestPartnerID1,
		Name:        TestPartnerName1,
		Pic:         TestPartnerPic1,
		Address:     TestPartnerAddress1,
		CallbackUrl: TestPartnerCallbackUrl1,
		IpWhitelist: TestPartnerIpwhitelist1,
		Status:      TestPartnerStatus1,
	}

	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestPartnerID1)
	service.On("ReadData", mock.Anything).Return(dataService, nil).Once()
	if assert.NoError(t, partner.ReadData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response partnerController.ResponsePartner
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
			assert.Equal(t, TestPartnerID1, response.ID)
			assert.Equal(t, TestPartnerName1, response.Name)
			assert.Equal(t, TestPartnerPic1, response.Pic)
			assert.Equal(t, TestPartnerAddress1, response.Address)
			assert.Equal(t, TestPartnerCallbackUrl1, response.CallbackUrl)
			assert.Equal(t, TestPartnerIpwhitelist1, response.IpWhitelist)
			assert.Equal(t, TestPartnerStatus1, response.Status)
		}
	}

	// 400 empty ID
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("  ")
	if assert.NoError(t, partner.ReadData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), partnerController.ErrRequiredID)
	}

	// 404 partner not found
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestPartnerID1)
	service.On("ReadData", mock.Anything).Return(partnerPort.PartnerService{}, errors.New(partnerController.ErrPartnerNotFound)).Once()
	if assert.NoError(t, partner.ReadData(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), partnerController.ErrPartnerNotFound)
	}

	// 422 err service
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestPartnerID1)
	service.On("ReadData", mock.Anything).Return(partnerPort.PartnerService{}, errors.New("")).Once()
	if assert.NoError(t, partner.ReadData(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestUpdate(t *testing.T) {
	e := echo.New()

	service := partnerService.New()
	partner := partnerController.New(service)
	endpoint := `/api/v1/partner/`
	b, _ := json.Marshal(TestPartnerIpwhitelist1)

	// 200
	updDataSuccess := fmt.Sprintf(`{"id":"%s","code":"%s","name":"%s","pic":"%s","address":"%s","callback_url":"%s","ip_whitelist":%s,"status":"%s","secret_key":"%s"}`, TestPartnerID1, TestPartnerCode1, TestPartnerName1, TestPartnerPic1, TestPartnerAddress1, TestPartnerCallbackUrl1, string(b), TestPartnerStatus1, TestPartnerSecretKey1)
	req := httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updDataSuccess))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestPartnerID1)
	service.On("UpdateData", mock.Anything).Return(nil).Once()
	if assert.NoError(t, partner.UpdateData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 400 empty ID
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updDataSuccess))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("")
	if assert.NoError(t, partner.UpdateData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), partnerController.ErrRequiredID)
	}

	// 400 bind
	updData := fmt.Sprintf(`{"id":"%s","name":"%s","pic":"%s","address":"%s","callback_url":"%s","ip_whitelist":"%s","status":"%s","secret_key":"%s"}`, TestPartnerID1, TestPartnerName1, TestPartnerPic1, TestPartnerAddress1, TestPartnerCallbackUrl1, TestPartnerIpwhitelist1, TestPartnerStatus1, TestPartnerSecretKey1)
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestPartnerID1)
	if assert.NoError(t, partner.UpdateData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	// 400 - validate
	updData = `{"Name":"*"}`
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestPartnerID1)
	if assert.NoError(t, partner.UpdateData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	// 404 - partner not found
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updDataSuccess))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestPartnerID1)
	service.On("UpdateData", mock.Anything).Return(errors.New(partnerController.ErrPartnerNotFound)).Once()
	if assert.NoError(t, partner.UpdateData(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), partnerController.ErrPartnerNotFound)
	}

	// 422 - err service
	req = httptest.NewRequest(http.MethodPut, endpoint, strings.NewReader(updDataSuccess))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestPartnerID1)
	service.On("UpdateData", mock.Anything).Return(errors.New("")).Once()
	if assert.NoError(t, partner.UpdateData(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestDelete(t *testing.T) {
	e := echo.New()

	service := partnerService.New()
	partner := partnerController.New(service)
	endpoint := `/api/v1/partner/`

	// 200
	req := httptest.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestPartnerID1)
	service.On("DeleteData", mock.Anything).Return(nil).Once()
	if assert.NoError(t, partner.DeleteData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 400 empty ID
	req = httptest.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("  ")
	if assert.NoError(t, partner.DeleteData(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), partnerController.ErrRequiredID)
	}

	// 404 partner not found
	req = httptest.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestPartnerID1)
	service.On("DeleteData", mock.Anything).Return(errors.New(partnerController.ErrPartnerNotFound)).Once()
	if assert.NoError(t, partner.DeleteData(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), partnerController.ErrPartnerNotFound)
	}

	// 422 error service
	req = httptest.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(TestPartnerID1)
	service.On("DeleteData", mock.Anything).Return(errors.New("")).Once()
	if assert.NoError(t, partner.DeleteData(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestList(t *testing.T) {
	e := echo.New()

	service := partnerService.New()
	partner := partnerController.New(service)
	endpoint := `/api/v1/partner`

	// 200
	partners := []partnerPort.PartnerService{
		{
			ID:          TestPartnerID1,
			Name:        TestPartnerName1,
			Pic:         TestPartnerPic1,
			Address:     TestPartnerAddress1,
			CallbackUrl: TestPartnerCallbackUrl1,
			IpWhitelist: TestPartnerIpwhitelist1,
			Status:      TestPartnerStatus1,
		},
	}
	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	service.On("ListData").Return(partners, nil).Once()
	if assert.NoError(t, partner.ListData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string][]partnerController.ResponsePartner
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
			assert.Equal(t, TestPartnerID1, response["data"][0].ID)
			assert.Equal(t, TestPartnerName1, response["data"][0].Name)
			assert.Equal(t, TestPartnerPic1, response["data"][0].Pic)
			assert.Equal(t, TestPartnerAddress1, response["data"][0].Address)
			assert.Equal(t, TestPartnerCallbackUrl1, response["data"][0].CallbackUrl)
			assert.Equal(t, TestPartnerIpwhitelist1, response["data"][0].IpWhitelist)
			assert.Equal(t, TestPartnerStatus1, response["data"][0].Status)
		}
	}

	// 200 empty data
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	service.On("ListData").Return([]partnerPort.PartnerService{}, nil).Once()
	if assert.NoError(t, partner.ListData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string][]partnerController.ResponsePartner
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
			assert.Equal(t, 0, len(response["data"]))
		}
	}

	// 422 err service
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	service.On("ListData").Return([]partnerPort.PartnerService{}, errors.New("")).Once()
	if assert.NoError(t, partner.ListData(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}
