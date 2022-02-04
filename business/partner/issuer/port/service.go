package port

import "time"

type (
	PartnerIssuerService struct {
		ID        string    `json:"id"`
		PartnerId string    `json:"partner_id"`
		IssuerId  string    `json:"issuer_id"`
		Config    string    `json:"config"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		DeletedAt time.Time `json:"deleted_at"`
	}
)

// Service is inbound port
type Service interface {
	//CreateData insert new data
	CreateData(partnerIssuer PartnerIssuerService) error

	//ReadData get data by ID
	ReadData(ID string) (PartnerIssuerService, error)

	//UpdateData update new data
	UpdateData(partnerIssuer PartnerIssuerService) error

	//DeleteData delete data
	DeleteData(ID string) error

	// //ListData get list data
	ListData() ([]PartnerIssuerService, error)
}
