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

// workflowRequest represents the data expected in workflow API requests
type workflowRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Triggers    []json.RawMessage `json:"triggers"`
	Actions     []json.RawMessage `json:"actions"`
	Active      bool              `json:"active"`
}

// WorkflowHandler handles workflow-related API endpoints
type WorkflowHandler struct {
	db *db.Manager
}

// NewWorkflowHandler creates a new workflow handler
func NewWorkflowHandler(dbManager *db.Manager) *WorkflowHandler {
	return &WorkflowHandler{
		db: dbManager,
	}
}

// ListWorkflows handles GET /api/v1/workflows
func (h *WorkflowHandler) ListWorkflows(w http.ResponseWriter, r *http.Request) {
	// Check for active filter
	activeOnly := r.URL.Query().Get("active") == "true"
	
	var workflows []models.Workflow
	var err error
	
	if activeOnly {
		workflows, err = h.db.Repos.Workflow().GetActive()
	} else {
		workflows, err = h.db.Repos.Workflow().GetAll()
	}
	
	if err != nil {
		response.Error(w, "Failed to retrieve workflows", http.StatusInternalServerError)
		return
	}
	
	// Convert to API response format
	result := make([]interface{}, len(workflows))
	for i, workflow := range workflows {
		result[i] = workflow.ToAPIWorkflow()
	}
	
	response.JSON(w, result, http.StatusOK)
}

// GetWorkflow handles GET /api/v1/workflows/{workflowID}
func (h *WorkflowHandler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowID := chi.URLParam(r, "workflowID")
	
	workflow, err := h.db.Repos.Workflow().GetByID(workflowID)
	if err != nil {
		response.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}
	
	response.JSON(w, workflow.ToAPIWorkflow(), http.StatusOK)
}

// CreateWorkflow handles POST /api/v1/workflows
func (h *WorkflowHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var req workflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.Name == "" {
		response.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	// Validate triggers and actions
	if len(req.Triggers) == 0 {
		response.Error(w, "At least one trigger is required", http.StatusBadRequest)
		return
	}
	
	if len(req.Actions) == 0 {
		response.Error(w, "At least one action is required", http.StatusBadRequest)
		return
	}
	
	// Convert triggers and actions to JSON
	triggersJSON, err := json.Marshal(req.Triggers)
	if err != nil {
		response.Error(w, "Invalid triggers format", http.StatusBadRequest)
		return
	}
	
	actionsJSON, err := json.Marshal(req.Actions)
	if err != nil {
		response.Error(w, "Invalid actions format", http.StatusBadRequest)
		return
	}
	
	// Create workflow model
	workflow := models.Workflow{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Triggers:  triggersJSON,
		Actions:   actionsJSON,
		Active:    req.Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Set description if provided
	if req.Description != "" {
		workflow.Description.String = req.Description
		workflow.Description.Valid = true
	}
	
	// Save to database
	if err := h.db.Repos.Workflow().Create(workflow); err != nil {
		response.Error(w, "Failed to create workflow", http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, workflow.ToAPIWorkflow(), http.StatusCreated)
}

// UpdateWorkflow handles PUT /api/v1/workflows/{workflowID}
func (h *WorkflowHandler) UpdateWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowID := chi.URLParam(r, "workflowID")
	
	// Check if workflow exists
	workflow, err := h.db.Repos.Workflow().GetByID(workflowID)
	if err != nil {
		response.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}
	
	// Parse request body
	var req workflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.Name == "" {
		response.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	// Validate triggers and actions
	if len(req.Triggers) == 0 {
		response.Error(w, "At least one trigger is required", http.StatusBadRequest)
		return
	}
	
	if len(req.Actions) == 0 {
		response.Error(w, "At least one action is required", http.StatusBadRequest)
		return
	}
	
	// Convert triggers and actions to JSON
	triggersJSON, err := json.Marshal(req.Triggers)
	if err != nil {
		response.Error(w, "Invalid triggers format", http.StatusBadRequest)
		return
	}
	
	actionsJSON, err := json.Marshal(req.Actions)
	if err != nil {
		response.Error(w, "Invalid actions format", http.StatusBadRequest)
		return
	}
	
	// Update workflow fields
	workflow.Name = req.Name
	workflow.Triggers = triggersJSON
	workflow.Actions = actionsJSON
	workflow.Active = req.Active
	workflow.UpdatedAt = time.Now()
	
	// Update description if provided
	if req.Description != "" {
		workflow.Description.String = req.Description
		workflow.Description.Valid = true
	} else {
		workflow.Description.Valid = false
	}
	
	// Save to database
	if err := h.db.Repos.Workflow().Update(workflow); err != nil {
		response.Error(w, "Failed to update workflow", http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, workflow.ToAPIWorkflow(), http.StatusOK)
}

// DeleteWorkflow handles DELETE /api/v1/workflows/{workflowID}
func (h *WorkflowHandler) DeleteWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowID := chi.URLParam(r, "workflowID")
	
	// Delete from database
	if err := h.db.Repos.Workflow().Delete(workflowID); err != nil {
		response.Error(w, "Failed to delete workflow", http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, map[string]interface{}{
		"success": true,
		"message": "Workflow deleted successfully",
	}, http.StatusOK)
}

// ToggleWorkflow handles PATCH /api/v1/workflows/{workflowID}/toggle
func (h *WorkflowHandler) ToggleWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowID := chi.URLParam(r, "workflowID")
	
	// Get current workflow
	workflow, err := h.db.Repos.Workflow().GetByID(workflowID)
	if err != nil {
		response.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}
	
	// Toggle active status
	newStatus := !workflow.Active
	
	// Update in database
	if err := h.db.Repos.Workflow().SetActive(workflowID, newStatus); err != nil {
		response.Error(w, "Failed to update workflow status", http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, map[string]interface{}{
		"success": true,
		"active":  newStatus,
		"message": "Workflow status updated successfully",
	}, http.StatusOK)
}

// ExecuteWorkflow handles POST /api/v1/workflows/{workflowID}/execute
func (h *WorkflowHandler) ExecuteWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowID := chi.URLParam(r, "workflowID")
	
	// Get workflow
	workflow, err := h.db.Repos.Workflow().GetByID(workflowID)
	if err != nil {
		response.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}
	
	// Check if workflow is active
	if !workflow.Active {
		response.Error(w, "Cannot execute inactive workflow", http.StatusBadRequest)
		return
	}
	
	// In a real implementation, you would:
	// 1. Create an execution record
	// 2. Start the workflow processing
	// 3. Update the execution record with results
	
	// For now, just create an execution record
	execution := models.Execution{
		ID:         uuid.New().String(),
		WorkflowID: workflowID,
		Status:     "in_progress",
		StartTime:  time.Now(),
		CreatedAt:  time.Now(),
	}
	
	// Save to database
	if err := h.db.Repos.Execution().Create(execution); err != nil {
		response.Error(w, "Failed to create execution record", http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, map[string]interface{}{
		"success":     true,
		"executionId": execution.ID,
		"message":     "Workflow execution started",
	}, http.StatusAccepted)
}
