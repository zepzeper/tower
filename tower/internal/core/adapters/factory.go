package adapters

import (
	"context"
	"fmt"
	"sync"
	
	"github.com/zepzeper/tower/internal/core/connectors"
)

// AdapterFactory creates and manages API adapters
type AdapterFactory struct {
	adapters map[string]Adapter
	mu       sync.RWMutex
}

// NewAdapterFactory creates a new adapter factory
func NewAdapterFactory() *AdapterFactory {
	return &AdapterFactory{
		adapters: make(map[string]Adapter),
	}
}

// Register registers an adapter with the factory
func (f *AdapterFactory) Register(adapter Adapter) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.adapters[adapter.Name()] = adapter
}

// GetAdapter retrieves an adapter by name
func (f *AdapterFactory) GetAdapter(name string) (Adapter, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	adapter, exists := f.adapters[name]
	if !exists {
		return nil, fmt.Errorf("adapter not found: %s", name)
	}
	
	return adapter, nil
}

// GetAdapterForConnector tries to find an adapter for a specific connector
func (f *AdapterFactory) GetAdapterForConnector(connector connectors.Connector) (Adapter, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	schema := connector.GetSchema()
	
	// Try to find adapter by entity name
	for _, adapter := range f.adapters {
		if adapter.SourceType() == schema.EntityName {
			return adapter, nil
		}
	}
	
	return nil, fmt.Errorf("no adapter found for entity type: %s", schema.EntityName)
}

// CreateAdapter creates a new adapter for a connector by analyzing its schema
func (f *AdapterFactory) CreateAdapter(ctx context.Context, connector connectors.Connector, sampleData []connectors.DataPayload) (Adapter, error) {
	schema := connector.GetSchema()
	
	// Create a base adapter with automatic name
	baseAdapter := NewBaseAdapter(
		fmt.Sprintf("%s-adapter", schema.EntityName),
		schema.EntityName,
	)
	
	// Try to create field mappings automatically
	// This is a simplified version - a real implementation would be more sophisticated
	mappings := make([]FieldMapping, 0)
	
	for fieldName, fieldDef := range schema.Fields {
		canonicalField := fieldName // Simple 1:1 mapping for demonstration
		
		mapping := FieldMapping{
			SourceField:    fieldDef.Path,
			CanonicalField: canonicalField,
			IsRequired:     fieldDef.Required,
		}
		
		mappings = append(mappings, mapping)
	}
	
	// Register the mappings
	baseAdapter.RegisterMappings(mappings)
	
	// Register the adapter
	f.Register(baseAdapter)
	
	return baseAdapter, nil
}

// ListAdapters returns the names of all registered adapters
func (f *AdapterFactory) ListAdapters() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	adapterNames := make([]string, 0, len(f.adapters))
	for name := range f.adapters {
		adapterNames = append(adapterNames, name)
	}
	
	return adapterNames
}
