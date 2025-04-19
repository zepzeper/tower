package woocommerce

import (
	"fmt"
	"os"

	"github.com/zepzeper/tower/internal/connectors/client/auth"
	"github.com/zepzeper/tower/internal/connectors/client/core"
)

// Client for interacting with the WooCommerce API
type Client struct {
	core *core.BaseClient
}

// NewClient creates a new WooCommerce API client
func NewClient(demo bool) (*Client, error) {
	apiURL := os.Getenv("WOOCOMMERCE_API_URL")
	if apiURL == "" {
		apiURL = "https://example.com/wp-json/wc/v3" // Default URL, should be overridden
	}
	
	consumerKey := os.Getenv("WOOCOMMERCE_CONSUMER_KEY")
	consumerSecret := os.Getenv("WOOCOMMERCE_CONSUMER_SECRET")
	
	if consumerKey == "" || consumerSecret == "" {
		return nil, fmt.Errorf("missing required WooCommerce credentials in environment variables")
	}
	
	// If demo mode is enabled, use a demo URL if specified
	if demo && os.Getenv("WOOCOMMERCE_DEMO_API_URL") != "" {
		apiURL = os.Getenv("WOOCOMMERCE_DEMO_API_URL")
	}
	
	// Create API key auth config
	// WooCommerce uses consumer_key and consumer_secret as query parameters
	apiKeyConfig := auth.APIKeyConfig{
		PublicKey:  consumerKey,
		PrivateKey: consumerSecret,
		KeyName:    "consumer_key",
		SecretName: "consumer_secret",
		InHeader:   false, // WooCommerce uses query parameters for authentication
	}
	
	// Create auth method
	apiKeyAuth := auth.NewAPIKeyAuth(apiKeyConfig)
	
	// Create base client
	baseClient := core.NewBaseClient(apiURL, apiKeyAuth)
	
	return &Client{core: baseClient}, nil
}

// BaseURL returns the base URL of the API
func (c *Client) BaseURL() string {
	return c.core.BaseURL()
}

// SetBaseURL sets the base URL of the API
func (c *Client) SetBaseURL(url string) {
	c.core.SetBaseURL(url)
}

// Request forwards the request to the core client
func (c *Client) Request(method, path string, body, result interface{}) error {
	return c.core.Execute(method, path, body, result)
}

// Execute forwards execution to the core client
func (c *Client) Execute(method, path string, body, result interface{}) error {
	return c.core.Execute(method, path, body, result)
}

// TestRequest makes a simple test request to the WooCommerce API
func (c *Client) TestRequest() (interface{}, error) {
	var result interface{}
	
	// Make a simple request to the products endpoint
	err := c.Execute("GET", "/", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("test request failed: %w", err)
	}
	
	return result, nil
}
