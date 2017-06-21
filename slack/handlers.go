package slack

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/c0nrad/mongobucks/models"
)

type Handler struct {
	Re          *regexp.Regexp
	HandlerFunc func(command string, vars map[string]string) string
}

var Handlers []Handler

func init() {
	Handlers = BuildHandlers()
}

// mongobucks: balance
//   # Respond to use with current balance
// mongobucks: give <username> 10
//   # Transfer 10 mongobucks from sender to <username>
// mongobucks: rain 10
//   # Evenly split 10 mongobucks amongst chat room
// mongobucks: help
//   # Display help room

func BuildHandlers() []Handler {
	handler := []Handler{}
	handler = append(handler, Handler{regexp.MustCompile("^(balance|b)$"), BalanceHandler})
	handler = append(handler, Handler{regexp.MustCompile("^(give|g) (?P<to>.*) (?P<amount>[0-9]*) (?P<memo>.*)$"), TransferHandler})
	handler = append(handler, Handler{regexp.MustCompile("^(give|g) (?P<to>.*) (?P<amount>[0-9]*)$"), TransferHandler})
	handler = append(handler, Handler{regexp.MustCompile("^(balance|b) all$"), AllBalanceHandler})
	handler = append(handler, Handler{regexp.MustCompile("^(gamble) (?P<amount>[0-9]*) (?P<orientation>.*)$"), GambleHandler})
	handler = append(handler, Handler{regexp.MustCompile("^(gamble) (?P<amount>[0-9]*)$"), GambleHandler})

	handler = append(handler, Handler{regexp.MustCompile("^(ticker|t) (?P<tickername>.*)$"), TickerHandler})

	handler = append(handler, Handler{regexp.MustCompile("^(sell) (?P<id>[0-9a-f]*)$"), SellInvestmentHandler})
	handler = append(handler, Handler{regexp.MustCompile("^(buy) (?P<amount>[0-9]*) (?P<tickername>.*) (?P<leverage>[0-9]*)x$"), BuyInvestmentHandler})
	handler = append(handler, Handler{regexp.MustCompile("^(buy) (?P<amount>[0-9]*) (?P<tickername>.*)$"), BuyInvestmentNoLeverageHandler})

	handler = append(handler, Handler{regexp.MustCompile("^(investments|i)$"), ShowInvestmentsHandler})

	handler = append(handler, Handler{regexp.MustCompile("^(redeem|r) (?P<token>.*)$"), RedeemHandler})

	handler = append(handler, Handler{regexp.MustCompile("^(help|h)$"), HelpHandler})
	handler = append(handler, Handler{regexp.MustCompile("^(help|h) (investments|i)$"), InvestmentHelpHandler})

	return handler
}

func HandleMessage(message Message) string {

	text := strings.Join(strings.Fields(message.Text)[1:], " ")

	for _, Handler := range Handlers {
		if Handler.Re.MatchString(text) {
			names := Handler.Re.SubexpNames()
			values := Handler.Re.FindAllStringSubmatch(text, -1)[0]
			values = values[1:]

			vars := map[string]string{}
			for i, value := range values {
				name := names[i+1]
				if name != "" {
					vars[names[i+1]] = value
				}
			}

			username, err := GetUsername(message.User)
			if err != nil {
				return err.Error()
			}
			vars["user"] = username
			vars["channel"] = message.Channel
			return Handler.HandlerFunc(text, vars)

		}
	}

	return "[-] Command not recognized. Use 'help' for available commands."
}

func TrimUsername(username string) string {
	username = strings.Replace(username, "<@", "", -1)
	username = strings.Replace(username, ">", "", -1)
	return username
}

func HelpHandler(command string, vars map[string]string) string {
	out := "```"
	out += `
Mongobucks is a slack based peer recognition tool. Everyone starts with 100mongobucks.  

Balance:
-------
@mongobucks: balance

Giving/Tipping:
--------------
Did a co-worker do something awesome? Give them some mongobucks to show your appreciation! They can use these mongobucks to buy rewards at http://mongobucks.mongobdb.cc/#/rewards.
@mongobucks: give @stuart 10 for being a rockstar

Gambling:
--------
You can gamble your hard earned mongobucks! The odds are exactly 50% for doubling your bet.
@mongobucks: gamble 5
@mongobucks: gamble 5 tails

Investments (beta):
-----------
Not the gambling type? Don't worry, you can invest your hard earned mongobucks into stocks. (It tracks stock prices at an increased rate, you don't actually buy them).

@mongobucks help investments
@mongobucks buy 10 AAPL 10x
@mongobucks investments 
@mongobucks sell <id>

http://mongobucks.mongodb.cc

Redeem (beta):
------
If you've been rewarded a Mongobucks voucher, you can redeem them using:

@mongobucks redeem <id>
`

	out += "```"

	return out
}

