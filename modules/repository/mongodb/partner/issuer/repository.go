package issuer

import (
	"encoding/json"
	"errors"
	"time"

	partnerIssuerPort "github.com/sepulsa/teleco/business/partner/issuer/port"
	mongo "github.com/sepulsa/teleco/utils/mgo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Repository struct {
		mongo.Collection
	}

	PartnerIssuer struct {
		ID        bson.ObjectId `bson:"_id,omitempty"`
		PartnerId string        `bson:"partner_id" json:"partner_id"`
		IssuerId  string        `bson:"issuer_id" json:"issuer_id"`
		Config    string        `bson:"config" json:"config"`
		CreatedAt time.Time     `bson:"created_at" json:"created_at"`
		UpdatedAt time.Time     `bson:"updated_at" json:"update_id"`
		DeletedAt time.Time     `bson:"-,omitempty" json:"deleted_at"`
	}
)

var (
	ErrInvalidID             = "Invalid ID"
	ErrPartnerIssuerNotFound = "Partner Issuer not found"
)

func New(Mgo *mongo.MongoDatabase) *Repository {
	return &Repository{
		Mgo.C("partnerIssuer"),
	}
}

func (db *Repository) CreateData(partnerIssuer partnerIssuerPort.PartnerIssuerRepo) error {
	data := PartnerIssuer{
		PartnerId: partnerIssuer.PartnerId,
		IssuerId:  partnerIssuer.IssuerId,
		Config:    partnerIssuer.Config,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Insert(data); err != nil {
		return err
	}
	return nil
}

func (db *Repository) FindByPartnerIssuerID(partnerId string, issuerId string) (partnerIssuer partnerIssuerPort.PartnerIssuerRepo, err error) {
	var data PartnerIssuer
	filterByPartnerIssuerID := bson.M{
		"partner_id": partnerId,
		"issuer_id":  issuerId,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	if err = db.Find(filterByPartnerIssuerID).One(&data); err != nil {
		if err == mgo.ErrNotFound {
			err = errors.New(ErrPartnerIssuerNotFound)
		}
		return
	}
	b, _ := json.Marshal(data)
	json.Unmarshal(b, &partnerIssuer)

	return
}

func (db *Repository) ReadData(ID string) (partnerIssuer partnerIssuerPort.PartnerIssuerRepo, err error) {
	if !bson.IsObjectIdHex(ID) {
		err = errors.New(ErrInvalidID)
		return
	}

	var data PartnerIssuer
	filterByID := bson.M{
		"_id": bson.ObjectIdHex(ID),
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	if err = db.Find(filterByID).One(&data); err != nil {
		if err == mgo.ErrNotFound {
			err = errors.New(ErrPartnerIssuerNotFound)
		}
		return
	}
	b, _ := json.Marshal(data)
	json.Unmarshal(b, &partnerIssuer)

	return
}

func (db *Repository) UpdateData(partnerIssuer partnerIssuerPort.PartnerIssuerRepo) error {
	data := bson.M{
		"partner_id": partnerIssuer.PartnerId,
		"issuer_id":  partnerIssuer.IssuerId,
		"config":     partnerIssuer.Config,
		"updated_at": time.Now(),
	}
	if err := db.Update(bson.M{"_id": bson.ObjectIdHex(partnerIssuer.ID)}, bson.M{"$set": data}); err != nil {
		return err
	}
	return nil
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
			err = errors.New(ErrPartnerIssuerNotFound)
		}
		return err
	}
	return nil
}

func (db *Repository) ListData() (partnerIssuers []partnerIssuerPort.PartnerIssuerRepo, err error) {
	var data []PartnerIssuer

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
	json.Unmarshal(d, &partnerIssuers)

	return
}
