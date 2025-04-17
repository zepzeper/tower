package dto

// TransformerRequest represents the request for creating/updating a transformer
type TransformerRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Mappings    []FieldMapping `json:"mappings"`
	Functions   []Function     `json:"functions,omitempty"`
}

// FieldMapping represents a field mapping
type FieldMapping struct {
	SourceField string `json:"sourceField"`
	TargetField string `json:"targetField"`
}

// Function represents a transformation function
type Function struct {
	Name        string   `json:"name"`
	TargetField string   `json:"targetField"`
	Args        []string `json:"args"`
}

// TransformerResponse represents the response for a transformer
type TransformerResponse struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Mappings    []FieldMapping `json:"mappings"`
	Functions   []Function     `json:"functions"`
	CreatedAt   string         `json:"createdAt"`
	UpdatedAt   string         `json:"updatedAt"`
}

// TransformerListResponse represents the response for listing transformers
type TransformerListResponse struct {
	Transformers []TransformerResponse `json:"transformers"`
}

// GenerateTransformerRequest represents the request for generating a transformer
type GenerateTransformerRequest struct {
	SourceID    string  `json:"sourceId"`
	TargetID    string  `json:"targetId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Threshold   float64 `json:"threshold,omitempty"`
}

// TransformDataRequest represents the request for transforming data
type TransformDataRequest struct {
	TransformerID string                 `json:"transformerId"`
	Data          map[string]interface{} `json:"data"`
}

// TransformDataResponse represents the response for transforming data
type TransformDataResponse struct {
	Result map[string]interface{} `json:"result"`
}
