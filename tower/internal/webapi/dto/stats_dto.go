// internal/webapi/dto/stats_dto.go
package dto

import "time"

// ConnectionStatsResponse represents statistics about connections
type ConnectionStatsResponse struct {
	TotalCount            int                     `json:"totalCount"`
	ActiveCount           int                     `json:"activeCount"`
	InactiveCount         int                     `json:"inactiveCount"`
	ConnectorDistribution map[string]int          `json:"connectorDistribution"`
	RecentlyCreated       []ConnectionStatSummary `json:"recentlyCreated"`
}

// ConnectionStatSummary represents summary data for a connection in stats
type ConnectionStatSummary struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	SourceID     string     `json:"sourceId"`
	TargetID     string     `json:"targetId"`
	Active       bool       `json:"active"`
	LastRun      *time.Time `json:"lastRun,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
}

// ExecutionStatsResponse represents statistics about executions
type ExecutionStatsResponse struct {
	TotalCount        int            `json:"totalCount"`
	SuccessCount      int            `json:"successCount"`
	FailureCount      int            `json:"failureCount"`
	PendingCount      int            `json:"pendingCount"`
	TimeDistribution  map[string]int `json:"timeDistribution"`
}
