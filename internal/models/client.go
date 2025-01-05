package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	. "github.com/sh3lwan/gosocket/internal/models/messages"
)

type Client struct {
	Id       string
	Username string
	*websocket.Conn
}

func (client *Client) Read(chat *Chat) {

	defer client.Conn.Close()

	defer chat.Remove(client)

	for {
		_, message, err := client.Conn.ReadMessage()

		if err != nil {
			log.Println("Error reading message:\n", err)
			break
		}

		receivedMessage := ReceivedMessage{
			Id: client.Id,
		}

		err = json.Unmarshal(message, &receivedMessage)

		if err != nil {
			log.Println("Error unmarshing message:\n", err)
			continue
		}

		fmt.Printf("received message: %v\n", receivedMessage)

        // set client's username from first message
        if client.Username == "" && receivedMessage.Username != "" {
            receivedMessage.Body = " joined chat!"
            client.Username = receivedMessage.Username
            chat.Update(client)
        }

		chat.broadcast <- &receivedMessage
	}
}
