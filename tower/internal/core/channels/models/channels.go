package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Channel represents an API integration channel in the database
type Channel struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	Description sql.NullString  `json:"description"`
	Config      json.RawMessage `json:"config"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// ToAPIChannel converts a database Channel to an API response Channel
func (c *Channel) ToAPIChannel() interface{} {
	var config map[string]interface{}
	json.Unmarshal(c.Config, &config)
	
	return map[string]interface{}{
		"id":          c.ID,
		"name":        c.Name,
		"type":        c.Type,
		"description": c.Description.String,
		"config":      config,
		"createdAt":   c.CreatedAt,
		"updatedAt":   c.UpdatedAt,
	}
}
