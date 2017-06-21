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

func Redeem(redemption string, user *User) error {
	session := Session.Copy()
	defer session.Close()

	ticket, err := GetTicketByToken(redemption)
	if err != nil {
		return err
	}

	info, err := session.DB(DB).C(TicketCollection).UpdateAll(bson.M{"redemption": redemption, "isused": false}, bson.M{"$set": bson.M{"isused": true}})
	if err != nil {
		return err
	}

	if info.Matched != 1 {
		return errors.New("ticket already used")
	}

	reward, err := GetRewardById(ticket.Reward)
	if err != nil {
		return err
	}

	err = reward.Redeem(user)
	return err
}

func GetTicketByToken(token string) (*Ticket, error) {
	session := Session.Copy()
	defer session.Close()

	var ticket Ticket
	err := session.DB(DB).C(TicketCollection).Find(bson.M{"redemption": token}).One(&ticket)

	return &ticket, err
}

func PurchaseTicket(user *User, reward *Reward) (*Ticket, error) {
	session := Session.Copy()
	defer session.Close()

	if user.Balance < reward.Price {
		return nil, errors.New("Insufficent funds")
	}

	user.Balance -= reward.Price
	err := user.Update()
	if err != nil {
		return nil, err
	}

	return NewTicket(reward.ID, user.Username, reward.Name, uuid.New())
}

func NewTicket(rewardId bson.ObjectId, username, rewardName, redemption string) (*Ticket, error) {
	session := Session.Copy()
	defer session.Close()

	t := Ticket{ID: bson.NewObjectId(), TS: time.Now(), Reward: rewardId, Username: username,
		Name: rewardName, Redemption: redemption, IsUsed: false}

	err := session.DB(DB).C(TicketCollection).Insert(t)
	return &t, err
}

func GetTicketsByUsername(username string) ([]*Ticket, error) {
	session := Session.Copy()
	defer session.Close()

	var out []*Ticket
	err := session.DB(DB).C(TicketCollection).Find(bson.M{"username": username}).All(&out)
	return out, err
}
