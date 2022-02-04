package port

import "time"

type (
	PartnerIssuerRepo struct {
		ID        string    `json:"id"`
		PartnerId string    `json:"partner_id"`
		IssuerId  string    `json:"issuer_id"`
		Config    string    `json:"config"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

// Repository is outbound port
type Repository interface {
	//CreateData insert new data
	CreateData(partnerIssuer PartnerIssuerRepo) error

	//ReadData get data by ID
	FindByPartnerIssuerID(partnerId string, issuerId string) (PartnerIssuerRepo, error)

	//ReadData get data by ID
	ReadData(ID string) (PartnerIssuerRepo, error)

	//UpdateData update new data
	UpdateData(partnerIssuer PartnerIssuerRepo) error

	//DeleteData delete data
	DeleteData(ID string) error

	//ListData get list data
	ListData() ([]PartnerIssuerRepo, error)
}
