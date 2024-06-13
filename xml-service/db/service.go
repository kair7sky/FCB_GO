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

    _, err := DB.Exec(query, serviceId, name, description)
    if err != nil {
        log.Printf("Error adding service: %v", err)
        return err
    }

    return nil
}
