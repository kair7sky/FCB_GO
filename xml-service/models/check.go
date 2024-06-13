package models

import "time"

// AutoCheck represents the structure for auto checks
type AutoCheck struct {
    ID        uint      `gorm:"primaryKey"`
    FilePath  string    `gorm:"not null"`
    Status    string    `gorm:"not null"`
    Result    string
    CreatedAt time.Time `gorm:"default:current_timestamp"`
}

type Check struct {
	ID                  uint   `gorm:"primaryKey" json:"id"`
	ServiceID           string `json:"service_id"`
	Request             string `json:"request"`
	Code                int    `json:"code"`
	ResponseExpectation string `json:"response_expectation"`
}
