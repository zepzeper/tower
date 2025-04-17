package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	
	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/api/dto"
	"github.com/zepzeper/tower/internal/api/response"
	"github.com/zepzeper/tower/internal/core/transformers"
	"github.com/zepzeper/tower/internal/services"
)

// TransformerHandler handles transformer-related API endpoints
type TransformerHandler struct {
	transformerService *services.TransformerService
}

// NewTransformerHandler creates a new transformer handler
func NewTransformerHandler(transformerService *services.TransformerService) *TransformerHandler {
	return &TransformerHandler{
		transformerService: transformerService,
	}
}

// ListTransformers handles GET /api/v1/transformers
func (h *TransformerHandler) ListTransformers(w http.ResponseWriter, r *http.Request) {
	// Get transformers from service
	transformers, err := h.transformerService.ListTransformers()
	if err != nil {
		response.Error(w, "Failed to list transformers: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Convert to response format
	transformerResponses := make([]dto.TransformerResponse, len(transformers))
	for i, t := range transformers {
		transformerResponses[i] = transformerToResponse(&t)
	}
	
	response.JSON(w, dto.TransformerListResponse{
		Transformers: transformerResponses,
	}, http.StatusOK)
}

// GetTransformer handles GET /api/v1/transformers/{transformerID}
func (h *TransformerHandler) GetTransformer(w http.ResponseWriter, r *http.Request) {
	transformerID := chi.URLParam(r, "transformerID")
	
	// Get transformer from service
	transformer, err := h.transformerService.GetTransformer(transformerID)
	if err != nil {
		response.Error(w, "Failed to get transformer: "+err.Error(), http.StatusNotFound)
		return
	}
	
	// Convert to response format
	resp := transformerToResponse(transformer)
	
	response.JSON(w, resp, http.StatusOK)
}

// CreateTransformer handles POST /api/v1/transformers
func (h *TransformerHandler) CreateTransformer(w http.ResponseWriter, r *http.Request) {
	var req dto.TransformerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.Name == "" {
		response.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	// Convert mappings to service format
	mappings := make([]transformers.FieldMapping, len(req.Mappings))
	for i, mapping := range req.Mappings {
		mappings[i] = transformers.FieldMapping{
			SourceField: mapping.SourceField,
			TargetField: mapping.TargetField,
		}
	}
	
	// Create transformer
	id, err := h.transformerService.CreateTransformer(req.Name, req.Description, mappings)
	if err != nil {
		response.Error(w, "Failed to create transformer: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get created transformer
	transformer, err := h.transformerService.GetTransformer(id)
	if err != nil {
		response.Error(w, "Transformer created but failed to retrieve: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Convert to response format
	resp := transformerToResponse(transformer)
	
	response.JSON(w, resp, http.StatusCreated)
}

// UpdateTransformer handles PUT /api/v1/transformers/{transformerID}
func (h *TransformerHandler) UpdateTransformer(w http.ResponseWriter, r *http.Request) {
	transformerID := chi.URLParam(r, "transformerID")
	
	var req dto.TransformerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.Name == "" {
		response.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	// Convert mappings to service format
	mappings := make([]transformers.FieldMapping, len(req.Mappings))
	for i, mapping := range req.Mappings {
		mappings[i] = transformers.FieldMapping{
			SourceField: mapping.SourceField,
			TargetField: mapping.TargetField,
		}
	}
	
	// Update transformer
	err := h.transformerService.UpdateTransformer(transformerID, req.Name, req.Description, mappings)
	if err != nil {
		response.Error(w, "Failed to update transformer: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get updated transformer
	transformer, err := h.transformerService.GetTransformer(transformerID)
	if err != nil {
		response.Error(w, "Transformer updated but failed to retrieve: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Convert to response format
	resp := transformerToResponse(transformer)
	
	response.JSON(w, resp, http.StatusOK)
}

// DeleteTransformer handles DELETE /api/v1/transformers/{transformerID}
func (h *TransformerHandler) DeleteTransformer(w http.ResponseWriter, r *http.Request) {
	transformerID := chi.URLParam(r, "transformerID")
	
	// Delete transformer
	err := h.transformerService.DeleteTransformer(transformerID)
	if err != nil {
		response.Error(w, "Failed to delete transformer: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, map[string]string{
		"message": "Transformer deleted successfully",
	}, http.StatusOK)
}

// GenerateTransformer handles POST /api/v1/transformers/generate
func (h *TransformerHandler) GenerateTransformer(w http.ResponseWriter, r *http.Request) {
	var req dto.GenerateTransformerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.SourceID == "" || req.TargetID == "" {
		response.Error(w, "Source and target IDs are required", http.StatusBadRequest)
		return
	}
	
	// Generate transformer
	id, err := h.transformerService.GenerateTransformer(req.SourceID, req.TargetID, req.Name, req.Description)
	if err != nil {
		response.Error(w, "Failed to generate transformer: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get created transformer
	transformer, err := h.transformerService.GetTransformer(id)
	if err != nil {
		response.Error(w, "Transformer generated but failed to retrieve: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Convert to response format
	resp := transformerToResponse(transformer)
	
	response.JSON(w, resp, http.StatusCreated)
}

// TransformData handles POST /api/v1/transformers/transform
func (h *TransformerHandler) TransformData(w http.ResponseWriter, r *http.Request) {
	var req dto.TransformDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Transform data
	result, err := h.transformerService.TransformData(r.Context(), req.TransformerID, req.Data)
	if err != nil {
		response.Error(w, "Failed to transform data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, dto.TransformDataResponse{
		Result: result,
	}, http.StatusOK)
}

// Helper function to convert transformer to response
func transformerToResponse(transformer *transformers.Transformer) dto.TransformerResponse {
	mappings := make([]dto.FieldMapping, len(transformer.Mappings))
	for i, mapping := range transformer.Mappings {
		mappings[i] = dto.FieldMapping{
			SourceField: mapping.SourceField,
			TargetField: mapping.TargetField,
		}
	}
	
	functions := make([]dto.Function, len(transformer.Functions))
	for i, fn := range transformer.Functions {
		functions[i] = dto.Function{
			Name:        fn.Name,
			TargetField: fn.TargetField,
			Args:        fn.Args,
		}
	}
	
	return dto.TransformerResponse{
		ID:          transformer.ID,
		Name:        transformer.Name,
		Description: transformer.Description,
		Mappings:    mappings,
		Functions:   functions,
		CreatedAt:   transformer.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   transformer.UpdatedAt.Format(time.RFC3339),
	}
}
