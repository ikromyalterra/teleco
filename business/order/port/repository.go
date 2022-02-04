package port

type (
	OrderRepo struct {
		ID                   string `json:"id"`
		CommandType          string `json:"command_type"`
		TransactionId        string `json:"transaction_id"`
		IssuerProductId      string `json:"issuer_product_id"`
		CustomerNumber       string `json:"customer_number"`
		PartnerId            string `json:"partner_id"`
		IssuerId             string `json:"issuer_id"`
		IssuerTransactionId  string `json:"issuer_transaction_id"`
		RequestData          string `json:"request_data"`
		ResponseData         string `json:"response_data"`
		CallbackRequestData  string `json:"callback_request_data"`
		CallbackResponseData string `json:"callback_response_data"`
	}
)

// Repository is outbound port
type Repository interface {
	//CreateData insert new data
	CreateData(issuer OrderRepo) error
}
