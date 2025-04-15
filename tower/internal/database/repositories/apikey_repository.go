package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zepzeper/tower/internal/database/models"
)

// APIKeyRepository handles database operations for API keys
type APIKeyRepository struct {
	db *sql.DB
}

// NewAPIKeyRepository creates a new API key repository
func NewAPIKeyRepository(db *sql.DB) *APIKeyRepository {
	return &APIKeyRepository{
		db: db,
	}
}

// GetByID retrieves an API key by ID
func (r *APIKeyRepository) GetByID(id string) (models.APIKey, error) {
	query := `
		SELECT id, user_id, key, name, created_at, expires_at, last_used
		FROM api_keys
		WHERE id = $1
	`
	
	var apiKey models.APIKey
	err := r.db.QueryRow(query, id).Scan(
		&apiKey.ID,
		&apiKey.UserID,
		&apiKey.Key,
		&apiKey.Name,
		&apiKey.CreatedAt,
		&apiKey.ExpiresAt,
		&apiKey.LastUsed,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return apiKey, fmt.Errorf("API key not found: %s", id)
		}
		return apiKey, err
	}
	
	return apiKey, nil
}

// GetByKey retrieves an API key by its key value
func (r *APIKeyRepository) GetByKey(key string) (models.APIKey, error) {
	query := `
		SELECT id, user_id, key, name, created_at, expires_at, last_used
		FROM api_keys
		WHERE key = $1
	`
	
	var apiKey models.APIKey
	err := r.db.QueryRow(query, key).Scan(
		&apiKey.ID,
		&apiKey.UserID,
		&apiKey.Key,
		&apiKey.Name,
		&apiKey.CreatedAt,
		&apiKey.ExpiresAt,
		&apiKey.LastUsed,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return apiKey, fmt.Errorf("API key not found with value: %s", key)
		}
		return apiKey, err
	}
	
	return apiKey, nil
}

// GetByUserID retrieves all API keys for a user
func (r *APIKeyRepository) GetByUserID(userID string) ([]models.APIKey, error) {
	query := `
		SELECT id, user_id, key, name, created_at, expires_at, last_used
		FROM api_keys
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var apiKeys []models.APIKey
	for rows.Next() {
		var apiKey models.APIKey
		if err := rows.Scan(
			&apiKey.ID,
			&apiKey.UserID,
			&apiKey.Key,
			&apiKey.Name,
			&apiKey.CreatedAt,
			&apiKey.ExpiresAt,
			&apiKey.LastUsed,
		); err != nil {
			return nil, err
		}
		apiKeys = append(apiKeys, apiKey)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return apiKeys, nil
}

// Create inserts a new API key into the database
func (r *APIKeyRepository) Create(apiKey models.APIKey) error {
	query := `
		INSERT INTO api_keys (id, user_id, key, name, created_at, expires_at, last_used)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	now := time.Now()
	apiKey.CreatedAt = now
	apiKey.LastUsed = now
	
	_, err := r.db.Exec(
		query,
		apiKey.ID,
		apiKey.UserID,
		apiKey.Key,
		apiKey.Name,
		apiKey.CreatedAt,
		apiKey.ExpiresAt,
		apiKey.LastUsed,
	)
	
	return err
}

// Delete removes an API key from the database
func (r *APIKeyRepository) Delete(id string) error {
	query := `DELETE FROM api_keys WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("API key not found: %s", id)
	}
	
	return nil
}

// UpdateLastUsed updates the last_used timestamp for an API key
func (r *APIKeyRepository) UpdateLastUsed(id string) error {
	query := `
		UPDATE api_keys
		SET last_used = $1
		WHERE id = $2
	`
	
	now := time.Now()
	
	result, err := r.db.Exec(
		query,
		now,
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
		return fmt.Errorf("API key not found: %s", id)
	}
	
	return nil
}

// DeleteExpired removes all expired API keys
func (r *APIKeyRepository) DeleteExpired() (int64, error) {
	query := `DELETE FROM api_keys WHERE expires_at < $1`
	
	now := time.Now()
	
	result, err := r.db.Exec(query, now)
	if err != nil {
		return 0, err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	
	return rowsAffected, nil
}

// DeleteByUserID removes all API keys for a user
func (r *APIKeyRepository) DeleteByUserID(userID string) (int64, error) {
	query := `DELETE FROM api_keys WHERE user_id = $1`
	
	result, err := r.db.Exec(query, userID)
	if err != nil {
		return 0, err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	
	return rowsAffected, nil
}
