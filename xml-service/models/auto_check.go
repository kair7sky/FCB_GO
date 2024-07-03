package models

import "time"

// AutoCheck represents the structure for auto checks
type AutoCheck struct {
	ID        uint   `gorm:"primaryKey"`
	URL       string `gorm:"not null"`
	Status    string `gorm:"not null"`
	Result    string
	Request   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
}

// Change represents the structure for changes in auto checks
type Change struct {
	URL       string
	Request   string
	OldResult string
	NewResult string
}
