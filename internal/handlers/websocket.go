package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	. "github.com/sh3lwan/gosocket/internal/models"
	"github.com/sh3lwan/gosocket/internal/server"
)

var s *server.Server = server.Get()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
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

	s.Chat.Join(client)

	go client.Receive(s.Chat)

	w.Write([]byte("OK"))
}
