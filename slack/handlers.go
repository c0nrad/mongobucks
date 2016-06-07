package slack

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

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

	handler = append(handler, Handler{regexp.MustCompile("^(help|h)$"), HelpHandler})

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
	out := `>>> Usage: 

@mongobucks: balance
@mongobucks: give @stuart 10 for being a rockstar
@mongobucks: help
@mongobucks: gamble 5 heads
http://mongobucks.mongodb.cc
`

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
