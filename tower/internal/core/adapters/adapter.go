package adapters

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	
	"github.com/zepzeper/tower/internal/core/cdm"
	"github.com/zepzeper/tower/internal/core/connectors"
)

// Adapter converts between API-specific formats and the canonical data model
type Adapter interface {
	// Name returns the name of the adapter
	Name() string
	
	// SourceType returns the type of the source system
	SourceType() string
	
	// ToCanonical converts data from source format to canonical format
	ToCanonical(ctx context.Context, sourceData connectors.DataPayload) (interface{}, error)
	
	// FromCanonical converts data from canonical format to source format
	FromCanonical(ctx context.Context, canonicalData interface{}) (connectors.DataPayload, error)
	
	// DiscoverSourceSchema discovers the schema of the source system
	DiscoverSourceSchema(ctx context.Context, sampleData []connectors.DataPayload) (connectors.Schema, error)
}

// BaseAdapter provides common functionality for adapters
type BaseAdapter struct {
	name       string
	sourceType string
	mappings   map[string]FieldMapping
}

// FieldMapping defines how data should be mapped between API and canonical model
type FieldMapping struct {
	SourceField      string
	CanonicalField   string
	IsRequired       bool
	DefaultValue     interface{}
	TransformToCanon func(interface{}) (interface{}, error)
	TransformFromCanon func(interface{}) (interface{}, error)
}

// NewBaseAdapter creates a new base adapter
func NewBaseAdapter(name, sourceType string) *BaseAdapter {
	return &BaseAdapter{
		name:       name,
		sourceType: sourceType,
		mappings:   make(map[string]FieldMapping),
	}
}

// Name returns the adapter name
func (a *BaseAdapter) Name() string {
	return a.name
}

// SourceType returns the type of source system
func (a *BaseAdapter) SourceType() string {
	return a.sourceType
}

// RegisterMapping registers a field mapping
func (a *BaseAdapter) RegisterMapping(mapping FieldMapping) {
	a.mappings[mapping.CanonicalField] = mapping
}

// RegisterMappings registers multiple field mappings
func (a *BaseAdapter) RegisterMappings(mappings []FieldMapping) {
	for _, mapping := range mappings {
		a.RegisterMapping(mapping)
	}
}

// ToCanonical converts from source to canonical format with registered mappings
func (a *BaseAdapter) ToCanonical(ctx context.Context, sourceData connectors.DataPayload) (interface{}, error) {
	// Determine entity type
	entityType := a.determineEntityType(sourceData)
	
	// Create the result map
	result := make(map[string]interface{})
	result["entityType"] = entityType
	
	// Apply mappings
	for canonicalField, mapping := range a.mappings {
		sourceValue, exists := getNestedValue(sourceData, mapping.SourceField)
		
		if !exists {
			if mapping.IsRequired {
				return nil, fmt.Errorf("required field %s missing in source data", mapping.SourceField)
			}
			
			if mapping.DefaultValue != nil {
				result[canonicalField] = mapping.DefaultValue
			}
			continue
		}
		
		// Apply transformation if provided
		if mapping.TransformToCanon != nil {
			transformedValue, err := mapping.TransformToCanon(sourceValue)
			if err != nil {
				return nil, fmt.Errorf("error transforming field %s: %w", mapping.SourceField, err)
			}
			result[canonicalField] = transformedValue
		} else {
			result[canonicalField] = sourceValue
		}
	}
	
	// Convert to appropriate entity type
	return cdm.Convert(entityType, result)
}

// FromCanonical converts from canonical to source format with registered mappings
func (a *BaseAdapter) FromCanonical(ctx context.Context, canonicalData interface{}) (connectors.DataPayload, error) {
	// Convert entity to map if it's not already a map
	var canonicalMap map[string]interface{}
	
	if reflect.TypeOf(canonicalData).Kind() == reflect.Map {
		var ok bool
		canonicalMap, ok = canonicalData.(map[string]interface{})
		if !ok {
			return nil, errors.New("failed to convert canonical data to map")
		}
	} else {
		canonicalMap = cdm.ConvertFromEntity(canonicalData)
	}
	
	// Create result map
	result := make(connectors.DataPayload)
	
	// Apply reverse mappings
	for canonicalField, mapping := range a.mappings {
		canonicalValue, exists := canonicalMap[canonicalField]
		if !exists {
			continue
		}
		
		// Apply transformation if provided
		var sourceValue interface{}
		var err error
		
		if mapping.TransformFromCanon != nil {
			sourceValue, err = mapping.TransformFromCanon(canonicalValue)
			if err != nil {
				return nil, fmt.Errorf("error reverse transforming field %s: %w", canonicalField, err)
			}
		} else {
			sourceValue = canonicalValue
		}
		
		// Set in result using source field path
		if err := setNestedValue(result, mapping.SourceField, sourceValue); err != nil {
			return nil, fmt.Errorf("error setting field %s: %w", mapping.SourceField, err)
		}
	}
	
	return result, nil
}

// Helper functions

// determineEntityType tries to determine the entity type from source data
func (a *BaseAdapter) determineEntityType(sourceData connectors.DataPayload) string {
	// Implementation will depend on your specific requirements
	// This is a placeholder
	
	// Try to get explicit entity type if available
	if entityType, ok := sourceData["entity_type"].(string); ok {
		return entityType
	}
	
	// Try to infer from data structure
	// E.g., presence of certain fields
	if _, ok := sourceData["sku"]; ok {
		return "product"
	}
	
	if _, ok := sourceData["email"]; ok {
		return "customer"
	}
	
	if _, ok := sourceData["order_number"]; ok {
		return "order"
	}
	
	// Default to generic entity
	return "entity"
}

// getNestedValue gets a value from a potentially nested map using dot notation
func getNestedValue(data map[string]interface{}, path string) (interface{}, bool) {
	// Implementation from your existing code in transformers
	// Not duplicating the full implementation here
	return nil, false
}

// setNestedValue sets a value in a potentially nested map using dot notation
func setNestedValue(data map[string]interface{}, path string, value interface{}) error {
	// Implementation from your existing code in transformers
	// Not duplicating the full implementation here
	return nil
}

// DiscoverSourceSchema discovers the schema of the source system
func (a *BaseAdapter) DiscoverSourceSchema(ctx context.Context, sampleData []connectors.DataPayload) (connectors.Schema, error) {
	discoverer := connectors.NewSchemaDiscoverer(10)
	return discoverer.DiscoverSchema(a.sourceType, sampleData), nil
}
