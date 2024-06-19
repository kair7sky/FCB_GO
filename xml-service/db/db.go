package db

import (
	"database/sql"
	"fmt"
	"log"

	"xml-service/config"
	"xml-service/models"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	SQL *sql.DB
)

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg *config.Config) {
	var err error

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	SQL, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = SQL.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: SQL,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to initialize Gorm: %v", err)
	}

	log.Println("Connected to the database")

	// Автоматически мигрировать схему базы данных
	err = DB.AutoMigrate(&models.AutoCheck{})
	if err != nil {
		log.Fatalf("Ошибка миграции схемы базы данных: %v", err)
	}
}
