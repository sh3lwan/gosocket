package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sh3lwan/gosocket/repositories"
)

type Client struct {
	Id       string
	Username string
	*websocket.Conn
}

type Chat struct {
	mu        sync.Mutex
	clients   map[string]*Client
	broadcast chan *ReceivedMessage
}

func NewChat() *Chat {
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

func (client *Client) Receive(chat *Chat) {

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

func (c *Chat) Receive() {
	for {
		received := <-c.broadcast

		if received.Username == "" {
			client, err := c.Get(received.Id)

			if err != nil {
				log.Printf(err.Error())
				return
			}

			received.Username = client.Username
		}

		id, err := repositories.InsertMessage(*received)

		if err != nil {
			fmt.Printf("Error inserting message: %v\n", err)
			return
		}

		received.Id = fmt.Sprint(id)

		received.BroadCast(c)
	}

}
