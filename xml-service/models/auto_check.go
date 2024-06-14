package models

import "time"

// AutoCheck represents the structure for auto checks
type AutoCheck struct {
    ID        uint      `gorm:"primaryKey"`
    URL       string    `gorm:"not null"`
    Status    string    `gorm:"not null"`
    Result    string
    CreatedAt time.Time `gorm:"default:current_timestamp"`
}
