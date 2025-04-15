package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Execution represents a record of a workflow execution in the database
type Execution struct {
	ID         string          `json:"id"`
	WorkflowID string          `json:"workflowId"`
	Status     string          `json:"status"` // "success", "failed", "in_progress"
	StartTime  time.Time       `json:"startTime"`
	EndTime    sql.NullTime    `json:"endTime"`
	Result     json.RawMessage `json:"result"`
	Error      sql.NullString  `json:"error"`
	CreatedAt  time.Time       `json:"createdAt"`
}

// ToAPIExecution converts a database Execution to an API response Execution
func (e *Execution) ToAPIExecution() interface{} {
	var result interface{}
	json.Unmarshal(e.Result, &result)
	
	execution := map[string]interface{}{
		"id":         e.ID,
		"workflowId": e.WorkflowID,
		"status":     e.Status,
		"startTime":  e.StartTime,
		"result":     result,
		"createdAt":  e.CreatedAt,
	}
	
	if e.EndTime.Valid {
		execution["endTime"] = e.EndTime.Time
	}
	
	if e.Error.Valid {
		execution["error"] = e.Error.String
	}
	
	return execution
}
