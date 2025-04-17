package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Execution represents a record of a connection execution
type Execution struct {
	ID           string          `json:"id"`
	ConnectionID string          `json:"connectionId"`
	Status       string          `json:"status"` // "success", "failed", "in_progress"
	StartTime    time.Time       `json:"startTime"`
	EndTime      sql.NullTime    `json:"endTime"`
	SourceData   json.RawMessage `json:"sourceData,omitempty"`
	TargetData   json.RawMessage `json:"targetData,omitempty"`
	Error        sql.NullString  `json:"error"`
	CreatedAt    time.Time       `json:"createdAt"`
}
