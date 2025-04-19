package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Token represents OAuth token data
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	ExpiresAt   time.Time
}

// OAuth2Config holds configuration for OAuth2 client credentials flow
type OAuth2Config struct {
	TokenURL     string
	ClientID     string
	ClientSecret string
	TenantID     string
	Scopes       []string
}

// OAuth2ClientCredentials implements OAuth2 client credentials flow
type OAuth2ClientCredentials struct {
	config     OAuth2Config
	Token      *Token
	HTTPClient *http.Client
}

// NewOAuth2ClientCredentials creates a new OAuth2 client credentials auth method
func NewOAuth2ClientCredentials(config OAuth2Config) *OAuth2ClientCredentials {
	return &OAuth2ClientCredentials{
		config:     config,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Authenticate adds the OAuth2 token to the request
func (o *OAuth2ClientCredentials) Authenticate(req *http.Request) error {
	if o.Token == nil {
		return fmt.Errorf("no token available")
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", o.Token.TokenType, o.Token.AccessToken))
	return nil
}

// IsValid checks if the token is still valid
func (o *OAuth2ClientCredentials) IsValid() bool {
	return o.Token != nil && time.Now().Before(o.Token.ExpiresAt)
}

// Refresh gets a new access token
func (o *OAuth2ClientCredentials) Refresh() error {
	// Prepare form data
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("scope", strings.Join(o.config.Scopes, " "))
	
	if o.config.TenantID != "" {
		formData.Set("tenant_id", o.config.TenantID)
	}
	
	formData.Set("client_id", o.config.ClientID)
	formData.Set("client_secret", o.config.ClientSecret)
	
	// Create request
	req, err := http.NewRequest("POST", o.config.TokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("error creating token request: %w", err)
	}
	
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	
	// Send request
	resp, err := o.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending token request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading token response: %w", err)
	}
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error getting token: status %d, response: %s", resp.StatusCode, string(body))
	}
	
	// Parse response
	var tokenResp Token
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("error parsing token response: %w", err)
	}
	
	// Set expiration time
	tokenResp.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	o.Token = &tokenResp
	
	return nil
}
