package woocommerce

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	
	"github.com/zepzeper/tower/internal/core/connectors"
)

// Connector implements the Connector interface for WooCommerce API
type Connector struct {
	connectors.BaseConnector
	ConsumerKey    string
	ConsumerSecret string
}

// NewConnector creates a new WooCommerce connector
func NewConnector(config map[string]interface{}) (*Connector, error) {
	// Extract configuration
	apiURL, ok := config["api_url"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid api_url in configuration")
	}
	
	consumerKey, ok := config["consumer_key"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid consumer_key in configuration")
	}
	
	consumerSecret, ok := config["consumer_secret"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid consumer_secret in configuration")
	}
	
	// Set up headers
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	
	// Create the connector
	connector := &Connector{
		BaseConnector: connectors.NewBaseConnector(apiURL, headers, GetSchema()),
		ConsumerKey:   consumerKey,
		ConsumerSecret: consumerSecret,
	}
	
	return connector, nil
}

// Connect implements the Connector interface
func (c *Connector) Connect(ctx context.Context) error {
	// Just test with a simple request
	_, err := c.makeRequest(ctx, "GET", "products", nil, nil)
	return err
}

// Fetch implements the Connector interface
func (c *Connector) Fetch(ctx context.Context, query map[string]interface{}) ([]connectors.DataPayload, error) {
	// Determine endpoint based on entity type
	endpoint := "products"
	if entityType, ok := query["entity_type"].(string); ok {
		if entityType == "order" {
			endpoint = "orders"
		} else if entityType == "customer" {
			endpoint = "customers"
		}
	}
	
	// Convert query to URL parameters
	params := url.Values{}
	for key, value := range query {
		// Skip entity_type as it's not a WooCommerce parameter
		if key == "entity_type" {
			continue
		}
		
		// Handle specific parameter types
		switch v := value.(type) {
		case map[string]interface{}:
			// Skip complex objects
			continue
		default:
			params.Set(key, fmt.Sprintf("%v", v))
		}
	}
	
	// Make the request
	responseData, err := c.makeRequest(ctx, "GET", endpoint, params, nil)
	if err != nil {
		return nil, err
	}
	
	// Parse the response
	var items []connectors.DataPayload
	err = json.Unmarshal(responseData, &items)
	if err != nil {
		// Try parsing as a single object
		var item connectors.DataPayload
		err = json.Unmarshal(responseData, &item)
		if err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		items = append(items, item)
	}
	
	return items, nil
}

// Push implements the Connector interface
func (c *Connector) Push(ctx context.Context, data []connectors.DataPayload) error {
	for _, item := range data {
		// Determine endpoint and ID for updates
		endpoint := "products"
		idField := "id"
		
		if entityType, ok := item["entity_type"].(string); ok {
			if entityType == "order" {
				endpoint = "orders"
			} else if entityType == "customer" {
				endpoint = "customers"
			}
			delete(item, "entity_type") // Remove helper field
		}
		
		// Check if this is an update (has ID) or create
		var (
			method string
			path   string
		)
		
		if id, ok := item[idField]; ok {
			method = "PUT"
			path = fmt.Sprintf("%s/%v", endpoint, id)
		} else {
			method = "POST"
			path = endpoint
		}
		
		// Convert data to JSON
		bodyJSON, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("error serializing data: %w", err)
		}
		
		// Make the request
		_, err = c.makeRequest(ctx, method, path, nil, strings.NewReader(string(bodyJSON)))
		if err != nil {
			return err
		}
	}
	
	return nil
}

// makeRequest makes a request to the WooCommerce API with authentication
func (c *Connector) makeRequest(ctx context.Context, method, endpoint string, params url.Values, body io.Reader) ([]byte, error) {
	// Create URL
	apiURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid API URL: %w", err)
	}
	
	// Add path
	apiURL.Path = strings.TrimSuffix(apiURL.Path, "/") + "/" + strings.TrimPrefix(endpoint, "/")
	
	// Create new params if nil
	if params == nil {
		params = url.Values{}
	}
	
	// Add authentication
	params.Set("consumer_key", c.ConsumerKey)
	params.Set("consumer_secret", c.ConsumerSecret)
	
	// Add params to URL
	apiURL.RawQuery = params.Encode()
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, method, apiURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	// Set headers
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
	
	// Perform request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check for error response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned non-success status: %d", resp.StatusCode)
	}
	
	// Read response body
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	
	return responseData, nil
}
