package xml

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

// ServiceReturn represents the structure of the parsed XML data within CDATA
type ServiceReturn struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		ActionResponse struct {
			Return struct {
				CDATAContent string `xml:",innerxml"`
			} `xml:",innerxml"`
		} `xml:"Body>actionResponse"`
	} `xml:"Body"`
}

// ParsedServiceReturn represents the structure after parsing CDATA content
type ParsedServiceReturn struct {
	XMLName       xml.Name `xml:"com.fcb.vshep.services.main.ServiceReturn"`
	ReturnMessage struct {
		ResponseInfo struct {
			MessageID    string `xml:"messageId"`
			ResponseDate struct {
				Year     int `xml:"year"`
				Month    int `xml:"month"`
				Day      int `xml:"day"`
				Hour     int `xml:"hour"`
				Minute   int `xml:"minute"`
				Second   int `xml:"second"`
				Timezone int `xml:"timezone"`
			} `xml:"responseDate"`
			Status struct {
				Code    string `xml:"code"`
				Message string `xml:"message"`
			} `xml:"status"`
			SessionID string `xml:"sessionId"`
		} `xml:"responseInfo"`
		ResponseData struct {
			Data string `xml:"data"`
		} `xml:"responseData"`
	} `xml:"returnMessage"`
	ErrorCode int `xml:"errorCode"`
}

// ParseXMLFile parses an XML file from the given filepath
func ParseXMLFile(filepath string) (*ParsedServiceReturn, error) {
	xmlFile, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open XML file: %w", err)
	}
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file: %w", err)
	}

	var result ServiceReturn
	err = xml.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML content: %w", err)
	}

	// Extract CDATA content
	cdataContent := result.Body.ActionResponse.Return.CDATAContent
	if cdataContent == "" {
		return nil, fmt.Errorf("no CDATA content found")
	}

	// Parse the CDATA content
	var parsedResult ParsedServiceReturn
	err = xml.Unmarshal([]byte(cdataContent), &parsedResult)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal CDATA content: %w", err)
	}

	return &parsedResult, nil
}
