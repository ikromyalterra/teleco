package mock

import (
	userPort "github.com/sepulsa/teleco/business/user/port"

	"github.com/stretchr/testify/mock"
)

type service struct {
	mock.Mock
}

func New() *service {
	return &service{}
}

func (s *service) CreateData(user userPort.UserService) error {
	result := s.Called(user)
	return result.Error(0)
}

func (s *service) ReadData(ID string) (userPort.UserService, error) {
	result := s.Called(ID)
	return result.Get(0).(userPort.UserService), result.Error(1)
}

func (s *service) UpdateData(user userPort.UserService) error {
	result := s.Called(user)
	return result.Error(0)
}

func (s *service) DeleteData(ID string) error {
	result := s.Called(ID)
	return result.Error(0)
}

func (s *service) ListData() ([]userPort.UserService, error) {
	result := s.Called()
	return result.Get(0).([]userPort.UserService), result.Error(1)
}
