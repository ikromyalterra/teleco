package mock

import (
	issuerPort "github.com/sepulsa/teleco/business/issuer/port"

	"github.com/stretchr/testify/mock"
)

type service struct {
	mock.Mock
}

func New() *service {
	return &service{}
}

func (s *service) CreateData(issuer issuerPort.IssuerService) error {
	result := s.Called(issuer)
	return result.Error(0)
}

func (s *service) ReadData(ID string) (issuerPort.IssuerService, error) {
	result := s.Called(ID)
	return result.Get(0).(issuerPort.IssuerService), result.Error(1)
}

func (s *service) UpdateData(issuer issuerPort.IssuerService) error {
	result := s.Called(issuer)
	return result.Error(0)
}

func (s *service) DeleteData(ID string) error {
	result := s.Called(ID)
	return result.Error(0)
}

func (s *service) ListData() ([]issuerPort.IssuerService, error) {
	result := s.Called()
	return result.Get(0).([]issuerPort.IssuerService), result.Error(1)
}
