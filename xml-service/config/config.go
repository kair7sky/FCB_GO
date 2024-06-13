package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

// Config holds the database configuration values
type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
}

// LoadConfig loads the configuration from the .env file and environment variables
func LoadConfig() (*Config, error) {
    err := godotenv.Load()
    if err != nil {
        log.Printf("Error loading .env file: %v", err)
    }

    config := &Config{
        DBHost:     os.Getenv("DB_HOST"),
        DBPort:     os.Getenv("DB_PORT"),
        DBUser:     os.Getenv("DB_USER"),
        DBPassword: os.Getenv("DB_PASSWORD"),
        DBName:     os.Getenv("DB_NAME"),
    }

    if config.DBHost == "" || config.DBPort == "" || config.DBUser == "" || config.DBPassword == "" || config.DBName == "" {
        log.Fatalf("Missing required environment variables")
        return nil, err
    }

    return config, nil
}
