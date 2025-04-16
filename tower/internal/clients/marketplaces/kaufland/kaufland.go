package kaufland

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/zepzeper/tower/internal/clients"
	"github.com/zepzeper/tower/internal/clients/marketplaces/mirakl/data"
	"github.com/zepzeper/tower/internal/helpers"
)

type ClientConfig struct {
	PublicKey     string
  PrivateKey    string
	UseBasicAuth  bool
	Debug         bool
}

type Client struct {
	config ClientConfig
}

// NewClient constructs a new woocommerce.Client object and sets the given configuration.
func NewClient(cc ClientConfig) Client {
	return Client{
		config: cc,
	}
}

// GetAPIEndpoint returns the API endpoint URL for the configured marketplace
func (cc ClientConfig) GetAPIEndpoint() string {
  return "https://sellerapi.kaufland.com/v2/"
}

func (c Client) Request(method, endpoint string, params url.Values, response interface{}) error {
  uri := c.config.GetAPIEndpoint()
	if uri == "" {
		return errors.New("invalid marketplace or missing API host")
	}

	// Create request
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return err
	}

  time := time.Now();

	// Set API key in the header if provided
	if c.config.PublicKey != "" {
		req.Header.Set("Shop-Client-Key", c.config.PublicKey)
	}
	
	// Add params to URL
	req.URL.RawQuery = params.Encode()

	// Debug
	if c.config.Debug {
		log.Printf("NEW REQUEST TO %s", uri)
	}

	// Perform request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	// Debug
	if c.config.Debug {
		helpers.Entry{
			URI:      req.URL.String(),
			Response: body,
		}.PrintPretty()
	}
	
	// Handle API errors
	if resp.StatusCode >= 400 {
		errResp := clients.ErrorResponse{}
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return err
		}
		err = errors.New(errResp.Code + ": " + errResp.Message)
	} else {
		err = json.Unmarshal(body, response)
	}
	
	return err
}

// Customers returns a Customers API client
func (c Client) Orders() data.Orders {
	return data.Orders{Client: c}
}

// Orders returns a Orders API client
func (c Client) Products() data.Products {
	return data.Products{Client: c}
}
