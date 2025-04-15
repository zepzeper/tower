package models

import (
	"database/sql"
	"time"
)

// User represents a user of the system in the database
type User struct {
	ID        string       `json:"id"`
	Email     string       `json:"email"`
	Password  string       `json:"-"` // Never expose in JSON
	Name      string       `json:"name"`
	Role      string       `json:"role"` // "admin", "user", etc.
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	LastLogin sql.NullTime `json:"lastLogin"`
}

// ToAPIUser converts a database User to an API response User
func (u *User) ToAPIUser() interface{} {
	user := map[string]interface{}{
		"id":        u.ID,
		"email":     u.Email,
		"name":      u.Name,
		"role":      u.Role,
		"createdAt": u.CreatedAt,
		"updatedAt": u.UpdatedAt,
	}
	
	if u.LastLogin.Valid {
		user["lastLogin"] = u.LastLogin.Time
	}
	
	return user
}
