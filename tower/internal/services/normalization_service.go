package services

import (
	"context"
	"fmt"
	
	"github.com/zepzeper/tower/internal/core/adapters"
	"github.com/zepzeper/tower/internal/core/connectors"
	"github.com/zepzeper/tower/internal/core/registry"
	"github.com/zepzeper/tower/internal/database"
)

// NormalizationService handles data normalization between APIs
type NormalizationService struct {
	dbManager       *database.Manager
	adapterFactory  *adapters.AdapterFactory
	schemaRegistry  *registry.SchemaRegistry
	connectorService *ConnectorService
}

// NewNormalizationService creates a new normalization service
func NewNormalizationService(
	dbManager *database.Manager,
	adapterFactory *adapters.AdapterFactory,
	schemaRegistry *registry.SchemaRegistry,
	connectorService *ConnectorService,
) *NormalizationService {
	return &NormalizationService{
		dbManager:       dbManager,
		adapterFactory:  adapterFactory,
		schemaRegistry:  schemaRegistry,
		connectorService: connectorService,
	}
}

// NormalizeForTarget normalizes source data for a specific target system
func (s *NormalizationService) NormalizeForTarget(
	ctx context.Context,
	sourceID string,
	targetID string,
	sourceData []connectors.DataPayload,
) ([]connectors.DataPayload, error) {
	// Get source connector
	sourceConnector, err := s.connectorService.GetConnector(sourceID)
	if err != nil {
		return nil, fmt.Errorf("source connector not found: %w", err)
	}
	
	// Get target connector
	targetConnector, err := s.connectorService.GetConnector(targetID)
	if err != nil {
		return nil, fmt.Errorf("target connector not found: %w", err)
	}
	
	// Get source adapter
	sourceAdapter, err := s.adapterFactory.GetAdapterForConnector(sourceConnector)
	if sourceAdapter == nil {
		// Try to create adapter on-the-fly
		sourceAdapter, err = s.adapterFactory.CreateAdapter(ctx, sourceConnector, sourceData)
		if err != nil {
			return nil, fmt.Errorf("error creating source adapter: %w", err)
		}
	}
	
	// Get target adapter
	targetAdapter, err := s.adapterFactory.GetAdapterForConnector(targetConnector)
	if targetAdapter == nil {
		// We need sample data for the target adapter
		// This is a bit of a chicken-and-egg problem, but we can use normalized data as a sample
		// First normalize one record from source
		if len(sourceData) == 0 {
			return nil, fmt.Errorf("no source data to normalize")
		}
		
		// Convert to canonical format
		canonical, err := sourceAdapter.ToCanonical(ctx, sourceData[0])
		if err != nil {
			return nil, fmt.Errorf("error converting to canonical format: %w", err)
		}
		
		// Create a sample target data from the canonical format (simple passthrough for now)
		targetSample := []connectors.DataPayload{
			connectors.DataPayload{
				"sample": "data",
				"entityType": "unknown",
			},
		}
		
		// Create target adapter
		targetAdapter, err = s.adapterFactory.CreateAdapter(ctx, targetConnector, targetSample)
		if err != nil {
			return nil, fmt.Errorf("error creating target adapter: %w", err)
		}
	}
	
	// Process each source data item
	normalizedData := make([]connectors.DataPayload, 0, len(sourceData))
	
	for _, item := range sourceData {
		// Step 1: Convert to canonical format
		canonical, err := sourceAdapter.ToCanonical(ctx, item)
		if err != nil {
			return nil, fmt.Errorf("error converting to canonical format: %w", err)
		}
		
		// Step 2: Convert from canonical to target format
		targetData, err := targetAdapter.FromCanonical(ctx, canonical)
		if err != nil {
			return nil, fmt.Errorf("error converting from canonical format: %w", err)
		}
		
		normalizedData = append(normalizedData, targetData)
	}
	
	return normalizedData, nil
}

// DetectSchema attempts to detect schema from sample data
func (s *NormalizationService) DetectSchema(
	ctx context.Context,
	connectorID string,
	sampleData []connectors.DataPayload,
) (connectors.Schema, error) {
	// Get connector
	connector, err := s.connectorService.GetConnector(connectorID)
	if err != nil {
		return connectors.Schema{}, fmt.Errorf("connector not found: %w", err)
	}
	
	// Try to get existing schema
	if schema, exists := s.schemaRegistry.GetSchema(connectorID); exists {
		return schema, nil
	}
	
	// No existing schema, try to discover it
	discoverer := connectors.NewSchemaDiscoverer(10)
	schema := discoverer.DiscoverSchema(connector.GetSchema().EntityName, sampleData)
	
	// Register the schema
	s.schemaRegistry.RegisterSchema(connectorID, schema)
	
	// Save to database for future reference
	err = s.schemaRegistry.SaveSchemaToDatabase(connectorID, schema)
	if err != nil {
		// Log the error but continue
		fmt.Printf("Error saving schema to database: %v\n", err)
	}
	
	return schema, nil
}

// GenerateMappings tries to generate mappings between source and target schemas
func (s *NormalizationService) GenerateMappings(
	ctx context.Context,
	sourceID string,
	targetID string,
) (string, error) {
	// Use the transformer service to generate transformer
	transformerService, ok := s.getTransformerService()
	if !ok {
		return "", fmt.Errorf("transformer service not available")
	}
	
	// Generate a name for the mapping
	mappingName := fmt.Sprintf("Auto-generated: %s to %s", sourceID, targetID)
	
	// Use the existing transformer service to create the mapping
	transformerID, err := transformerService.GenerateTransformer(
		sourceID,
		targetID,
		mappingName,
		"Automatically generated by the normalization service",
	)
	
	if err != nil {
		return "", fmt.Errorf("error generating transformer: %w", err)
	}
	
	// Register the mapping in the schema registry
	err = s.schemaRegistry.RegisterMapping(sourceID, targetID, transformerID)
	if err != nil {
		return "", fmt.Errorf("error registering mapping: %w", err)
	}
	
	return transformerID, nil
}

// Helper method to get transformer service
func (s *NormalizationService) getTransformerService() (*TransformerService, bool) {
	// In a real implementation, you'd want to inject this dependency
	// This is a workaround for the example
	
	// Check if the service is registered with the service registry
	// For now, just return false to indicate not available
	return nil, false
}
