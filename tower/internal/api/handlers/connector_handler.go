package handlers

import (
	"encoding/json"
	"net/http"
	
	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/api/dto"
	"github.com/zepzeper/tower/internal/api/response"
	"github.com/zepzeper/tower/internal/services"
)

// ConnectorHandler handles connector-related API endpoints
type ConnectorHandler struct {
	connectorService *services.ConnectorService
}

// NewConnectorHandler creates a new connector handler
func NewConnectorHandler(connectorService *services.ConnectorService) *ConnectorHandler {
	return &ConnectorHandler{
		connectorService: connectorService,
	}
}

// ListConnectors handles GET /api/v1/connectors
func (h *ConnectorHandler) ListConnectors(w http.ResponseWriter, r *http.Request) {
	// Get connectors from service
	connectorIDs := h.connectorService.ListConnectors()
	
	// Convert to response format
	connectors := make([]dto.ConnectorInfo, 0, len(connectorIDs))
	for _, id := range connectorIDs {
		connector, err := h.connectorService.GetConnector(id)
		if err != nil {
			continue // Skip if connector not found
		}
		
		schema := connector.GetSchema()
		connectors = append(connectors, dto.ConnectorInfo{
			ID:          id,
			Name:        schema.EntityName,
			Type:        id, // Use ID as type for now
			Description: "Connector for " + schema.EntityName,
		})
	}
	
	response.JSON(w, dto.ConnectorListResponse{
		Connectors: connectors,
	}, http.StatusOK)
}

// GetConnectorSchema handles GET /api/v1/connectors/{connectorID}/schema
func (h *ConnectorHandler) GetConnectorSchema(w http.ResponseWriter, r *http.Request) {
	connectorID := chi.URLParam(r, "connectorID")
	
	// Get schema from service
	schema, err := h.connectorService.GetSchema(connectorID)
	if err != nil {
		response.Error(w, "Failed to get schema: "+err.Error(), http.StatusNotFound)
		return
	}
	
	// Convert to response format
	fields := make(map[string]dto.FieldDefinitionInfo)
	for name, def := range schema.Fields {
		fields[name] = dto.FieldDefinitionInfo{
			Type:        def.Type,
			Required:    def.Required,
			Path:        def.Path,
			EnumValues:  def.EnumValues,
			Description: def.Description,
		}
	}
	
	response.JSON(w, dto.ConnectorSchemaResponse{
		ID:         connectorID,
		EntityName: schema.EntityName,
		Fields:     fields,
	}, http.StatusOK)
}

// TestConnector handles POST /api/v1/connectors/test
func (h *ConnectorHandler) TestConnector(w http.ResponseWriter, r *http.Request) {
	var req dto.ConnectorTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Test connector
	err := h.connectorService.TestConnector(r.Context(), req.ID)
	if err != nil {
		response.JSON(w, dto.ConnectorTestResponse{
			Success: false,
			Message: "Connection test failed: " + err.Error(),
		}, http.StatusOK)
		return
	}
	
	response.JSON(w, dto.ConnectorTestResponse{
		Success: true,
		Message: "Connection successful",
	}, http.StatusOK)
}

// FetchData handles POST /api/v1/connectors/fetch
func (h *ConnectorHandler) FetchData(w http.ResponseWriter, r *http.Request) {
	var req dto.DataFetchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Fetch data
	data, err := h.connectorService.FetchData(r.Context(), req.ID, req.Query)
	if err != nil {
		response.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, dto.DataFetchResponse{
		Success: true,
		Count:   len(data),
		Data:    data,
	}, http.StatusOK)
}
