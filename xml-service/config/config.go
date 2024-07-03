package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    DBHost           string
    DBPort           string
    DBUser           string
    DBPassword       string
    DBName           string
    RabbitMQURL      string
    SMTPHost         string
    SMTPPort         string
    SMTPUser         string
    SMTPPassword     string
    NotificationEmail string // добавляем это поле
}

func LoadConfig() (*Config, error) {
    viper.SetConfigFile(".env")
    err := viper.ReadInConfig()
    if err != nil {
        return nil, err
    }

    config := &Config{
        DBHost:           viper.GetString("DB_HOST"),
        DBPort:           viper.GetString("DB_PORT"),
        DBUser:           viper.GetString("DB_USER"),
        DBPassword:       viper.GetString("DB_PASSWORD"),
        DBName:           viper.GetString("DB_NAME"),
        RabbitMQURL:      viper.GetString("RABBITMQ_URL"),
        SMTPHost:         viper.GetString("SMTP_HOST"),
        SMTPPort:         viper.GetString("SMTP_PORT"),
        SMTPUser:         viper.GetString("SMTP_USER"),
        SMTPPassword:     viper.GetString("SMTP_PASSWORD"),
        NotificationEmail: viper.GetString("NOTIFICATION_EMAIL"), // инициализируем это поле
    }

    return config, nil
}
