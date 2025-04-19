package auth

import (
	"net/http"
)

// APIKeyConfig holds configuration for API key authentication
type APIKeyConfig struct {
	PublicKey    string // Consumer Key for WooCommerce
	PrivateKey   string // Consumer Secret for WooCommerce
	KeyName      string // Header or query parameter name for the public key
	SecretName   string // Header or query parameter name for the private key
	InHeader     bool   // Whether the keys go in header or query parameters
	ApiSignature string // Signature format for some APIs (default is empty)
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
		// Use header-based authentication
		if a.config.KeyName != "" && a.config.PublicKey != "" {
			req.Header.Set(a.config.KeyName, a.config.PublicKey)
		}
		
		if a.config.SecretName != "" && a.config.PrivateKey != "" {
			req.Header.Set(a.config.SecretName, a.config.PrivateKey)
		}
	} else {
		// Use query parameter-based authentication
		q := req.URL.Query()
		
		if a.config.KeyName != "" && a.config.PublicKey != "" {
			q.Add(a.config.KeyName, a.config.PublicKey)
		}
		
		if a.config.SecretName != "" && a.config.PrivateKey != "" {
			q.Add(a.config.SecretName, a.config.PrivateKey)
		}
		
		req.URL.RawQuery = q.Encode()
	}
	
	return nil
}

// IsValid checks if the API key is valid
func (a *APIKeyAuth) IsValid() bool {
	// API keys typically don't expire
	return a.config.PublicKey != "" && a.config.PrivateKey != ""
}

// Refresh refreshes the API key (typically not needed)
func (a *APIKeyAuth) Refresh() error {
	// API keys typically don't need refreshing
	return nil
}
