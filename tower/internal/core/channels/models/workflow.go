package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Workflow represents a data processing workflow in the database
type Workflow struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description sql.NullString  `json:"description"`
	Triggers    json.RawMessage `json:"triggers"`
	Actions     json.RawMessage `json:"actions"`
	Active      bool            `json:"active"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// ToAPIWorkflow converts a database Workflow to an API response Workflow
func (w *Workflow) ToAPIWorkflow() interface{} {
	var triggers []interface{}
	var actions []interface{}
	
	json.Unmarshal(w.Triggers, &triggers)
	json.Unmarshal(w.Actions, &actions)
	
	return map[string]interface{}{
		"id":          w.ID,
		"name":        w.Name,
		"description": w.Description.String,
		"triggers":    triggers,
		"actions":     actions,
		"active":      w.Active,
		"createdAt":   w.CreatedAt,
		"updatedAt":   w.UpdatedAt,
	}
}
