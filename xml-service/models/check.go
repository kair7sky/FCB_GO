package models

type Check struct {
	ID                  uint   `gorm:"primaryKey" json:"id"`
	ServiceID           string `json:"service_id"`
	Request             string `json:"request"`
	Code                int    `json:"code"`
	ResponseExpectation string `json:"response_expectation"`
}
