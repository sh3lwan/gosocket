package internal

import (
	. "github.com/sh3lwan/gosocket/internal/handlers"
    "github.com/sh3lwan/gosocket/pkg/utils"
	"net/http"
)

func HandleRoutes() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        config.WriteJSON(w, http.StatusOK, map[string]any{
            "message": "Welcome to GoSocket",
        })
    })

    http.HandleFunc("/ws", SocketHandler)

	http.HandleFunc("GET /api/connect", HandleConnect)
}
