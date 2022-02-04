package mock

import (
	partnerIssuerPort "github.com/sepulsa/teleco/business/partner/issuer/port"

	"github.com/stretchr/testify/mock"
)

type service struct {
	mock.Mock
}

func New() *service {
	return &service{}
}

func (s *service) CreateData(partnerIssuer partnerIssuerPort.PartnerIssuerService) error {
	result := s.Called(partnerIssuer)
	return result.Error(0)
}

func (s *service) ReadData(ID string) (partnerIssuerPort.PartnerIssuerService, error) {
	result := s.Called(ID)
	return result.Get(0).(partnerIssuerPort.PartnerIssuerService), result.Error(1)
}

func (s *service) UpdateData(partnerIssuer partnerIssuerPort.PartnerIssuerService) error {
	result := s.Called(partnerIssuer)
	return result.Error(0)
}

func (s *service) DeleteData(ID string) error {
	result := s.Called(ID)
	return result.Error(0)
}

func (s *service) ListData() ([]partnerIssuerPort.PartnerIssuerService, error) {
	result := s.Called()
	return result.Get(0).([]partnerIssuerPort.PartnerIssuerService), result.Error(1)
}
