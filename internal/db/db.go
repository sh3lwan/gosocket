package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"github.com/joho/godotenv"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/mattn/go-sqlite3"
)

var DBConnection *sql.DB = nil

type DBConfig struct {
    Driver string
    DSN string
    TestMode bool
}
func init() {
    // Load SQLite driver if in test mode
    _ = func() error {
        _ = "github.com/mattn/go-sqlite3" // Force import in test mode
        return nil
    }()
}

func LoadConfig(testMode bool) (*DBConfig, error) {

    if testMode {
        err := godotenv.Load( "../../.env.test")
        if err != nil {
            fmt.Printf("Error loading .env file")
        }
        return nil, err
    }

    err := godotenv.Load()

    if err != nil {
        fmt.Printf("Error loading .env file")
        return nil, err
    }

    dsn := os.Getenv("DB_DSN")

    if dsn == "" {
        dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
    }

    return &DBConfig{
        Driver: os.Getenv("DB_DRIVER"),
        DSN: dsn,
        TestMode: testMode,
    }, nil
}

func Init(config *DBConfig) {
	conn, err := sql.Open(config.Driver, config.DSN)

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
    // initalize once singilton
	if DBConnection == nil {
        config, err := LoadConfig(false)
        if err != nil {
            fmt.Printf("Error loading config: %v\n", err)
            return nil
        }

        Init(config)
	}

	return DBConnection
}
