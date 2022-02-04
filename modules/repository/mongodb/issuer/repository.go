package issuer

import (
	"encoding/json"
	"errors"
	"time"

	issuerPort "github.com/sepulsa/teleco/business/issuer/port"
	mongo "github.com/sepulsa/teleco/utils/mgo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Repository struct {
		mongo.Collection
	}

	Issuer struct {
		ID               bson.ObjectId `bson:"_id,omitempty"`
		Code             string        `bson:"code"`
		Label            string        `bson:"label"`
		Config           string        `bson:"config"`
		ThreadNum        int           `bson:"thread_num" json:"thread_num"`
		ThreadTimeout    int           `bson:"thread_timeout" json:"thread_timeout"`
		QueueWorkerLimit int           `bson:"queue_worker_limit" json:"queue_worker_limit"`
		CreatedAt        time.Time     `bson:"created_at"`
		UpdatedAt        time.Time     `bson:"updated_at"`
		DeletedAt        time.Time     `bson:"-,omitempty"`
	}
)

var (
	ErrInvalidID      = "Invalid ID"
	ErrIssuerNotFound = "Issuer not found"
)

func New(Mgo *mongo.MongoDatabase) *Repository {
	return &Repository{
		Mgo.C("issuer"),
	}
}

func (db *Repository) FindByCode(code string) (issuer issuerPort.IssuerRepo) {
	var data Issuer
	filterByCode := bson.M{
		"code": code,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	if err := db.Find(filterByCode).One(&data); err != nil {
		return
	}
	b, _ := json.Marshal(data)
	json.Unmarshal(b, &issuer)

	return
}

func (db *Repository) CreateData(issuer issuerPort.IssuerRepo) error {
	data := Issuer{
		Code:             issuer.Code,
		Label:            issuer.Label,
		Config:           issuer.Config,
		ThreadNum:        issuer.ThreadNum,
		ThreadTimeout:    issuer.ThreadTimeout,
		QueueWorkerLimit: issuer.QueueWorkerLimit,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return db.Insert(data)
}

func (db *Repository) ReadData(ID string) (issuer issuerPort.IssuerRepo, err error) {
	if !bson.IsObjectIdHex(ID) {
		err = errors.New(ErrInvalidID)
		return
	}

	var data Issuer
	filterByID := bson.M{
		"_id": bson.ObjectIdHex(ID),
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	if err = db.Find(filterByID).One(&data); err != nil {
		if err == mgo.ErrNotFound {
			err = errors.New(ErrIssuerNotFound)
		}
		return
	}
	b, _ := json.Marshal(data)
	json.Unmarshal(b, &issuer)

	return
}

func (db *Repository) UpdateData(issuer issuerPort.IssuerRepo) error {
	data := bson.M{
		"code":               issuer.Code,
		"label":              issuer.Label,
		"config":             issuer.Config,
		"thread_num":         issuer.ThreadNum,
		"thread_timeout":     issuer.ThreadTimeout,
		"queue_worker_limit": issuer.QueueWorkerLimit,
		"updated_at":         time.Now(),
	}
	return db.Update(bson.M{"_id": bson.ObjectIdHex(issuer.ID)}, bson.M{"$set": data})
}

func (db *Repository) DeleteData(ID string) error {
	if !bson.IsObjectIdHex(ID) {
		return errors.New(ErrInvalidID)
	}

	filter := bson.M{
		"_id": bson.ObjectIdHex(ID),
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	data := bson.M{
		"deleted_at": time.Now(),
	}
	if err := db.Update(filter, bson.M{"$set": data}); err != nil {
		if err == mgo.ErrNotFound {
			err = errors.New(ErrIssuerNotFound)
		}
		return err
	}
	return nil
}

func (db *Repository) ListData() (issuers []issuerPort.IssuerRepo, err error) {
	var data []Issuer

	filter := bson.M{
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	if err = db.Find(filter).All(&data); err != nil {
		if err == mgo.ErrNotFound {
			err = nil
		}
		return
	}

	d, _ := json.Marshal(data)
	json.Unmarshal(d, &issuers)

	return
}
