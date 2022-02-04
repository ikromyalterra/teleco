package issuer

type ResponsePartnerIssuer struct {
	ID        string `json:"id"`
	PartnerId string `json:"partner_id"`
	IssuerId  string `json:"issuer_id"`
	Config    string `json:"config"`
}
