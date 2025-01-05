package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	. "github.com/sh3lwan/gosocket/internal/models"
	"github.com/sh3lwan/gosocket/internal/server"
)


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

    chat := server.Get().Chat

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("something is wrong %v\n", err)
		return
	}

	client := &Client{
		Id:   uuid.New().String(),
		Conn: conn,
	}

	chat.Join(client)

    go chat.Receive()

	go client.Read(chat)

	w.Write([]byte("OK"))
}
