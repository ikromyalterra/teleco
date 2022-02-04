package partner

import (
	partnerPort "github.com/sepulsa/teleco/business/partner/port"

	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func New() *Repository {
	return &Repository{}
}

func (db *Repository) FindByCode(code string) partnerPort.PartnerRepo {
	result := db.Called(code)
	return result.Get(0).(partnerPort.PartnerRepo)
}

func (db *Repository) CreateData(partnerIssuer partnerPort.PartnerRepo) error {
	result := db.Called(partnerIssuer)
	return result.Error(0)
}

func (db *Repository) ReadData(ID string) (partnerPort.PartnerRepo, error) {
	result := db.Called(ID)
	return result.Get(0).(partnerPort.PartnerRepo), result.Error(1)
}

func (db *Repository) UpdateData(partnerIssuer partnerPort.PartnerRepo) error {
	result := db.Called(partnerIssuer)
	return result.Error(0)
}

func (db *Repository) DeleteData(ID string) error {
	result := db.Called(ID)
	return result.Error(0)
}

func (db *Repository) ListData() ([]partnerPort.PartnerRepo, error) {
	result := db.Called()
	return result.Get(0).([]partnerPort.PartnerRepo), result.Error(1)
}
