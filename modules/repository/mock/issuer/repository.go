package issuer

import (
	issuerPort "github.com/sepulsa/teleco/business/issuer/port"

	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func New() *Repository {
	return &Repository{}
}

func (db *Repository) FindByCode(code string) issuerPort.IssuerRepo {
	result := db.Called(code)
	return result.Get(0).(issuerPort.IssuerRepo)
}

func (db *Repository) CreateData(issuer issuerPort.IssuerRepo) error {
	result := db.Called(issuer)
	return result.Error(0)
}

func (db *Repository) ReadData(ID string) (issuerPort.IssuerRepo, error) {
	result := db.Called(ID)
	return result.Get(0).(issuerPort.IssuerRepo), result.Error(1)
}

func (db *Repository) UpdateData(issuer issuerPort.IssuerRepo) error {
	result := db.Called(issuer)
	return result.Error(0)
}

func (db *Repository) DeleteData(ID string) error {
	result := db.Called(ID)
	return result.Error(0)
}

func (db *Repository) ListData() ([]issuerPort.IssuerRepo, error) {
	result := db.Called()
	return result.Get(0).([]issuerPort.IssuerRepo), result.Error(1)
}
