package brincr

import (
	"fmt"
	"net/url"
	"os"

	"github.com/zepzeper/tower/internal/connectors/client/auth"
	"github.com/zepzeper/tower/internal/connectors/client/core"
)

// Client for interacting with the Brincr API
type Client struct {
	core *core.BaseClient
}

// NewClient creates a new Brincr API client
func NewClient(demo bool) (*Client, error) {
	tenantID := os.Getenv("BRINCR_TENANT_ID")
	clientID := os.Getenv("BRINCR_CLIENT_ID")
	clientSecret := os.Getenv("BRINCR_CLIENT_SECRET")
	
	// Create OAuth2 config
	oauth2Config := auth.OAuth2Config{
		TokenURL:     "https://connect.identity.stagaws.visma.com/connect/token",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TenantID:     tenantID,
		Scopes:       []string{
                            "brincr:customer",
                            "brincr:customerContact",
                            "brincr:customerShippingAddress",
                            "brincr:order",
                            "brincr:orderConnected",
                            "brincr:orderLine",
                            "brincr:product",
                            "brincr:productImage",
                            "brincr:productVariant",
                            "brincr:productWarehouseStock",
                            "brincr:salesOrder",
                            "brincr:status",
                    },
	}
	
	// Create auth method
	oauth2Auth := auth.NewOAuth2ClientCredentials(oauth2Config)
	
	// Initialize with a token
	if err := oauth2Auth.Refresh(); err != nil {
		return nil, fmt.Errorf("failed to get initial token: %w", err)
	}
    
var baseClient *core.BaseClient

  if demo {
        // Create base client
        baseClient = core.NewBaseClient("https://sandbox-api.brincr.com", oauth2Auth)
    } else {
        // Create base client
        baseClient = core.NewBaseClient("https://api.brincr.com", oauth2Auth)
    }

	
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

// TestRequest makes a simple test request to the Brincr API
func (c *Client) TestRequest() (interface{}, error) {
	var result interface{}

	q := url.Values{}
	q.Add("filter[product_type]", "default") // change "123" to whatever value you want

	url := "/api/v1/products?" + q.Encode()

	err := c.Execute("GET", url, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("test request failed: %w", err)
	}

	return result, nil
}
