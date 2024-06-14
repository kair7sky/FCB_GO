package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"notification-service/config"
	emailPkg "notification-service/email" // Renaming the import to avoid conflict
)

var db *gorm.DB

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MessageTo string    `json:"messageTo"`
	Content   string    `json:"content"`
	SentAt    time.Time `json:"sentAt"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Set up logging to a file
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("Application started")

	cfg := config.LoadConfig()

	db = initDB(cfg.DSN)

	// Loading .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("Loaded .env file successfully")

	// Applying migrations
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Ошибка при получении объекта базы данных: %v", err)
	}
	err = goose.Up(sqlDB, "./migrations")
	if err != nil {
		log.Fatalf("Ошибка при применении миграций: %v", err)
	}
	log.Println("Database migrations applied successfully")

	// Auto migrate Notification model
	db.AutoMigrate(&Notification{})
	log.Println("Database migrated")

	// RabbitMQ
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	log.Println("Connected to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	log.Println("RabbitMQ channel opened")

	q, err := ch.QueueDeclare(
		"notification_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a notification: %s", d.Body)
			var notification Notification
			err := json.Unmarshal(d.Body, &notification)
			if err != nil {
				log.Printf("Error unmarshaling JSON: %v", err)
				continue
			}

			// Save notification to database
			notification.SentAt = time.Now()
			db.Create(&notification)

			// Send email notification
			err = emailPkg.SendEmail(notification.MessageTo, notification.Content)
			if err != nil {
				log.Printf("Error sending email: %v", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for notifications. To exit press CTRL+C")
	<-forever

	// HTTP
	log.Println("Сервер запущен на :8081")
	http.ListenAndServe(":8081", nil)
}

func initDB(dsn string) *gorm.DB {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	return db
}
