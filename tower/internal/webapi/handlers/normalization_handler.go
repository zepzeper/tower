package handlers

import (
	"encoding/json"
	"net/http"
	
	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/core/connectors"
	"github.com/zepzeper/tower/internal/response"
	"github.com/zepzeper/tower/internal/services"
)

// NormalizationHandler handles API endpoints for data normalization
type NormalizationHandler struct {
	normalizationService *services.NormalizationService
	integrationService  *services.IntegrationService
}

// NewNormalizationHandler creates a new normalization handler
func NewNormalizationHandler(
	normalizationService *services.NormalizationService,
	integrationService *services.IntegrationService,
) *NormalizationHandler {
	return &NormalizationHandler{
		normalizationService: normalizationService,
		integrationService:  integrationService,
	}
}

// RegisterRoutes registers the handler's routes
func (h *NormalizationHandler) RegisterRoutes(router chi.Router) {
	router.Route("/normalize", func(r chi.Router) {
		r.Post("/{sourceID}/to/{targetID}", h.NormalizeData)
		r.Post("/discover-schema/{connectorID}", h.DiscoverSchema)
		r.Post("/auto-configure", h.AutoConfigureIntegration)
	})
}

// NormalizeData normalizes data from a source to a target format
func (h *NormalizationHandler) NormalizeData(w http.ResponseWriter, r *http.Request) {
	// Extract path parameters
	sourceID := chi.URLParam(r, "sourceID")
	targetID := chi.URLParam(r, "targetID")
	
	// Parse request body
	var sourceData []connectors.DataPayload
	err := json.NewDecoder(r.Body).Decode(&sourceData)
	if err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Normalize the data
	normalizedData, err := h.normalizationService.NormalizeForTarget(
		r.Context(),
		sourceID,
		targetID,
		sourceData,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return the normalized data
	response.JSON(w, normalizedData, http.StatusOK)
}

// DiscoverSchema discovers the schema for a connector
func (h *NormalizationHandler) DiscoverSchema(w http.ResponseWriter, r *http.Request) {
	// Extract path parameter
	connectorID := chi.URLParam(r, "connectorID")
	
	// Parse request body
	var sampleData []connectors.DataPayload
	err := json.NewDecoder(r.Body).Decode(&sampleData)
	if err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Discover the schema
	schema, err := h.normalizationService.DetectSchema(
		r.Context(),
		connectorID,
		sampleData,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return the schema
	response.JSON(w, schema, http.StatusOK)
}

// AutoConfigureRequest is the request body for auto-configuration
type AutoConfigureRequest struct {
	SourceID    string `json:"sourceId"`
	TargetID    string `json:"targetId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Schedule    string `json:"schedule"`
}

// AutoConfigureIntegration automatically configures an integration between systems
func (h *NormalizationHandler) AutoConfigureIntegration(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req AutoConfigureRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate required fields
	if req.SourceID == "" || req.TargetID == "" || req.Name == "" {
		response.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	
	// Auto-configure the integration
	connectionID, err := h.integrationService.AutoConfigureIntegration(
		r.Context(),
		req.SourceID,
		req.TargetID,
		req.Name,
		req.Description,
		req.Schedule,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return the connection ID
	response.JSON(w, map[string]string{"connectionId": connectionID}, http.StatusOK)
}
