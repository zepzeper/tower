package mapping

// TestRequest represents the data needed to test a mapping
type TestRequest struct {
	ConnectionID string            `json:"connectionId"`
	SourceType   string            `json:"sourceType"`
	TargetType   string            `json:"targetType"`
	Mappings     []MappingMetadata `json:"mappings"`
}

// MappingMetadata represents the metadata for a mapping
type MappingMetadata struct {
	ID            string  `json:"id"`
	SourcePath    string  `json:"sourcePath"`
	TargetPath    string  `json:"targetPath"`
	Transform     *string `json:"transform"`
	SourceFieldID string  `json:"sourceFieldId"`
	TargetFieldID string  `json:"targetFieldId"`
}

// TestResponse is the result of testing a mapping
type TestResponse struct {
	SourceData      map[string]interface{} `json:"sourceData"`
	TransformedData map[string]interface{} `json:"transformedData"`
}

// FieldDefinition represents a single mappable field
type FieldDefinition struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Path     string `json:"path"`
	Sample   string `json:"sample,omitempty"`
	Required bool   `json:"required,omitempty"`
}

// MappingDefinition represents a field mapping
type MappingDefinition struct {
	ID          string  `json:"id"`
	SourceField string  `json:"sourceField"`
	TargetField string  `json:"targetField"`
	Transform   *string `json:"transform"`
}

// MappingData holds all fields and suggested mappings
type MappingData struct {
	SourceFields []FieldDefinition   `json:"sourceFields"`
	TargetFields []FieldDefinition   `json:"targetFields"`
	Mappings     []MappingDefinition `json:"mappings"`
}

type schemaTransformOptions struct {
	Prefix          string
	IncludeRequired bool
}

