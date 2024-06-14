package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	// "xml-service/db"
	"xml-service/email"
	"xml-service/soap"
	"xml-service/utils"
)

// RequestPayload represents the payload for the auto check request
type RequestPayload struct {
	URL      string   `json:"url"`
	Requests []string `json:"requests"`
	Email    string   `json:"email"`
}

// AutoCheckHandler handles automatic checking of SOAP requests
func AutoCheckHandler(w http.ResponseWriter, r *http.Request) {
	var request RequestPayload

	// Decode the JSON request payload
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Error decoding JSON request payload: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check the SOAP requests
	results, err := soap.CheckSOAPRequests(request.URL, request.Requests)
	if err != nil {
		log.Printf("Error checking SOAP requests: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check SOAP requests")
		return
	}

	// // Store the auto check results in the database
	// for _, result := range results {
	// 	db.DB.Create(&result)
	// }

	// Format the results for email
	var emailBody strings.Builder
	emailBody.WriteString("Auto Check Results:\n\n")
	for _, result := range results {
		emailBody.WriteString(fmt.Sprintf("URL: %s\nStatus: %s\nResult: %s\nCreatedAt: %s\n\n",
			result.URL, result.Status, result.Result, result.CreatedAt))
	}

	// Send the results via email
	err = email.SendEmail(request.Email, "Auto Check Results", emailBody.String())
	if err != nil {
		log.Printf("Error sending email: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to send email")
		return
	}

	// Respond with the results
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"results": results})
}
