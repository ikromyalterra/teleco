package order

import (
	"errors"

	issuerPort "github.com/sepulsa/teleco/business/issuer/port"
	orderPort "github.com/sepulsa/teleco/business/order/port"
	partnerIssuerPort "github.com/sepulsa/teleco/business/partner/issuer/port"
	partnerPort "github.com/sepulsa/teleco/business/partner/port"
)

type (
	service struct {
		issuerRepository        issuerPort.Repository
		partnerRepository       partnerPort.Repository
		partnerIssuerRepository partnerIssuerPort.Repository
		orderRepository         orderPort.Repository
		issuerApi               orderPort.IssuerApi
	}
)

var (
	ErrConfigNotFound = "Partner Issuer Config Not Found"
)

func New(issuerRepository issuerPort.Repository, partnerRepository partnerPort.Repository, partnerIssuerRepository partnerIssuerPort.Repository, orderRepository orderPort.Repository, issuerApi orderPort.IssuerApi) orderPort.Service {
	return &service{
		issuerRepository,
		partnerRepository,
		partnerIssuerRepository,
		orderRepository,
		issuerApi,
	}
}

func (s *service) Purchase(order orderPort.OrderService) (orderPort.OrderServiceResult, error) {
	// Get Data Partner Issuer
	issuerData := s.issuerRepository.FindByCode(order.IssuerCode)
	partnerData := s.partnerRepository.FindByCode(order.PartnerCode)
	partnerIssuerData, err := s.partnerIssuerRepository.FindByPartnerIssuerID(partnerData.ID, issuerData.ID)

	if err != nil {
		return orderPort.OrderServiceResult{}, errors.New(ErrConfigNotFound)
	}

	// Do Purchase
	orderIssuer := orderPort.OrderIssuerApi{
		ID:                  order.ID,
		CommandType:         orderPort.Purchase,
		IssuerCode:          order.IssuerCode,
		TransactionId:       order.TransactionId,
		IssuerProductId:     order.IssuerProductId,
		CustomerNumber:      order.CustomerNumber,
		PartnerIssuerConfig: partnerIssuerData.Config,
		IssuerConfig:        issuerData.Config,
		IssuerThreadNum:     issuerData.ThreadNum,
		IssuerThreadTimeout: issuerData.ThreadTimeout,
		PartnerCallbackUrl:  partnerData.CallbackUrl,
		PartnerId:           partnerData.ID,
		IssuerId:            issuerData.ID,
	}
	issuerResult, errApi := s.issuerApi.Do(orderIssuer)

	// Store Log Data
	orderData := orderPort.OrderRepo{
		CommandType:         orderPort.Purchase,
		TransactionId:       order.TransactionId,
		IssuerProductId:     order.IssuerProductId,
		CustomerNumber:      order.CustomerNumber,
		PartnerId:           partnerData.ID,
		IssuerId:            issuerData.ID,
		IssuerTransactionId: issuerResult.IssuerTransactionId,
		RequestData:         issuerResult.RequestData,
		ResponseData:        issuerResult.ResponseData,
	}
	s.orderRepository.CreateData(orderData)

	if errApi != nil {
		return orderPort.OrderServiceResult{}, errApi
	}

	result := orderPort.OrderServiceResult{
		IssuerTransactionId: issuerResult.IssuerTransactionId,
		SerialNumber:        issuerResult.SerialNumber,
		IssuerRescode:       issuerResult.IssuerRescode,
		Message:             issuerResult.Message,
		RawData:             issuerResult.RawData,
	}

	return result, nil
}

