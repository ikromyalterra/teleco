package order

import (
	orderPort "github.com/sepulsa/teleco/business/order/port"

	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func New() *Repository {
	return &Repository{}
}

func (db *Repository) CreateData(order orderPort.OrderRepo) error {
	result := db.Called(order)
	return result.Error(0)
}

func (db *Repository) UpdateData(order orderPort.OrderRepo) error {
	result := db.Called(order)
	return result.Error(0)
}
