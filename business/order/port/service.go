package port

type OrderService struct {
	ID                  string `json:"id"`
	TransactionId       string `json:"transaction_id"`
	IssuerProductId     string `json:"issuer_product_id"`
	CustomerNumber      string `json:"customer_number"`
	PartnerCode         string `json:"partner_code"`
	IssuerCode          string `json:"issuer_code"`
	IssuerTransactionId string `json:"issuer_transaction_id"`
}

type OrderServiceResult struct {
	IssuerTransactionId string `json:"issuer_transaction_id"`
	SerialNumber        string `json:"serial_number"`
	IssuerRescode       string `json:"issuer_rescode"`
	Message             string `json:"message"`
	RawData             string `json:"rawdata"`
}

// Service is inbound port
type Service interface {
	//Purchase ...
	Purchase(order OrderService) (OrderServiceResult, error)

	//Advise ...
	Advise(order OrderService) (OrderServiceResult, error)

	//Reversal ...
	Reversal(order OrderService) (OrderServiceResult, error)
}
