package main

import "github.com/gorilla/websocket"

type Response struct {
	StatusCode int16 `json:"status_code"`
	Data       map[string]any
}

type AuthBody struct {
	Username string `json:"username"`
}

type Client struct {
	Id       string
	Username string
	*websocket.Conn
}

type ReceivedMessage struct {
	Id       string `json:"id"`
	Message  string `json:"message"`
	Username string `json:"username"`
	IsNew    bool   `json:"is_new"`
}
