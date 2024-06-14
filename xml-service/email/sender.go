package email

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendEmail sends an email with the specified subject and body to the specified recipient.
func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("SMTP_MAIL")
	pass := os.Getenv("SMTP_KEY")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, pass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
