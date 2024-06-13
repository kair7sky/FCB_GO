package xml

import (
    "fmt"
    
)



// CheckXMLQueries checks the XML queries and returns a list of errors if any
func CheckXMLQueries(xmlData *ParsedServiceReturn) ([]string, error) {
    var errors []string

    if xmlData.ReturnMessage.ResponseInfo.MessageID == "" {
        errors = append(errors, "MessageID is missing")
    }
    if xmlData.ReturnMessage.ResponseInfo.Status.Code == "" {
        errors = append(errors, "Status code is missing")
    }

    // Add more checks as needed

    if len(errors) > 0 {
        return errors, fmt.Errorf("validation failed with %d error(s)", len(errors))
    }

    return nil, nil
}
