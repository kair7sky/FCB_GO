package db

import (
	"log"
)

// AddService adds a new service to the database
func AddService(serviceId, name, description string) error {
	query := `
        INSERT INTO services (serviceId, name, description)
        VALUES ($1, $2, $3)
    `

	result := DB.Exec(query, serviceId, name, description)
	if result.Error != nil {
		log.Printf("Error adding service: %v", result.Error)
		return result.Error
	}

	return nil
}