func InvestmentHelpHandler(command string, vars map[string]string) string {
	out := "```"
	out += `
Investments:
-----------
You can use your hard earned mongobucks to buy stocks. As the price of stocks goes up and down, so do your mongobucks, but at an increased rate (called leverage). Whether you prefer steady blue chip stocks or penny stocks, it's all up to you.

Each Mongobuck is worth $100. To see the current price of a stock type:

@mongobucks ticker TWLO
> TWLO trading at $26.87 (0.2687 mongobucks each)

If you decide to buy some Twilio stocks:
@mongobucks buy 20 TWLO 5x
> Buy order placed for 20xTWLO at 0.27mongobucks/per, total price of 5mongobucks.

This will place an order for 20 TWLO stocks with 5x leverage.

5x leverage means that if the stock goes up 1%, your mongobucks will instead go up 5%. But if the stock goes down 1%, your investment goes down 5%. If your investment ever hits a sell value of 0 with leverage it will be marked as "overextened", and disabled. You can't have a negative investment.

To check on your investments:
@mongobucks investments
> Investments: (1 mongobucks == $100)
> ID                       | Stock | Quantity | Buy Price | Total Buy | Leverage | Current Price | Sell Value /w Leverage | Current ROI
> 5945693fac2685e6dcf36496   AAPL    10         1.42        14          10x        1.42            14                       0 
> 59456b76ac2685e8890ca4f1   TIME    10         0.14        1           20x        0.14            1                        0 
> 59456baeac2685e8890ca4f2   TIME    1          0.14        0           10x        0.14            0                        0 

> Total Investments Value: 15
> Total Investments and Balance: 100

And then you can sell any investment based on ID:
@mongobucks sell 5945693fac2685e6dcf36496`

	out += "```"

	return out
}

func BalanceHandler(command string, vars map[string]string) string {
	fmt.Println("BalanceHandler", vars)

	balance, err := models.GetBalance(vars["user"])
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%d mongobucks. http://mongobucks.mongodb.cc/#/u/%s", balance, vars["user"])
}

func TickerHandler(command string, vars map[string]string) string {
	fmt.Println("TickerHandler", vars)

	vars["tickername"] = strings.ToUpper(vars["tickername"])

	ticker, err := models.GetLastTickerByName(vars["tickername"])
	if err != nil {
		return err.Error()
	}

	// username := vars["user"]

	// user, err := models.FindUser(username)
	// if err != nil {
	// 	return "invalid user: " + err.Error()
	// }

	fmt.Println(ticker.Last)
	fmt.Println(models.PenniesToMongobucks)
	return fmt.Sprintf("%s trading at $%.2f (%.4f mongobucks each)", ticker.Name, float64(ticker.Last)/100, float64(ticker.Last)*models.PenniesToMongobucks)
}

func BuyInvestmentNoLeverageHandler(command string, vars map[string]string) string {
	vars["leverage"] = "1"
	return BuyInvestmentHandler(command, vars)
}

func BuyInvestmentHandler(command string, vars map[string]string) string {
	fmt.Println("BuyInvestmentHandler", vars)
	vars["tickername"] = strings.ToUpper(vars["tickername"])

	// @mongobucks buy 10 AAPL 10x
	amount, err := strconv.Atoi(vars["amount"])
	if err != nil {
		return "invalid amount: " + err.Error()
	}

	ticker, err := models.GetLastTickerByName(vars["tickername"])
	if err != nil {
		return "invalid ticker:" + err.Error()
	}

	leverage, err := strconv.Atoi(vars["leverage"])
	if err != nil {
		return "invalid leverage: " + err.Error()
	}

	username := vars["user"]

	user, err := models.FindUser(username)
	if err != nil {
		return "invalid user: " + err.Error()
	}

	i, err := models.BuyInvestment(user, ticker.Name, amount, leverage)
	if err != nil {
		return "Could not buy investment: " + err.Error()
	}

	// balance, err := models.GetBalance(vars["user"])
	// if err != nil {
	// 	return "Could not load balance :" + err.Error()
	// }

	return fmt.Sprintf("Buy order placed for %dx%s at %.2fmongobucks/per, total price of %dmongobucks.", i.Amount, i.TickerName, float64(i.BuyPrice)*models.PenniesToMongobucks, i.TotalBuyValue())
}

func SellInvestmentHandler(command string, vars map[string]string) string {
	fmt.Println("SellInvestmentHandler", vars)

	id := vars["id"]
	fmt.Println()
	if !bson.IsObjectIdHex(id) {
		return "not a valid investment id, should be an objectid"
	}

	investment, err := models.GetInvestment(id)
	if err != nil {
		return "Could not find investment."
	}

	if investment.Username != vars["user"] {
		return "This investment does not belong to you"
	}

	if !investment.IsValid {
		return "This investment is no longer valid"
	}

	username := vars["user"]

	user, err := models.FindUser(username)
	if err != nil {
		return "invalid user: " + err.Error()
	}

	sellValue, err := models.SellInvestment(user, investment)

	balance, err := models.GetBalance(vars["user"])
	if err != nil {
		return "Could not load balance :" + err.Error()
	}

	return fmt.Sprintf("Sold investment for %d. New Balance: %d", sellValue, balance)
}

