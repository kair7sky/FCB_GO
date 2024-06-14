package soap

import (
    "bytes"
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "net/http"
    "xml-service/models"
)

type SOAPEnvelope struct {
    XMLName xml.Name `xml:"soapenv:Envelope"`
    SOAPEnv string   `xml:"xmlns:soapenv,attr"`
    Body    SOAPBody
}

type SOAPBody struct {
    XMLName xml.Name `xml:"soapenv:Body"`
    Content string   `xml:",innerxml"`
}

// SendSOAPRequest sends a SOAP request and returns the response
func SendSOAPRequest(url string, xmlData string) (string, error) {
    soapEnvelope := SOAPEnvelope{
        SOAPEnv: "http://schemas.xmlsoap.org/soap/envelope/",
        Body: SOAPBody{
            Content: xmlData,
        },
    }

    soapBytes, err := xml.MarshalIndent(soapEnvelope, "", "  ")
    if err != nil {
        return "", fmt.Errorf("failed to marshal SOAP request: %w", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(soapBytes))
    if err != nil {
        return "", fmt.Errorf("failed to create HTTP request: %w", err)
    }
    req.Header.Set("Content-Type", "text/xml")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to send HTTP request: %w", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read HTTP response: %w", err)
    }

    return string(body), nil
}

// CheckSOAPRequests sends SOAP requests and checks the responses
func CheckSOAPRequests(url string, requests []string) ([]models.AutoCheck, error) {
    var results []models.AutoCheck

    for _, request := range requests {
        response, err := SendSOAPRequest(url, request)
        status := "success"
        result := response
        if err != nil {
            status = "failure"
            result = err.Error()
        }

        autoCheck := models.AutoCheck{
            URL:    url,
            Status: status,
            Result: result,
        }

        results = append(results, autoCheck)
    }

    return results, nil
}
