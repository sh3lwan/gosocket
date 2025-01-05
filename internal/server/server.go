package server

import (
	"fmt"
	. "github.com/sh3lwan/gosocket/internal/models"
	. "github.com/sh3lwan/gosocket/pkg/utils"
	"net/http"
)

var server *Server

type Logger struct {
	handler http.Handler
}

func NewLogger(handler http.Handler) *Logger {
	return &Logger{handler}
}

func (l *Logger) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	l.handler.ServeHTTP(w, r)
}

type Server struct {
	port string
	Chat *Chat
}

func Get() *Server {
    if server == nil {
        panic("Server Not Found")
    }

	return server
}

func NewServer(port string) *Server {
	fmt.Println("Server Created..")

	server = &Server{
		port: port,
		Chat: NewChat(),
	}
    
    return server
}

func (s *Server) Start() error {
	fmt.Println("Server is starting...")

    fmt.Println(fmt.Sprintf("localhost:%s", s.port))

	return http.ListenAndServe(
		fmt.Sprintf(":%s", s.port),
        nil,
	)

}
