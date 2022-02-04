package port

type (
	OrderIssuerApi struct {
		ID                  string `json:"id"`
		CommandType         string `json:"command_type"`
		IssuerCode          string `json:"issuer_code"`
		TransactionId       string `json:"transaction_id"`
		IssuerProductId     string `json:"issuer_product_id"`
		CustomerNumber      string `json:"customer_number"`
		IssuerConfig        string `json:"issuer_config"`
		PartnerIssuerConfig string `json:"partner_issuer_config"`
		IssuerTransactionId string `json:"issuer_transaction_id"`
		IssuerThreadNum     int    `json:"issuer_thread_num"`
		IssuerThreadTimeout int    `json:"issuer_thread_timeout"`
		PartnerCallbackUrl  string `json:"partner_callback_url"`
		PartnerId           string `json:"partner_id"`
		IssuerId            string `json:"issuer_id"`
	}

	OrderIssuerApiResult struct {
		IssuerTransactionId string `json:"issuer_transaction_id"`
		SerialNumber        string `json:"serial_number"`
		IssuerRescode       string `json:"issuer_rescode"`
		Message             string `json:"message"`
		RawData             string `json:"rawdata"`
		RequestData         string `json:"request_data"`
		ResponseData        string `json:"response_data"`
	}

	Error struct {
		Err error
	}
)

const (
	Purchase string = "purchase"
	Advise   string = "advise"
	Reversal string = "reversal"
)

// IssuerApi is outbound port
type IssuerApi interface {
	//Do ...
	Do(order OrderIssuerApi) (OrderIssuerApiResult, error)
}

// Issuer is outbound port
type Issuer interface {
	//Purchase ...
	Purchase(order OrderIssuerApi, orderResult *OrderIssuerApiResult, err *Error)

	//Advise ...
	Advise(order OrderIssuerApi, orderResult *OrderIssuerApiResult, err *Error)

	//Reversal ...
	Reversal(order OrderIssuerApi, orderResult *OrderIssuerApiResult, err *Error)
}
