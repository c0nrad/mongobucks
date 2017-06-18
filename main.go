package main

import "github.com/c0nrad/mongobucks/slack"

func main() {

	// if !strings.Contains(api.GetReturnHost(), "localhost") {
	go slack.StartSlackListener()
	// }
	StartServer()
}
