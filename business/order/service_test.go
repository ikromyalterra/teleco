package order_test

import (
	"errors"
	"testing"

	orderService "github.com/sepulsa/teleco/business/order"
	orderPort "github.com/sepulsa/teleco/business/order/port"
	issuerRepo "github.com/sepulsa/teleco/modules/repository/mock/issuer"
	orderRepo "github.com/sepulsa/teleco/modules/repository/mock/order"
	partnerRepo "github.com/sepulsa/teleco/modules/repository/mock/partner"
	partnerIssuerRepo "github.com/sepulsa/teleco/modules/repository/mock/partner/issuer"

	issuerPort "github.com/sepulsa/teleco/business/issuer/port"
	partnerIssuerPort "github.com/sepulsa/teleco/business/partner/issuer/port"
	partnerPort "github.com/sepulsa/teleco/business/partner/port"

	issuerApi "github.com/sepulsa/teleco/modules/issuerapi/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	ErrIssuerCodeNotFound = "Issuer Code Not Found"
	ErrConfigNotFound     = "Partner Issuer Config Not Found"
)

func TestPurchase(t *testing.T) {
	orderRepository := orderRepo.New()
	issuerRepository := issuerRepo.New()
	partnerRepository := partnerRepo.New()
	partnerIssuerRepository := partnerIssuerRepo.New()
	issuerApi := issuerApi.New()
	service := orderService.New(issuerRepository, partnerRepository, partnerIssuerRepository, orderRepository, issuerApi)

	// Error Partner Issuer Not found
	partnerRepository.On("FindByCode", mock.Anything).Return(partnerPort.PartnerRepo{ID: "12345"}).Once()
	issuerRepository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: "12345", Config: "\"\":\"\""}).Once()
	partnerIssuerRepository.On("FindByPartnerIssuerID", mock.Anything, mock.Anything).Return(partnerIssuerPort.PartnerIssuerRepo{}, errors.New(ErrConfigNotFound)).Once()
	_, err := service.Purchase(orderPort.OrderService{})
	assert.NotNil(t, err)
	assert.Equal(t, ErrConfigNotFound, err.Error())

	partnerRepository.On("FindByCode", mock.Anything).Return(partnerPort.PartnerRepo{ID: "12345"}).Once()
	issuerRepository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: "12345", Config: "\"\":\"\""}).Once()
	partnerIssuerRepository.On("FindByPartnerIssuerID", mock.Anything, mock.Anything).Return(partnerIssuerPort.PartnerIssuerRepo{Config: "{\"\":\"\"}"}, nil).Once()
	issuerApi.On("Do", mock.Anything).Return(orderPort.OrderIssuerApiResult{}, errors.New("test error")).Once()
	orderRepository.On("CreateData", mock.Anything).Return(nil).Once()
	_, err = service.Purchase(orderPort.OrderService{})
	assert.NotNil(t, err)

	partnerRepository.On("FindByCode", mock.Anything).Return(partnerPort.PartnerRepo{ID: "12345"}).Once()
	issuerRepository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: "12345", Config: "\"\":\"\""}).Once()
	partnerIssuerRepository.On("FindByPartnerIssuerID", mock.Anything, mock.Anything).Return(partnerIssuerPort.PartnerIssuerRepo{Config: "{\"\":\"\"}"}, nil).Once()
	issuerApi.On("Do", mock.Anything).Return(orderPort.OrderIssuerApiResult{}, nil).Once()
	orderRepository.On("CreateData", mock.Anything).Return(nil).Once()
	_, err = service.Purchase(orderPort.OrderService{})
	assert.Nil(t, err)
}

