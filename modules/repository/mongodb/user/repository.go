package user

import (
	"errors"
	"time"

	userPort "github.com/sepulsa/teleco/business/user/port"
	mongo "github.com/sepulsa/teleco/utils/mgo"
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

type (
	Repository struct {
		mongo.Collection
	}

	User struct {
		ID        bson.ObjectId `bson:"_id,omitempty"`
		Email     string        `bson:"email"`
		Password  string        `bson:"password"`
		Fullname  string        `bson:"fullname"`
		CreatedAt time.Time     `bson:"created_at"`
		UpdatedAt time.Time     `bson:"updated_at"`
		DeletedAt time.Time     `bson:"-,omitempty"`
	}
)

var (
	ErrUserNotFound error = errors.New("user not found")
	ErrInvalidID    error = errors.New("invalid id")
)

func New(Mgo *mongo.MongoDatabase) *Repository {
	return &Repository{
		Mgo.C("user"),
	}
}

func (db *Repository) FindByEmail(email string) userPort.UserRepo {
	var data User
	var user userPort.UserRepo

	filterByEmail := bson.M{
		"email": email,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	if err := db.Find(filterByEmail).One(&data); err != nil {
		return user
	}

	user.ID = data.ID.Hex()
	user.Email = data.Email
	user.Fullname = data.Fullname
	user.Password = data.Password

	return user
}

func (db *Repository) ReadData(ID string) (userPort.UserRepo, error) {
	var data User
	var user userPort.UserRepo

	if !bson.IsObjectIdHex(ID) {
		return user, ErrInvalidID
	}

	filterByID := bson.M{
		"_id": bson.ObjectIdHex(ID),
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	if err := db.Find(filterByID).One(&data); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrUserNotFound
		}
		return user, err
	}

	user.ID = data.ID.Hex()
	user.Email = data.Email
	user.Fullname = data.Fullname
	user.Password = data.Password

	return user, nil
}

func (db *Repository) CreateData(user userPort.UserRepo) error {
	var data User
	data.Email = user.Email
	data.Fullname = user.Fullname
	data.Password = user.Password
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	return db.Insert(data)
}

func (db *Repository) UpdateData(user userPort.UserRepo) error {
	data := make(bson.M)
	data["email"] = user.Email
	data["fullname"] = user.Fullname
	data["updated_at"] = time.Now()

	if user.Password != "" {
		data["password"] = user.Password
	}

	return db.Update(bson.M{"_id": bson.ObjectIdHex(user.ID)}, bson.M{"$set": data})
}

func (db *Repository) DeleteData(ID string) error {
	if !bson.IsObjectIdHex(ID) {
		return ErrInvalidID
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
			err = ErrUserNotFound
		}
		return err
	}
	return nil
}

func (db *Repository) ListData() ([]userPort.UserRepo, error) {
	var datas []User
	users := make([]userPort.UserRepo, 0)

	filter := bson.M{
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	if err := db.Find(filter).All(&datas); err != nil {
		if err == mgo.ErrNotFound {
			err = nil
		}
		return users, err
	}

	var user userPort.UserRepo

	for i := range datas {
		user.ID = datas[i].ID.Hex()
		user.Email = datas[i].Email
		user.Fullname = datas[i].Fullname
		users = append(users, user)
	}

	return users, nil
}
