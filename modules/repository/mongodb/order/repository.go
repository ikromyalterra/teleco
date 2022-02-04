package order

import (
	"time"

	orderPort "github.com/sepulsa/teleco/business/order/port"
	mongo "github.com/sepulsa/teleco/utils/mgo"
	"gopkg.in/mgo.v2/bson"
)

type (
	Repository struct {
		mongo.Collection
	}

	Order struct {
		ID                   bson.ObjectId `bson:"_id,omitempty"`
		CommandType          string        `bson:"command_type" json:"command_type"`
		TransactionId        string        `bson:"transaction_id" json:"transaction_id"`
		IssuerProductId      string        `bson:"issuer_product_id" json:"issuer_product_id"`
		CustomerNumber       string        `bson:"customer_number" json:"customer_number"`
		PartnerId            string        `bson:"partner_id" json:"partner_id"`
		IssuerId             string        `bson:"issuer_id" json:"issuer_id"`
		IssuerTransactionId  string        `bson:"issuer_transaction_id" json:"issuer_transaction_id"`
		RequestData          string        `bson:"request_data" json:"request_data"`
		ResponseData         string        `bson:"response_data" json:"response_data"`
		CallbackRequestData  string        `bson:"callback_request_data"`
		CallbackResponseData string        `bson:"callback_response_data"`
		CreatedAt            time.Time     `bson:"created_at" json:"created_at"`
		UpdatedAt            time.Time     `bson:"updated_at" json:"update_id"`
		DeletedAt            time.Time     `bson:"-,omitempty" json:"deleted_at"`
	}
)

var (
	ErrInvalidID     = "Invalid ID"
	ErrOrderNotFound = "Partner Issuer not found"
)

func New(Mgo *mongo.MongoDatabase) *Repository {
	return &Repository{
		Mgo.C("order"),
	}
}

func (db *Repository) CreateData(order orderPort.OrderRepo) error {
	data := Order{
		CommandType:          order.CommandType,
		IssuerTransactionId:  order.IssuerTransactionId,
		TransactionId:        order.TransactionId,
		IssuerProductId:      order.IssuerProductId,
		CustomerNumber:       order.CustomerNumber,
		PartnerId:            order.PartnerId,
		IssuerId:             order.IssuerId,
		RequestData:          order.RequestData,
		ResponseData:         order.ResponseData,
		CallbackRequestData:  order.CallbackRequestData,
		CallbackResponseData: order.CallbackResponseData,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	if err := db.Insert(data); err != nil {
		return err
	}
	return nil
}
