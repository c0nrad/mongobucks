package main

import (
	"net/http"

	"github.com/c0nrad/mongobucks/api"
	"github.com/codegangsta/negroni"
)

var (
	Port = "8081"
)

func StartServer() {
	r := api.BuildRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	n := negroni.Classic()

	n.UseHandler(r)
	n.Run(":" + Port)
}
