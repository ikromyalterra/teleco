package task

import (
	"encoding/json"

	orderPort "github.com/sepulsa/teleco/business/order/port"
	"github.com/sepulsa/teleco/modules/callback"
	orderRepository "github.com/sepulsa/teleco/modules/repository/mongodb/order"
	"github.com/sepulsa/teleco/utils/config"
	log "github.com/sepulsa/teleco/utils/logger"
)

var (
	packageLog = "teleco/modules/issuerapi/task"
)

type WorkerTask struct{}

func (t *WorkerTask) Run(payload string) {
	var order orderPort.OrderIssuerApi
	var orderResult orderPort.OrderIssuerApiResult
	var errs orderPort.Error
	json.Unmarshal([]byte(payload), &order)

	issuerApi, err := getIssuerAPI(order.IssuerCode)
	if err != nil {
		errs.Err = err
		return
	}
	switch order.CommandType {
	case orderPort.Purchase:
		issuerApi.Purchase(order, &orderResult, &errs)
	case orderPort.Advise:
		issuerApi.Advise(order, &orderResult, &errs)
	case orderPort.Reversal:
		issuerApi.Reversal(order, &orderResult, &errs)
	}

	callbackPort := callback.New()
	callBackResult := callbackPort.Do(order, orderResult)
	db := config.Mgo
	orderRepo := orderRepository.New(db)
	// Store Log Data
	orderData := orderPort.OrderRepo{
		CommandType:          order.CommandType,
		TransactionId:        order.TransactionId,
		IssuerProductId:      order.IssuerProductId,
		CustomerNumber:       order.CustomerNumber,
		PartnerId:            order.PartnerId,
		IssuerId:             order.IssuerId,
		IssuerTransactionId:  orderResult.IssuerTransactionId,
		RequestData:          orderResult.RequestData,
		ResponseData:         orderResult.ResponseData,
		CallbackRequestData:  callBackResult.RequestData,
		CallbackResponseData: callBackResult.ResponseData,
	}
	orderRepo.CreateData(orderData)
	log.Info().Str("event", "queue.executed").Str("package", packageLog).Msgf("Payload: %s", payload)
}
