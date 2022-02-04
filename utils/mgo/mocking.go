package mgo

import "gopkg.in/mgo.v2"

type (
	// MockQuery satisfies Query and act as a mock.
	MockQuery struct{}

	// MockCollection satisfies Collection and act as a mock.
	MockCollection struct{}

	// MockDatabase satisfies DataLayer and act as a mock.
	MockDatabase struct{}

	// MockSession satisfies Session and act as a mock of *mgo.session.
	MockSession struct{}
)

// All mock.
func (fq MockQuery) All(result interface{}) error {
	return nil
}

// One mock.
func (fq MockQuery) One(result interface{}) error {
	return nil
}

func (fq MockQuery) Sort(fields ...string) Query {
	return fq
}

func (fq MockQuery) Limit(n int) Query {
	return fq
}

// Find mock.
func (fc MockCollection) Find(query interface{}) Query {
	return MockQuery{}
}

// Insert mock.
func (fc MockCollection) Insert(docs ...interface{}) error {
	return nil
}

// Update mock.
func (fc MockCollection) Update(selector interface{}, update interface{}) error {
	return nil
}

func (fc MockCollection) UpdateAll(selector interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	return nil, nil
}

func (fc MockCollection) Remove(selector interface{}) error {
	return nil
}

func (fc MockCollection) DropIndexName(name string) error {
	return nil
}

func (c MockCollection) EnsureIndex(index mgo.Index) error {
	return nil
}

func (fc MockCollection) EnsureIndexKey(key ...string) error {
	return nil
}

// C mocks mgo.Database(name).Collection(name).
func (db MockDatabase) C(name string) Collection {
	return MockCollection{}
}

// DB mocks mgo.Session.DB().
func (fs MockSession) DB(name string) DataLayer {
	return MockDatabase{}
}

// NewMockSession mock NewSession.
func NewMockSession() Session {
	mockSession := MockSession{}
	return mockSession
}
