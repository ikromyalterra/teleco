package usertoken

import (
	"errors"
	"time"

	authPort "github.com/sepulsa/teleco/business/auth/port"
	mongo "github.com/sepulsa/teleco/utils/mgo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Repository struct {
		mongo.Collection
	}

	UserToken struct {
		ID        bson.ObjectId `bson:"_id,omitempty"`
		UserID    string        `bson:"user_id"`
		TokenID   string        `bson:"token_id"`
		CreatedAt time.Time     `bson:"created_at"`
		UpdatedAt time.Time     `bson:"updated_at"`
		DeletedAt time.Time     `bson:"-,omitempty"`
	}
)

var (
	ErrInvalidID         = "Invalid ID"
	ErrUserTokenNotFound = "Token not found"
)

func New(Mgo *mongo.MongoDatabase) *Repository {
	return &Repository{
		Mgo.C("usertoken"),
	}
}

func (db *Repository) CreateData(jwt authPort.UserTokenRepo) error {
	var data UserToken
	data.UserID = jwt.UserID
	data.TokenID = jwt.TokenID
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	return db.Insert(data)
}

func (db *Repository) FindByTokenID(tokenID string) authPort.UserTokenRepo {
	var userToken authPort.UserTokenRepo
	var data UserToken
	filterByTokenID := bson.M{
		"token_id": tokenID,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	if err := db.Find(filterByTokenID).One(&data); err == nil {
		userToken.UserID = data.UserID
		userToken.TokenID = data.TokenID
	}

	return userToken
}

func (db *Repository) DeleteData(tokenID string) error {
	filter := bson.M{
		"token_id": tokenID,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	data := bson.M{
		"deleted_at": time.Now(),
	}
	if err := db.Update(filter, bson.M{"$set": data}); err != nil {
		if err == mgo.ErrNotFound {
			err = errors.New(ErrUserTokenNotFound)
		}
		return err
	}
	return nil
}

func (db *Repository) DeleteDataByUserID(userID string) error {
	filter := bson.M{
		"user_id": userID,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	data := bson.M{
		"deleted_at": time.Now(),
	}
	if _, err := db.UpdateAll(filter, bson.M{"$set": data}); err != nil {
		if err == mgo.ErrNotFound {
			err = errors.New(ErrUserTokenNotFound)
		}
		return err
	}
	return nil
}
