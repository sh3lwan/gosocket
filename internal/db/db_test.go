package db

import (
    "testing"
)

func TestLoadConfig(t *testing.T) {
    config := LoadConfig(true)

    if config.Driver != "sqlite3" {
        t.Fatalf("Driver should be sqlite3, got %s", config.Driver)
    }

    if config.DSN != ":memory:" {
        t.Fatalf("DSN should be ::memory:, got %s", config.DSN)
    }

    if config.TestMode != true {
        t.Fatalf("TestMode should be true, got %v", config.TestMode)
    }
}

func TestInit(t *testing.T) {
    config := LoadConfig(true)

	Init(config)

	if DBConnection == nil {
		t.Fatal("DB connection should not be nil after initialization")
	}

	err := DBConnection.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestDBConnection(t *testing.T) {
    config := LoadConfig(true)

    Init(config)

    db := DB()
    if db == nil {
        t.Fatalf("Failed to fetch database: %v", db)
    }

    err := db.Ping()
    if err != nil {
        t.Fatalf("Failed to ping database: %v", err)
    }
}
