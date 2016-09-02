package models

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"

	"gopkg.in/mgo.v2"
)

const (
	UserCollection        = "users"
	TransactionCollection = "transactions"
	GambleCollection      = "gambles"
	RewardCollection      = "rewards"
	TicketCollection      = "tickets"

	DB = "mongobucks"

	StartingBalance = 100
)

var MongoUri = ""
var Session *mgo.Session

func init() {
	MongoUri = os.Getenv("MONGO_URI")
	Session = ConnectToMongoTLS(MongoUri)
	//SeedData()
}

func ConnectToMongoTLS(uri string) *mgo.Session {
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	dialInfo, err := mgo.ParseURL(uri)
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	fmt.Println("attemptung to connect", dialInfo)
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("COnnect to mongo")

	EnsureUserIndex(session)

	return session
}

func ConnectToMongo(uri string) *mgo.Session {
	session, err := mgo.Dial(uri)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connect to mongo")
	EnsureUserIndex(session)

	return session
}

func EnsureUserIndex(session *mgo.Session) {
	index := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := session.DB(DB).C(UserCollection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
