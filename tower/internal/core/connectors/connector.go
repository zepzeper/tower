package connectors

import (
	"context"
)

// DataPayload represents generic data exchanged between systems
type DataPayload map[string]interface{}

// Connector defines the interface that all API connectors must implement
type Connector interface {
	// Connect establishes the connection with the API
	Connect(ctx context.Context) error
	
	// Fetch retrieves data from the API
	Fetch(ctx context.Context, query map[string]interface{}) ([]DataPayload, error)
	
	// Push sends data to the API
	Push(ctx context.Context, data []DataPayload) error
	
	// GetSchema returns the data schema for this connector
	GetSchema() Schema
}

// ConnectorInfo contains metadata about a connector
type ConnectorInfo struct {
	ID          string
	Name        string
	Description string
	Type        string
}
