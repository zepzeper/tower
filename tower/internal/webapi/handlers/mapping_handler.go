package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zepzeper/tower/internal/response"
	"github.com/zepzeper/tower/internal/services/mapping"
)

type MappingHandler struct {
	mappingService *mapping.Service
}

func NewMappingHandler(mappingService *mapping.Service) *MappingHandler {
	return &MappingHandler{
		mappingService: mappingService,
	}
}

// Generate handles GET /api/mappings/schema?source={source}&target={target}
func (h *MappingHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	source := r.URL.Query().Get("source")
	target := r.URL.Query().Get("target")
	operation := r.URL.Query().Get("operation")

	if source == "" || target == "" || operation == "" {
		response.Error(w, "Missing 'source' or 'target' query parameters", http.StatusBadRequest)
		return
	}

	data, err := h.mappingService.GenerateMapping(source, target, operation)

	if err != nil {
		response.Error(w, "Failed to generate mapping: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.JSON(w, data, http.StatusOK)
}

// handleTestMapping handles POST /api/mappings/test requests
func (h *MappingHandler) HandleTestMapping(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req mapping.TestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
		return
	}

	// Apply transformations based on mappings
	response, err := h.mappingService.ApplyMappings(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to apply mappings: %v", err), http.StatusBadRequest)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
