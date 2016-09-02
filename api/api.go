package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"strconv"

	"github.com/c0nrad/mongobucks/models"
	"github.com/c0nrad/mongobucks/ticket"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var CookieStore = sessions.NewCookieStore([]byte("i4masup3rs3cret!!!13373!??swagswagswag"))
var CookieName = "msession"

func BuildRouter() *mux.Router {
	r := mux.NewRouter()

	// User
	r.HandleFunc("/api/users", GetUsersHandler).Methods("GET")
	r.HandleFunc("/api/users/me", GetMeHandler).Methods("GET")
	r.HandleFunc("/api/users/{username}", GetUserHandler).Methods("GET")
	r.HandleFunc("/api/users/{username}/transactions", GetUserTransactionsHandler).Methods("GET")
	r.HandleFunc("/api/users/{username}/gambles", GetUserGamblesHandler).Methods("GET")
	r.HandleFunc("/api/users/me/tickets", GetMyTicketsHandler).Methods("GET")

	// r.HandleFunc("/api/transactions", NewTransactionHandler).Methods("POST")
	r.HandleFunc("/api/transactions/recent", GetRecentTransactionsHandler).Methods("GET")
	r.HandleFunc("/api/transactions/{id}", GetTransactionHandler).Methods("GET")

	r.HandleFunc("/api/gambles/recent", GetRecentGamblesHandler).Methods("GET")
	r.HandleFunc("/api/gambles/{id}", GetGambleHandler).Methods("GET")

	r.HandleFunc("/api/rewards", GetAllRewardsHandler).Methods("GET")

	r.HandleFunc("/api/tickets", BuyTicketHandler).Methods("POST")
	r.HandleFunc("/api/tickets/{token}", GetTicketHandler).Methods("GET")
	r.HandleFunc("/api/tickets/{token}/render", GenerateTicketHandler).Methods("GET")
	r.HandleFunc("/api/tickets/{token}/redeem", RedeemTicketHandler).Methods("POST")

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

// r.HandleFunc("/api/rewards", GetAllRewardsHandler).Methods("GET")

// r.HandleFunc("/api/ticket", BuyTicketHandler).Methods("POST")
// r.HandleFunc("/api/users/me/tickets", GetMyTicketsHandler).Methods("GET")
// r.HandleFunc("/api/ticket/:token/redeem", RedeemTicketHandler).Methods("POST")

func GetAllRewardsHandler(w http.ResponseWriter, r *http.Request) {
	rewards, err := models.GetRewards()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(rewards)
}

func BuyTicketHandler(w http.ResponseWriter, r *http.Request) {
	username := context.Get(r, "username").(string)

	user, err := models.FindUser(username)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var t models.Ticket
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Println(t)

	reward, err := models.GetRewardById(t.Reward)
	if err != nil {
		http.Error(w, "Could not find reward", 400)
		return
	}

	ticket, err := models.PurchaseTicket(user, reward)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(ticket)
}

func GetMyTicketsHandler(w http.ResponseWriter, r *http.Request) {
	username := context.Get(r, "username").(string)

	tickets, err := models.GetTicketsByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(tickets)
}

func RedeemTicketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	err := models.Redeem(token)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write([]byte("redeemed"))
}

func GetTicketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	ticket, err := models.GetTicketByToken(token)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(ticket)
}

func GenerateTicketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	t, err := models.GetTicketByToken(token)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	ticketImg := ticket.GenerateTicketImage(t)

	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, ticketImg)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	_, err = w.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}
