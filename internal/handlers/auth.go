package handlers

import (
	. "github.com/sh3lwan/gosocket/internal/repositories"
	. "github.com/sh3lwan/gosocket/pkg/utils"
	"net/http"
)

func HandleConnect(w http.ResponseWriter, r *http.Request) {
	receiver := r.URL.Query().Get("receiver")

	messages, err := GetMessages(receiver)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
