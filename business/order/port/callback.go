package port

type (
	CallbackResult struct {
		RequestData  string `json:"request_data"`
		ResponseData string `json:"response_data"`
	}
)

// Callback is outbound port
type Callback interface {
	//Do ...
	Do(order OrderIssuerApi, orderResult OrderIssuerApiResult) CallbackResult
}
