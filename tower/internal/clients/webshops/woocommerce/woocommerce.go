package woocommerce

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/zepzeper/tower/internal/clients"
	"github.com/zepzeper/tower/internal/helpers"
	"github.com/zepzeper/tower/internal/clients/webshops/woocommerce/data"
)

type ClientConfig struct {
	APIHost        string // Complete REST API base URL: "https://example.com/wp-json/wc/v3"
	ConsumerKey    string
	ConsumerSecret string
	UseBasicAuth   bool
	Debug          bool
}

func (cc ClientConfig) GetAPIEndpoint(method string) string {
	base, err := url.Parse(cc.APIHost)
	if err != nil {
		return ""
	}

  base.Path = path.Join(base.Path, method)
	return base.String()
}

// Client is the entry point for all methods.
type Client struct {
	config ClientConfig
}

// NewClient constructs a new woocommerce.Client object and sets the given configuration.
func NewClient(cc ClientConfig) Client {
	return Client{
		config: cc,
	}
}

// Request performs a request and unmarshals JSON response into given response object.
func (c Client) Request(method, endpoint string, params url.Values, response interface{}) error {
	uri := c.config.GetAPIEndpoint(endpoint)

	// Parse params
	if params == nil {
		params = url.Values{}
	}

	// Set Auth
	if !c.config.UseBasicAuth {
		params.Set("consumer_key", c.config.ConsumerKey)
		params.Set("consumer_secret", c.config.ConsumerSecret)
	}

	// Add params to URI
	encodedParams := params.Encode()
	if encodedParams != "" {
		uri += "?" + params.Encode()
	}

	// Create request
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return err
	}

	// Set Auth
	if c.config.UseBasicAuth {
		req.SetBasicAuth(c.config.ConsumerKey, c.config.ConsumerSecret)
	}

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
            URI:      uri,
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
func (c Client) Customers() data.Customers {
	return data.Customers{Client: c}
}
