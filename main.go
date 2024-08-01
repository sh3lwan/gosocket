package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sh3lwan/gosocket/config"
	. "github.com/sh3lwan/gosocket/repositories"
	. "github.com/sh3lwan/gosocket/services"
	. "github.com/sh3lwan/gosocket/types"
)

var chat *Chat

var broadcast = make(chan *ReceivedMessage)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	srv := config.NewServer("80")

	err := srv.Start()

	log.Fatal(err.Error())

	chat = NewChat()

	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("GET /api/connect", handleConnect)
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
	config.EnableCors(&w)

	receiver := r.URL.Query().Get("receiver")

	messages, err := GetMessages(receiver)

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

	client := &Client{
		Id:   uuid.New().String(),
		Conn: conn,
	}

	chat.Join(client)

	go client.Receive(chat)
}
