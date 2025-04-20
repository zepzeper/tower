package mapping

import (
	"fmt"

	"github.com/zepzeper/tower/internal/database"
)

// SchemaFetcher is a dependency interface for retrieving schemas
type SchemaFetcher interface {
	GetSchema(schemaType, operation, sourceortarget string) (map[string]interface{}, error)
	GetApis() (map[string]interface{}, error)
}

// Service provides mapping-related functionality
type Service struct {
	schemaFetcher   SchemaFetcher
	transformer     *Transformer
	databaseManager *database.Manager
}

// NewService creates a new instance with the given fetcher
func NewService(fetcher SchemaFetcher, databaseManager *database.Manager) *Service {
	return &Service{
		schemaFetcher:   fetcher,
		transformer:     NewTransformer(),
		databaseManager: databaseManager,
	}
}

// GenerateMapping generates flattened field data + auto mappings
func (s *Service) GenerateMapping(sourceType, targetType, operation string) (*MappingData, error) {
	sourceSchema, err := s.schemaFetcher.GetSchema(sourceType, operation, "source")
	if err != nil {
		return nil, fmt.Errorf("failed to get source schema (%s): %w", sourceType, err)
	}

	targetSchema, err := s.schemaFetcher.GetSchema(targetType, operation, "target")
	if err != nil {
		return nil, fmt.Errorf("failed to get target schema (%s): %w", targetType, err)
	}

	mapping, err := GenerateMappingData(sourceSchema, targetSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to generate mappings: %w", err)
	}

	return &mapping, nil
}

// ApplyMappings tests mappings by applying them to sample data
func (s *Service) ApplyMappings(req TestRequest) (*TestResponse, error) {
	// Get source data
	sourceData, err := s.schemaFetcher.GetSchema(req.SourceType, "products", "source")
	if err != nil {
		return nil, fmt.Errorf("failed to get source schema (%s): %w", req.SourceType, err)
	}

	// Get target schema for type information
	targetSchema, err := s.schemaFetcher.GetSchema(req.TargetType, "products", "target")
	if err != nil {
		return nil, fmt.Errorf("failed to get target schema (%s): %w", req.TargetType, err)
	}

	// Use the transformer to apply mappings and convert to target format
	transformedData := s.transformer.TransformFields(sourceData, targetSchema, req.Mappings)

	return &TestResponse{
		SourceData:      sourceData,
		TransformedData: transformedData,
	}, nil
}