func (s *service) Advise(order orderPort.OrderService) (orderPort.OrderServiceResult, error) {
	// Get Data Partner Issuer
	issuerData := s.issuerRepository.FindByCode(order.IssuerCode)
	partnerData := s.partnerRepository.FindByCode(order.PartnerCode)
	partnerIssuerData, err := s.partnerIssuerRepository.FindByPartnerIssuerID(partnerData.ID, issuerData.ID)

	if err != nil {
		return orderPort.OrderServiceResult{}, errors.New(ErrConfigNotFound)
	}

	// Do Advise
	orderIssuer := orderPort.OrderIssuerApi{
		ID:                  order.ID,
		CommandType:         orderPort.Advise,
		IssuerCode:          order.IssuerCode,
		IssuerTransactionId: order.IssuerTransactionId,
		PartnerIssuerConfig: partnerIssuerData.Config,
		IssuerConfig:        issuerData.Config,
		IssuerThreadNum:     issuerData.ThreadNum,
		IssuerThreadTimeout: issuerData.ThreadTimeout,
		PartnerCallbackUrl:  partnerData.CallbackUrl,
		PartnerId:           partnerData.ID,
		IssuerId:            issuerData.ID,
	}
	issuerResult, errApi := s.issuerApi.Do(orderIssuer)

	// Store Log Data
	orderData := orderPort.OrderRepo{
		CommandType:         orderPort.Advise,
		IssuerTransactionId: issuerResult.IssuerTransactionId,
		PartnerId:           partnerData.ID,
		IssuerId:            issuerData.ID,
		RequestData:         issuerResult.RequestData,
		ResponseData:        issuerResult.ResponseData,
	}
	s.orderRepository.CreateData(orderData)

	if errApi != nil {
		return orderPort.OrderServiceResult{}, errApi
	}

	result := orderPort.OrderServiceResult{
		IssuerTransactionId: issuerResult.IssuerTransactionId,
		SerialNumber:        issuerResult.SerialNumber,
		IssuerRescode:       issuerResult.IssuerRescode,
		Message:             issuerResult.Message,
		RawData:             issuerResult.RawData,
	}

	return result, nil
}

func (s *service) Reversal(order orderPort.OrderService) (orderPort.OrderServiceResult, error) {
	// Get Data Partner Issuer
	issuerData := s.issuerRepository.FindByCode(order.IssuerCode)
	partnerData := s.partnerRepository.FindByCode(order.PartnerCode)
	partnerIssuerData, err := s.partnerIssuerRepository.FindByPartnerIssuerID(partnerData.ID, issuerData.ID)

	if err != nil {
		return orderPort.OrderServiceResult{}, errors.New(ErrConfigNotFound)
	}

	// Do Reversal
	orderIssuer := orderPort.OrderIssuerApi{
		ID:                  order.ID,
		CommandType:         orderPort.Reversal,
		PartnerIssuerConfig: partnerIssuerData.Config,
		IssuerConfig:        issuerData.Config,
		IssuerCode:          order.IssuerCode,
		IssuerTransactionId: order.IssuerTransactionId,
		IssuerThreadNum:     issuerData.ThreadNum,
		IssuerThreadTimeout: issuerData.ThreadTimeout,
		PartnerCallbackUrl:  partnerData.CallbackUrl,
		PartnerId:           partnerData.ID,
		IssuerId:            issuerData.ID,
	}
	issuerResult, errApi := s.issuerApi.Do(orderIssuer)

	// Store Log Data
	orderData := orderPort.OrderRepo{
		CommandType:         orderPort.Reversal,
		RequestData:         issuerResult.RequestData,
		ResponseData:        issuerResult.ResponseData,
		PartnerId:           partnerData.ID,
		IssuerId:            issuerData.ID,
		IssuerTransactionId: issuerResult.IssuerTransactionId,
	}
	s.orderRepository.CreateData(orderData)

	if errApi != nil {
		return orderPort.OrderServiceResult{}, errApi
	}

	result := orderPort.OrderServiceResult{
		Message:             issuerResult.Message,
		RawData:             issuerResult.RawData,
		IssuerTransactionId: issuerResult.IssuerTransactionId,
		IssuerRescode:       issuerResult.IssuerRescode,
	}

	return result, nil
}
