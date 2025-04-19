package auth

import (
	"net/http"
)

// APIKeyConfig holds configuration for API key authentication
type APIKeyConfig struct {
	PublicKey  string
	PrivateKey string
	KeyName    string // Header or query parameter name
	InHeader   bool   // Whether the key goes in header or query
}

// APIKeyAuth implements API key authentication
type APIKeyAuth struct {
	config APIKeyConfig
}

// NewAPIKeyAuth creates a new API key auth method
func NewAPIKeyAuth(config APIKeyConfig) *APIKeyAuth {
	return &APIKeyAuth{
		config: config,
	}
}

// Authenticate adds the API key to the request
func (a *APIKeyAuth) Authenticate(req *http.Request) error {
	if a.config.InHeader {
		req.Header.Set(a.config.KeyName, a.config.PrivateKey)
		// Some APIs might need both keys
		if a.config.PublicKey != "" {
			req.Header.Set("X-Public-Key", a.config.PublicKey)
		}
	} else {
		q := req.URL.Query()
		q.Add(a.config.KeyName, a.config.PrivateKey)
		if a.config.PublicKey != "" {
			q.Add("public_key", a.config.PublicKey)
		}
		req.URL.RawQuery = q.Encode()
	}
	return nil
}

// IsValid checks if the API key is valid
func (a *APIKeyAuth) IsValid() bool {
	// API keys typically don't expire
	return a.config.PrivateKey != ""
}

// Refresh refreshes the API key (typically not needed)
func (a *APIKeyAuth) Refresh() error {
	// API keys typically don't need refreshing
	return nil
}
