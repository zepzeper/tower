package dto

// ConnectorListResponse represents the response for listing connectors
type ConnectorListResponse struct {
	Connectors []ConnectorInfo `json:"connectors"`
}

// ConnectorInfo represents information about a connector
type ConnectorInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// ConnectorSchemaResponse represents the response for a connector schema
type ConnectorSchemaResponse struct {
	ID          string                         `json:"id"`
	EntityName  string                         `json:"entityName"`
	Fields      map[string]FieldDefinitionInfo `json:"fields"`
}

// FieldDefinitionInfo represents information about a field definition
type FieldDefinitionInfo struct {
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Path        string   `json:"path"`
	EnumValues  []string `json:"enumValues,omitempty"`
	Description string   `json:"description,omitempty"`
}

// ConnectorTestRequest represents the request for testing a connector
type ConnectorTestRequest struct {
	ID     string                 `json:"id"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// ConnectorTestResponse represents the response for testing a connector
type ConnectorTestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// DataFetchRequest represents the request for fetching data from a connector
type DataFetchRequest struct {
	ID    string                 `json:"id"`
	Query map[string]interface{} `json:"query"`
}

// DataFetchResponse represents the response for fetching data
type DataFetchResponse struct {
	Success bool        `json:"success"`
	Count   int         `json:"count"`
	Data    interface{} `json:"data"`
}
