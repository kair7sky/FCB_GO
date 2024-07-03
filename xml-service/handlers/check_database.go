package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"xml-service/db"
	"xml-service/models"
)

// CheckDatabaseForFailures checks the database for failures and sends email notifications
func CheckDatabaseForFailures(email string) {
	var failures []models.AutoCheck
	db.DB.Where("status = ?", "failure").Find(&failures)

	if len(failures) > 0 {
		notificationPayload := map[string]interface{}{
			"messageTo": email,
			"content":   generateEmailContentFailures(failures),
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
			return
		}
	}
}

func generateEmailContentFailures(failures []models.AutoCheck) string {
	const tmpl = `
    <html>
    <body>
        <h2>Auto Check Failure Results</h2>
        <table border="1">
            <tr>
                <th>URL</th>
                <th>Status</th>
                <th>Result</th>
                <th>Request</th>
            </tr>
            {{ range . }}
            <tr>
                <td>{{ .URL }}</td>
                <td>{{ .Status }}</td>
                <td><pre>{{ .Result }}</pre></td>
                <td><pre>{{ .Request }}</pre></td>
            </tr>
            {{ end }}
        </table>
    </body>
    </html>
    `
	t := template.Must(template.New("emailTemplate").Parse(tmpl))

	var content bytes.Buffer
	err := t.Execute(&content, failures)
	if err != nil {
		log.Printf("Error generating email content: %v", err)
		return "Failed to generate email content"
	}

	return content.String()
}

// CheckDatabaseForChanges checks the database for changes in successful requests and sends email notifications
func CheckDatabaseForChanges(email string) {
	var autoChecks []models.AutoCheck
	db.DB.Where("status = ?", "success").Find(&autoChecks)

	// Group the autoChecks by URL and check for changes in requests
	changes := make(map[string][]models.AutoCheck)
	for _, autoCheck := range autoChecks {
		changes[autoCheck.URL] = append(changes[autoCheck.URL], autoCheck)
	}

	var differences []models.Change
	for _, checks := range changes {
		requests := make(map[string]string)
		for _, check := range checks {
			if oldResult, exists := requests[check.Request]; exists {
				if oldResult != check.Result {
					differences = append(differences, models.Change{
						URL:       check.URL,
						Request:   check.Request,
						OldResult: oldResult,
						NewResult: check.Result,
					})
				}
			} else {
				requests[check.Request] = check.Result
			}
		}
	}

	if len(differences) > 0 {
		sendChangeNotifications(email, differences)
	}
}

func sendChangeNotifications(email string, differences []models.Change) {
	notificationPayload := map[string]interface{}{
		"messageTo": email,
		"content":   generateChangeEmailContent(differences),
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

func generateChangeEmailContent(changes []models.Change) string {
	const tmpl = `
    <html>
    <body>
        <h2>Auto Check Changes</h2>
        <table border="1">
            <tr>
                <th>URL</th>
                <th>Request</th>
                <th>Old Result</th>
                <th>New Result</th>
            </tr>
            {{ range . }}
            <tr>
                <td>{{ .URL }}</td>
                <td><pre>{{ .Request }}</pre></td>
                <td><pre>{{ .OldResult }}</pre></td>
                <td><pre>{{ .NewResult }}</pre></td>
            </tr>
            {{ end }}
        </table>
    </body>
    </html>
    `
	t := template.Must(template.New("emailTemplate").Parse(tmpl))

	var content bytes.Buffer
	err := t.Execute(&content, changes)
	if err != nil {
		log.Printf("Error generating email content: %v", err)
		return "Failed to generate email content"
	}

	return content.String()
}
