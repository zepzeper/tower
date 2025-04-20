package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Connection represents a data integration connection
type Connection struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Description   sql.NullString  `json:"description"`
	SourceID      string          `json:"sourceId"`
	TargetID      string          `json:"targetId"`
	TransformerID string          `json:"transformerId"`
	Config        json.RawMessage `json:"config"`
	Schedule      string          `json:"schedule"`
	Active        bool            `json:"active"`
	LastRun       sql.NullTime    `json:"lastRun"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

// APIConnection represents an API connection in the database
type APIConnection struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description sql.NullString    `json:"description"`
	Type        string            `json:"type"`
	Active      bool              `json:"active"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Config      map[string]string `json:"config,omitempty"` // Not stored in DB directly
}

// APIConnectionConfig represents a key-value configuration for a connection
type APIConnectionConfig struct {
	ConnectionID string    `json:"connectionId"`
	Key          string    `json:"key"`
	Value        string    `json:"value"`
	IsSecret     bool      `json:"isSecret"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
