package issuer

type RequestPartnerIssuer struct {
	PartnerId string `json:"partner_id" validate:"required"`
	IssuerId  string `json:"issuer_id" validate:"required"`
	Config    string `json:"config" validate:"required"`
}
