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
	go receive()

	http.HandleFunc("/", wsHandler)

	http.HandleFunc("/api/connect", handleConnect)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func receive() {
	for {
		message := <-broadcast

		recevied := ReceivedMessage{
			Id:       message.Id,
			Message:  message.Message,
			IsNew:    message.IsNew,
			Username: clients[message.Id].Username,
		}

        fmt.Printf("received message %v\n", recevied)
		if first := checkfirst(recevied); first != nil {
			recevied = *first


		}

		sendToClients(recevied)
	}

}

func checkfirst(message ReceivedMessage) *ReceivedMessage {

	var firstMessage ReceivedMessage

	err := json.Unmarshal([]byte(message.Message), &firstMessage)

	if err != nil || !firstMessage.IsNew {
		return nil
	}

	id := message.Id
	client := clients[id]
	client.Username = firstMessage.Username
	clients[id] = client

	return &ReceivedMessage{
		Id:       client.Id,
		Message:  "Connection established!",
		IsNew:    true,
		Username: client.Username,
	}

}

func sendToClients(received ReceivedMessage) {
	for _, client := range clients {
		conn := client.Conn

		id, err := insertMessage(received)

		if err != nil {
			log.Printf("Error writing message: %v\n", err)
			return
		}

		received.Id = fmt.Sprint(id)

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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]any{
		"status":   http.StatusOK,
		"messages": getMessages(),
	})

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

		fmt.Printf("Message Received: %v", message)

		broadcast <- &ReceivedMessage{
			Id:      client.Id,
			Message: string(message),
		}
	}

	client.Conn.Close()

	delete(clients, client.Id)
}
