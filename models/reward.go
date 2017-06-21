package models

import "gopkg.in/mgo.v2/bson"

const (
	NullReward   int = iota
	RedeemReward int = iota
)

type Reward struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	Name        string
	Description string

	IsHidden   bool
	RewardType int

	Price int
}

func GetRewards() ([]*Reward, error) {
	session := Session.Copy()
	defer session.Close()

	var out []*Reward

	err := session.DB(DB).C(RewardCollection).Find(bson.M{"ishidden": false}).All(&out)
	return out, err
}

func NewReward(name, description string, hidden bool, rewardType, price int) (*Reward, error) {
	session := Session.Copy()
	defer session.Close()

	reward := &Reward{ID: bson.NewObjectId(), Name: name, Description: description, IsHidden: hidden, RewardType: rewardType, Price: price}
	err := session.DB(DB).C(RewardCollection).Insert(reward)

	if err != nil {
		return nil, err
	}
	return reward, err
}

func GetRewardById(id bson.ObjectId) (*Reward, error) {
	session := Session.Copy()
	defer session.Close()

	var reward Reward
	err := session.DB(DB).C(RewardCollection).FindId(id).One(&reward)
	return &reward, err
}

func (r Reward) Redeem(user *User) error {
	if r.RewardType == RedeemReward {
		user.Balance += r.Price
		return user.Update()
	}

	return nil
}