func AllBalanceHandler(command string, vars map[string]string) string {
	users, err := models.GetUsers()
	if err != nil {
		return err.Error()
	}

	out := "Balances: \n"
	for _, u := range users {
		out += "@" + u.Username + ": " + strconv.Itoa(u.Balance) + "\n"
	}

	return out
}

func ShowInvestmentsHandler(command string, vars map[string]string) string {
	fmt.Println("[+] ShowInvestmentsHandler", vars)

	balance, err := models.GetBalance(vars["user"])
	if err != nil {
		return "Could not load balance :" + err.Error()
	}

	investments, err := models.GetInvestmentsForUser(vars["user"])
	if err != nil {
		return "unable to get investments: " + err.Error()
	}

	out := "```Investments: (1 mongobucks == $100)\nID                       | Stock | Quantity | Buy Price | Total Buy | Leverage | Current Price | Sell Value /w Leverage | Current ROI\n"
	total := 0

	for _, i := range investments {

		ticker, err := models.GetLastTickerByName(i.TickerName)
		if err != nil {
			return err.Error()
		}

		totalBuyValue := i.TotalBuyValue()
		totalSellValue := i.TotalSellValue(*ticker)
		totalDelta := i.ReturnOnInvestment(*ticker)

		// Stock Quantity Buy Price Total Amount  Leverage Current Sell Value /w Leverage Current ROI
		// ID APPL  5        154.44    13mongobucks  10x      150.44  6mongobucks            (7 mongobucks)
		//s   q    b    t
		out += fmt.Sprintf("%-26s %-7s %-10d %-11.2f %-11d %-10s %-15.2f %-24d %d %s\n", i.ID.Hex(), i.TickerName, i.Amount, float64(i.BuyPrice)*models.PenniesToMongobucks, totalBuyValue, strconv.Itoa(i.Leverage)+"x", float64(ticker.Last)*models.PenniesToMongobucks, totalSellValue, totalDelta, i.Error)
		// out += fmt.Sprintf("%dx%s bought at %.2f with %dx leverage. Ticker Value: %.2f. Current Sell Value with leverage: %d\n", i.Amount, i.TickerName, float64(i.BuyPrice)/100, i.Leverage, float64(ticker.Last)/100, roi/10)

		total += totalSellValue
	}
	out += fmt.Sprintf("\nTotal Investments Value: %d\n", total)
	out += fmt.Sprintf("Total Investments and Balance: %d\n```", total+balance)
	return out

}

func TransferHandler(command string, vars map[string]string) string {
	fmt.Println("[+] TransferHandler", vars)

	if !strings.HasPrefix(vars["to"], "<@") {
		return "[-] Prefix the username with '@', for example '@stuart'"
	}

	from := vars["user"]
	to, err := GetUsername(TrimUsername(vars["to"]))
	if err != nil {
		return "invalid user: " + err.Error()
	}

	amount, err := strconv.Atoi(vars["amount"])
	if err != nil {
		return "invalid amount: " + err.Error()
	}

	t, err := models.ExecuteTransfer(from, to, amount, vars["memo"])
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Transfer complete! http://mongobucks.mongodb.cc/#/t/" + t.ID.Hex())
}

func RedeemHandler(command string, vars map[string]string) string {
	fmt.Println("[+] RedeemHandler", vars)

	username := vars["user"]

	user, err := models.FindUser(username)
	if err != nil {
		return "invalid user: " + err.Error()
	}

	token := vars["token"]
	err = models.Redeem(token, user)

	if err != nil {
		return err.Error()
	}

	return "redeemed"
}

func GambleHandler(command string, vars map[string]string) string {
	fmt.Println("[+] GambleHandler", vars)

	username := vars["user"]

	balance, err := models.GetBalance(vars["user"])
	if err != nil {
		return err.Error()
	}

	amount, err := strconv.Atoi(vars["amount"])
	if err != nil {
		return "invalid amount: " + err.Error()
	}

	if amount > balance {
		return "insufficent funds"
	}

	orientation := strings.ToLower(vars["orientation"])

	if orientation != "heads" && orientation != "tails" {
		orientation = "heads"
	}

	user, err := models.FindUser(username)
	if err != nil {
		return err.Error()
	}

	g, err := models.ExecuteGamble(username, amount, orientation == "heads")
	if g.IsWinner {
		user.Balance += amount
	} else {
		user.Balance -= amount
	}

	err = user.Update()
	if err != nil {
		return err.Error()
	}

	if g.IsWinner {
		return fmt.Sprintf("You win! New Balance: %d. http://mongobucks.mongodb.cc/#/g/%s", user.Balance, g.ID.Hex())
	} else {
		return fmt.Sprintf("Better luck next time. New Balance: %d. http://mongobucks.mongodb.cc/#/g/%s", user.Balance, g.ID.Hex())
	}
}
