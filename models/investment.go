package models

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const MongobucksToUSD = float64(100)
const MongobucksToPennies = MongobucksToUSD * 100
const USDToMongobucks = 1 / MongobucksToUSD
const PenniesToMongobucks = USDToMongobucks / 100

type Investment struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	TS time.Time

	Username   string
	TickerName string
	Amount     int
	BuyPrice   int // stored as USD pennies
	Leverage   int

	IsValid bool
	Error   string

	LastPrice int
}

func (i Investment) Pennies() int {
	return i.BuyPrice
}

// Returns in MongoBucks
func (i Investment) TotalBuyValue() int {
	return int(float64(i.Amount) * float64(i.BuyPrice) * PenniesToMongobucks)
}

func (i Investment) TotalSellValue(t Ticker) int {
	change := (t.Last - i.BuyPrice) * i.Leverage

	sqewedValue := t.Last + change

	return int(float64(sqewedValue*i.Amount) * PenniesToMongobucks)
}

func (i Investment) ReturnOnInvestment(t Ticker) int {
	return i.TotalSellValue(t) - i.TotalBuyValue()
}

func (i *Investment) MarkOverextended() error {
	session := Session.Copy()
	defer session.Close()

	err := session.DB(DB).C(InvestmentCollection).UpdateId(i.ID, bson.M{"error": "Overextened leverage.", "amount": 0})
	return err
}

func BuyInvestment(user *User, name string, amount, leverage int) (*Investment, error) {
	session := Session.Copy()
	defer session.Close()

	investments, err := GetInvestmentsForUser(user.Username)
	if err != nil {
		return nil, err
	}

	if len(investments) >= 5 {
		return nil, errors.New("you can only ahve 5 active investments")
	}

	lastTicker, err := GetLastTickerByName(name)
	if err != nil {
		return nil, err
	}

	if int(float64(lastTicker.Last*amount)*PenniesToMongobucks) < 1 {
		return nil, errors.New(fmt.Sprintf("The buy order must be at least 1mongobuck, currently: %.2f", float64(lastTicker.Last*amount)*PenniesToMongobucks))
	}

	if lastTicker.Last*amount > int(float64(user.Balance)*MongobucksToPennies) {
		return nil, errors.New("insufficent funds")
	}

	if leverage > 10 || leverage < 1 {
		return nil, errors.New("leverage must be between 1 and 50")
	}

	user.Balance -= int(float64(lastTicker.Last*amount) * PenniesToMongobucks)
	err = user.Update()
	if err != nil {
		return nil, err
	}

	i := Investment{ID: bson.NewObjectId(), Username: user.Username, TickerName: name, Amount: amount, BuyPrice: lastTicker.Last, Leverage: leverage, TS: time.Now(), IsValid: true, Error: ""}

	err = session.DB(DB).C(InvestmentCollection).Insert(i)
	return &i, err
}

func SellInvestment(user *User, investment *Investment) (int, error) {
	session := Session.Copy()
	defer session.Close()

	lastTicker, err := GetLastTickerByName(investment.TickerName)
	if err != nil {
		return 0, err
	}

	err = session.DB(DB).C(InvestmentCollection).UpdateId(investment.ID, bson.M{"$set": bson.M{"isvalid": false}})
	if err != nil {
		return 0, err
	}

	user.Balance += investment.TotalSellValue(*lastTicker)
	err = user.Update()
	if err != nil {
		return 0, err
	}
	sellValue := investment.TotalSellValue(*lastTicker)

	return sellValue, err
}

func GetInvestment(id string) (*Investment, error) {
	session := Session.Copy()
	defer session.Close()

	var investment *Investment
	err := session.DB(DB).C(InvestmentCollection).FindId(bson.ObjectIdHex(id)).One(&investment)
	return investment, err
}

func GetInvestmentsForUser(username string) ([]*Investment, error) {
	session := Session.Copy()
	defer session.Close()

	var investments []*Investment
	err := session.DB(DB).C(InvestmentCollection).Find(bson.M{"username": username, "isvalid": true}).All(&investments)
	return investments, err
}

func GetAllInvestmentsByTicker(ticker string) ([]*Investment, error) {
	session := Session.Copy()
	defer session.Close()

	var investments []*Investment
	err := session.DB(DB).C(InvestmentCollection).Find(bson.M{"tickername": ticker, "isvalid": true}).All(&investments)
	return investments, err
}

func PruneInvestments(tickers []*Ticker) []error {
	errors := []error{}

	for _, ticker := range tickers {
		investments, err := GetAllInvestmentsByTicker(ticker.Name)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		for _, investment := range investments {
			if investment.TotalSellValue(*ticker) <= 0 {
				err = investment.MarkOverextended()
				if err != nil {
					errors = append(errors, err)
				}
			}
		}
	}
	return errors
}
