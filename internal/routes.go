package internal

import (
	. "github.com/sh3lwan/gosocket/internal/handlers"
	"net/http"
)

func HandleRoutes() {

	http.HandleFunc("/ws", SocketHandler)

	http.HandleFunc("GET /api/connect", HandleConnect)
}
