package main

import (
	"strings"

	"github.com/c0nrad/ctf/api"
	"github.com/c0nrad/mongobucks/slack"
)

func main() {

	if !strings.Contains(api.GetReturnHost(), "localhost") {
		go slack.StartSlackListener()
	}
	StartServer()
}
