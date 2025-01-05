package handlers

import (
	. "github.com/sh3lwan/gosocket/internal/repositories"
	. "github.com/sh3lwan/gosocket/pkg/utils"
	"net/http"
)

func HandleConnect(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	receiver := r.URL.Query().Get("receiver")

	messages, err := GetMessages(receiver)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]any{
            "error": err.Error(),
        })
		return
	}

	WriteJSON(
		w,
		http.StatusOK,
		map[string]any{
			"messages": messages,
			"receiver": receiver,
		},
	)
}
