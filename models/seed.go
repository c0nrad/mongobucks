package models

import "gopkg.in/mgo.v2/bson"

func SeedData() {
	session := Session.Clone()
	defer session.Close()

	session.DB(DB).C(RewardCollection).DropCollection()

	r := Reward{ID: bson.NewObjectId(), Name: "1 Highfive from Andrew Erlichson", Description: "Recieve a ticket that is redeemable for one highfive from Andrew Erlichson", Price: 10}
	session.DB(DB).C(RewardCollection).Insert(r)

	r = Reward{ID: bson.NewObjectId(), Name: "Lunch with Dev", Description: "Recieve a ticket that is redeemable for one solo lunch with Dev (only one a quarter)", Price: 1000}
	session.DB(DB).C(RewardCollection).Insert(r)

	r = Reward{ID: bson.NewObjectId(), Name: "Free Snacks from Kitchenette", Description: "Recieve a ticket that is redeemable for free snack from the kitchenette", Price: 5}
	session.DB(DB).C(RewardCollection).Insert(r)

	r = Reward{ID: bson.NewObjectId(), Name: "Free Vacation Day", Description: "Recieve a ticket that is redeemable for a vacation day (pending your managers approval)", Price: 5}
	session.DB(DB).C(RewardCollection).Insert(r)

	r = Reward{ID: bson.NewObjectId(), Name: "Free security consultation", Description: "Recieve a ticket that is redeemable for a security consultation via the Security Team (for mongodb related work)", Price: 10}
	session.DB(DB).C(RewardCollection).Insert(r)
}
