package slack

import (
	"fmt"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	Text    string `json:"text"`
}

var Counter uint64

func WriteMessage(ws *websocket.Conn, channel, text string) error {
	m := Message{}
	m.Id = atomic.AddUint64(&Counter, 1)
	m.Channel = channel
	m.Text = text
	m.Type = "message"

	fmt.Printf("Out: %+v\n", m)

	return websocket.JSON.Send(ws, m)
}

func ReadMessage(ws *websocket.Conn) (m Message, err error) {
	err = websocket.JSON.Receive(ws, &m)
	fmt.Printf("In:  %+v\n", m)
	return
}