func TestAdvise(t *testing.T) {
	orderRepository := orderRepo.New()
	issuerRepository := issuerRepo.New()
	partnerRepository := partnerRepo.New()
	partnerIssuerRepository := partnerIssuerRepo.New()
	issuerApi := issuerApi.New()
	service := orderService.New(issuerRepository, partnerRepository, partnerIssuerRepository, orderRepository, issuerApi)

	// Error Partner Issuer Not found
	partnerRepository.On("FindByCode", mock.Anything).Return(partnerPort.PartnerRepo{ID: "12345"}).Once()
	issuerRepository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: "12345", Config: "\"\":\"\""}).Once()
	partnerIssuerRepository.On("FindByPartnerIssuerID", mock.Anything, mock.Anything).Return(partnerIssuerPort.PartnerIssuerRepo{}, errors.New(ErrConfigNotFound)).Once()
	_, err := service.Advise(orderPort.OrderService{})
	assert.NotNil(t, err)
	assert.Equal(t, ErrConfigNotFound, err.Error())

	partnerRepository.On("FindByCode", mock.Anything).Return(partnerPort.PartnerRepo{ID: "12345"}).Once()
	issuerRepository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: "12345", Config: "\"\":\"\""}).Once()
	partnerIssuerRepository.On("FindByPartnerIssuerID", mock.Anything, mock.Anything).Return(partnerIssuerPort.PartnerIssuerRepo{Config: "{\"\":\"\"}"}, nil).Once()
	issuerApi.On("Do", mock.Anything).Return(orderPort.OrderIssuerApiResult{}, errors.New("test error")).Once()
	orderRepository.On("CreateData", mock.Anything).Return(nil).Once()
	_, err = service.Advise(orderPort.OrderService{})
	assert.NotNil(t, err)

	partnerRepository.On("FindByCode", mock.Anything).Return(partnerPort.PartnerRepo{ID: "12345"}).Once()
	issuerRepository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: "12345", Config: "\"\":\"\""}).Once()
	partnerIssuerRepository.On("FindByPartnerIssuerID", mock.Anything, mock.Anything).Return(partnerIssuerPort.PartnerIssuerRepo{Config: "{\"\":\"\"}"}, nil).Once()
	issuerApi.On("Do", mock.Anything).Return(orderPort.OrderIssuerApiResult{}, nil).Once()
	orderRepository.On("CreateData", mock.Anything).Return(nil).Once()
	_, err = service.Advise(orderPort.OrderService{})
	assert.Nil(t, err)
}

func TestReversal(t *testing.T) {
	orderRepository := orderRepo.New()
	issuerRepository := issuerRepo.New()
	partnerRepository := partnerRepo.New()
	partnerIssuerRepository := partnerIssuerRepo.New()
	issuerApi := issuerApi.New()
	service := orderService.New(issuerRepository, partnerRepository, partnerIssuerRepository, orderRepository, issuerApi)

	// Error Partner Issuer Not found
	partnerRepository.On("FindByCode", mock.Anything).Return(partnerPort.PartnerRepo{ID: "12345"}).Once()
	issuerRepository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: "12345", Config: "\"\":\"\""}).Once()
	partnerIssuerRepository.On("FindByPartnerIssuerID", mock.Anything, mock.Anything).Return(partnerIssuerPort.PartnerIssuerRepo{}, errors.New(ErrConfigNotFound)).Once()
	_, err := service.Reversal(orderPort.OrderService{})
	assert.NotNil(t, err)
	assert.Equal(t, ErrConfigNotFound, err.Error())

	partnerRepository.On("FindByCode", mock.Anything).Return(partnerPort.PartnerRepo{ID: "12345"}).Once()
	issuerRepository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: "12345", Config: "\"\":\"\""}).Once()
	partnerIssuerRepository.On("FindByPartnerIssuerID", mock.Anything, mock.Anything).Return(partnerIssuerPort.PartnerIssuerRepo{Config: "{\"\":\"\"}"}, nil).Once()
	issuerApi.On("Do", mock.Anything).Return(orderPort.OrderIssuerApiResult{}, errors.New("test error")).Once()
	orderRepository.On("CreateData", mock.Anything).Return(nil).Once()
	_, err = service.Reversal(orderPort.OrderService{})
	assert.NotNil(t, err)

	partnerRepository.On("FindByCode", mock.Anything).Return(partnerPort.PartnerRepo{ID: "12345"}).Once()
	issuerRepository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: "12345", Config: "\"\":\"\""}).Once()
	partnerIssuerRepository.On("FindByPartnerIssuerID", mock.Anything, mock.Anything).Return(partnerIssuerPort.PartnerIssuerRepo{Config: "{\"\":\"\"}"}, nil).Once()
	issuerApi.On("Do", mock.Anything).Return(orderPort.OrderIssuerApiResult{}, nil).Once()
	orderRepository.On("CreateData", mock.Anything).Return(nil).Once()
	_, err = service.Reversal(orderPort.OrderService{})
	assert.Nil(t, err)
}
