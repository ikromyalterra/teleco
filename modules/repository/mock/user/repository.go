package user

import (
	userPort "github.com/sepulsa/teleco/business/user/port"

	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func New() *Repository {
	return &Repository{}
}

func (db *Repository) FindByEmail(email string) userPort.UserRepo {
	result := db.Called(email)
	return result.Get(0).(userPort.UserRepo)
}

func (db *Repository) CreateData(user userPort.UserRepo) error {
	result := db.Called(user)
	return result.Error(0)
}

func (db *Repository) ReadData(ID string) (userPort.UserRepo, error) {
	result := db.Called(ID)
	return result.Get(0).(userPort.UserRepo), result.Error(1)
}

func (db *Repository) UpdateData(user userPort.UserRepo) error {
	result := db.Called(user)
	return result.Error(0)
}

func (db *Repository) DeleteData(ID string) error {
	result := db.Called(ID)
	return result.Error(0)
}

func (db *Repository) ListData() ([]userPort.UserRepo, error) {
	result := db.Called()
	return result.Get(0).([]userPort.UserRepo), result.Error(1)
}
