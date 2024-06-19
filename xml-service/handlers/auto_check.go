package handlers

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "xml-service/db"
    "xml-service/models"
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

    // Store the auto check results in the database
    for _, result := range results {
        autoCheck := models.AutoCheck{
            URL:    result.URL,
            Status: result.Status,
            Result: result.Result,
        }
        db.DB.Create(&autoCheck)
    }

    // Prepare the payload for notification service
    notificationPayload := map[string]interface{}{
        "messageTo": request.Email,
        "content":   generateEmailContent(results),
    }

    payloadBytes, err := json.Marshal(notificationPayload)
    if err != nil {
        log.Printf("Error marshaling JSON payload: %v", err)
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to prepare notification payload")
        return
    }

    // Send the results to the notification service
    resp, err := http.Post("http://localhost:8081/send-report", "application/json", bytes.NewBuffer(payloadBytes))
    if err != nil {
        log.Printf("Error sending report to notification service: %v", err)
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to send report to notification service")
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Notification service responded with status: %v", resp.Status)
        utils.RespondWithError(w, http.StatusInternalServerError, "Notification service error")
        return
    }

    // Respond with the results
    utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"results": results})
}

func generateEmailContent(results []soap.CheckResult) string {
    content := "Here are the results of the auto checks:\n\n"
    for _, result := range results {
        content += "URL: " + result.URL + "\nStatus: " + result.Status + "\nResult: " + result.Result + "\n\n"
    }
    return content
}
