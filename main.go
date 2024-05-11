package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sh3lwan/gosocket/config"
	"github.com/sh3lwan/gosocket/repositories"
	. "github.com/sh3lwan/gosocket/types"
)

var clients = make(map[string]Client)

var broadcast = make(chan *ReceivedMessage)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var p *repositories.Producer

func main() {
	p = repositories.NewProducer()
	// Handle sent client messages
	go send()

	//handle routes
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("GET /api/connect", handleConnect)

	// start server
	fmt.Println("Started server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func receive(client *Client) {

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
			clients[client.Id] = *client
		}

		broadcast <- &receivedMessage
	}

	client.Conn.Close()

	delete(clients, client.Id)
}

func send() {
	for {
		received := <-broadcast

		if received.Username == "" {
			received.Username = clients[received.Id].Username
		}

		id, err := repositories.InsertMessage(*received)

		if err != nil {
			fmt.Printf("Error inserting message: %v\n", err)
			return
		}

		p.WriteNotification(*received)

		received.Id = fmt.Sprint(id)

		sendToClients(*received)
	}

}

func sendToClients(received ReceivedMessage) {

	for _, client := range clients {

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

func handleConnect(w http.ResponseWriter, r *http.Request) {
	config.EnableCors(&w)

	receiver := r.URL.Query().Get("receiver")

	messages, err := repositories.GetMessages(receiver)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	config.WriteJSON(
		w,
		http.StatusOK,
		map[string]any{
			"messages": messages,
			"receiver": receiver,
		},
	)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("something went wrong %v\n", err)
		return
	}

	uuid := uuid.New().String()

	client := &Client{
		Conn: conn,
		Id:   uuid,
	}

	clients[uuid] = *client

	go receive(client)
}
