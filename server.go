package main

import (
	"net/http"

	"github.com/c0nrad/mongobucks/api"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
)

var (
	Port = "8080"
)

func StartServer() {
	r := api.BuildRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(api.CookieAuthentication))
	//n.Use(negroni.HandlerFunc(LocalAuthentication))

	n.UseHandler(r)
	n.Run(":" + Port)
}

func LocalAuthentication(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "username", "stuart.larsen")
	next(w, r)
	return
}
