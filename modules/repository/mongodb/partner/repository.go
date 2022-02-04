package partner

import (
	"encoding/json"
	"errors"
	"time"

	partnerPort "github.com/sepulsa/teleco/business/partner/port"
	mongo "github.com/sepulsa/teleco/utils/mgo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Repository struct {
		mongo.Collection
	}

	Partner struct {
		ID          bson.ObjectId `bson:"_id,omitempty"`
		Code        string        `bson:"code"`
		Name        string        `bson:"name"`
		Pic         string        `bson:"pic"`
		Address     string        `bson:"address"`
		CallbackUrl string        `json:"callback_url" bson:"callback_url"`
		IpWhitelist []string      `json:"ip_whitelist" bson:"ip_whitelist"`
		Status      string        `bson:"status"`
		SecretKey   string        `json:"secret_key" bson:"secret_key"`
		CreatedAt   time.Time     `bson:"created_at"`
		UpdatedAt   time.Time     `bson:"updated_at"`
		DeletedAt   time.Time     `bson:"-,omitempty"`
	}
)

var (
	ErrInvalidID       = "Invalid ID"
	ErrPartnerNotFound = "Partner not found"
)

func New(Mgo *mongo.MongoDatabase) *Repository {
	return &Repository{
		Mgo.C("partner"),
	}
}

func (db *Repository) FindByCode(code string) (issuer partnerPort.PartnerRepo) {
	var data Partner
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

func (db *Repository) CreateData(partner partnerPort.PartnerRepo) error {
	insertData := Partner{
		Code:        partner.Code,
		Name:        partner.Name,
		Pic:         partner.Pic,
		Address:     partner.Address,
		CallbackUrl: partner.CallbackUrl,
		IpWhitelist: partner.IpWhitelist,
		Status:      partner.Status,
		SecretKey:   partner.SecretKey,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := db.Insert(insertData); err != nil {
		return err
	}
	return nil
}

func (db *Repository) ReadData(ID string) (partner partnerPort.PartnerRepo, err error) {
	if !bson.IsObjectIdHex(ID) {
		err = errors.New(ErrInvalidID)
		return
	}

	var data Partner
	filterByID := bson.M{
		"_id": bson.ObjectIdHex(ID),
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	if err = db.Find(filterByID).One(&data); err != nil {
		if err == mgo.ErrNotFound {
			err = errors.New(ErrPartnerNotFound)
		}
		return
	}
	b, _ := json.Marshal(data)
	json.Unmarshal(b, &partner)

	return
}

func (db *Repository) UpdateData(partner partnerPort.PartnerRepo) error {
	data := bson.M{
		"code":         partner.Code,
		"name":         partner.Name,
		"pic":          partner.Pic,
		"address":      partner.Address,
		"callback_url": partner.CallbackUrl,
		"ip_whitelist": partner.IpWhitelist,
		"status":       partner.Status,
		"secret_key":   partner.SecretKey,
		"updated_at":   time.Now(),
	}
	if err := db.Update(bson.M{"_id": bson.ObjectIdHex(partner.ID)}, bson.M{"$set": data}); err != nil {
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
			err = errors.New(ErrPartnerNotFound)
		}
		return err
	}
	return nil
}

func (db *Repository) ListData() (partners []partnerPort.PartnerRepo, err error) {
	var data []Partner

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
	json.Unmarshal(d, &partners)

	return
}
