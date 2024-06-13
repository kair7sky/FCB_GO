package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "xml-service/config"
    "xml-service/db"
    "xml-service/handlers"
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
    defer db.DB.Close()

    // Set up router
    router := mux.NewRouter()
    router.HandleFunc("/auto-check", handlers.AutoCheckHandler).Methods("POST")
    router.HandleFunc("/manual-check", handlers.ManualCheckHandler).Methods("POST")
    router.HandleFunc("/add-service", handlers.AddServiceHandler).Methods("POST")

    // Start the server
    log.Println("Server starting at :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
