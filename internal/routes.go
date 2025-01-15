package internal

import (
	"net/http"

	"github.com/sh3lwan/gosocket/internal/db"
	. "github.com/sh3lwan/gosocket/internal/handlers"
	"github.com/sh3lwan/gosocket/pkg/utils"
)

func HandleRoutes() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        config.WriteJSON(w, http.StatusOK, map[string]any{
            "status": "success",
            "message": "Welcome to GoSocket",
        })
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        error := db.DB().Ping()

        if error != nil {
            config.WriteJSON(w, http.StatusOK, map[string]any{
                "status": "fail",
                "message": "DB isn't connected",
            })
        }

        config.WriteJSON(w, http.StatusOK, map[string]any{
            "status": "success",
            "message": "Welcome to GoSocket",
        })
    })
    http.HandleFunc("/ws", SocketHandler)

	http.HandleFunc("GET /api/connect", HandleConnect)
}
