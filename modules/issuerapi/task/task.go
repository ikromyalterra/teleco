package task

import (
	"encoding/json"
	"errors"

	orderPort "github.com/sepulsa/teleco/business/order/port"
	"github.com/sepulsa/teleco/modules/callback"
	"github.com/sepulsa/teleco/modules/issuerapi/issuer/dummy"
	orderRepository "github.com/sepulsa/teleco/modules/repository/mongodb/order"
	"github.com/sepulsa/teleco/utils/config"
	"github.com/sepulsa/teleco/utils/queue/producer"
)

type OrderTask struct {
	Order       orderPort.OrderIssuerApi
	OrderResult *orderPort.OrderIssuerApiResult
	Err         *orderPort.Error
}

var (
	ErrTimeout            = "Transaction still in progress"
	ErrConcurrentLimit    = "Concurrent limit reached. Transaction added to queue process."
	ErrIssuerCodeNotFound = "Issuer API Not Found"
)

func getIssuerAPI(issuerCode string) (orderPort.Issuer, error) {
	switch issuerCode {
	case "dummy":
		return dummy.New(), nil
	default:
		return nil, errors.New(ErrIssuerCodeNotFound)
	}
}

func (t *OrderTask) Run() {
	issuerApi, err := getIssuerAPI(t.Order.IssuerCode)
	if err != nil {
		t.Err.Err = err
		return
	}
	switch t.Order.CommandType {
	case orderPort.Purchase:
		issuerApi.Purchase(t.Order, t.OrderResult, t.Err)
	case orderPort.Advise:
		issuerApi.Advise(t.Order, t.OrderResult, t.Err)
	case orderPort.Reversal:
		issuerApi.Reversal(t.Order, t.OrderResult, t.Err)
	}
}

func (t *OrderTask) RunWhenTimeout() {
	t.OrderResult.IssuerRescode = "10"
	t.OrderResult.Message = ErrTimeout
}

func (t *OrderTask) RunAfterTimeout() {
	callbackPort := callback.New()
	callBackResult := callbackPort.Do(t.Order, *t.OrderResult)
	db := config.Mgo
	orderRepo := orderRepository.New(db)
	// Store Log Data
	orderData := orderPort.OrderRepo{
		CommandType:          t.Order.CommandType,
		TransactionId:        t.Order.TransactionId,
		IssuerProductId:      t.Order.IssuerProductId,
		CustomerNumber:       t.Order.CustomerNumber,
		PartnerId:            t.Order.PartnerId,
		IssuerId:             t.Order.IssuerId,
		IssuerTransactionId:  t.OrderResult.IssuerTransactionId,
		RequestData:          t.OrderResult.RequestData,
		ResponseData:         t.OrderResult.ResponseData,
		CallbackRequestData:  callBackResult.RequestData,
		CallbackResponseData: callBackResult.ResponseData,
	}
	orderRepo.CreateData(orderData)
}

func (t *OrderTask) RunWhenFull() {
	t.OrderResult.IssuerRescode = "10"
	t.OrderResult.Message = ErrConcurrentLimit
	js, _ := json.Marshal(t.Order)
	str := string(js)
	queueName := "teleco_" + t.Order.IssuerCode
	producer.Queue.CreateItem(queueName, str)
}
