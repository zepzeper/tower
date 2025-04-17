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
