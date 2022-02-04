package issuer_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	partnerIssuerService "github.com/sepulsa/teleco/business/partner/issuer"
	partnerIssuerPort "github.com/sepulsa/teleco/business/partner/issuer/port"
	partnerIssuerRepo "github.com/sepulsa/teleco/modules/repository/mock/partner/issuer"
)

var (
	TestID        = "6138813fb95630b0b528b160"
	TestPartnerID = "6138813fb95630b0b528b162"
	TestIssuerID  = "6138813fb95630b0b528b163"
	TestConfig    = "config"

	TestErrInvalidID = "Invalid ID"
)

func TestCreateData(t *testing.T) {
	repository := partnerIssuerRepo.New()

	dataService := partnerIssuerPort.PartnerIssuerService{
		ID:        "",
		PartnerId: TestPartnerID,
		IssuerId:  TestIssuerID,
		Config:    TestConfig,
	}

	// success
	repository.On("CreateData", mock.Anything).Return(nil).Once()
	service := partnerIssuerService.New(repository)
	err := service.CreateData(dataService)
	assert.Nil(t, err)

	// error mongo
	repository.On("CreateData", mock.Anything).Return(errors.New("")).Once()
	service = partnerIssuerService.New(repository)
	err = service.CreateData(dataService)
	assert.NotNil(t, err)
}

func TestReadData(t *testing.T) {
	repository := partnerIssuerRepo.New()

	id := TestID
	dataRepo := partnerIssuerPort.PartnerIssuerRepo{
		ID:        TestID,
		PartnerId: TestPartnerID,
		IssuerId:  TestIssuerID,
		Config:    TestConfig,
	}

	// success
	repository.On("ReadData", mock.Anything).Return(dataRepo, nil).Once()
	service := partnerIssuerService.New(repository)
	partnerIssuer, err := service.ReadData(id)
	if assert.Nil(t, err) {
		assert.Equal(t, dataRepo.ID, partnerIssuer.ID)
		assert.Equal(t, dataRepo.PartnerId, partnerIssuer.PartnerId)
		assert.Equal(t, dataRepo.IssuerId, partnerIssuer.IssuerId)
	}

	// error
	repository.On("ReadData", mock.Anything).Return(partnerIssuerPort.PartnerIssuerRepo{}, errors.New(TestErrInvalidID)).Once()
	service = partnerIssuerService.New(repository)
	_, err = service.ReadData(id)
	assert.Equal(t, TestErrInvalidID, err.Error())
}

func TestUpdateData(t *testing.T) {
	repository := partnerIssuerRepo.New()

	dataService := partnerIssuerPort.PartnerIssuerService{
		ID:        TestID,
		PartnerId: TestPartnerID + TestPartnerID,
		IssuerId:  TestIssuerID,
		Config:    TestConfig,
	}

	dataRepo := partnerIssuerPort.PartnerIssuerRepo{
		ID:        TestID,
		PartnerId: TestPartnerID,
		IssuerId:  TestIssuerID,
		Config:    TestConfig,
	}

	// success
	repository.On("ReadData", mock.Anything).Return(dataRepo, nil).Once()
	repository.On("UpdateData", mock.Anything).Return(nil).Once()
	service := partnerIssuerService.New(repository)
	err := service.UpdateData(dataService)
	assert.Nil(t, err)

	// error invalid id
	repository.On("ReadData", mock.Anything).Return(partnerIssuerPort.PartnerIssuerRepo{}, errors.New(TestErrInvalidID)).Once()
	service = partnerIssuerService.New(repository)
	err = service.UpdateData(dataService)
	assert.Equal(t, TestErrInvalidID, err.Error())

	// error mongo
	repository.On("ReadData", mock.Anything).Return(dataRepo, nil).Once()
	repository.On("UpdateData", mock.Anything).Return(errors.New("")).Once()
	service = partnerIssuerService.New(repository)
	err = service.UpdateData(dataService)
	assert.NotNil(t, err)
}

func TestDeleteData(t *testing.T) {
	repository := partnerIssuerRepo.New()

	id := TestID

	// success
	repository.On("DeleteData", mock.Anything).Return(nil).Once()
	service := partnerIssuerService.New(repository)
	err := service.DeleteData(id)
	assert.Nil(t, err)

	// error
	repository.On("DeleteData", mock.Anything).Return(errors.New(TestErrInvalidID)).Once()
	service = partnerIssuerService.New(repository)
	err = service.DeleteData(id)
	assert.Equal(t, TestErrInvalidID, err.Error())
}

func TestListData(t *testing.T) {
	repository := partnerIssuerRepo.New()

	dataRepo := []partnerIssuerPort.PartnerIssuerRepo{
		{
			ID:        TestID,
			PartnerId: TestPartnerID,
			IssuerId:  TestIssuerID,
			Config:    TestConfig,
		},
	}

	// success
	repository.On("ListData").Return(dataRepo, nil).Once()
	service := partnerIssuerService.New(repository)
	partnerIssuers, err := service.ListData()
	if assert.Nil(t, err) {
		assert.Equal(t, TestID, partnerIssuers[0].ID)
		assert.Equal(t, TestPartnerID, partnerIssuers[0].PartnerId)
		assert.Equal(t, TestIssuerID, partnerIssuers[0].IssuerId)
		assert.Equal(t, TestConfig, partnerIssuers[0].Config)
	}

	// error
	repository.On("ListData").Return([]partnerIssuerPort.PartnerIssuerRepo{}, errors.New("")).Once()
	service = partnerIssuerService.New(repository)
	_, err = service.ListData()
	assert.NotNil(t, err)
}
