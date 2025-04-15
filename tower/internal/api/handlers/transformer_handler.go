package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/zepzeper/tower/internal/api/response"
	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/database/models"
)

// transformerRequest represents the data expected in transformer API requests
type transformerRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Mappings    []json.RawMessage `json:"mappings"`
	Functions   []json.RawMessage `json:"functions"`
}

// TransformerHandler handles transformer-related API endpoints
type TransformerHandler struct {
	db *db.Manager
}

// NewTransformerHandler creates a new transformer handler
func NewTransformerHandler(dbManager *db.Manager) *TransformerHandler {
	return &TransformerHandler{
		db: dbManager,
	}
}

// ListTransformers handles GET /api/v1/transformers
func (h *TransformerHandler) ListTransformers(w http.ResponseWriter, r *http.Request) {
	// Check for name search query
	nameQuery := r.URL.Query().Get("name")
	
	var transformers []models.Transformer
	var err error
	
	if nameQuery != "" {
		transformers, err = h.db.Repos.Transformer().FindByName(nameQuery)
	} else {
		transformers, err = h.db.Repos.Transformer().GetAll()
	}
	
	if err != nil {
		response.Error(w, "Failed to retrieve transformers", http.StatusInternalServerError)
		return
	}
	
	// Convert to API response format
	result := make([]interface{}, len(transformers))
	for i, transformer := range transformers {
		result[i] = transformer.ToAPITransformer()
	}
	
	response.JSON(w, result, http.StatusOK)
}

// GetTransformer handles GET /api/v1/transformers/{transformerID}
func (h *TransformerHandler) GetTransformer(w http.ResponseWriter, r *http.Request) {
	transformerID := chi.URLParam(r, "transformerID")
	
	transformer, err := h.db.Repos.Transformer().GetByID(transformerID)
	if err != nil {
		response.Error(w, "Transformer not found", http.StatusNotFound)
		return
	}
	
	response.JSON(w, transformer.ToAPITransformer(), http.StatusOK)
}

// CreateTransformer handles POST /api/v1/transformers
func (h *TransformerHandler) CreateTransformer(w http.ResponseWriter, r *http.Request) {
	var req transformerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.Name == "" {
		response.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	// Convert mappings and functions to JSON
	mappingsJSON, err := json.Marshal(req.Mappings)
	if err != nil {
		response.Error(w, "Invalid mappings format", http.StatusBadRequest)
		return
	}
	
	functionsJSON, err := json.Marshal(req.Functions)
	if err != nil {
		response.Error(w, "Invalid functions format", http.StatusBadRequest)
		return
	}
	
	// Create transformer model
	transformer := models.Transformer{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Mappings:  mappingsJSON,
		Functions: functionsJSON,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Set description if provided
	if req.Description != "" {
		transformer.Description.String = req.Description
		transformer.Description.Valid = true
	}
	
	// Save to database
	if err := h.db.Repos.Transformer().Create(transformer); err != nil {
		response.Error(w, "Failed to create transformer", http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, transformer.ToAPITransformer(), http.StatusCreated)
}

// UpdateTransformer handles PUT /api/v1/transformers/{transformerID}
func (h *TransformerHandler) UpdateTransformer(w http.ResponseWriter, r *http.Request) {
	transformerID := chi.URLParam(r, "transformerID")
	
	// Check if transformer exists
	transformer, err := h.db.Repos.Transformer().GetByID(transformerID)
	if err != nil {
		response.Error(w, "Transformer not found", http.StatusNotFound)
		return
	}
	
	// Parse request body
	var req transformerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.Name == "" {
		response.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	// Convert mappings and functions to JSON
	if len(req.Mappings) > 0 {
		mappingsJSON, err := json.Marshal(req.Mappings)
		if err != nil {
			response.Error(w, "Invalid mappings format", http.StatusBadRequest)
			return
		}
		transformer.Mappings = mappingsJSON
	}
	
	if len(req.Functions) > 0 {
		functionsJSON, err := json.Marshal(req.Functions)
		if err != nil {
			response.Error(w, "Invalid functions format", http.StatusBadRequest)
			return
		}
		transformer.Functions = functionsJSON
	}
	
	// Update transformer fields
	transformer.Name = req.Name
	transformer.UpdatedAt = time.Now()
	
	// Update description if provided
	if req.Description != "" {
		transformer.Description.String = req.Description
		transformer.Description.Valid = true
	} else {
		transformer.Description.Valid = false
	}
	
	// Save to database
	if err := h.db.Repos.Transformer().Update(transformer); err != nil {
		response.Error(w, "Failed to update transformer", http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, transformer.ToAPITransformer(), http.StatusOK)
}

// DeleteTransformer handles DELETE /api/v1/transformers/{transformerID}
func (h *TransformerHandler) DeleteTransformer(w http.ResponseWriter, r *http.Request) {
	transformerID := chi.URLParam(r, "transformerID")
	
	// Delete from database
	if err := h.db.Repos.Transformer().Delete(transformerID); err != nil {
		response.Error(w, "Failed to delete transformer", http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, map[string]interface{}{
		"success": true,
		"message": "Transformer deleted successfully",
	}, http.StatusOK)
}

// TestTransformer handles POST /api/v1/transformers/{transformerID}/test
func (h *TransformerHandler) TestTransformer(w http.ResponseWriter, r *http.Request) {
	transformerID := chi.URLParam(r, "transformerID")
	
	// Get transformer
	transformer, err := h.db.Repos.Transformer().GetByID(transformerID)
	if err != nil {
		response.Error(w, "Transformer not found", http.StatusNotFound)
		return
	}
	
	// Parse input data from request
	var inputData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		response.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}
	
	// In a real implementation, you would:
	// 1. Apply the transformer mappings to the input data
	// 2. Apply the transformer functions to the intermediate data
	// 3. Return the transformed data
	
	// For now, just return example output
	response.JSON(w, map[string]interface{}{
		"success": true,
		"input":   inputData,
		"output": map[string]interface{}{
			"transformedData": "This is a placeholder for the transformed data",
			"transformerId":   transformer.ID,
			"timestamp":       time.Now(),
		},
	}, http.StatusOK)
}
