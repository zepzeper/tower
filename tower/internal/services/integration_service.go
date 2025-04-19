package services

import (
	"context"
	"fmt"
	
	"github.com/zepzeper/tower/internal/config"
	"github.com/zepzeper/tower/internal/core/adapters"
	"github.com/zepzeper/tower/internal/core/registry"
	"github.com/zepzeper/tower/internal/database"
)

// IntegrationService orchestrates the integration between systems
type IntegrationService struct {
	dbManager         *database.Manager
	connectorService  *ConnectorService
	transformerService *TransformerService
	connectionService *ConnectionService
	normalizationService *NormalizationService
	adapterFactory    *adapters.AdapterFactory
	schemaRegistry    *registry.SchemaRegistry
	configDir         string
}

// NewIntegrationService creates a new integration service
func NewIntegrationService(
	dbManager *database.Manager,
	connectorService *ConnectorService,
	transformerService *TransformerService,
	connectionService *ConnectionService,
	normalizationService *NormalizationService,
	adapterFactory *adapters.AdapterFactory,
	schemaRegistry *registry.SchemaRegistry,
	configDir string,
) *IntegrationService {
	return &IntegrationService{
		dbManager:         dbManager,
		connectorService:  connectorService,
		transformerService: transformerService,
		connectionService: connectionService,
		normalizationService: normalizationService,
		adapterFactory:    adapterFactory,
		schemaRegistry:    schemaRegistry,
		configDir:         configDir,
	}
}

// InitializeFromConfig initializes the system from configuration files
func (s *IntegrationService) InitializeFromConfig(ctx context.Context) error {
	// Load integration configuration
	integrationConfig, err := config.LoadIntegrationConfig(s.configDir)
	if err != nil {
		return fmt.Errorf("error loading integration config: %w", err)
	}
	
	// Process adapter mappings
	for _, adapterMapping := range integrationConfig.AdapterMappings {
		// Create a base adapter
		baseAdapter := adapters.NewBaseAdapter(adapterMapping.AdapterName, adapterMapping.EntityType)
		
		// Convert mappings
		mappings := make([]adapters.FieldMapping, 0, len(adapterMapping.Mappings))
		for _, mapping := range adapterMapping.Mappings {
			adapterMapping := adapters.FieldMapping{
				SourceField:    mapping.SourceField,
				CanonicalField: mapping.CanonicalField,
				IsRequired:     mapping.IsRequired,
			}
			
			// Set default value if specified
			if mapping.DefaultValue != "" {
				adapterMapping.DefaultValue = mapping.DefaultValue
			}
			
			// Add transformation functions if specified
			// In a real implementation, these would be parsed from strings
			// to actual functions, possibly using a registry of transformations
			
			mappings = append(mappings, adapterMapping)
		}
		
		// Register mappings with the adapter
		baseAdapter.RegisterMappings(mappings)
		
		// Register adapter with the factory
		s.adapterFactory.Register(baseAdapter)
	}
	
	// Process transformers
	for _, transformerConfig := range integrationConfig.Transformers {
		// Convert mappings
		mappings := make([]config.FieldMapping, 0, len(transformerConfig.Mappings))
		for _, mapping := range transformerConfig.Mappings {
			mappings = append(mappings, config.FieldMapping{
				SourceField: mapping.SourceField,
				TargetField: mapping.TargetField,
			})
		}
		
		// Check if transformer already exists
		_, err := s.transformerService.GetTransformer(transformerConfig.ID)
		if err != nil {
			// Create transformer
			_, err = s.transformerService.CreateTransformer(
				transformerConfig.Name,
				transformerConfig.Description,
				mappings,
			)
			if err != nil {
				return fmt.Errorf("error creating transformer %s: %w", transformerConfig.ID, err)
			}
		} else {
			// Update transformer
			err = s.transformerService.UpdateTransformer(
				transformerConfig.ID,
				transformerConfig.Name,
				transformerConfig.Description,
				mappings,
			)
			if err != nil {
				return fmt.Errorf("error updating transformer %s: %w", transformerConfig.ID, err)
			}
		}
	}
	
	// Process connections
	for _, connectionConfig := range integrationConfig.Connections {
		// Check if connection already exists
		_, err := s.connectionService.GetConnection(connectionConfig.ID)
		if err != nil {
			// Create connection
			_, err = s.connectionService.CreateConnection(
				ctx,
				connectionConfig.Name,
				connectionConfig.Description,
				connectionConfig.SourceID,
				connectionConfig.TargetID,
				connectionConfig.TransformerID,
				connectionConfig.Query,
				connectionConfig.Schedule,
			)
			if err != nil {
				return fmt.Errorf("error creating connection %s: %w", connectionConfig.ID, err)
			}
		} else {
			// Update connection
			err = s.connectionService.UpdateConnection(
				ctx,
				connectionConfig.ID,
				connectionConfig.Name,
				connectionConfig.Description,
				connectionConfig.SourceID,
				connectionConfig.TargetID,
				connectionConfig.TransformerID,
				connectionConfig.Query,
				connectionConfig.Schedule,
				connectionConfig.Active,
			)
			if err != nil {
				return fmt.Errorf("error updating connection %s: %w", connectionConfig.ID, err)
			}
		}
	}
	
	return nil
}

