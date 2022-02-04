package usertoken

import (
	authPort "github.com/sepulsa/teleco/business/auth/port"

	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func New() *Repository {
	return &Repository{}
}

func (db *Repository) FindByTokenID(tokenID string) authPort.UserTokenRepo {
	result := db.Called(tokenID)
	return result.Get(0).(authPort.UserTokenRepo)
}

func (db *Repository) CreateData(userToken authPort.UserTokenRepo) error {
	result := db.Called(userToken)
	return result.Error(0)
}

func (db *Repository) DeleteData(tokenID string) error {
	result := db.Called(tokenID)
	return result.Error(0)
}

func (db *Repository) DeleteDataByUserID(userID string) error {
	result := db.Called(userID)
	return result.Error(0)
}
