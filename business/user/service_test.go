package user_test

import (
	"errors"
	"testing"

	mockUserRepo "github.com/sepulsa/teleco/modules/repository/mock/user"
	mockUserTokenRepo "github.com/sepulsa/teleco/modules/repository/mock/usertoken"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	userServ "github.com/sepulsa/teleco/business/user"
	"github.com/sepulsa/teleco/business/user/port"
)

var (
	TestID       = "6138813fb95630b0b528b160"
	TestIDErr    = "1 2 3"
	TestEmail    = "test@test.com"
	TestFullname = "Test Fullname"
	TestPassword = "test"

	UserID  = "615321f9b95630b0b529867c"
	TokenID = "6328402206636078"

	TestErrInvalidID      error = errors.New("invalid id")
	TestErrDuplicateEmail error = errors.New(userServ.ErrDuplicateEmail)
)

func TestCreateData(t *testing.T) {
	userRepo := mockUserRepo.New()
	userTokenRepo := mockUserTokenRepo.New()

	// payload success
	dataUserServ := port.UserService{
		Email:    TestEmail,
		Fullname: TestFullname,
		Password: TestPassword,
	}

	// start populate mock data
	dataUserRepo := port.UserRepo{
		ID:       TestID,
		Email:    TestEmail,
		Fullname: TestFullname,
	}

	userRepo.On("FindByEmail", dataUserServ.Email).Return(dataUserRepo).Once()
	userRepo.On("FindByEmail", dataUserServ.Email).Return(port.UserRepo{})
	userRepo.On("CreateData", mock.Anything).Return(nil)
	// end populate

	// inject
	s := userServ.New(userRepo, userTokenRepo)

	// err duplicate email
	err := s.CreateData(dataUserServ)
	assert.Equal(t, TestErrDuplicateEmail.Error(), err.Error())

	// success
	assert.Nil(t, s.CreateData(dataUserServ))
}

func TestUpdateData(t *testing.T) {
	userRepo := mockUserRepo.New()
	userTokenRepo := mockUserTokenRepo.New()

	// payload success
	dataUserServ := port.UserService{
		ID:       TestID,
		Email:    "test@email.new",
		Fullname: TestFullname,
		Password: TestPassword,
	}

	// start populate mock data
	dataUserRepo := port.UserRepo{
		ID:       TestID,
		Email:    TestEmail,
		Fullname: TestFullname,
	}

	userRepo.On("ReadData", dataUserServ.ID).Return(dataUserRepo, errors.New("")).Once()
	userRepo.On("ReadData", dataUserServ.ID).Return(dataUserRepo, nil)
	userRepo.On("FindByEmail", dataUserServ.Email).Return(dataUserRepo).Once()
	userRepo.On("FindByEmail", dataUserServ.Email).Return(port.UserRepo{})
	userRepo.On("UpdateData", mock.Anything).Return(nil)
	// end populate

	// inject
	s := userServ.New(userRepo, userTokenRepo)

	// err repo, invalid id or user not found
	assert.NotNil(t, s.UpdateData(dataUserServ))

	// err duplicate email
	err := s.UpdateData(dataUserServ)
	assert.Equal(t, TestErrDuplicateEmail.Error(), err.Error())

	// success
	assert.Nil(t, s.UpdateData(dataUserServ))
}

func TestReadData(t *testing.T) {
	userRepo := mockUserRepo.New()
	userTokenRepo := mockUserTokenRepo.New()

	// start populate mock data
	dataUserRepo := port.UserRepo{
		ID:       TestID,
		Email:    TestEmail,
		Fullname: TestFullname,
	}
	userRepo.On("ReadData", TestID).Return(dataUserRepo, nil)
	userRepo.On("ReadData", TestIDErr).Return(dataUserRepo, TestErrInvalidID)
	// end populate

	// inject
	s := userServ.New(userRepo, userTokenRepo)

	// success
	user, err := s.ReadData(TestID)
	if assert.Nil(t, err) {
		assert.Equal(t, dataUserRepo.ID, user.ID)
		assert.Equal(t, dataUserRepo.Email, user.Email)
		assert.Equal(t, dataUserRepo.Fullname, user.Fullname)
	}

	// failed
	_, err = s.ReadData(TestIDErr)
	assert.Equal(t, TestErrInvalidID.Error(), err.Error())
}

func TestDeleteData(t *testing.T) {
	// populate data user
	userRepo := mockUserRepo.New()
	userRepo.On("DeleteData", TestID).Return(nil)
	userRepo.On("DeleteData", TestIDErr).Return(TestErrInvalidID)
	// populate data user token
	userTokenRepo := mockUserTokenRepo.New()
	userTokenRepo.On("DeleteDataByUserID", TestID).Return(nil)
	// end populate

	// inject
	s := userServ.New(userRepo, userTokenRepo)

	// success
	err := s.DeleteData(TestID)
	assert.Nil(t, err)

	// failed
	err = s.DeleteData(TestIDErr)
	assert.Equal(t, TestErrInvalidID.Error(), err.Error())
}

func TestListData(t *testing.T) {
	// populate data user
	datasUserRepo := []port.UserRepo{
		{
			ID:       TestID,
			Email:    TestEmail,
			Fullname: TestFullname,
		},
	}
	userRepo := mockUserRepo.New()
	userRepo.On("ListData").Return(datasUserRepo, nil).Once()
	userRepo.On("ListData").Return(datasUserRepo, errors.New("")).Once()

	userTokenRepo := mockUserTokenRepo.New()
	// end populate

	// inject
	s := userServ.New(userRepo, userTokenRepo)

	// success
	users, err := s.ListData()
	if assert.Nil(t, err) {
		assert.Equal(t, TestID, users[0].ID)
		assert.Equal(t, TestEmail, users[0].Email)
		assert.Equal(t, TestFullname, users[0].Fullname)
	}

	// failed
	_, err = s.ListData()
	assert.NotNil(t, err)
}
