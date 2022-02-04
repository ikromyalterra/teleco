package callback

import (
	"encoding/json"

	orderPort "github.com/sepulsa/teleco/business/order/port"
	log "github.com/sepulsa/teleco/utils/logger"
	"github.com/sepulsa/teleco/utils/net/httpclient"
	"github.com/sepulsa/teleco/utils/net/httpdump"
)

type (
	Callback struct{}

	RequestData struct {
		IssuerTransactionId string `json:"issuer_transaction_id"`
		SerialNumber        string `json:"serial_number"`
		IssuerRescode       string `json:"issuer_rescode"`
		Message             string `json:"message"`
		RawData             string `json:"rawdata"`
	}
)

var (
	packageLog = "teleco/modules/callback"
)

func New() *Callback {
	return &Callback{}
}

func (c *Callback) Do(order orderPort.OrderIssuerApi, orderResult orderPort.OrderIssuerApiResult) (result orderPort.CallbackResult) {
	reqData := RequestData{
		IssuerTransactionId: orderResult.IssuerTransactionId,
		SerialNumber:        orderResult.SerialNumber,
		IssuerRescode:       orderResult.IssuerRescode,
		Message:             orderResult.Message,
		RawData:             orderResult.RawData,
	}
	b, _ := json.Marshal(reqData)

	// Set HTTP Parameters
	var httpParam httpclient.HttpParam
	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	httpParam.Url = order.PartnerCallbackUrl
	httpParam.Method = "post"
	httpParam.Header = header
	httpParam.Body = string(b)
	httpParam.Timeout = 30

	// log request
	b, _ = json.Marshal(httpParam)
	result.RequestData = string(b)

	// Hit API
	res, err := httpParam.HttpDo()

	if err != nil {
		result.ResponseData = err.Error()
		return
	}
	// // log response
	result.ResponseData = httpdump.DumpResponse(res)

	defer res.Body.Close()

	log.Info().
		Str("event", "callback.executed").
		Str("package", packageLog).
		Str("request_data", result.RequestData).
		Str("response_data", result.ResponseData).
		Msgf("")

	return result
}
