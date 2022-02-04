package mgo

import (
	"gopkg.in/mgo.v2"
)

type (
	// Query is an interface to access to the database struct
	Query interface {
		All(result interface{}) error
		One(result interface{}) error
		Sort(fields ...string) Query
		Limit(n int) Query
	}

	// Collection is an interface to access to the collection struct.
	Collection interface {
		Find(query interface{}) Query
		Insert(docs ...interface{}) error
		Update(selector interface{}, update interface{}) error
		UpdateAll(selector interface{}, update interface{}) (*mgo.ChangeInfo, error)
		Remove(selector interface{}) error
		DropIndexName(name string) error
		EnsureIndex(index mgo.Index) error
		EnsureIndexKey(key ...string) error
	}

	// DataLayer is an interface to access to the database struct.
	DataLayer interface {
		C(name string) Collection
	}

	// Session is an interface to access to the Session struct.
	Session interface {
		DB(name string) DataLayer
	}
)

type (
	// MongoQuery wraps a mgo.Query to embed methods in models.
	MongoQuery struct {
		*mgo.Query
	}

	// MongoCollection wraps a mgo.Collection to embed methods in models.
	MongoCollection struct {
		*mgo.Collection
	}

	// MongoDatabase wraps a mgo.Database to embed methods in models.
	MongoDatabase struct {
		*mgo.Database
	}

	// MongoSession is currently a Mongo session.
	MongoSession struct {
		*mgo.Session
	}
)

func (q MongoQuery) All(result interface{}) error {
	return q.Query.All(result)
}

func (q MongoQuery) One(result interface{}) error {
	return q.Query.One(result)
}

func (q MongoQuery) Sort(fields ...string) Query {
	return MongoQuery{Query: q.Query.Sort(fields...)}
}

func (q MongoQuery) Limit(n int) Query {
	return MongoQuery{Query: q.Query.Limit(n)}
}

// Find shadows *mgo.Collection to returns a Query interface instead of *mgo.Query.
func (c MongoCollection) Find(query interface{}) Query {
	return MongoQuery{Query: c.Collection.Find(query)}
}

func (c MongoCollection) Insert(docs ...interface{}) error {
	return c.Collection.Insert(docs...)
}

func (c MongoCollection) Update(selector interface{}, update interface{}) error {
	return c.Collection.Update(selector, update)
}

func (c MongoCollection) UpdateAll(selector interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	return c.Collection.UpdateAll(selector, update)
}

func (c MongoCollection) Remove(selector interface{}) error {
	return c.Collection.Remove(selector)
}

func (c MongoCollection) DropIndexName(name string) error {
	return c.Collection.DropIndexName(name)
}

func (c MongoCollection) EnsureIndex(index mgo.Index) error {
	return c.Collection.EnsureIndex(index)
}

func (c MongoCollection) EnsureIndexKey(key ...string) error {
	return c.Collection.EnsureIndexKey(key...)
}

// C shadows *mgo.DB to returns a DataLayer interface instead of *mgo.Database.
func (d MongoDatabase) C(name string) Collection {
	return &MongoCollection{Collection: d.Database.C(name)}
}

// DB shadows *mgo.DB to returns a DataLayer interface instead of *mgo.Database.
func (s MongoSession) DB(name string) DataLayer {
	return &MongoDatabase{Database: s.Session.DB(name)}
}
