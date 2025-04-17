package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Connection represents a data integration connection
type Connection struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Description   sql.NullString  `json:"description"`
	SourceID      string          `json:"sourceId"`
	TargetID      string          `json:"targetId"`
	TransformerID string          `json:"transformerId"`
	Config        json.RawMessage `json:"config"`
	Schedule      string          `json:"schedule"`
	Active        bool            `json:"active"`
	LastRun       sql.NullTime    `json:"lastRun"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}
