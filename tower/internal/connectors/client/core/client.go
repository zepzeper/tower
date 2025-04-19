package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIClient defines the base API client interface
type APIClient interface {
	Execute(method, path string, body, result interface{}) error
	BaseURL() string
	SetBaseURL(url string)
	TestRequest() (interface{}, error)
	Request(method, path string, body, result interface{}) error
}

// BaseClient implements the APIClient interface
type BaseClient struct {
	baseURL    string
	httpClient *http.Client
	auth       AuthMethod
}

// AuthMethod defines how different authentication strategies work
type AuthMethod interface {
	// Authenticate adds authentication to an HTTP request
	Authenticate(req *http.Request) error
	
	// IsValid checks if authentication credentials are still valid
	IsValid() bool
	
	// Refresh attempts to refresh the authentication if needed
	Refresh() error
}

// NewBaseClient creates a new API client with the specified authentication method
func NewBaseClient(baseURL string, auth AuthMethod) *BaseClient {
	return &BaseClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		auth:       auth,
	}
}

// BaseURL returns the base URL of the API
func (c *BaseClient) BaseURL() string {
	return c.baseURL
}

// SetBaseURL sets the base URL of the API
func (c *BaseClient) SetBaseURL(url string) {
	c.baseURL = url
}

// Execute performs an HTTP request with authentication
func (c *BaseClient) Execute(method, path string, body, result interface{}) error {
	// Prepare URL
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	
	// Prepare request body
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Create request
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	// Check if auth is valid, refresh if needed
	if c.auth != nil {
		if !c.auth.IsValid() {
			if err := c.auth.Refresh(); err != nil {
				return fmt.Errorf("error refreshing authentication: %w", err)
			}
		}
		
		// Apply authentication to request
		if err := c.auth.Authenticate(req); err != nil {
			return fmt.Errorf("error applying authentication: %w", err)
		}
	}
	
	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read the full response body for better error messages
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	
	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}
	
	// Parse response if result is expected
	if result != nil {
		if err := json.Unmarshal(bodyBytes, result); err != nil {
			return fmt.Errorf("error parsing response: %w, body: %s", err, string(bodyBytes))
		}
	}
	
	return nil
}

// Request method that forwards to Execute
func (c *BaseClient) Request(method, path string, body, result interface{}) error {
	return c.Execute(method, path, body, result)
}

// TestRequest makes a test request to the API
func (c *BaseClient) TestRequest() (interface{}, error) {
	var result interface{}
	err := c.Execute("GET", "/", nil, &result)
	return result, err
}
