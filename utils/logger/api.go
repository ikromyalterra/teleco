package logger

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// APILogHandler : handle something who need to do
func APILogHandler(c echo.Context, req, res []byte) {
	appLog := AppLogger{}
	appLog.ConnectionType = "http"
	appLog.Package = "api"

	c.Response().Header().Set("X-Api-ResponseTime", time.Now().Format(time.RFC3339))
	reqTime, err := time.Parse(time.RFC3339, c.Request().Header.Get("X-Api-RequestTime"))
	if err == nil {
		appLog.TimerStart = reqTime
	}

	reqURL := c.Request().RequestURI
	appLog.Event = "request.http"
	appLog.RequestURL = reqURL
	appLog.RequestHost = c.Request().Host
	appLog.RequestMethod = c.Request().Method
	appLog.RequestIP = c.RealIP()
	appLog.RequestTime = c.Request().Header.Get("X-Api-RequestTime")

	var reqHeadIntf interface{}
	reqHeader, _ := json.Marshal(c.Request().Header)
	lerr := json.Unmarshal(reqHeader, &reqHeadIntf)
	if lerr == nil {
		appLog.RequestHeader = reqHeadIntf
	}

	var reqIntf interface{}
	lerr = json.Unmarshal(req, &reqIntf)
	if lerr == nil {
		appLog.Request = reqIntf
	}

	var data interface{}
	lerr = json.Unmarshal(req, &data)
	if lerr == nil {
		appLog.Data = data
	} else {
		appLog.Data = appLog.Request
	}

	appLog.ResponseTime = c.Response().Header().Get("X-Api-ResponseTime")
	appLog.RawResponse = string(res)

	var resHeadIntf interface{}
	respHeader, _ := json.Marshal(c.Response().Header())
	lerr = json.Unmarshal(respHeader, &resHeadIntf)
	if lerr == nil {
		appLog.ResponseHeader = resHeadIntf
	}

	var resIntf interface{}
	lerr = json.Unmarshal([]byte(string(res)), &resIntf)
	if lerr == nil {
		appLog.Response = resIntf
	}

	appLog.HTTPCode = c.Response().Status
	appLog.Info().Msg("")

}

// APILogSkipper : rules for APILogHandler
func APILogSkipper(c echo.Context) bool {
	// bool, is this url request include "/api"?
	rules1 := strings.Contains(c.Request().RequestURI, "/api")

	if rules1 {
		return false
	}

	return true
}
