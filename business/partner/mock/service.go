package mock

import (
	partnerPort "github.com/sepulsa/teleco/business/partner/port"

	"github.com/stretchr/testify/mock"
)

type service struct {
	mock.Mock
}

func New() *service {
	return &service{}
}

func (s *service) CreateData(partner partnerPort.PartnerService) error {
	result := s.Called(partner)
	return result.Error(0)
}

func (s *service) ReadData(ID string) (partnerPort.PartnerService, error) {
	result := s.Called(ID)
	return result.Get(0).(partnerPort.PartnerService), result.Error(1)
}

func (s *service) UpdateData(partner partnerPort.PartnerService) error {
	result := s.Called(partner)
	return result.Error(0)
}

func (s *service) DeleteData(ID string) error {
	result := s.Called(ID)
	return result.Error(0)
}

func (s *service) ListData() ([]partnerPort.PartnerService, error) {
	result := s.Called()
	return result.Get(0).([]partnerPort.PartnerService), result.Error(1)
}
