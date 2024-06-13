package db

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
    "xml-service/config"
)

// DB is a global variable to hold the database connection
var DB *sql.DB

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg *config.Config) {
    var err error

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatalf("Error pinging database: %v", err)
    }

    log.Println("Connected to the database")
}
