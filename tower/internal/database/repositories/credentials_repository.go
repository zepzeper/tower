package repositories

import (
	"database/sql"
	"fmt"

	"github.com/zepzeper/tower/internal/database/models"
)

type CredentialsRepository struct {
	db *sql.DB
}

func NewCredentialsRepository(db *sql.DB) *CredentialsRepository {
	return &CredentialsRepository{db: db}
}

func (r *CredentialsRepository) GetAll() ([]models.Credentials, error) {
	query := `
		SELECT id, name, description, type, active, created_at, updated_at
		FROM credentials
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []models.Credentials
	for rows.Next() {
		var conn models.Credentials
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
		connections = append(connections, conn)
	}

	return connections, rows.Err()
}

func (r *CredentialsRepository) GetByID(id string) (*models.Credentials, error) {
	query := `
		SELECT id, name, description, type, active, created_at, updated_at
		FROM credentials
		WHERE id = $1
	`

	var conn models.Credentials
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
			return nil, fmt.Errorf("API connection not found: %s", id)
		}
		return nil, err
	}

	return &conn, nil
}

func (r *CredentialsRepository) GetConfigsByConnectionID(connectionID string) ([]models.CredentialsConfig, error) {
	query := `
		SELECT connection_id, key, value, is_secret, created_at, updated_at
		FROM credentials_configs
		WHERE connection_id = $1
	`

	rows, err := r.db.Query(query, connectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []models.CredentialsConfig
	for rows.Next() {
		var cfg models.CredentialsConfig
		if err := rows.Scan(
			&cfg.ConnectionID,
			&cfg.Key,
			&cfg.Value,
			&cfg.IsSecret,
			&cfg.CreatedAt,
			&cfg.UpdatedAt,
		); err != nil {
			return nil, err
		}
		configs = append(configs, cfg)
	}

	return configs, rows.Err()
}

func (r *CredentialsRepository) Create(conn *models.Credentials, configs []models.CredentialsConfig) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO credentials (id, name, description, type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

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

	if len(configs) > 0 {
		for _, config := range configs {
			configQuery := `
				INSERT INTO credentials_configs (connection_id, key, value, is_secret, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6)
			`
			_, err := tx.Exec(
				configQuery,
				config.ConnectionID,
				config.Key,
				config.Value,
				config.IsSecret,
				config.CreatedAt,
				config.UpdatedAt,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *CredentialsRepository) Update(conn *models.Credentials, configs []models.CredentialsConfig) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE credentials
		SET name = $1, description = $2, type = $3, active = $4, updated_at = $5
		WHERE id = $6
	`

	_, err = tx.Exec(
		query,
		conn.Name,
		conn.Description,
		conn.Type,
		conn.Active,
		conn.UpdatedAt,
		conn.ID,
	)
	if err != nil {
		return err
	}

	// Remove existing configs
	_, err = tx.Exec("DELETE FROM credentials_configs WHERE connection_id = $1", conn.ID)
	if err != nil {
		return err
	}

	for _, config := range configs {
		insert := `
			INSERT INTO credentials_configs (connection_id, key, value, is_secret, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err := tx.Exec(
			insert,
			config.ConnectionID,
			config.Key,
			config.Value,
			config.IsSecret,
			config.CreatedAt,
			config.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *CredentialsRepository) Delete(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM credentials_configs WHERE connection_id = $1", id)
	if err != nil {
		return err
	}

	res, err := tx.Exec("DELETE FROM credentials WHERE id = $1", id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("no connection found with ID: %s", id)
	}

	return tx.Commit()
}
