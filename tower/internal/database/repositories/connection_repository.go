package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zepzeper/tower/internal/database/models"
)

// ConnectionRepository handles database operations for connections
type ConnectionRepository struct {
	db *sql.DB
}

// NewConnectionRepository creates a new connection repository
func NewConnectionRepository(db *sql.DB) *ConnectionRepository {
	return &ConnectionRepository{
		db: db,
	}
}

// GetAll retrieves all connections
func (r *ConnectionRepository) GetAll() ([]models.Connection, error) {
	query := `
		SELECT id, name, description, source_id, target_id, transformer_id, config, schedule, active, last_run, created_at, updated_at
		FROM connections
		ORDER BY name ASC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var connections []models.Connection
	for rows.Next() {
		var connection models.Connection
		if err := rows.Scan(
			&connection.ID,
			&connection.Name,
			&connection.Description,
			&connection.SourceID,
			&connection.TargetID,
			&connection.TransformerID,
			&connection.Config,
			&connection.Schedule,
			&connection.Active,
			&connection.LastRun,
			&connection.CreatedAt,
			&connection.UpdatedAt,
		); err != nil {
			return nil, err
		}
		connections = append(connections, connection)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return connections, nil
}

// GetByID retrieves a connection by ID
func (r *ConnectionRepository) GetByID(id string) (models.Connection, error) {
	query := `
		SELECT id, name, description, source_id, target_id, transformer_id, config, schedule, active, last_run, created_at, updated_at
		FROM connections
		WHERE id = $1
	`
	
	var connection models.Connection
	err := r.db.QueryRow(query, id).Scan(
		&connection.ID,
		&connection.Name,
		&connection.Description,
		&connection.SourceID,
		&connection.TargetID,
		&connection.TransformerID,
		&connection.Config,
		&connection.Schedule,
		&connection.Active,
		&connection.LastRun,
		&connection.CreatedAt,
		&connection.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return connection, fmt.Errorf("connection not found: %s", id)
		}
		return connection, err
	}
	
	return connection, nil
}

// Create inserts a new connection
func (r *ConnectionRepository) Create(connection models.Connection) error {
	query := `
		INSERT INTO connections (id, name, description, source_id, target_id, transformer_id, config, schedule, active, last_run, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	
	now := time.Now()
	connection.CreatedAt = now
	connection.UpdatedAt = now
	
	_, err := r.db.Exec(
		query,
		connection.ID,
		connection.Name,
		connection.Description,
		connection.SourceID,
		connection.TargetID,
		connection.TransformerID,
		connection.Config,
		connection.Schedule,
		connection.Active,
		connection.LastRun,
		connection.CreatedAt,
		connection.UpdatedAt,
	)
	
	return err
}

// Update updates an existing connection
func (r *ConnectionRepository) Update(connection models.Connection) error {
	query := `
		UPDATE connections
		SET name = $1, description = $2, source_id = $3, target_id = $4, transformer_id = $5, 
			config = $6, schedule = $7, active = $8, last_run = $9, updated_at = $10
		WHERE id = $11
	`
	
	connection.UpdatedAt = time.Now()
	
	result, err := r.db.Exec(
		query,
		connection.Name,
		connection.Description,
		connection.SourceID,
		connection.TargetID,
		connection.TransformerID,
		connection.Config,
		connection.Schedule,
		connection.Active,
		connection.LastRun,
		connection.UpdatedAt,
		connection.ID,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("connection not found: %s", connection.ID)
	}
	
	return nil
}

// Delete removes a connection
func (r *ConnectionRepository) Delete(id string) error {
	query := `DELETE FROM connections WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("connection not found: %s", id)
	}
	
	return nil
}

// UpdateLastRun updates the last run timestamp for a connection
func (r *ConnectionRepository) UpdateLastRun(id string, lastRun time.Time) error {
	query := `
		UPDATE connections
		SET last_run = $1
		WHERE id = $2
	`
	
	result, err := r.db.Exec(query, lastRun, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("connection not found: %s", id)
	}
	
	return nil
}

// SetActive sets the active status of a connection
func (r *ConnectionRepository) SetActive(id string, active bool) error {
	query := `
		UPDATE connections
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
		return fmt.Errorf("connection not found: %s", id)
	}
	
	return nil
}
