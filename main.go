package main

import "github.com/c0nrad/mongobucks/slack"

func main() {
	go slack.StartSlackListener()
	StartServer()
}
