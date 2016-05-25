package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	Username string
	Name     string
	Picture  string

	Balance int
}

func (user *User) Update() error {
	session := Session.Clone()
	defer session.Close()

	err := session.DB(DB).C(UserCollection).Update(bson.M{"username": user.Username}, bson.M{"$set": bson.M{"balance": user.Balance}})
	return err
}

func FindOrCreateUser(username string) (*User, error) {
	session := Session.Clone()
	defer session.Close()

	// Atomically safe, but this will probably fail a lot. index on username
	session.DB(DB).C(UserCollection).Insert(bson.M{"username": username, "balance": StartingBalance})

	var user User
	err := session.DB(DB).C(UserCollection).Find(bson.M{"username": username}).One(&user)
	return &user, err
}

func FindUser(username string) (*User, error) {
	session := Session.Clone()
	defer session.Close()

	var user User
	err := session.DB(DB).C(UserCollection).Find(bson.M{"username": username}).One(&user)
	return &user, err
}

func GetUsers() ([]*User, error) {
	session := Session.Clone()
	defer session.Close()

	var out []*User

	err := session.DB(DB).C(UserCollection).Find(nil).All(&out)
	return out, err

}

func GetBalance(username string) (int, error) {
	user, err := FindOrCreateUser(username)
	if err != nil {
		return 0, err
	}
	return user.Balance, nil
}

func (user *User) UpdateProfile(name, picture string) error {
	user.Name = name
	user.Picture = picture

	err := Session.DB(DB).C(UserCollection).UpdateId(user.ID, bson.M{"$set": bson.M{"name": name, "picture": picture}})
	return err
}
