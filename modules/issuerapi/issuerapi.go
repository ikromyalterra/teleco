package issuerapi

import (
	orderPort "github.com/sepulsa/teleco/business/order/port"
	"github.com/sepulsa/teleco/modules/issuerapi/task"
	"github.com/sepulsa/teleco/utils/threadpool"
)

type (
	issuerApi struct{}
)

func New() *issuerApi {
	return &issuerApi{}
}

func (is *issuerApi) Do(order orderPort.OrderIssuerApi) (orderResult orderPort.OrderIssuerApiResult, err error) {
	var errPort orderPort.Error
	ot := &task.OrderTask{Order: order, OrderResult: &orderResult, Err: &errPort}
	threadpool.Run(order.IssuerCode, order.IssuerThreadNum, order.IssuerThreadTimeout, ot)
	return orderResult, errPort.Err
}
