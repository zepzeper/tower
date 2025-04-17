package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/zepzeper/tower/internal/core/connectors"
	"github.com/zepzeper/tower/internal/core/registry"
	"github.com/zepzeper/tower/internal/core/transformers"
	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/database/models"
)

// TransformerService handles transformer-related business logic
type TransformerService struct {
	dbManager *database.Manager
	registry  *registry.ConnectorRegistry
}

// NewTransformerService creates a new transformer service
func NewTransformerService(dbManager *database.Manager, registry *registry.ConnectorRegistry) *TransformerService {
	return &TransformerService{
		dbManager: dbManager,
		registry:  registry,
	}
}

// ListTransformers retrieves all transformers
func (s *TransformerService) ListTransformers() ([]transformers.Transformer, error) {
	// Get from database
	transformerModels, err := s.dbManager.Repos.Transformer().GetAll()
	if err != nil {
		return nil, fmt.Errorf("error retrieving transformers: %w", err)
	}
	
	// Convert to domain objects
	result := make([]transformers.Transformer, len(transformerModels))
	for i, model := range transformerModels {
		transformer, err := s.modelToDomain(model)
		if err != nil {
			return nil, err
		}
		result[i] = *transformer
	}
	
	return result, nil
}

// CreateTransformer creates a new transformer
func (s *TransformerService) CreateTransformer(name, description string, mappings []transformers.FieldMapping) (string, error) {
	// Generate ID
	transformerID := generateID()
	
	// Convert mappings to JSON
	mappingsJSON, err := json.Marshal(mappings)
	if err != nil {
		return "", fmt.Errorf("error serializing mappings: %w", err)
	}
	
	// Create empty functions array
	functionsJSON := []byte("[]")
	
	// Create transformer model
	transformer := models.Transformer{
		ID:          transformerID,
		Name:        name,
		Description: sqlNullString(description),
		Mappings:    mappingsJSON,
		Functions:   functionsJSON,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// Save to database
	if err := s.dbManager.Repos.Transformer().Create(transformer); err != nil {
		return "", fmt.Errorf("error saving transformer: %w", err)
	}
	
	return transformerID, nil
}

// GetTransformer gets a transformer by ID
func (s *TransformerService) GetTransformer(id string) (*transformers.Transformer, error) {
	// Get from database
	transformerModel, err := s.dbManager.Repos.Transformer().GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transformer: %w", err)
	}
	
	// Convert to domain object
	return s.modelToDomain(transformerModel)
}

// UpdateTransformer updates a transformer
func (s *TransformerService) UpdateTransformer(id, name, description string, mappings []transformers.FieldMapping) error {
	// Get existing transformer
	transformerModel, err := s.dbManager.Repos.Transformer().GetByID(id)
	if err != nil {
		return fmt.Errorf("error retrieving transformer: %w", err)
	}
	
	// Convert mappings to JSON
	mappingsJSON, err := json.Marshal(mappings)
	if err != nil {
		return fmt.Errorf("error serializing mappings: %w", err)
	}
	
	// Update fields
	transformerModel.Name = name
	transformerModel.Description = sqlNullString(description)
	transformerModel.Mappings = mappingsJSON
	transformerModel.UpdatedAt = time.Now()
	
	// Save to database
	if err := s.dbManager.Repos.Transformer().Update(transformerModel); err != nil {
		return fmt.Errorf("error updating transformer: %w", err)
	}
	
	return nil
}

// DeleteTransformer deletes a transformer
func (s *TransformerService) DeleteTransformer(id string) error {
	return s.dbManager.Repos.Transformer().Delete(id)
}

// GenerateTransformer automatically creates a transformer between two connectors
func (s *TransformerService) GenerateTransformer(sourceID, targetID string, name, description string) (string, error) {
	// Get source connector
	sourceConn, exists := s.registry.Get(sourceID)
	if !exists {
		return "", fmt.Errorf("source connector not found: %s", sourceID)
	}
	
	// Get target connector
	targetConn, exists := s.registry.Get(targetID)
	if !exists {
		return "", fmt.Errorf("target connector not found: %s", targetID)
	}
	
	// Get schemas
	sourceSchema := sourceConn.GetSchema()
	targetSchema := targetConn.GetSchema()
	
	// Generate mappings
	autoMapper := transformers.NewAutoMapper(0.7) // 0.7 similarity threshold
	mappings := autoMapper.GenerateMappings(sourceSchema, targetSchema)
	
	// Create transformer name if not provided
	if name == "" {
		name = fmt.Sprintf("Auto-generated: %s to %s", sourceID, targetID)
	}
	
	// Create transformer description if not provided
	if description == "" {
		description = fmt.Sprintf("Automatically generated transformer from %s to %s", sourceID, targetID)
	}
	
	// Create and save the transformer
	return s.CreateTransformer(name, description, mappings)
}

// TransformData transforms data using a specific transformer
func (s *TransformerService) TransformData(ctx context.Context, transformerID string, data connectors.DataPayload) (connectors.DataPayload, error) {
	// Get transformer
	transformer, err := s.GetTransformer(transformerID)
	if err != nil {
		return nil, err
	}
	
	// Transform data
	result, err := transformer.Transform(data)
	if err != nil {
		return nil, fmt.Errorf("error transforming data: %w", err)
	}
	
	return result, nil
}

// Helper function to convert model to domain object
func (s *TransformerService) modelToDomain(model models.Transformer) (*transformers.Transformer, error) {
	// Convert mappings from JSON
	var mappings []transformers.FieldMapping
	if err := json.Unmarshal(model.Mappings, &mappings); err != nil {
		return nil, fmt.Errorf("error parsing mappings: %w", err)
	}
	
	// Convert functions from JSON
	var functions []transformers.Function
	if err := json.Unmarshal(model.Functions, &functions); err != nil {
		return nil, fmt.Errorf("error parsing functions: %w", err)
	}
	
	// Create transformer object
	transformer := &transformers.Transformer{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description.String,
		Mappings:    mappings,
		Functions:   functions,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
	
	return transformer, nil
}

// Helper function to generate SQL null string
func sqlNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// Helper function to generate ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
