package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	. "github.com/sh3lwan/gosocket/internal"
	. "github.com/sh3lwan/gosocket/internal/server"
)

func main() {
    err := godotenv.Load()

    if err != nil {
        log.Fatal("Error loading .env file")
    }

    port := os.Getenv("APP_PORT")

	srv := NewServer(port)

	HandleRoutes()

	err = srv.Start()

	log.Fatal(err.Error())
}
