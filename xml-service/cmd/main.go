package main

import (
    "log"
    "net/http"
    "time"

    "xml-service/config"
    "xml-service/db"
    "xml-service/handlers"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    // Connect to the database
    db.Connect(cfg)

    // Defer closing the database connection
    sqlDB, err := db.DB.DB()
    if err != nil {
        log.Fatalf("Error getting database instance: %v", err)
    }
    defer sqlDB.Close()

    // Set up router
    router := mux.NewRouter()
    router.HandleFunc("/auto-check", handlers.AutoCheckHandler).Methods("POST")
    router.HandleFunc("/manual-check", handlers.ManualCheckHandler).Methods("POST")
    router.HandleFunc("/add-service", handlers.AddServiceHandler).Methods("POST")

    // Periodic checks
    go func() {
        for {
            handlers.CheckDatabaseForFailures(cfg.NotificationEmail)
            handlers.CheckDatabaseForChanges(cfg.NotificationEmail)
            time.Sleep(24 * time.Hour) // Run every 24 hours
        }
    }()

    // Start the server
    log.Println("Server starting at :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
