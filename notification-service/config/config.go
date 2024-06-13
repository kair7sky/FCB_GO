package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    SMTPMail     string
    SMTPKey      string
    APIURL       string
    RabbitMQURL  string
    DSN          string
}

func LoadConfig() *Config {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    return &Config{
        SMTPMail:    os.Getenv("SMTP_MAIL"),
        SMTPKey:     os.Getenv("SMTP_KEY"),
        APIURL:      os.Getenv("API_URL"),
        RabbitMQURL: os.Getenv("RABBITMQ_URL"),
        DSN:         os.Getenv("DSN"),
    }
}
