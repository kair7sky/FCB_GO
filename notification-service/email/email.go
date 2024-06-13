package email

import (
    "net/smtp"
    "os"

    "github.com/jordan-wright/email"
)

func SendEmail(to, content string) error {
    from := os.Getenv("SMTP_MAIL")
    pass := os.Getenv("SMTP_KEY")

    e := email.NewEmail()
    e.From = from
    e.To = []string{to}
    e.Subject = "Notification"
    e.Text = []byte(content)

    return e.Send("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"))
}
