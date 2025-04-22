package models

import (
	"database/sql"
	"time"
)

// Credentials represents an API connection in the database
type Credentials struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description sql.NullString    `json:"description"`
	Type        string            `json:"type"`
	Active      bool              `json:"active"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Config      map[string]string `json:"config,omitempty"` // Not stored in DB directly
}

// CredentialsConfig represents a key-value configuration for a connection
type CredentialsConfig struct {
	ConnectionID string    `json:"connectionId"`
	Key          string    `json:"key"`
	Value        string    `json:"value"`
	IsSecret     bool      `json:"isSecret"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
