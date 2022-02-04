package mock

import (
	orderPort "github.com/sepulsa/teleco/business/order/port"
	"github.com/stretchr/testify/mock"
)

type (
	issuerApi struct {
		mock.Mock
	}
)

func New() *issuerApi {
	return &issuerApi{}
}

func (is *issuerApi) Do(order orderPort.OrderIssuerApi) (orderPort.OrderIssuerApiResult, error) {
	result := is.Called(order)
	return result.Get(0).(orderPort.OrderIssuerApiResult), result.Error(1)
}
