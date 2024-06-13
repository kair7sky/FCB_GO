package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"xml-service/utils"
	"xml-service/xml"
)

// RequestPayload represents the payload for the auto check request
type RequestPayload struct {
	FilePath string `json:"file_path"`
}

// AutoCheckHandler handles automatic checking of XML files
func AutoCheckHandler(w http.ResponseWriter, r *http.Request) {
	var request RequestPayload

	// Decode the JSON request payload
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Error decoding JSON request payload: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Parse the XML file
	xmlData, err := xml.ParseXMLFile(request.FilePath)
	if err != nil {
		log.Printf("Error parsing XML file: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to parse XML file")
		return
	}

	// Check the XML queries
	errors, err := xml.CheckXMLQueries(xmlData)
	if err != nil {
		log.Printf("Error checking XML queries: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error checking XML queries")
		return
	}

	// Respond with the results
	if len(errors) > 0 {
		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"errors": errors})
	} else {
		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"message": "All queries are valid"})
	}
}
