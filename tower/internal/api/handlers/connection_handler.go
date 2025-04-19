package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	
	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/api/dto"
	"github.com/zepzeper/tower/internal/response"
	"github.com/zepzeper/tower/internal/services"
)

// ConnectionHandler handles connection-related API endpoints
type ConnectionHandler struct {
	connectionService *services.ConnectionService
}

// NewConnectionHandler creates a new connection handler
func NewConnectionHandler(connectionService *services.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{
		connectionService: connectionService,
	}
}

// ListConnections handles GET /api/v1/connections
func (h *ConnectionHandler) ListConnections(w http.ResponseWriter, r *http.Request) {
	// Get connections from service
	connections, err := h.connectionService.ListConnections()
	if err != nil {
		response.Error(w, "Failed to list connections: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Convert to response format
	connectionResponses := make([]dto.ConnectionResponse, len(connections))
	for i, conn := range connections {
		connectionResponses[i] = connectionToResponse(conn)
	}
	
	response.JSON(w, dto.ConnectionListResponse{
		Connections: connectionResponses,
	}, http.StatusOK)
}

// GetConnection handles GET /api/v1/connections/{connectionID}
func (h *ConnectionHandler) GetConnection(w http.ResponseWriter, r *http.Request) {
	connectionID := chi.URLParam(r, "connectionID")
	
	// Get connection from service
	connection, err := h.connectionService.GetConnection(connectionID)
	if err != nil {
		response.Error(w, "Failed to get connection: "+err.Error(), http.StatusNotFound)
		return
	}
	
	// Convert to response format
	resp := connectionToResponse(*connection)
	
	response.JSON(w, resp, http.StatusOK)
}

// CreateConnection handles POST /api/v1/connections
func (h *ConnectionHandler) CreateConnection(w http.ResponseWriter, r *http.Request) {
	var req dto.ConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.Name == "" {
		response.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	if req.SourceID == "" || req.TargetID == "" || req.TransformerID == "" {
		response.Error(w, "Source, target, and transformer IDs are required", http.StatusBadRequest)
		return
	}
	
	// Create connection
	id, err := h.connectionService.CreateConnection(
		r.Context(),
		req.Name,
		req.Description,
		req.SourceID,
		req.TargetID,
		req.TransformerID,
		req.Query,
		req.Schedule,
	)
	
	if err != nil {
		response.Error(w, "Failed to create connection: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get created connection
	connection, err := h.connectionService.GetConnection(id)
	if err != nil {
		response.Error(w, "Connection created but failed to retrieve: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Convert to response format
	resp := connectionToResponse(*connection)
	
	response.JSON(w, resp, http.StatusCreated)
}

// UpdateConnection handles PUT /api/v1/connections/{connectionID}
func (h *ConnectionHandler) UpdateConnection(w http.ResponseWriter, r *http.Request) {
	connectionID := chi.URLParam(r, "connectionID")
	
	var req dto.ConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.Name == "" {
		response.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	if req.SourceID == "" || req.TargetID == "" || req.TransformerID == "" {
		response.Error(w, "Source, target, and transformer IDs are required", http.StatusBadRequest)
		return
	}
	
	// Update connection
	err := h.connectionService.UpdateConnection(
		r.Context(),
		connectionID,
		req.Name,
		req.Description,
		req.SourceID,
		req.TargetID,
		req.TransformerID,
		req.Query,
		req.Schedule,
		req.Active,
	)
	
	if err != nil {
		response.Error(w, "Failed to update connection: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get updated connection
	connection, err := h.connectionService.GetConnection(connectionID)
	if err != nil {
		response.Error(w, "Connection updated but failed to retrieve: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Convert to response format
	resp := connectionToResponse(*connection)
	
	response.JSON(w, resp, http.StatusOK)
}

// DeleteConnection handles DELETE /api/v1/connections/{connectionID}
func (h *ConnectionHandler) DeleteConnection(w http.ResponseWriter, r *http.Request) {
	connectionID := chi.URLParam(r, "connectionID")
	
	// Delete connection
	err := h.connectionService.DeleteConnection(connectionID)
	if err != nil {
		response.Error(w, "Failed to delete connection: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, map[string]string{
		"message": "Connection deleted successfully",
	}, http.StatusOK)
}

// ExecuteConnection handles POST /api/v1/connections/{connectionID}/execute
func (h *ConnectionHandler) ExecuteConnection(w http.ResponseWriter, r *http.Request) {
	connectionID := chi.URLParam(r, "connectionID")
	
	// Execute connection
	executionID, err := h.connectionService.ExecuteConnection(r.Context(), connectionID)
	if err != nil {
		response.Error(w, "Failed to execute connection: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, dto.ExecuteConnectionResponse{
		ExecutionID: executionID,
		Message:     "Connection execution started",
	}, http.StatusAccepted)
}

// SetActive handles PATCH /api/v1/connections/{connectionID}/active
func (h *ConnectionHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	connectionID := chi.URLParam(r, "connectionID")
	
	var req dto.SetActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Set active status
	err := h.connectionService.SetActive(r.Context(), connectionID, req.Active)
	if err != nil {
		response.Error(w, "Failed to set active status: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get updated connection
	connection, err := h.connectionService.GetConnection(connectionID)
	if err != nil {
		response.Error(w, "Status updated but failed to retrieve connection: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Convert to response format
	resp := connectionToResponse(*connection)
	
	response.JSON(w, resp, http.StatusOK)
}

// GetExecutions handles GET /api/v1/connections/{connectionID}/executions
func (h *ConnectionHandler) GetExecutions(w http.ResponseWriter, r *http.Request) {
	connectionID := chi.URLParam(r, "connectionID")
	
	// Parse pagination parameters
	page := 1
	pageSize := 20
	
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}
	
	// Calculate offset
	offset := (page - 1) * pageSize
	
	// Get executions from service
	executions, err := h.connectionService.GetExecutions(connectionID, pageSize, offset)
	if err != nil {
		response.Error(w, "Failed to get executions: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Convert to response format
	executionResponses := make([]dto.ExecutionResponse, len(executions))
	for i, exec := range executions {
		var sourceData, targetData interface{}
		
		// Parse source data if available
		if exec.SourceData != nil {
			json.Unmarshal(exec.SourceData, &sourceData)
		}
		
		// Parse target data if available
		if exec.TargetData != nil {
			json.Unmarshal(exec.TargetData, &targetData)
		}
		
		// Create response object
		resp := dto.ExecutionResponse{
			ID:           exec.ID,
			ConnectionID: exec.ConnectionID,
			Status:       exec.Status,
			StartTime:    exec.StartTime,
			SourceData:   sourceData,
			TargetData:   targetData,
			CreatedAt:    exec.CreatedAt,
		}
		
		// Add end time if available
		if exec.EndTime.Valid {
			endTime := exec.EndTime.Time
			resp.EndTime = &endTime
		}
		
		// Add error if available
		if exec.Error.Valid {
			resp.Error = exec.Error.String
		}
		
		executionResponses[i] = resp
	}
	
	// Use pagination response
	response.Paginated(w, dto.ExecutionListResponse{
		Executions: executionResponses,
		Total:      len(executionResponses), // This should ideally be a count from the database
		Page:       page,
		PageSize:   pageSize,
	}, page, pageSize, len(executionResponses)) // Total should come from the database
}

// Helper function to convert domain object to response
func connectionToResponse(conn services.ConnectionInfo) dto.ConnectionResponse {
	return dto.ConnectionResponse{
		ID:            conn.ID,
		Name:          conn.Name,
		Description:   conn.Description,
		SourceID:      conn.SourceID,
		TargetID:      conn.TargetID,
		TransformerID: conn.TransformerID,
		Query:         conn.Query,
		Schedule:      conn.Schedule,
		Active:        conn.Active,
		LastRun:       conn.LastRun,
		Status:        conn.Status,
		CreatedAt:     conn.CreatedAt,
		UpdatedAt:     conn.UpdatedAt,
	}
}
