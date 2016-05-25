package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/c0nrad/mongobucks/models"

	"golang.org/x/net/websocket"
)

type RealTimeMessagingResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	Url   string `json:"url"`
	Self  struct {
		Id string `json:"id"`
	}
}

type SlackUser struct {
	Ok   bool
	User struct {
		Id      string
		Name    string
		Profile struct {
			Email string
		}
	}
}

var SlackUsernameMap map[string]string
var Token string
var Username string

func init() {
	SlackUsernameMap = make(map[string]string)
	Token = os.Getenv("SLACK_API_TOKEN")
	Username = os.Getenv("SLACK_USERNAME")
}

func GetUsername(id string) (string, error) {
	username, ok := SlackUsernameMap[id]
	if ok {
		return username, nil
	}

	fmt.Println("[-] Doing hard lookup on", id)

	base := "https://slack.com/api/users.info?token=" + Token + "&user=" + id
	u, _ := url.Parse(base)

	fmt.Println(u.String())

	res, err := http.Get(u.String())
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var user SlackUser
	err = decoder.Decode(&user)

	fmt.Printf("%+v\n", user)

	if err != nil {
		return "", err
	}

	username = strings.Split(user.User.Profile.Email, "@")[0]

	SlackUsernameMap[id] = username
	models.FindOrCreateUser(username)
	return username, nil

}

func BuildRealTImeMessageConnection() (conn *websocket.Conn, username string, err error) {
	rtm, err := GetRealTimeMessagingKey()
	if err != nil {
		return nil, "", err
	}
	username = rtm.Self.Id

	conn, err = websocket.Dial(rtm.Url, "", "https://api.slack.com/")
	return conn, username, err
}

func GetRealTimeMessagingKey() (*RealTimeMessagingResponse, error) {
	url := "https://slack.com/api/rtm.start?token=" + Token
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var out RealTimeMessagingResponse
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func StartSlackListener() {
	conn, id, err := BuildRealTImeMessageConnection()
	if err != nil {
		panic(err)
	}

	fmt.Println("[+] Starting with username:", Username, "and id:", id)

	for {
		message, err := ReadMessage(conn)
		if err != nil {
			fmt.Println("[-] Error", err)
			continue
		}

		if message.Type == "message" && IsMessageForBot(message.Text, Username, id) {
			out := HandleMessage(message)
			WriteMessage(conn, message.Channel, out)
		}
	}
}

func IsMessageForBot(message, Username, id string) bool {
	if strings.HasPrefix(message, "<@"+id) || strings.HasPrefix(message, Username+":") {
		return true
	}
	return false
}
