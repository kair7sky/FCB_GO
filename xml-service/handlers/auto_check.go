package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
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

	// Store the auto check results in the database and collect failures
	var failures []soap.CheckResult
	for _, result := range results {
		autoCheck := models.AutoCheck{
			URL:     result.URL,
			Status:  result.Status,
			Result:  result.Result,
			Request: result.Request,
		}
		db.DB.Create(&autoCheck)
		if result.Status == "failure" {
			failures = append(failures, result)
		}
	}

	// If there are failures, send an email notification
	if len(failures) > 0 {
		sendFailureNotifications(request.Email, failures)
	}

	// Check database for changes and send email notifications if any
	CheckDatabaseForChanges(request.Email)

	// Respond with the results
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"results": results})
}

func sendFailureNotifications(email string, failures []soap.CheckResult) {
	notificationPayload := map[string]interface{}{
		"messageTo": email,
		"content":   generateEmailContent(failures),
	}

	payloadBytes, err := json.Marshal(notificationPayload)
	if err != nil {
		log.Printf("Error marshaling JSON payload: %v", err)
		return
	}

	// Send the results to the notification service
	resp, err := http.Post("http://localhost:8081/send-report", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error sending report to notification service: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Notification service responded with status: %v", resp.Status)
	}
}

func generateEmailContent(results []soap.CheckResult) string {
	const tmpl = `
    <html>
    <body>
        <h2>Auto Check Results</h2>
        <table border="1">
            <tr>
                <th>URL</th>
                <th>Status</th>
                <th>Result</th>
            </tr>
            {{ range . }}
            <tr>
                <td>{{ .URL }}</td>
                <td>{{ .Status }}</td>
                <td><pre>{{ .Result }}</pre></td>
            </tr>
            {{ end }}
        </table>
    </body>
    </html>
    `
	t := template.Must(template.New("emailTemplate").Parse(tmpl))

	var content bytes.Buffer
	err := t.Execute(&content, results)
	if err != nil {
		log.Printf("Error generating email content: %v", err)
		return "Failed to generate email content"
	}

	return content.String()
}
