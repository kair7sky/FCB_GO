package soap

import (
    "bytes"
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "net/http"
)

// SOAPEnvelope and SOAPBody are structs for parsing SOAP responses
type SOAPEnvelope struct {
    XMLName xml.Name `xml:"Envelope"`
    Body    SOAPBody `xml:"Body"`
}

type SOAPBody struct {
    XMLName xml.Name `xml:"Body"`
    Content string   `xml:",innerxml"`
}

// CheckResult represents the result of a SOAP request check
type CheckResult struct {
    URL    string `json:"url"`
    Status string `json:"status"`
    Result string `json:"result"`
}

// SendSOAPRequest sends a SOAP request and returns the response
func SendSOAPRequest(url string, xmlData string) (string, error) {
    req, err := http.NewRequest("POST", url, bytes.NewBufferString(xmlData))
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
func CheckSOAPRequests(url string, requests []string) ([]CheckResult, error) {
    var results []CheckResult

    for _, request := range requests {
        response, err := SendSOAPRequest(url, request)
        status := "success"
        result := response
        if err != nil {
            status = "failure"
            result = err.Error()
        }

        checkResult := CheckResult{
            URL:    url,
            Status: status,
            Result: result,
        }

        results = append(results, checkResult)
    }

    return results, nil
}
