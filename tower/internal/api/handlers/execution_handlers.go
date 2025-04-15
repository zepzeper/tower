package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/api/response"
	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/database/models"
)

// ExecutionHandler handles execution-related API endpoints
type ExecutionHandler struct {
	db *db.Manager
}

// NewExecutionHandler creates a new execution handler
func NewExecutionHandler(dbManager *db.Manager) *ExecutionHandler {
	return &ExecutionHandler{
		db: dbManager,
	}
}

// ListExecutions handles GET /api/v1/executions
func (h *ExecutionHandler) ListExecutions(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	workflowID := r.URL.Query().Get("workflow_id")
	
	limit := 20 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	offset := 0 // Default offset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	
	var executions []models.Execution
	var err error
	
	// Check if filtering by workflow ID
	if workflowID != "" {
		executions, err = h.db.Repos.Execution().GetByWorkflowID(workflowID, limit, offset)
	} else {
		executions, err = h.db.Repos.Execution().GetRecentExecutions(limit, offset)
	}
	
	if err != nil {
		response.Error(w, "Failed to retrieve executions", http.StatusInternalServerError)
		return
	}
	
	// Convert to API response format
	result := make([]interface{}, len(executions))
	for i, execution := range executions {
		result[i] = execution.ToAPIExecution()
	}
	
	// In a real implementation, you would get the total count for pagination
	response.JSON(w, result, http.StatusOK)
}

// GetExecution handles GET /api/v1/executions/{executionID}
func (h *ExecutionHandler) GetExecution(w http.ResponseWriter, r *http.Request) {
	executionID := chi.URLParam(r, "executionID")
	
	execution, err := h.db.Repos.Execution().GetByID(executionID)
	if err != nil {
		response.Error(w, "Execution not found", http.StatusNotFound)
		return
	}
	
	response.JSON(w, execution.ToAPIExecution(), http.StatusOK)
}

// GetExecutionStats handles GET /api/v1/executions/stats
func (h *ExecutionHandler) GetExecutionStats(w http.ResponseWriter, r *http.Request) {
	workflowID := r.URL.Query().Get("workflow_id")
	
	if workflowID == "" {
		response.Error(w, "workflow_id parameter is required", http.StatusBadRequest)
		return
	}
	
	// Get execution counts by status
	counts, err := h.db.Repos.Execution().GetExecutionCounts(workflowID)
	if err != nil {
		response.Error(w, "Failed to retrieve execution statistics", http.StatusInternalServerError)
		return
	}
	
	// Ensure all statuses are represented
	if _, ok := counts["success"]; !ok {
		counts["success"] = 0
	}
	
	if _, ok := counts["failed"]; !ok {
		counts["failed"] = 0
	}
	
	if _, ok := counts["in_progress"]; !ok {
		counts["in_progress"] = 0
	}
	
	// Calculate total
	total := counts["success"] + counts["failed"] + counts["in_progress"]
	
	response.JSON(w, map[string]interface{}{
		"workflow_id": workflowID,
		"total":       total,
		"by_status":   counts,
	}, http.StatusOK)
}

// CancelExecution handles POST /api/v1/executions/{executionID}/cancel
func (h *ExecutionHandler) CancelExecution(w http.ResponseWriter, r *http.Request) {
	executionID := chi.URLParam(r, "executionID")
	
	// Get current execution
	execution, err := h.db.Repos.Execution().GetByID(executionID)
	if err != nil {
		response.Error(w, "Execution not found", http.StatusNotFound)
		return
	}
	
	// Check if execution is already completed
	if execution.Status == "success" || execution.Status == "failed" {
		response.Error(w, "Cannot cancel completed execution", http.StatusBadRequest)
		return
	}
	
	// Create cancellation result
	cancelResult := map[string]interface{}{
		"status":     "cancelled",
		"reason":     "User requested cancellation",
		"cancelled_at": execution.EndTime,
	}
	
	resultJSON, err := json.Marshal(cancelResult)
	if err != nil {
		response.Error(w, "Failed to create cancellation record", http.StatusInternalServerError)
		return
	}
	
	// In a real implementation, you would:
	// 1. Signal the running execution to stop
	// 2. Update the execution record with a cancelled status
	
	// Update the execution status
	if err := h.db.Repos.Execution().UpdateStatus(
		executionID,
		"failed",
		execution.EndTime.Time,
		resultJSON,
		"Execution cancelled by user",
	); err != nil {
		response.Error(w, "Failed to cancel execution", http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, map[string]interface{}{
		"success": true,
		"message": "Execution cancelled successfully",
	}, http.StatusOK)
}
