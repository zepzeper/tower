package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zepzeper/tower/internal/database/models"
)

// TransformerRepository handles database operations for transformers
type TransformerRepository struct {
	db *sql.DB
}

// NewTransformerRepository creates a new transformer repository
func NewTransformerRepository(db *sql.DB) *TransformerRepository {
	return &TransformerRepository{
		db: db,
	}
}

// GetAll retrieves all transformers
func (r *TransformerRepository) GetAll() ([]models.Transformer, error) {
	query := `
		SELECT id, name, description, mappings, functions, created_at, updated_at
		FROM transformers
		ORDER BY name ASC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var transformers []models.Transformer
	for rows.Next() {
		var transformer models.Transformer
		if err := rows.Scan(
			&transformer.ID,
			&transformer.Name,
			&transformer.Description,
			&transformer.Mappings,
			&transformer.Functions,
			&transformer.CreatedAt,
			&transformer.UpdatedAt,
		); err != nil {
			return nil, err
		}
		transformers = append(transformers, transformer)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return transformers, nil
}

// GetByID retrieves a transformer by ID
func (r *TransformerRepository) GetByID(id string) (models.Transformer, error) {
	query := `
		SELECT id, name, description, mappings, functions, created_at, updated_at
		FROM transformers
		WHERE id = $1
	`
	
	var transformer models.Transformer
	err := r.db.QueryRow(query, id).Scan(
		&transformer.ID,
		&transformer.Name,
		&transformer.Description,
		&transformer.Mappings,
		&transformer.Functions,
		&transformer.CreatedAt,
		&transformer.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return transformer, fmt.Errorf("transformer not found: %s", id)
		}
		return transformer, err
	}
	
	return transformer, nil
}

// Create inserts a new transformer
func (r *TransformerRepository) Create(transformer models.Transformer) error {
	query := `
		INSERT INTO transformers (id, name, description, mappings, functions, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	now := time.Now()
	transformer.CreatedAt = now
	transformer.UpdatedAt = now
	
	_, err := r.db.Exec(
		query,
		transformer.ID,
		transformer.Name,
		transformer.Description,
		transformer.Mappings,
		transformer.Functions,
		transformer.CreatedAt,
		transformer.UpdatedAt,
	)
	
	return err
}

// Update updates an existing transformer
func (r *TransformerRepository) Update(transformer models.Transformer) error {
	query := `
		UPDATE transformers
		SET name = $1, description = $2, mappings = $3, functions = $4, updated_at = $5
		WHERE id = $6
	`
	
	transformer.UpdatedAt = time.Now()
	
	result, err := r.db.Exec(
		query,
		transformer.Name,
		transformer.Description,
		transformer.Mappings,
		transformer.Functions,
		transformer.UpdatedAt,
		transformer.ID,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("transformer not found: %s", transformer.ID)
	}
	
	return nil
}

// Delete removes a transformer
func (r *TransformerRepository) Delete(id string) error {
	query := `DELETE FROM transformers WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("transformer not found: %s", id)
	}
	
	return nil
}
