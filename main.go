package main

import "github.com/c0nrad/mongobucks/slack"

func main() {
	//touch
	go slack.StartSlackListener()
	StartServer()
}
