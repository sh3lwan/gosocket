package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
)

var DBConnection *sql.DB = nil

func Init() {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	fmt.Printf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)

	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database))

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
