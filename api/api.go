package api

import (
	"encoding/json"
	"net/http"

	"github.com/c0nrad/mongobucks/models"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var CookieStore = sessions.NewCookieStore([]byte("i4masup3rs3cret!!!13373!??swagswagswag"))
var CookieName = "session"

func BuildRouter() *mux.Router {
	r := mux.NewRouter()

	// User
	r.HandleFunc("/api/users", GetUsersHandler).Methods("GET")
	r.HandleFunc("/api/users/me", GetMeHandler).Methods("GET")
	r.HandleFunc("/api/users/{username}", GetUserHandler).Methods("GET")
	r.HandleFunc("/api/users/{username}/transactions", GetUserTransactionsHandler).Methods("GET")
	r.HandleFunc("/api/users/{username}/gambles", GetUserGamblesHandler).Methods("GET")

	// r.HandleFunc("/api/transactions", NewTransactionHandler).Methods("POST")
	r.HandleFunc("/api/transactions/recent", GetRecentTransactionsHandler).Methods("GET")
	r.HandleFunc("/api/transactions/{id}", GetTransactionHandler).Methods("GET")

	r.HandleFunc("/api/gambles/recent", GetRecentGamblesHandler).Methods("GET")
	r.HandleFunc("/api/gambles/{id}", GetGambleHandler).Methods("GET")

	// // Login
	r.HandleFunc("/login/google", LoginGoogleHandler)
	r.HandleFunc("/oauth/google", GoogleOAuthCallbackHandler)

	return r
}

func GetMeHandler(w http.ResponseWriter, r *http.Request) {
	me := context.Get(r, "username").(string)

	user, err := models.FindUser(me)
	if err != nil {
		http.Error(w, "unable to find user", 400)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {

	users, err := models.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	user, err := models.FindUser(username)
	if err != nil {
		http.Error(w, "unable to find user", 400)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func GetUserTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	transactions, err := models.GetTransactionsForUser(username)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}

func GetRecentTransactionsHandler(w http.ResponseWriter, r *http.Request) {

	transactions, err := models.GetRecentTransactions()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}

func GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	transaction, err := models.FindTransaction(id)
	if err != nil {
		http.Error(w, "unable to find transaction", 400)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func GetUserGamblesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	gambles, err := models.GetGamblesForUser(username)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(gambles)
}

func GetRecentGamblesHandler(w http.ResponseWriter, r *http.Request) {

	gambles, err := models.GetRecentGambles()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(gambles)
}

func GetGambleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	gamble, err := models.FindGamble(id)
	if err != nil {
		http.Error(w, "unable to find gamble", 400)
		return
	}

	json.NewEncoder(w).Encode(gamble)
}
