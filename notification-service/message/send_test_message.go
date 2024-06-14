package main

import (
    "encoding/json"
    "log"

    "github.com/streadway/amqp"
)

type Notification struct {
    MessageTo string `json:"messageTo"`
    Content   string `json:"content"`
}

func main() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "notification_queue",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to declare a queue: %v", err)
    }

    notification := Notification{
        MessageTo: "kairkhanovabylai@gmail.com",
        Content:   "Отчет по проверке xml",
    }
    body, err := json.Marshal(notification)
    if err != nil {
        log.Fatalf("Failed to marshal JSON: %v", err)
    }

    err = ch.Publish(
        "",
        q.Name,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
    if err != nil {
        log.Fatalf("Failed to publish a message: %v", err)
    }

    log.Println("Test message sent successfully")
}
