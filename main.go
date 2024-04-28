package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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
    fmt.Println("Started server...")
	// Migrate databse

	// Handle client messages
	go receive()

	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("GET /api/connect", handleConnect)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func receive() {
	for {
		message := <-broadcast

		recevied := ReceivedMessage{
			Id:       message.Id,
			Body:     message.Body,
			IsNew:    message.IsNew,
			Username: clients[message.Id].Username,
		}

		if first := checkfirst(recevied); first != nil {
			recevied = *first
		}

		id, err := insertMessage(recevied)

		if err != nil {
			fmt.Printf("Error inserting message: %v\n", err)
			return
		}

		recevied.Id = fmt.Sprint(id)

		sendToClients(recevied)
	}

}

func checkfirst(message ReceivedMessage) *ReceivedMessage {

	var firstMessage ReceivedMessage

	err := json.Unmarshal([]byte(message.Body), &firstMessage)

	if err != nil || !firstMessage.IsNew {
		return nil
	}

	id := message.Id
	client := clients[id]
	client.Username = firstMessage.Username
	clients[id] = client

	return &ReceivedMessage{
		Id:       client.Id,
		Body:     "Connection established!",
		IsNew:    true,
		Username: client.Username,
	}

}

func sendToClients(received ReceivedMessage) {
	for _, client := range clients {
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
	enableCors(&w)

	messages, err := getMessages()

    fmt.Printf("Messages: %v", messages)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	WriteJSON(
		w,
		http.StatusOK,
		map[string]any{
			"messages": messages,
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

	go initClient(client)
}

func initClient(client *Client) {

	for {
		_, message, err := client.Conn.ReadMessage()

		if err != nil {
			log.Println("Error reading message:\n", err)
			break
		}

		broadcast <- &ReceivedMessage{
			Id:   client.Id,
			Body: string(message),
		}
	}

	client.Conn.Close()

	delete(clients, client.Id)
}
