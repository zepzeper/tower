package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zepzeper/tower/internal/response"
	"github.com/zepzeper/tower/internal/services/connection"
	"github.com/zepzeper/tower/internal/webapi/dto"
)

type ConnectionHandler struct {
	connectionService *connection.Service
}

func NewConnectionHandler(connectionService *connection.Service) *ConnectionHandler {
	return &ConnectionHandler{
		connectionService: connectionService,
	}
}

func (h *ConnectionHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	// h.connectionService.FetchConnections()

}

func (h *ConnectionHandler) Test(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var req dto.ApiConnectionCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Type == "" || len(req.Configs) == 0 {
		response.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Call the service to test the connection
	conn, err := h.connectionService.TestConnection(req)
	if err != nil {
		response.Error(w, "Failed to test connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the response data
	data := map[string]interface{}{
		"success": true,
		"message": "Connection test successful",
		"connection": map[string]interface{}{
			"id":          conn.ID,
			"name":        conn.Name,
			"description": conn.Description,
			"type":        conn.Type,
			"base_url":    conn.BaseURL,
			"active":      conn.Active,
		},
	}

	response.JSON(w, data, http.StatusOK)
}
