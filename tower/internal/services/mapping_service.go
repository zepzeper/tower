
package services

import (
	"fmt"

	"github.com/zepzeper/tower/internal/core/transformers"
)

// MappingService provides mapping-related functionality
type MappingService struct {
	schemaFetcher SchemaFetcher
}

// SchemaFetcher is a dependency interface for retrieving schemas
type SchemaFetcher interface {
	GetSchema(schemaType string) (map[string]interface{}, error)
}

// NewMappingService creates a new instance with the given fetcher
func NewMappingService(fetcher SchemaFetcher) *MappingService {
	return &MappingService{
		schemaFetcher: fetcher,
	}
}

// GenerateMapping generates flattened field data + auto mappings
func (s *MappingService) GenerateMapping(sourceType, targetType string) (*transformers.MappingData, error) {
	sourceSchema, err := s.schemaFetcher.GetSchema(sourceType)
	if err != nil {
		return nil, fmt.Errorf("failed to get source schema (%s): %w", sourceType, err)
	}

	targetSchema, err := s.schemaFetcher.GetSchema(targetType)
	if err != nil {
		return nil, fmt.Errorf("failed to get target schema (%s): %w", targetType, err)
	}

	mapping := transformers.GenerateMappingData(sourceSchema, targetSchema)
	return &mapping, nil
}
