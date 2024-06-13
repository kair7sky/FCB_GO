package handlers

import (
    "encoding/json"
    "net/http"
    "xml-service/db"
    "xml-service/utils"
)

// AddServiceRequest represents the JSON request body for adding a service
type AddServiceRequest struct {
    ServiceId   string `json:"serviceId"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

// AddServiceHandler handles adding a new service
func AddServiceHandler(w http.ResponseWriter, r *http.Request) {
    var req AddServiceRequest

    // Decode the JSON request body
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    // Validate the request payload
    if req.ServiceId == "" || req.Name == "" {
        utils.RespondWithError(w, http.StatusBadRequest, "Missing required fields")
        return
    }

    // Add the service to the database
    err = db.AddService(req.ServiceId, req.Name, req.Description)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Error adding service")
        return
    }

    // Respond with success
    utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Service added successfully"})
}
