package main

import (
	_ "github.com/go-sql-driver/mysql"
	. "github.com/sh3lwan/gosocket/internal"
	. "github.com/sh3lwan/gosocket/internal/server"
	"log"
)

const (
	PORT = "8080"
)

func main() {
	srv := NewServer(PORT)

	HandleRoutes()

	err := srv.Start()

	log.Fatal(err.Error())
}
