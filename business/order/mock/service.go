package mock

import (
	orderPort "github.com/sepulsa/teleco/business/order/port"

	"github.com/stretchr/testify/mock"
)

type service struct {
	mock.Mock
}

func New() *service {
	return &service{}
}

func (s *service) Purchase(order orderPort.OrderService) (orderPort.OrderServiceResult, error) {
	result := s.Called(order)
	return result.Get(0).(orderPort.OrderServiceResult), result.Error(1)
}

func (s *service) Advise(order orderPort.OrderService) (orderPort.OrderServiceResult, error) {
	result := s.Called(order)
	return result.Get(0).(orderPort.OrderServiceResult), result.Error(1)
}

func (s *service) Reversal(order orderPort.OrderService) (orderPort.OrderServiceResult, error) {
	result := s.Called(order)
	return result.Get(0).(orderPort.OrderServiceResult), result.Error(1)
}
