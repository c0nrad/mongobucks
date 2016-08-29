package models

import (
	"errors"
	"time"

	"github.com/pborman/uuid"

	"gopkg.in/mgo.v2/bson"
)

type Ticket struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	TS time.Time

	Reward   bson.ObjectId
	Username string

	Name string

	Redemption string
	IsUsed     bool
}

func Redeem(redemption string) error {
	session := Session.Clone()
	defer session.Close()

	err := session.DB(DB).C(TicketCollection).Update(bson.M{"redemption": redemption}, bson.M{"$set": bson.M{"isused": true}})

	return err
}

func GetTicketByToken(token string) (*Ticket, error) {
	session := Session.Clone()
	defer session.Close()

	var ticket Ticket
	err := session.DB(DB).C(TicketCollection).Find(bson.M{"redemption": token}).One(&ticket)

	return &ticket, err
}

func PurchaseTicket(user *User, reward *Reward) (*Ticket, error) {
	session := Session.Clone()
	defer session.Close()

	if user.Balance < reward.Price {
		return nil, errors.New("Insufficent funds")
	}

	user.Balance -= reward.Price
	err := user.Update()
	if err != nil {
		return nil, err
	}

	t := Ticket{ID: bson.NewObjectId(), TS: time.Now(), Reward: reward.ID, Username: user.Username,
		Name: reward.Name, Redemption: uuid.New(), IsUsed: false}

	err = session.DB(DB).C(TicketCollection).Insert(t)
	return &t, err
}

func GetTicketsByUsername(username string) ([]*Ticket, error) {
	session := Session.Clone()
	defer session.Close()

	var out []*Ticket
	err := session.DB(DB).C(TicketCollection).Find(bson.M{"username": username}).All(&out)
	return out, err
}
