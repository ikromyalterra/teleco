package config

import (
	"github.com/spf13/viper"

	//driver
	"gopkg.in/mgo.v2"

	mongo "github.com/sepulsa/teleco/utils/mgo"
)

var (
	//Mgo Mgo
	Mgo *mongo.MongoDatabase
)

//Database Database
type Database struct {
	Host              string
	User              string
	Password          string
	DBName            string
	DBNumber          int
	Port              int
	DebugMode         bool
	ReconnectRetry    int
	ReconnectInterval int64
}

// LoadDBConfig load database configuration
func LoadDBConfig(name string) Database {
	db := viper.Sub("database." + name)
	conf := Database{
		Host:      db.GetString("host"),
		User:      db.GetString("user"),
		Password:  db.GetString("password"),
		DBName:    db.GetString("db_name"),
		Port:      db.GetInt("port"),
		DebugMode: db.GetBool("debug"),
	}
	return conf
}

//MongoConnect MongoConnect
func MongoConnect() {
	if viper.Get("env") != "testing" {
		conf := LoadDBConfig("mongo")
		mongoConf := &mgo.DialInfo{
			Addrs:    []string{conf.Host},
			Username: conf.User,
			Password: conf.Password,
			Database: conf.DBName,
		}
		session, err := mgo.DialWithInfo(mongoConf)
		if err != nil {
			panic(err)
		}
		s := mongo.MongoSession{Session: session}
		Mgo = &mongo.MongoDatabase{Database: s.Session.DB(conf.DBName)}
	}
}
