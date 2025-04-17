package dto

import (
	"time"
)

// ConnectionRequest represents the request for creating/updating a connection
type ConnectionRequest struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	SourceID      string                 `json:"sourceId"`
	TargetID      string                 `json:"targetId"`
	TransformerID string                 `json:"transformerId"`
	Query         map[string]interface{} `json:"query"`
	Schedule      string                 `json:"schedule"`
	Active        bool                   `json:"active"`
}

// ConnectionResponse represents the response for a connection
type ConnectionResponse struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	SourceID      string                 `json:"sourceId"`
	TargetID      string                 `json:"targetId"`
	TransformerID string                 `json:"transformerId"`
	Query         map[string]interface{} `json:"query"`
	Schedule      string                 `json:"schedule"`
	Active        bool                   `json:"active"`
	LastRun       *time.Time             `json:"lastRun,omitempty"`
	Status        string                 `json:"status,omitempty"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

// ConnectionListResponse represents the response for listing connections
type ConnectionListResponse struct {
	Connections []ConnectionResponse `json:"connections"`
}

// ExecuteConnectionResponse represents the response when executing a connection
type ExecuteConnectionResponse struct {
	ExecutionID string `json:"executionId"`
	Message     string `json:"message"`
}

// SetActiveRequest represents the request for toggling a connection
type SetActiveRequest struct {
	Active bool `json:"active"`
}

// ExecutionResponse represents the response for an execution
type ExecutionResponse struct {
	ID           string                 `json:"id"`
	ConnectionID string                 `json:"connectionId"`
	Status       string                 `json:"status"`
	StartTime    time.Time              `json:"startTime"`
	EndTime      *time.Time             `json:"endTime,omitempty"`
	SourceData   interface{}            `json:"sourceData,omitempty"`
	TargetData   interface{}            `json:"targetData,omitempty"`
	Error        string                 `json:"error,omitempty"`
	CreatedAt    time.Time              `json:"createdAt"`
}

// ExecutionListResponse represents the response for listing executions
type ExecutionListResponse struct {
	Executions []ExecutionResponse `json:"executions"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"pageSize"`
}
