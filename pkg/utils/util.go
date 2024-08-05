package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	. "github.com/sh3lwan/gosocket/internal/models"
)

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func getEnv(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func WriteJSON(w http.ResponseWriter, statusCode int, data map[string]any) {

	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(
		Response{
			StatusCode: int16(statusCode),
			Data:       data,
		},
	)

	if err != nil {
		fmt.Printf("Error marsheling: %v", err)
		return
	}

	w.Write(response)
}
