package models

import (
	"time"
)

// APIKey represents an API key in the database
type APIKey struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	LastUsed  time.Time `json:"lastUsed"`
}

// ToAPIKey converts a database APIKey to an API response APIKey
func (a *APIKey) ToAPIKey() interface{} {
	// For security, don't include the actual key value in the response
	// unless specifically requested
	return map[string]interface{}{
		"id":        a.ID,
		"userId":    a.UserID,
		"name":      a.Name,
		"createdAt": a.CreatedAt,
		"expiresAt": a.ExpiresAt,
		"lastUsed":  a.LastUsed,
	}
}

// ToAPIKeyWithValue is similar to ToAPIKey but includes the key value
// Should only be used when first creating a key
func (a *APIKey) ToAPIKeyWithValue() interface{} {
	apiKey := a.ToAPIKey().(map[string]interface{})
	apiKey["key"] = a.Key
	return apiKey
}
