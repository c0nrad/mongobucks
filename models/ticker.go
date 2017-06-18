package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/TrevorDev/go-finance"

	"gopkg.in/mgo.v2/bson"
)

type Ticker struct {
	Name string
	Last int // in pennies
	TS   time.Time
}

func init() {
	go TickerUpdateJob()
}

var KnownTickers = []string{"AAPL", "GOOG", "YHOO", "MSFT"}

func IsKnownTicker(ticker string) bool {
	for _, t := range KnownTickers {
		if t == ticker {
			return true
		}
	}
	return false
}

func TickerUpdateJob() {

	time.Sleep(time.Second * 3)

	for {
		tickers, err := UpdateTickers(KnownTickers)
		if err != nil {
			fmt.Println("[-]", err.Error())
		}

		errors := PruneInvestments(tickers)
		if len(errors) != 0 {
			for _, err := range errors {
				fmt.Println("[-] ", err.Error())
			}
		}

		time.Sleep(5 * time.Minute)
	}
}

func NewTicker(name string, price int) (*Ticker, error) {
	session := Session.Copy()
	defer session.Close()

	t := &Ticker{name, price, time.Now()}
	err := session.DB(DB).C(TickerCollection).Insert(t)

	return t, err
}

func UpdateTickers(names []string) ([]*Ticker, error) {
	var tickers []*Ticker
	out, err := finance.GetStockInfo(names, []string{finance.Last_Trade_Price_Only})
	if err != nil {
		return nil, err
	}

	fmt.Println(out)
	for _, name := range names {
		priceStr := out[name][finance.Last_Trade_Price_Only]

		price, _ := strconv.ParseFloat(priceStr, 64)

		t, err := NewTicker(name, int(price*100))
		if err != nil {
			return nil, err
		}
		tickers = append(tickers, t)
	}

	return tickers, nil

}

func GetLastTickerByName(name string) (*Ticker, error) {
	session := Session.Copy()
	defer session.Close()

	name = strings.ToUpper(name)
	if !IsKnownTicker(name) {
		_, err := UpdateTickers([]string{name})
		if err != nil {
			return nil, err
		}

		KnownTickers = append(KnownTickers, name)
	}

	var t Ticker
	err := session.DB(DB).C(TickerCollection).Find(bson.M{"name": name}).Sort("-ts").One(&t)

	return &t, err
}

func GetTickerHistoryByName(name string, count int) ([]*Ticker, error) {
	session := Session.Copy()
	defer session.Close()

	var t []*Ticker
	err := session.DB(DB).C(TickerCollection).Find(bson.M{"name": name}).Sort("-ts").Limit(count).All(&t)

	return t, err
}
