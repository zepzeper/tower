package dto

import "time"

// SummaryResponse represents dashboard summary data
type SummaryResponse struct {
	TotalConnections    int `json:"totalConnections"`
	ActiveConnections   int `json:"activeConnections"`
	SuccessfulTransfers int `json:"successfulTransfers"`
	FailedTransfers     int `json:"failedTransfers"`
	PendingTransfers    int `json:"pendingTransfers"`
}

// RecentActivityResponse represents recent activity data
type RecentActivityResponse struct {
	Executions []ExecutionSummary `json:"executions"`
}

// ExecutionSummary represents summary data for an execution
type ExecutionSummary struct {
	ID              string     `json:"id"`
	ConnectionID    string     `json:"connectionId"`
	ConnectionName  string     `json:"connectionName"`
	Status          string     `json:"status"`
	StartTime       time.Time  `json:"startTime"`
	EndTime         *time.Time `json:"endTime,omitempty"`
}
