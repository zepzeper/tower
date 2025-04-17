package services

import (
	"context"
	"fmt"
	
	"github.com/zepzeper/tower/internal/core/connectors"
	"github.com/zepzeper/tower/internal/core/registry"
)

// ConnectorService handles connector-related business logic
type ConnectorService struct {
	registry *registry.ConnectorRegistry
}

// NewConnectorService creates a new connector service
func NewConnectorService(registry *registry.ConnectorRegistry) *ConnectorService {
	return &ConnectorService{
		registry: registry,
	}
}

// ListConnectors returns all registered connectors
func (s *ConnectorService) ListConnectors() []string {
	return s.registry.List()
}

// GetConnector returns a specific connector by ID
func (s *ConnectorService) GetConnector(id string) (connectors.Connector, error) {
	connector, exists := s.registry.Get(id)
	if !exists {
		return nil, fmt.Errorf("connector not found: %s", id)
	}
	return connector, nil
}

// TestConnector tests the connection to a specific connector
func (s *ConnectorService) TestConnector(ctx context.Context, id string) error {
	connector, err := s.GetConnector(id)
	if err != nil {
		return err
	}
	
	return connector.Connect(ctx)
}

// GetSchema returns the schema for a specific connector
func (s *ConnectorService) GetSchema(id string) (connectors.Schema, error) {
	connector, err := s.GetConnector(id)
	if err != nil {
		return connectors.Schema{}, err
	}
	
	return connector.GetSchema(), nil
}

// FetchData fetches data from a connector
func (s *ConnectorService) FetchData(ctx context.Context, id string, query map[string]interface{}) ([]connectors.DataPayload, error) {
	connector, err := s.GetConnector(id)
	if err != nil {
		return nil, err
	}
	
	return connector.Fetch(ctx, query)
}

// PushData pushes data to a connector
func (s *ConnectorService) PushData(ctx context.Context, id string, data []connectors.DataPayload) error {
	connector, err := s.GetConnector(id)
	if err != nil {
		return err
	}
	
	return connector.Push(ctx, data)
}
