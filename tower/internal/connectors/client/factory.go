package client

import (
	"fmt"
	
	"github.com/zepzeper/tower/internal/connectors/brincr"
	"github.com/zepzeper/tower/internal/connectors/client/core"
	// Import other clients as needed
)

// Factory creates API clients
type Factory struct{}

// NewFactory creates a new client factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateClient creates a client of the specified type
func (f *Factory) CreateClient(clientType string, demo bool) (core.APIClient, error) {
	switch clientType {
	case "brincr":
		client, err := brincr.NewClient(demo)
		if err != nil {
			return nil, err
		}
		return client, nil
	// Add cases for other client types
	default:
		return nil, fmt.Errorf("unknown client type: %s", clientType)
	}
}