// SaveConfigFromDatabase saves the current state to configuration files
func (s *IntegrationService) SaveConfigFromDatabase() error {
	// Create a new integration config
	integrationConfig := &config.IntegrationConfig{
		Connections:    make([]config.ConnectionConfig, 0),
		Transformers:   make([]config.TransformerConfig, 0),
		AdapterMappings: make([]config.AdapterMapping, 0),
	}
	
	// Get all connections
	connections, err := s.connectionService.ListConnections()
	if err != nil {
		return fmt.Errorf("error retrieving connections: %w", err)
	}
	
	// Convert connections to config
	for _, conn := range connections {
		connectionConfig := config.ConnectionConfig{
			ID:            conn.ID,
			Name:          conn.Name,
			Description:   conn.Description,
			SourceID:      conn.SourceID,
			TargetID:      conn.TargetID,
			TransformerID: conn.TransformerID,
			Query:         conn.Query,
			Schedule:      conn.Schedule,
			Active:        conn.Active,
		}
		
		integrationConfig.Connections = append(integrationConfig.Connections, connectionConfig)
	}
	
	// Get all transformers
	transformers, err := s.transformerService.ListTransformers()
	if err != nil {
		return fmt.Errorf("error retrieving transformers: %w", err)
	}
	
	// Convert transformers to config
	for _, tr := range transformers {
		mappings := make([]config.FieldMapping, 0, len(tr.Mappings))
		for _, m := range tr.Mappings {
			mappings = append(mappings, config.FieldMapping{
				SourceField: m.SourceField,
				TargetField: m.TargetField,
			})
		}
		
		functions := make([]config.Function, 0, len(tr.Functions))
		for _, f := range tr.Functions {
			functions = append(functions, config.Function{
				Name:        f.Name,
				TargetField: f.TargetField,
				Args:        f.Args,
			})
		}
		
		transformerConfig := config.TransformerConfig{
			ID:          tr.ID,
			Name:        tr.Name,
			Description: tr.Description,
			Mappings:    mappings,
			Functions:   functions,
		}
		
		integrationConfig.Transformers = append(integrationConfig.Transformers, transformerConfig)
	}
	
	// Save the config
	return config.SaveIntegrationConfig(s.configDir, integrationConfig)
}

// AutoConfigureIntegration attempts to automatically configure an integration between systems
func (s *IntegrationService) AutoConfigureIntegration(
	ctx context.Context,
	sourceID string,
	targetID string,
	name string,
	description string,
	schedule string,
) (string, error) {
	// Step 1: Get or generate mappings
	transformerID, err := s.normalizationService.GenerateMappings(ctx, sourceID, targetID)
	if err != nil {
		return "", fmt.Errorf("error generating mappings: %w", err)
	}
	
	// Step 2: Create a connection
	connectionID, err := s.connectionService.CreateConnection(
		ctx,
		name,
		description,
		sourceID,
		targetID,
		transformerID,
		make(map[string]interface{}), // Empty query for now
		schedule,
	)
	if err != nil {
		return "", fmt.Errorf("error creating connection: %w", err)
	}
	
	return connectionID, nil
}
