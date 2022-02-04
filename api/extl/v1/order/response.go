package order

type ResponseOrder struct {
	IssuerTransactionId string `json:"issuer_transaction_id"`
	SerialNumber        string `json:"serial_number"`
	IssuerRescode       string `json:"issuer_rescode"`
	Message             string `json:"message"`
	RawData             string `json:"rawdata"`
}
