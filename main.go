package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var clients = make(map[string]Client)

var broadcast = make(chan *ReceivedMessage)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	// Handle client messages
	go func() {
		for {
			sender := <-broadcast

			fmt.Printf("client count: %d\n", len(clients))

			for _, client := range clients {
				conn := client.Conn

				message, err := json.Marshal(ReceivedMessage{
					Id:      sender.Id,
					Message: sender.Message,
                    IsNew: sender.IsNew,
				})

				if err != nil {
					log.Printf("Error writing message: %v", err)
				}

				err = conn.WriteMessage(websocket.TextMessage, message)

				fmt.Printf("message sent: %s %s\n", client.Id, message)

				if err != nil {
					log.Printf("Error writing message: %v", err)
				}

			}
		}

	}()
	http.HandleFunc("/", wsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Client struct {
	Id string
	*websocket.Conn
}

type ReceivedMessage struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	IsNew   bool   `json:"is_new"`
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("something went wrong %v\n", err)
		conn.Close()
		return
	}

	uuid := uuid.New().String()

	client := &Client{
		Conn: conn,
		Id:   uuid,
	}

	clients[uuid] = *client

	fmt.Printf("new connection established: %s\n", uuid)

	fmt.Printf("clients count: %d\n", len(clients))

	go handleMessages(client)
}

func handleMessages(client *Client) {

	broadcast <- &ReceivedMessage{
		Id:      client.Id,
		Message: "Connection established",
        IsNew: true,
	}

	for {
		_, message, err := client.Conn.ReadMessage()

		if err != nil {
			log.Println("Error reading message:\n", err)
			break
		}

		broadcast <- &ReceivedMessage{
			Id:      client.Id,
			Message: string(message),
		}
	}

	fmt.Printf("Connection Closed: %s\n", client.Id)

	client.Conn.Close()

	delete(clients, client.Id)
}
