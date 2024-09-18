package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/sh3lwan/gosocket/internal"
	. "github.com/sh3lwan/gosocket/internal/models"
	. "github.com/sh3lwan/gosocket/internal/server"
)

const (
	PORT = "8080"
)

var chat *Chat

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	srv := NewServer(PORT)

	HandleRoutes()

	err := srv.Start()

	log.Fatal(err.Error())

}
