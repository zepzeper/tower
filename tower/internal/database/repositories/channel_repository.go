package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zepzeper/tower/internal/database/models"
)

// ChannelRepository handles database operations for channels
type ChannelRepository struct {
	db *sql.DB
}

// NewChannelRepository creates a new channel repository
func NewChannelRepository(db *sql.DB) *ChannelRepository {
	return &ChannelRepository{
		db: db,
	}
}

// GetAll retrieves all channels from the database
func (r *ChannelRepository) GetAll() ([]models.Channel, error) {
	query := `
		SELECT id, name, type, description, config, created_at, updated_at
		FROM channels
		ORDER BY name ASC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var channels []models.Channel
	for rows.Next() {
		var channel models.Channel
		if err := rows.Scan(
			&channel.ID,
			&channel.Name,
			&channel.Type,
			&channel.Description,
			&channel.Config,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		); err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return channels, nil
}

// GetByID retrieves a channel by ID
func (r *ChannelRepository) GetByID(id string) (models.Channel, error) {
	query := `
		SELECT id, name, type, description, config, created_at, updated_at
		FROM channels
		WHERE id = $1
	`
	
	var channel models.Channel
	err := r.db.QueryRow(query, id).Scan(
		&channel.ID,
		&channel.Name,
		&channel.Type,
		&channel.Description,
		&channel.Config,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return channel, fmt.Errorf("channel not found: %s", id)
		}
		return channel, err
	}
	
	return channel, nil
}

// Create inserts a new channel into the database
func (r *ChannelRepository) Create(channel models.Channel) error {
	query := `
		INSERT INTO channels (id, name, type, description, config, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	now := time.Now()
	channel.CreatedAt = now
	channel.UpdatedAt = now
	
	_, err := r.db.Exec(
		query,
		channel.ID,
		channel.Name,
		channel.Type,
		channel.Description,
		channel.Config,
		channel.CreatedAt,
		channel.UpdatedAt,
	)
	
	return err
}

// Update updates an existing channel
func (r *ChannelRepository) Update(channel models.Channel) error {
	query := `
		UPDATE channels
		SET name = $1, type = $2, description = $3, config = $4, updated_at = $5
		WHERE id = $6
	`
	
	channel.UpdatedAt = time.Now()
	
	result, err := r.db.Exec(
		query,
		channel.Name,
		channel.Type,
		channel.Description,
		channel.Config,
		channel.UpdatedAt,
		channel.ID,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("channel not found: %s", channel.ID)
	}
	
	return nil
}

// Delete removes a channel from the database
func (r *ChannelRepository) Delete(id string) error {
	query := `DELETE FROM channels WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("channel not found: %s", id)
	}
	
	return nil
}

// GetByType retrieves all channels of a specific type
func (r *ChannelRepository) GetByType(channelType string) ([]models.Channel, error) {
	query := `
		SELECT id, name, type, description, config, created_at, updated_at
		FROM channels
		WHERE type = $1
		ORDER BY name ASC
	`
	
	rows, err := r.db.Query(query, channelType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var channels []models.Channel
	for rows.Next() {
		var channel models.Channel
		if err := rows.Scan(
			&channel.ID,
			&channel.Name,
			&channel.Type,
			&channel.Description,
			&channel.Config,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		); err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return channels, nil
}

// UpdateConfig updates only the configuration of a channel
func (r *ChannelRepository) UpdateConfig(id string, config map[string]interface{}) error {
	query := `
		UPDATE channels
		SET config = $1, updated_at = $2
		WHERE id = $3
	`
	
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	
	result, err := r.db.Exec(
		query,
		configJSON,
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
		return fmt.Errorf("channel not found: %s", id)
	}
	
	return nil
}
