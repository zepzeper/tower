package registry

import (
	"sync"

	"github.com/zepzeper/tower/internal/core/connectors"
)

// SchemaRegistry manages API schemas and their mappings
type SchemaFetcher struct {
	schemas       map[string]connectors.Schema
	mappings      map[string][]MappingEntry
	mu            sync.RWMutex
}

// NewSchemaRegistry creates a new schema registry
func NewSchemaRegistry() *SchemaFetcher {
	return &SchemaFetcher{
		schemas:   make(map[string]connectors.Schema),
		mappings:  make(map[string][]MappingEntry),
	}
}

// GetSchema retrieves a schema by connector ID
func (r *SchemaFetcher) GetSchema(connectorID string) (connectors.Schema, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	schema, exists := r.schemas[connectorID]
	return schema, exists
}

