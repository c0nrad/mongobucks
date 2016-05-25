package models

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Transaction struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	From   string
	To     string
	Amount int

	Memo string

	TS time.Time
}

func (t *Transaction) Insert() error {
	session := Session.Clone()
	defer session.Close()

	err := session.DB(DB).C(TransactionCollection).Insert(t)
	return err
}

func FindTransaction(id string) (*Transaction, error) {
	session := Session.Clone()
	defer session.Close()

	var transaction Transaction
	err := session.DB(DB).C(TransactionCollection).FindId(bson.ObjectIdHex(id)).One(&transaction)
	return &transaction, err
}

func GetTransactionsForUser(username string) ([]*Transaction, error) {
	session := Session.Clone()
	defer session.Close()

	var transactions []*Transaction
	err := session.DB(DB).C(TransactionCollection).Find(bson.M{"$or": []bson.M{{"from": username}, {"to": username}}}).All(&transactions)
	return transactions, err
}

func GetRecentTransactions() ([]*Transaction, error) {
	session := Session.Clone()
	defer session.Close()

	var transactions []*Transaction
	err := session.DB(DB).C(TransactionCollection).Find(nil).Sort("-ts").Limit(5).All(&transactions)
	return transactions, err
}

func ExecuteTransfer(from, to string, amount int, memo string) (*Transaction, error) {
	fmt.Println("[+] Executing Transfer", from, to, amount, memo)

	if from == to {
		return nil, errors.New("can't give to self")
	}

	fromUser, err := FindUser(from)
	if err != nil {
		return nil, err
	}

	toUser, err := FindUser(to)
	if err != nil {
		return nil, err
	}

	if fromUser.Balance < amount {
		return nil, errors.New("insufficent funds")
	}

	// There's probably race conditions, so we'll just keep the audit log first, and if something bad happens we'll recalculate manually
	t := Transaction{ID: bson.NewObjectId(), From: from, To: to, Amount: amount, TS: time.Now(), Memo: memo}
	err = t.Insert()
	if err != nil {
		return nil, err
	}

	fromUser.Balance -= amount
	toUser.Balance += amount

	err = fromUser.Update()
	if err != nil {
		return nil, err
	}

	err = toUser.Update()
	if err != nil {
		return nil, err
	}

	return &t, nil

}
