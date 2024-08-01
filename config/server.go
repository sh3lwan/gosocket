package config

import (
	"fmt"
	"net/http"
)

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{
        port: port,
    }
}

func (s *Server) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), nil)
}
