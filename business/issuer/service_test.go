package issuer_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	issuerService "github.com/sepulsa/teleco/business/issuer"
	issuerPort "github.com/sepulsa/teleco/business/issuer/port"
	issuerRepo "github.com/sepulsa/teleco/modules/repository/mock/issuer"
)

var (
	TestID            = "6138813fb95630b0b528b160"
	TestCode          = "issuer"
	TestLabel         = "issuer testing"
	TestThreadNum     = 5
	TestThreadTimeout = 30

	TestErrInvalidID = "Invalid ID"
)

func TestCreateData(t *testing.T) {
	repository := issuerRepo.New()

	dataService := issuerPort.IssuerService{
		ID:            "",
		Code:          TestCode,
		Label:         TestLabel,
		ThreadNum:     TestThreadNum,
		ThreadTimeout: TestThreadTimeout,
	}

	// success
	repository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: ""}).Once()
	repository.On("CreateData", mock.Anything).Return(nil).Once()
	service := issuerService.New(repository)
	err := service.CreateData(dataService)
	assert.Nil(t, err)

	// duplicate code
	repository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: TestID}).Once()
	service = issuerService.New(repository)
	err = service.CreateData(dataService)
	assert.Equal(t, issuerService.ErrDuplicateCode, err.Error())

	// error mongo
	repository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: ""}).Once()
	repository.On("CreateData", mock.Anything).Return(errors.New("")).Once()
	service = issuerService.New(repository)
	err = service.CreateData(dataService)
	assert.NotNil(t, err)
}

func TestReadData(t *testing.T) {
	repository := issuerRepo.New()

	id := TestID
	dataRepo := issuerPort.IssuerRepo{
		ID:            TestID,
		Code:          TestCode,
		Label:         TestLabel,
		ThreadNum:     TestThreadNum,
		ThreadTimeout: TestThreadTimeout,
	}

	// success
	repository.On("ReadData", mock.Anything).Return(dataRepo, nil).Once()
	service := issuerService.New(repository)
	issuer, err := service.ReadData(id)
	if assert.Nil(t, err) {
		assert.Equal(t, dataRepo.ID, issuer.ID)
		assert.Equal(t, dataRepo.Code, issuer.Code)
		assert.Equal(t, dataRepo.Label, issuer.Label)
	}

	// error
	repository.On("ReadData", mock.Anything).Return(issuerPort.IssuerRepo{}, errors.New(TestErrInvalidID)).Once()
	service = issuerService.New(repository)
	_, err = service.ReadData(id)
	assert.Equal(t, TestErrInvalidID, err.Error())
}

func TestUpdateData(t *testing.T) {
	repository := issuerRepo.New()

	dataService := issuerPort.IssuerService{
		ID:            TestID,
		Code:          TestCode + TestCode,
		Label:         TestLabel,
		ThreadNum:     TestThreadNum,
		ThreadTimeout: TestThreadTimeout,
	}

	dataRepo := issuerPort.IssuerRepo{
		ID:            TestID,
		Code:          TestCode,
		Label:         TestLabel,
		ThreadNum:     TestThreadNum,
		ThreadTimeout: TestThreadTimeout,
	}

	// success
	repository.On("ReadData", mock.Anything).Return(dataRepo, nil).Once()
	repository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: ""}).Once()
	repository.On("UpdateData", mock.Anything).Return(nil).Once()
	service := issuerService.New(repository)
	err := service.UpdateData(dataService)
	assert.Nil(t, err)

	// error invalid id
	repository.On("ReadData", mock.Anything).Return(issuerPort.IssuerRepo{}, errors.New(TestErrInvalidID)).Once()
	service = issuerService.New(repository)
	err = service.UpdateData(dataService)
	assert.Equal(t, TestErrInvalidID, err.Error())

	// error duplicate
	repository.On("ReadData", mock.Anything).Return(dataRepo, nil).Once()
	repository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: TestID}).Once()
	service = issuerService.New(repository)
	err = service.UpdateData(dataService)
	assert.Equal(t, issuerService.ErrDuplicateCode, err.Error())

	// error mongo
	repository.On("ReadData", mock.Anything).Return(dataRepo, nil).Once()
	repository.On("FindByCode", mock.Anything).Return(issuerPort.IssuerRepo{ID: ""}).Once()
	repository.On("UpdateData", mock.Anything).Return(errors.New("")).Once()
	service = issuerService.New(repository)
	err = service.UpdateData(dataService)
	assert.NotNil(t, err)
}

func TestDeleteData(t *testing.T) {
	repository := issuerRepo.New()

	id := TestID

	// success
	repository.On("DeleteData", mock.Anything).Return(nil).Once()
	service := issuerService.New(repository)
	err := service.DeleteData(id)
	assert.Nil(t, err)

	// error
	repository.On("DeleteData", mock.Anything).Return(errors.New(TestErrInvalidID)).Once()
	service = issuerService.New(repository)
	err = service.DeleteData(id)
	assert.Equal(t, TestErrInvalidID, err.Error())
}

func TestListData(t *testing.T) {
	repository := issuerRepo.New()

	dataRepo := []issuerPort.IssuerRepo{
		{
			ID:            TestID,
			Code:          TestCode,
			Label:         TestLabel,
			ThreadNum:     TestThreadNum,
			ThreadTimeout: TestThreadTimeout,
		},
	}

	// success
	repository.On("ListData").Return(dataRepo, nil).Once()
	service := issuerService.New(repository)
	issuers, err := service.ListData()
	if assert.Nil(t, err) {
		assert.Equal(t, TestID, issuers[0].ID)
		assert.Equal(t, TestCode, issuers[0].Code)
		assert.Equal(t, TestLabel, issuers[0].Label)
	}

	// error
	repository.On("ListData").Return([]issuerPort.IssuerRepo{}, errors.New("")).Once()
	service = issuerService.New(repository)
	_, err = service.ListData()
	assert.NotNil(t, err)
}
