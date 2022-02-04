package partner_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	partnerService "github.com/sepulsa/teleco/business/partner"
	partnerPort "github.com/sepulsa/teleco/business/partner/port"
	partnerRepo "github.com/sepulsa/teleco/modules/repository/mock/partner"
)

var (
	TestPartnerID1          = "613f12af2edd3a56323f0d1d"
	TestPartnerCode1        = "p1"
	TestPartnerName1        = "partner1"
	TestPartnerPic1         = "pic1"
	TestPartnerAddress1     = "address1"
	TestPartnerCallbackUrl1 = "callback_url1"
	TestPartnerIpwhitelist1 = []string{"192.168.1.1"}
	TestPartnerStatus1      = "active"
	TestPartnerSecretKey1   = "SECRETKEY"

	TestErrInvalidID = "Invalid ID"
)

func TestCreateData(t *testing.T) {
	repository := partnerRepo.New()

	dataService := partnerPort.PartnerService{
		ID:          "",
		Code:        TestPartnerCode1,
		Name:        TestPartnerName1,
		Pic:         TestPartnerPic1,
		Address:     TestPartnerAddress1,
		CallbackUrl: TestPartnerCallbackUrl1,
		IpWhitelist: TestPartnerIpwhitelist1,
		Status:      TestPartnerStatus1,
		SecretKey:   TestPartnerSecretKey1,
	}

	// success
	repository.On("CreateData", mock.Anything).Return(nil).Once()
	service := partnerService.New(repository)
	err := service.CreateData(dataService)
	assert.Nil(t, err)

	// error mongo
	repository.On("CreateData", mock.Anything).Return(errors.New("")).Once()
	service = partnerService.New(repository)
	err = service.CreateData(dataService)
	assert.NotNil(t, err)
}

func TestReadData(t *testing.T) {
	repository := partnerRepo.New()

	id := TestPartnerID1
	dataRepo := partnerPort.PartnerRepo{
		ID:          TestPartnerID1,
		Code:        TestPartnerCode1,
		Name:        TestPartnerName1,
		Pic:         TestPartnerPic1,
		Address:     TestPartnerAddress1,
		CallbackUrl: TestPartnerCallbackUrl1,
		IpWhitelist: TestPartnerIpwhitelist1,
		Status:      TestPartnerStatus1,
		SecretKey:   TestPartnerSecretKey1,
	}

	t.Log(dataRepo)

	// success
	repository.On("ReadData", mock.Anything).Return(dataRepo, nil).Once()
	service := partnerService.New(repository)
	partner, err := service.ReadData(id)
	if assert.Nil(t, err) {
		assert.Equal(t, dataRepo.ID, partner.ID)
		assert.Equal(t, dataRepo.Code, partner.Code)
		assert.Equal(t, dataRepo.Name, partner.Name)
		assert.Equal(t, dataRepo.Pic, partner.Pic)
		assert.Equal(t, dataRepo.Address, partner.Address)
		assert.Equal(t, dataRepo.CallbackUrl, partner.CallbackUrl)
		assert.Equal(t, dataRepo.IpWhitelist, partner.IpWhitelist)
		assert.Equal(t, dataRepo.Status, partner.Status)
	}

	// error
	repository.On("ReadData", mock.Anything).Return(partnerPort.PartnerRepo{}, errors.New(TestErrInvalidID)).Once()
	service = partnerService.New(repository)
	_, err = service.ReadData(id)
	assert.Equal(t, TestErrInvalidID, err.Error())
}

func TestUpdateData(t *testing.T) {
	repository := partnerRepo.New()

	dataService := partnerPort.PartnerService{
		ID:          TestPartnerID1,
		Code:        TestPartnerCode1 + TestPartnerID1,
		Name:        TestPartnerName1,
		Pic:         TestPartnerPic1,
		Address:     TestPartnerAddress1,
		CallbackUrl: TestPartnerCallbackUrl1,
		IpWhitelist: TestPartnerIpwhitelist1,
		Status:      TestPartnerStatus1,
		SecretKey:   TestPartnerSecretKey1,
	}

	dataRepo := partnerPort.PartnerRepo{
		ID:          TestPartnerID1,
		Name:        TestPartnerName1,
		Pic:         TestPartnerPic1,
		Address:     TestPartnerAddress1,
		CallbackUrl: TestPartnerCallbackUrl1,
		IpWhitelist: TestPartnerIpwhitelist1,
		Status:      TestPartnerStatus1,
		SecretKey:   TestPartnerSecretKey1,
	}

	// success
	repository.On("ReadData", mock.Anything).Return(dataRepo, nil).Once()
	repository.On("UpdateData", mock.Anything).Return(nil).Once()
	service := partnerService.New(repository)
	err := service.UpdateData(dataService)
	assert.Nil(t, err)

	// error invalid id
	repository.On("ReadData", mock.Anything).Return(partnerPort.PartnerRepo{}, errors.New(TestErrInvalidID)).Once()
	service = partnerService.New(repository)
	err = service.UpdateData(dataService)
	assert.Equal(t, TestErrInvalidID, err.Error())

	// error mongo
	repository.On("ReadData", mock.Anything).Return(dataRepo, nil).Once()
	repository.On("UpdateData", mock.Anything).Return(errors.New("")).Once()
	service = partnerService.New(repository)
	err = service.UpdateData(dataService)
	assert.NotNil(t, err)
}

func TestDeleteData(t *testing.T) {
	repository := partnerRepo.New()

	id := TestPartnerID1

	// success
	repository.On("DeleteData", mock.Anything).Return(nil).Once()
	service := partnerService.New(repository)
	err := service.DeleteData(id)
	assert.Nil(t, err)

	// error
	repository.On("DeleteData", mock.Anything).Return(errors.New(TestErrInvalidID)).Once()
	service = partnerService.New(repository)
	err = service.DeleteData(id)
	assert.Equal(t, TestErrInvalidID, err.Error())
}

func TestListData(t *testing.T) {
	repository := partnerRepo.New()

	dataRepo := []partnerPort.PartnerRepo{
		{
			ID:          TestPartnerID1,
			Name:        TestPartnerName1,
			Pic:         TestPartnerPic1,
			Address:     TestPartnerAddress1,
			CallbackUrl: TestPartnerCallbackUrl1,
			IpWhitelist: TestPartnerIpwhitelist1,
			Status:      TestPartnerStatus1,
			SecretKey:   TestPartnerSecretKey1,
		},
	}

	// success
	repository.On("ListData").Return(dataRepo, nil).Once()
	service := partnerService.New(repository)
	partners, err := service.ListData()
	if assert.Nil(t, err) {
		assert.Equal(t, TestPartnerID1, partners[0].ID)
		assert.Equal(t, TestPartnerName1, partners[0].Name)
		assert.Equal(t, TestPartnerPic1, partners[0].Pic)
		assert.Equal(t, TestPartnerAddress1, partners[0].Address)
		assert.Equal(t, TestPartnerCallbackUrl1, partners[0].CallbackUrl)
		assert.Equal(t, TestPartnerIpwhitelist1, partners[0].IpWhitelist)
		assert.Equal(t, TestPartnerStatus1, partners[0].Status)
		assert.Equal(t, TestPartnerSecretKey1, partners[0].SecretKey)

	}

	// error
	repository.On("ListData").Return([]partnerPort.PartnerRepo{}, errors.New("")).Once()
	service = partnerService.New(repository)
	_, err = service.ListData()
	assert.NotNil(t, err)
}
