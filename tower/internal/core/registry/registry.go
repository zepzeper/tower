// tower/internal/core/registry/schema_registry.go
package registry

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/zepzeper/tower/internal/core/connectors"
	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/database/models"
)

// SchemaRegistry manages API schemas and their mappings
type SchemaRegistry struct {
	schemas       map[string]connectors.Schema
	mappings      map[string][]MappingEntry
	dbManager     *database.Manager
	mu            sync.RWMutex
}

// MappingEntry defines a mapping between source and target schemas
type MappingEntry struct {
	SourceSchema string
	TargetSchema string
	MappingID    string // ID of the transformer that handles this mapping
}

// NewSchemaRegistry creates a new schema registry
func NewSchemaRegistry(dbManager *database.Manager) *SchemaRegistry {
	return &SchemaRegistry{
		schemas:   make(map[string]connectors.Schema),
		mappings:  make(map[string][]MappingEntry),
		dbManager: dbManager,
	}
}

// RegisterSchema registers a schema with the registry
func (r *SchemaRegistry) RegisterSchema(connectorID string, schema connectors.Schema) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.schemas[connectorID] = schema
}

// GetSchema retrieves a schema by connector ID
func (r *SchemaRegistry) GetSchema(connectorID string) (connectors.Schema, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	schema, exists := r.schemas[connectorID]
	return schema, exists
}

// RegisterMapping registers a mapping between source and target schemas
func (r *SchemaRegistry) RegisterMapping(sourceID, targetID, mappingID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Check if both schemas exist
	_, sourceExists := r.schemas[sourceID]
	_, targetExists := r.schemas[targetID]
	
	if !sourceExists {
		return fmt.Errorf("source schema not found: %s", sourceID)
	}
	
	if !targetExists {
		return fmt.Errorf("target schema not found: %s", targetID)
	}
	
	// Create mapping key (source-to-target)
	mappingKey := fmt.Sprintf("%s-to-%s", sourceID, targetID)
	
	// Add mapping entry
	mapping := MappingEntry{
		SourceSchema: sourceID,
		TargetSchema: targetID,
		MappingID:    mappingID,
	}
	
	// Add to or update mappings
	existingMappings, exists := r.mappings[mappingKey]
	if !exists {
		r.mappings[mappingKey] = []MappingEntry{mapping}
	} else {
		// Check if mapping already exists
		for i, existing := range existingMappings {
			if existing.MappingID == mappingID {
				// Update existing mapping
				existingMappings[i] = mapping
				r.mappings[mappingKey] = existingMappings
				return nil
			}
		}
		
		// Add new mapping
		r.mappings[mappingKey] = append(existingMappings, mapping)
	}
	
	return nil
}

// GetMappings retrieves mappings between source and target schemas
func (r *SchemaRegistry) GetMappings(sourceID, targetID string) ([]MappingEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Create mapping key
	mappingKey := fmt.Sprintf("%s-to-%s", sourceID, targetID)
	
	mappings, exists := r.mappings[mappingKey]
	if !exists {
		return nil, fmt.Errorf("no mappings found for %s to %s", sourceID, targetID)
	}
	
	return mappings, nil
}

// LoadFromDatabase loads schemas and mappings from the database
func (r *SchemaRegistry) LoadFromDatabase() error {
	// Load transformers
	transformers, err := r.dbManager.Repos.Transformer().GetAll()
	if err != nil {
		return fmt.Errorf("error loading transformers: %w", err)
	}
	
	// Load connections to get source-target pairs
	connections, err := r.dbManager.Repos.Connection().GetAll()
	if err != nil {
		return fmt.Errorf("error loading connections: %w", err)
	}
	
	// Process connections to register mappings
	for _, conn := range connections {
		// Try to register mapping
		err := r.RegisterMapping(conn.SourceID, conn.TargetID, conn.TransformerID)
		if err != nil {
			// Log error but continue with other mappings
			fmt.Printf("Error registering mapping for connection %s: %v\n", conn.ID, err)
		}
	}
	
	// Also load transformers that aren't connected to any connection yet
	for _, transformer := range transformers {
		var mappings []map[string]string
		
		if err := json.Unmarshal(transformer.Mappings, &mappings); err != nil {
			fmt.Printf("Error parsing mappings for transformer %s: %v\n", transformer.ID, err)
			continue
		}
		
		// If we could determine source and target from the mappings, we could register them
		// This would depend on your specific mapping structure
	}
	
	return nil
}

// SaveSchemaToDatabase saves a connector schema to the database
func (r *SchemaRegistry) SaveSchemaToDatabase(connectorID string, schema connectors.Schema) error {
	// Convert schema to JSON
	_, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("error serializing schema: %w", err)
	}
	
	// Create a transformer model as a placeholder for the schema
	// You might want to create a dedicated schema table in a production system
	transformer := models.Transformer{
		ID:          fmt.Sprintf("schema-%s", connectorID),
		Name:        fmt.Sprintf("Schema for %s", connectorID),
		Description: sql.NullString{String: fmt.Sprintf("Auto-generated schema for connector %s", connectorID), Valid: true},
		Mappings:    []byte("[]"), // Empty mappings
		Functions:   []byte("[]"), // Empty functions
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// You'd probably want to store the schema in a custom field
	// This is a workaround using the existing structure
	
	// Try to create, update if exists
	err = r.dbManager.Repos.Transformer().Create(transformer)
	if err != nil {
		// Try to update
		err = r.dbManager.Repos.Transformer().Update(transformer)
		if err != nil {
			return fmt.Errorf("error saving schema: %w", err)
		}
	}
	
	return nil
}

// Create custom database model and repo for schemas in a real implementation
