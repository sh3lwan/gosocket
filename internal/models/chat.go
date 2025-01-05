package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	. "github.com/sh3lwan/gosocket/internal/models/messages"
	. "github.com/sh3lwan/gosocket/internal/repositories"
	"log"
	"sync"
)

type Chat struct {
	mu        sync.Mutex
	clients   map[string]*Client
	broadcast chan *ReceivedMessage
}

func NewChat() *Chat {
    fmt.Println("Chat Created..")

	return &Chat{
		clients:   make(map[string]*Client),
		broadcast: make(chan *ReceivedMessage),
	}
}

func (c *Chat) Get(key string) (*Client, error) {
	client, ok := c.clients[key]

	if ok {
		return client, nil
	}

	return nil, errors.New("Client Not Found")
}

func (c *Chat) List() map[string]*Client {
	return c.clients
}

func (c *Chat) Join(client *Client) {
	c.mu.Lock()
	c.clients[client.Id] = client
	c.mu.Unlock()
}

func (c *Chat) Remove(client *Client) {
	c.mu.Lock()
	delete(c.clients, client.Id)
	c.mu.Unlock()
}

func (c *Chat) Update(client *Client) error {
	client, err := c.Get(client.Id)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.clients[client.Id] = client
	c.mu.Unlock()

	return nil
}

func (c *Chat) Receive() {
	for {

        fmt.Println("waiting for message..")

		received := <-c.broadcast

        fmt.Printf("received on chat: %v\n", received)

		if received.Username == "" {
			client, err := c.Get(received.Id)

			if err != nil {
				log.Printf(err.Error())
				return
			}

			received.Username = client.Username
		}

		id, err := InsertMessage(*received)

		if err != nil {
			fmt.Printf("Error inserting message: %v\n", err)
			return
		}

		received.Id = fmt.Sprint(id)

		c.Broadcast(received)
	}

}

func (c *Chat) Broadcast(received *ReceivedMessage) {
	for _, client := range c.List() {

		if received.Receiver != "" {
			if received.Username != client.Username && received.Receiver != client.Username {
				continue
			}
		}

		conn := client.Conn

		message, err := json.Marshal(received)

		if err != nil {
			log.Printf("Error marshing message: %v\n", err)
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, message)

		if err != nil {
			log.Printf("Error writing message: %v\n", err)
			return
		}
	}
}
