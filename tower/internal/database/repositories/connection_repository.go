package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zepzeper/tower/internal/database/models"
)

// APIConnectionRepository handles database operations for API connections
type APIConnectionRepository struct {
	db *sql.DB
}

// NewAPIConnectionRepository creates a new API connection repository
func NewAPIConnectionRepository(db *sql.DB) *APIConnectionRepository {
	return &APIConnectionRepository{
		db: db,
	}
}

// GetAll retrieves all API connections with their configurations
func (r *APIConnectionRepository) GetAll() ([]models.APIConnection, error) {
	// Get all connections
	query := `
		SELECT id, name, description, type, active, created_at, updated_at
		FROM api_connections
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []models.APIConnection
	for rows.Next() {
		var conn models.APIConnection
		if err := rows.Scan(
			&conn.ID,
			&conn.Name,
			&conn.Description,
			&conn.Type,
			&conn.Active,
			&conn.CreatedAt,
			&conn.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// Load config for this connection
		config, err := r.GetConnectionConfig(conn.ID, false)
		if err != nil {
			return nil, err
		}
		conn.Config = config

		connections = append(connections, conn)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return connections, nil
}

// GetByID retrieves an API connection by ID
func (r *APIConnectionRepository) GetByID(id string, includeSecrets bool) (models.APIConnection, error) {
	query := `
		SELECT id, name, description, type, active, created_at, updated_at
		FROM api_connections
		WHERE id = $1
	`

	var conn models.APIConnection
	err := r.db.QueryRow(query, id).Scan(
		&conn.ID,
		&conn.Name,
		&conn.Description,
		&conn.Type,
		&conn.Active,
		&conn.CreatedAt,
		&conn.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return conn, fmt.Errorf("API connection not found: %s", id)
		}
		return conn, err
	}

	// Load config for this connection
	config, err := r.GetConnectionConfig(id, includeSecrets)
	if err != nil {
		return conn, err
	}
	conn.Config = config

	return conn, nil
}

// Create inserts a new API connection with its configuration
func (r *APIConnectionRepository) Create(conn models.APIConnection) error {
	// Start a transaction for atomicity
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert connection
	query := `
		INSERT INTO api_connections (id, name, description, type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	conn.CreatedAt = now
	conn.UpdatedAt = now

	_, err = tx.Exec(
		query,
		conn.ID,
		conn.Name,
		conn.Description,
		conn.Type,
		conn.Active,
		conn.CreatedAt,
		conn.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// Insert configuration
	if len(conn.Config) > 0 {
		for key, value := range conn.Config {
			// Determine if this is a secret (you might want to have a more sophisticated approach)
			isSecret := key == "api_key" || key == "client_secret" || key == "password" || key == "token"

			configQuery := `
				INSERT INTO api_connection_configs (connection_id, key, value, is_secret, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6)
			`

			_, err = tx.Exec(
				configQuery,
				conn.ID,
				key,
				value,
				isSecret,
				now,
				now,
			)

			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// Update updates an existing API connection and its configuration
func (r *APIConnectionRepository) Update(conn models.APIConnection) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update connection
	query := `
		UPDATE api_connections
		SET name = $1, description = $2, type = $3, active = $4, updated_at = $5
		WHERE id = $6
	`

	now := time.Now()
	conn.UpdatedAt = now

	result, err := tx.Exec(
		query,
		conn.Name,
		conn.Description,
		conn.Type,
		conn.Active,
		now,
		conn.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API connection not found: %s", conn.ID)
	}

	// For configuration, we'll use a merge approach:
	// 1. Delete any config items not in the new config
	// 2. Update existing items
	// 3. Insert new items

	if len(conn.Config) > 0 {
		// Get existing config keys
		existingKeys := make(map[string]bool)
		existingQuery := `
			SELECT key FROM api_connection_configs
			WHERE connection_id = $1
		`

		rows, err := tx.Query(existingQuery, conn.ID)
		if err != nil {
			return err
		}

		for rows.Next() {
			var key string
			if err := rows.Scan(&key); err != nil {
				rows.Close()
				return err
			}
			existingKeys[key] = true
		}
		rows.Close()

		// Update or insert config items
		for key, value := range conn.Config {
			isSecret := key == "api_key" || key == "client_secret" || key == "password" || key == "token"

			if existingKeys[key] {
				// Update
				updateQuery := `
					UPDATE api_connection_configs
					SET value = $1, is_secret = $2, updated_at = $3
					WHERE connection_id = $4 AND key = $5
				`

				_, err = tx.Exec(updateQuery, value, isSecret, now, conn.ID, key)
				if err != nil {
					return err
				}

				delete(existingKeys, key)
			} else {
				// Insert
				insertQuery := `
					INSERT INTO api_connection_configs (connection_id, key, value, is_secret, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6)
				`

				_, err = tx.Exec(insertQuery, conn.ID, key, value, isSecret, now, now)
				if err != nil {
					return err
				}
			}
		}

		// Delete any remaining keys that aren't in the new config
		for key := range existingKeys {
			deleteQuery := `
				DELETE FROM api_connection_configs
				WHERE connection_id = $1 AND key = $2
			`

			_, err = tx.Exec(deleteQuery, conn.ID, key)
			if err != nil {
				return err
			}
		}
	} else {
		// If the new config is empty, delete all config items
		deleteQuery := `
			DELETE FROM api_connection_configs
			WHERE connection_id = $1
		`

		_, err = tx.Exec(deleteQuery, conn.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetConnectionConfig retrieves all configuration for a connection
func (r *APIConnectionRepository) GetConnectionConfig(connectionID string, includeSecrets bool) (map[string]string, error) {
	query := `
		SELECT key, value, is_secret
		FROM api_connection_configs
		WHERE connection_id = $1
	`

	rows, err := r.db.Query(query, connectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	config := make(map[string]string)
	for rows.Next() {
		var key, value string
		var isSecret bool

		if err := rows.Scan(&key, &value, &isSecret); err != nil {
			return nil, err
		}

		// Skip secrets if not requested
		if isSecret && !includeSecrets {
			continue
		}

		config[key] = value
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

// SetActive sets the active status of a connection
func (r *APIConnectionRepository) SetActive(id string, active bool) error {
	query := `
		UPDATE api_connections
		SET active = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.Exec(
		query,
		active,
		time.Now(),
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API connection not found: %s", id)
	}

	return nil
}
