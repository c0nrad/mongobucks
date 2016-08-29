package models

import "gopkg.in/mgo.v2/bson"

type Reward struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	Name        string
	Description string

	Price int
}

func GetRewards() ([]*Reward, error) {
	session := Session.Clone()
	defer session.Close()

	var out []*Reward

	err := session.DB(DB).C(RewardCollection).Find(nil).All(&out)
	return out, err
}

func GetRewardById(id bson.ObjectId) (*Reward, error) {
	session := Session.Clone()
	defer session.Close()

	var reward Reward
	err := session.DB(DB).C(RewardCollection).FindId(id).One(&reward)
	return &reward, err
}
