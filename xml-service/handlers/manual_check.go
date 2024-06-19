package handlers

import (
	"encoding/json"
	"net/http"
	"xml-service/db"
	"xml-service/models"
	"xml-service/utils"
)

// ManualCheckHandler handles manual checks and stores them in the database
func ManualCheckHandler(w http.ResponseWriter, r *http.Request) {
	var check models.Check

	// Decode the JSON request payload
	err := json.NewDecoder(r.Body).Decode(&check)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the incoming check data
	if check.ServiceID == "" || check.Request == "" || check.Code == 0 || check.ResponseExpectation == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Insert the check into the database
	result := db.DB.Exec("INSERT INTO checks (service_id, request, code, response_expectation) VALUES ($1, $2, $3, $4)",
		check.ServiceID, check.Request, check.Code, check.ResponseExpectation)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error inserting manual check")
		return
	}

	// Respond with a success message
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"message": "Manual check added successfully"})
}
