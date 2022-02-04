package dummy

import (
	"encoding/json"
	"io/ioutil"

	orderPort "github.com/sepulsa/teleco/business/order/port"
	"github.com/sepulsa/teleco/utils/net/httpclient"
	"github.com/sepulsa/teleco/utils/net/httpdump"
)

type (
	Issuer struct{}

	IssuerConfig struct {
		Url string `json:"url"`
	}

	PartnerIssuerConfig struct {
		ID       string `json:"id"`
		PIN      string `json:"pin"`
		User     string `json:"user"`
		Password string `json:"password"`
	}

	RequestData struct {
		ID         string `json:"id"`
		PIN        string `json:"pin"`
		User       string `json:"user"`
		Password   string `json:"password"`
		KodeProduk string `json:"kodeproduk"`
		IdTrx      string `json:"Idtrx"`
		Tujuan     string `json:"tujuan"`
	}

	ResponseData struct {
		Success      bool   `json:"success"`
		Produk       string `json:"produk"`
		Tujuan       string `json:"tujuan"`
		ReffId       string `json:"reffid"`
		Rc           string `json:"rc"`
		Msg          string `json:"msg"`
		Status       string `json:"status"`
		SerialNumber string `json:"serial_number"`
	}
)

func New() *Issuer {
	return &Issuer{}
}

func (is *Issuer) Purchase(order orderPort.OrderIssuerApi, result *orderPort.OrderIssuerApiResult, errOrder *orderPort.Error) {

	// Parse json string issuer config
	var issuerConfig IssuerConfig
	json.Unmarshal([]byte(order.IssuerConfig), &issuerConfig)

	// Parse json string config
	var partnerIssuerConfig PartnerIssuerConfig
	json.Unmarshal([]byte(order.PartnerIssuerConfig), &partnerIssuerConfig)
	reqData := RequestData{
		ID:         partnerIssuerConfig.ID,
		PIN:        partnerIssuerConfig.PIN,
		User:       partnerIssuerConfig.User,
		Password:   partnerIssuerConfig.Password,
		KodeProduk: order.IssuerProductId,
		IdTrx:      order.TransactionId,
		Tujuan:     order.CustomerNumber,
	}
	b, _ := json.Marshal(reqData)

	// Set HTTP Parameters
	var httpParam httpclient.HttpParam
	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	httpParam.Url = issuerConfig.Url + "/purchase"
	httpParam.Method = "post"
	httpParam.Header = header
	httpParam.Body = string(b)
	httpParam.Timeout = 30

	// log request
	b, _ = json.Marshal(httpParam)
	result.RequestData = string(b)

	// Hit API
	res, err := httpParam.HttpDo()

	// log response
	result.ResponseData = httpdump.DumpResponse(res)

	if err != nil {
		errOrder.Err = err
		return
	}
	defer res.Body.Close()

	// read response
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		errOrder.Err = err
		return
	}

	// parse response
	var resData ResponseData
	json.Unmarshal(response, &resData)

	result.SerialNumber = resData.SerialNumber
	result.IssuerTransactionId = resData.ReffId
	result.Message = resData.Msg
	result.IssuerRescode = resData.Rc
	result.RawData = string(response)
}

func (is *Issuer) Advise(order orderPort.OrderIssuerApi, result *orderPort.OrderIssuerApiResult, errOrder *orderPort.Error) {
	// Parse json string issuer config
	var issuerConfig IssuerConfig
	json.Unmarshal([]byte(order.IssuerConfig), &issuerConfig)

	// Parse json string config
	var partnerIssuerConfig PartnerIssuerConfig
	json.Unmarshal([]byte(order.PartnerIssuerConfig), &partnerIssuerConfig)
	reqData := RequestData{
		ID:       partnerIssuerConfig.ID,
		PIN:      partnerIssuerConfig.PIN,
		User:     partnerIssuerConfig.User,
		Password: partnerIssuerConfig.Password,
		IdTrx:    order.IssuerTransactionId,
	}
	b, _ := json.Marshal(reqData)

	// Set HTTP Parameters
	var httpParam httpclient.HttpParam
	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	httpParam.Url = issuerConfig.Url + "/advise"
	httpParam.Method = "post"
	httpParam.Header = header
	httpParam.Body = string(b)
	httpParam.Timeout = 30

	// log request
	b, _ = json.Marshal(httpParam)
	result.RequestData = string(b)

	// Hit API
	res, err := httpParam.HttpDo()

	// log response
	result.ResponseData = httpdump.DumpResponse(res)

	if err != nil {
		errOrder.Err = err
		return
	}
	defer res.Body.Close()

	// read response
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		errOrder.Err = err
		return
	}

	// parse response
	var resData ResponseData
	json.Unmarshal(response, &resData)

	result.SerialNumber = resData.SerialNumber
	result.IssuerTransactionId = resData.ReffId
	result.Message = resData.Msg
	result.IssuerRescode = resData.Rc
	result.RawData = string(response)

}

func (is *Issuer) Reversal(order orderPort.OrderIssuerApi, result *orderPort.OrderIssuerApiResult, errOrder *orderPort.Error) {

	// Parse json string issuer config
	var issuerConfig IssuerConfig
	json.Unmarshal([]byte(order.IssuerConfig), &issuerConfig)

	// Parse json string config
	var partnerIssuerConfig PartnerIssuerConfig
	json.Unmarshal([]byte(order.PartnerIssuerConfig), &partnerIssuerConfig)
	reqData := RequestData{
		ID:       partnerIssuerConfig.ID,
		PIN:      partnerIssuerConfig.PIN,
		User:     partnerIssuerConfig.User,
		Password: partnerIssuerConfig.Password,
		IdTrx:    order.IssuerTransactionId,
	}
	b, _ := json.Marshal(reqData)

	// Set HTTP Parameters
	var httpParam httpclient.HttpParam
	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	httpParam.Url = issuerConfig.Url + "/reversal"
	httpParam.Method = "post"
	httpParam.Header = header
	httpParam.Body = string(b)
	httpParam.Timeout = 30

	// log request
	b, _ = json.Marshal(httpParam)
	result.RequestData = string(b)

	// Hit API
	res, err := httpParam.HttpDo()

	// log response
	result.ResponseData = httpdump.DumpResponse(res)

	if err != nil {
		errOrder.Err = err
		return
	}
	defer res.Body.Close()

	// read response
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		errOrder.Err = err
		return
	}

	// parse response
	var resData ResponseData
	json.Unmarshal(response, &resData)

	result.SerialNumber = resData.SerialNumber
	result.IssuerTransactionId = resData.ReffId
	result.Message = resData.Msg
	result.IssuerRescode = resData.Rc
	result.RawData = string(response)
}
