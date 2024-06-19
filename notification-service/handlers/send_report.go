package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/jordan-wright/email"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MessageTo string    `json:"messageTo"`
	Content   string    `json:"content"`
	SentAt    time.Time `json:"sentAt"`
}

var db *gorm.DB

func InitHandler(database *gorm.DB) {
	db = database
}

func SendReportHandler(w http.ResponseWriter, r *http.Request) {
	var notification Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		log.Printf("Error decoding JSON request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Save notification to database
	notification.SentAt = time.Now()
	db.Create(&notification)

	// Send email notification
	err = sendEmail(notification.MessageTo, notification.Content)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Report sent successfully"})
}

func sendEmail(to, content string) error {
	from := os.Getenv("SMTP_MAIL")
	if from == "" {
		return fmt.Errorf("SMTP_MAIL is not set")
	}
	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = "Auto Check Report"
	e.HTML = []byte(content)

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_MAIL")
	smtpPass := os.Getenv("SMTP_KEY")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP configuration is not set properly")
	}

	return e.Send(smtpHost+":"+smtpPort, smtp.PlainAuth("", smtpUser, smtpPass, smtpHost))
}
