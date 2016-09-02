package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/c0nrad/mongobucks/models"
	"github.com/gorilla/context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConf = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_SECRET"),
	RedirectURL:  "http://" + GetReturnHost() + "/oauth/google",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	},
	Endpoint: google.Endpoint,
}

func GetReturnHost() string {
	h := os.Getenv("OAUTH_RETURN")
	if h == "" {
		return "localhost:8080"
	}
	return h
}

func LoginGoogleHandler(w http.ResponseWriter, r *http.Request) {
	tmpConf := GoogleOAuthConf

	redirect := "http://" + r.Host + "/oauth/google"

	tmpConf.RedirectURL = redirect
	url := tmpConf.AuthCodeURL("state")
	http.Redirect(w, r, url, 307)
}

type GoogleProfileStruct struct {
	Id      string
	Email   string
	Name    string
	Picture string
	Hd      string
	// given_name, locale, hd,
}

func GoogleOAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	authcode := r.FormValue("code")

	tok, err := GoogleOAuthConf.Exchange(oauth2.NoContext, authcode)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var profile GoogleProfileStruct
	err = json.Unmarshal(contents, &profile)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if profile.Hd != "10gen.com" {
		http.Error(w, "must use a 10gen.com email account", 500)
		return
	}

	username := strings.Split(profile.Email, "@")[0]

	user, err := models.FindOrCreateUser(username)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = user.UpdateProfile(profile.Name, profile.Picture)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	session, err := CookieStore.Get(r, CookieName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	session.Values["username"] = username
	session.Save(r, w)

	fmt.Println("[+] ", username, "just logged in")

	http.Redirect(w, r, "/", 301)
}

func CookieAuthentication(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, err := CookieStore.Get(r, CookieName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	username := session.Values["username"]
	if username == nil {
		fmt.Println("[-] Missing Cookie, using anonymous")
		context.Set(r, "username", "anonymous")
	} else {
		context.Set(r, "username", username)
	}

	next(w, r)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := CookieStore.Get(r, CookieName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	session.Values["username"] = ""
	session.Save(r, w)

	http.Redirect(w, r, "/", 301)
}
