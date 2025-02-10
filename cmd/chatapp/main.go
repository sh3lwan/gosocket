package main

import (
	"log"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/sh3lwan/gosocket/internal"
	. "github.com/sh3lwan/gosocket/internal/server"
)

const (
    defaultPort = "8080"
)
func main() {
	srv := NewServer(defaultPort)

	HandleRoutes()

    err := srv.Start()

	log.Fatal(err.Error())
}
