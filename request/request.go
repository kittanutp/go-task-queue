package request

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type RequestSchema struct {
	Header      map[string]interface{} `json:"header"`
	RequestType string                 `json:"type"`
	Url         string                 `json:"url"`
	Payload     map[string]interface{} `json:"payload"`
	Query       map[string]interface{} `json:"query"`
}

func (r *RequestSchema) MakeRequest() error {
	body, err := handlePayload(r.Payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(r.RequestType, r.Url, body)
	if err != nil {
		return err
	}
	handleHeader(req, r.Header)

	err = handleQuery(req, r.Query)
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(req)
	log.Printf("sending request to %v", req.URL)
	return err
}

func handlePayload(payload map[string]interface{}) (io.Reader, error) {
	if payload == nil {
		return nil, nil
	}
	// Marshal the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonData), nil
}
func handleQuery(req *http.Request, query map[string]interface{}) error {
	if query == nil {
		return nil
	}

	// Get the URL's query values and append new parameters
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, toString(value))
	}

	// Set the encoded query string on the request URL
	req.URL.RawQuery = q.Encode()
	return nil
}

func handleHeader(req *http.Request, header map[string]interface{}) {
	for key, value := range header {
		req.Header.Add(key, toString(value))
	}
	req.Header.Add("Task-Queue-Request", time.Now().Format(time.RFC3339))
}

// Helper function to convert interface{} to string
func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}
