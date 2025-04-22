package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zepzeper/tower/internal/response"
	"github.com/zepzeper/tower/internal/services/relations"
	"github.com/zepzeper/tower/internal/webapi/dto"
)

type RelationsHandler struct {
	relationsService *relations.Service
}

func NewRelationsHandler(relationsService *relations.Service) *RelationsHandler {
	return &RelationsHandler{
		relationsService: relationsService,
	}
}

func (h *RelationsHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	connectionID := getInitiatorIdFromUrl(r)

	json, err := h.relationsService.Fetch(connectionID)
	if err != nil {
		response.Error(w, "Failed to fetch relations: "+err.Error(), http.StatusInternalServerError)
		return
	}
	response.JSON(w, json, http.StatusOK)
}

func (h *RelationsHandler) FetchLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	connectionID := getInitiatorIdFromUrl(r)

	json, err := h.relationsService.FetchLogs(connectionID)
	if err != nil {
		response.Error(w, "Failed to fetch relations: "+err.Error(), http.StatusInternalServerError)
		return
	}
	response.JSON(w, json, http.StatusOK)
}

func (h *RelationsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var req dto.CredentialsRelationCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Call the service to create the connection
	err := h.relationsService.Create(req)
	if err != nil {
		response.Error(w, "Failed to create connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func getInitiatorIdFromUrl(r *http.Request) string {

	// Extract connectionID from URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")

	// Find the index of "connections" in the path
	connectionIndex := -1
	for i, part := range parts {
		if part == "connections" && i+1 < len(parts) {
			connectionIndex = i + 1
			break
		}
	}

	if connectionIndex == -1 || connectionIndex >= len(parts)-1 {
		return "0"
	}

	connectionID := parts[connectionIndex]

	return connectionID
}
