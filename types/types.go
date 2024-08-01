package types

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Response struct {
	StatusCode int16 `json:"status_code"`
	Data       map[string]any
}

type AuthBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ReceivedMessage struct {
	Id       string `json:"id"`
	Body     string `json:"body"`
	Username string `json:"username"`
	Receiver string `json:"receiver"`
	IsNew    bool   `json:"is_new"`
}

func (received *ReceivedMessage) BroadCast(chat *Chat) {

	for _, client := range chat.List() {

		if received.Receiver != "" {
			if received.Username != client.Username && received.Receiver != client.Username {
				continue
			}
		}

		conn := client.Conn

		message, err := json.Marshal(received)

		if err != nil {
			log.Printf("Error writing message: %v\n", err)
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, message)

		if err != nil {
			log.Printf("Error writing message: %v\n", err)
			return
		}
	}
}
