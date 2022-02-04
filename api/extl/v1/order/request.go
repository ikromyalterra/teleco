package order

type PurchaseRequestOrder struct {
	TransactionId   string `json:"transaction_id" validate:"required"`
	IssuerProductId string `json:"issuer_product_id" validate:"required"`
	CustomerNumber  string `json:"customer_number" validate:"required"`
	IssuerCode      string `json:"issuer_code" validate:"required"`
}

type AdviseRequestOrder struct {
	IssuerTransactionId string `json:"issuer_transaction_id"`
	IssuerCode          string `json:"issuer_code" validate:"required"`
}

type ReversalRequestOrder struct {
	IssuerTransactionId string `json:"issuer_transaction_id"`
	IssuerCode          string `json:"issuer_code" validate:"required"`
}
