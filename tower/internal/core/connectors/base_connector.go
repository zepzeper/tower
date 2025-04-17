package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// BaseConnector provides common functionality for API connectors
type BaseConnector struct {
	BaseURL     string
	Headers     map[string]string
	Schema      Schema
	HTTPClient  *http.Client
}

// NewBaseConnector creates a new base connector
func NewBaseConnector(baseURL string, headers map[string]string, schema Schema) BaseConnector {
	return BaseConnector{
		BaseURL:    baseURL,
		Headers:    headers,
		Schema:     schema,
		HTTPClient: &http.Client{},
	}
}

// MakeRequest performs an HTTP request and returns the response
func (bc *BaseConnector) MakeRequest(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
	url := bc.BaseURL
	if !strings.HasSuffix(url, "/") && !strings.HasPrefix(path, "/") {
		url += "/"
	}
	url += path
	
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	// Set headers
	for key, value := range bc.Headers {
		req.Header.Set(key, value)
	}
	
	resp, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned non-success status: %d", resp.StatusCode)
	}
	
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	
	return responseBody, nil
}

// GetSchema returns the connector's schema
func (bc *BaseConnector) GetSchema() Schema {
	return bc.Schema
}

// ParseJSON parses JSON data into a map
func (bc *BaseConnector) ParseJSON(data []byte) (DataPayload, error) {
	var result DataPayload
	err := json.Unmarshal(data, &result)
	return result, err
}
