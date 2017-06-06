package models

import (
	"math/rand"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Gamble struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	User       string
	GuessHeads bool
	Amount     int

	IsWinner bool

	TS time.Time
}

func (g *Gamble) Insert() error {
	session := Session.Copy()
	defer session.Close()

	err := session.DB(DB).C(GambleCollection).Insert(g)
	return err
}

func ExecuteGamble(user string, amount int, guessHeads bool) (*Gamble, error) {
	isHeads := rand.Intn(2) == 1

	g := &Gamble{ID: bson.NewObjectId(), User: user, GuessHeads: guessHeads, Amount: amount, IsWinner: isHeads == guessHeads, TS: time.Now()}

	err := g.Insert()
	return g, err
}

func FindGamble(id string) (*Gamble, error) {
	session := Session.Copy()
	defer session.Close()

	var gamble Gamble
	err := session.DB(DB).C(GambleCollection).FindId(bson.ObjectIdHex(id)).One(&gamble)
	return &gamble, err
}

func GetGamblesForUser(username string) ([]*Gamble, error) {
	session := Session.Copy()
	defer session.Close()

	var gambles []*Gamble
	err := session.DB(DB).C(GambleCollection).Find(bson.M{"user": username}).All(&gambles)
	return gambles, err
}

func GetRecentGambles() ([]*Gamble, error) {
	session := Session.Copy()
	defer session.Close()

	var gambles []*Gamble
	err := session.DB(DB).C(GambleCollection).Find(nil).Sort("-ts").Limit(5).All(&gambles)
	return gambles, err
}
