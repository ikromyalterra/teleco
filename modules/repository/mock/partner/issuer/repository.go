package issuer

import (
	partnerIssuerPort "github.com/sepulsa/teleco/business/partner/issuer/port"

	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func New() *Repository {
	return &Repository{}
}

func (db *Repository) CreateData(partnerIssuer partnerIssuerPort.PartnerIssuerRepo) error {
	result := db.Called(partnerIssuer)
	return result.Error(0)
}

func (db *Repository) FindByPartnerIssuerID(partnerId string, issuerId string) (partnerIssuerPort.PartnerIssuerRepo, error) {
	result := db.Called(partnerId, issuerId)
	return result.Get(0).(partnerIssuerPort.PartnerIssuerRepo), result.Error(1)
}

func (db *Repository) ReadData(ID string) (partnerIssuerPort.PartnerIssuerRepo, error) {
	result := db.Called(ID)
	return result.Get(0).(partnerIssuerPort.PartnerIssuerRepo), result.Error(1)
}

func (db *Repository) UpdateData(partnerIssuer partnerIssuerPort.PartnerIssuerRepo) error {
	result := db.Called(partnerIssuer)
	return result.Error(0)
}

func (db *Repository) DeleteData(ID string) error {
	result := db.Called(ID)
	return result.Error(0)
}

func (db *Repository) ListData() ([]partnerIssuerPort.PartnerIssuerRepo, error) {
	result := db.Called()
	return result.Get(0).([]partnerIssuerPort.PartnerIssuerRepo), result.Error(1)
}
