package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

var DBConnection *sql.DB = nil

func getEnv(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func Init() {
	conn, err := sql.Open("mysql", "root:password@tcp(mysql_db:3306)/chat")

	if err != nil {
        fmt.Printf("Error creating connection: %v\n", err)
		return
	}

	// See "Important settings" section.
	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)
	DBConnection = conn
}

func DB() *sql.DB {
	if DBConnection == nil {
		Init()
	}

	return DBConnection
}

func Migrate() error {
	driver, err := mysql.WithInstance(DB(), &mysql.Config{})

	if err != nil {
		fmt.Printf("Error creating migrations: %v\n", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations/",
		"mysql", driver)

	if err != nil {
		fmt.Printf("Error : %v\n", err)
		return err
	}

	return m.Steps(1)
}

func WriteJSON(w http.ResponseWriter, statusCode int, data map[string]any) {

	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(Response{
		StatusCode: int16(statusCode),
		Data:       data,
	})

	if err != nil {
		fmt.Printf("Error marsheling: %v", err)
		return
	}

	w.Write(response)
}
