package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Transformer represents a data transformation definition in the database
type Transformer struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description sql.NullString  `json:"description"`
	Mappings    json.RawMessage `json:"mappings"`
	Functions   json.RawMessage `json:"functions"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// ToAPITransformer converts a database Transformer to an API response Transformer
func (t *Transformer) ToAPITransformer() interface{} {
	var mappings []interface{}
	var functions []interface{}
	
	json.Unmarshal(t.Mappings, &mappings)
	json.Unmarshal(t.Functions, &functions)
	
	return map[string]interface{}{
		"id":          t.ID,
		"name":        t.Name,
		"description": t.Description.String,
		"mappings":    mappings,
		"functions":   functions,
		"createdAt":   t.CreatedAt,
		"updatedAt":   t.UpdatedAt,
	}
}
