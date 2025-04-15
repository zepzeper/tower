package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zepzeper/tower/internal/database/models"
)

// WorkflowRepository handles database operations for workflows
type WorkflowRepository struct {
	db *sql.DB
}

// NewWorkflowRepository creates a new workflow repository
func NewWorkflowRepository(db *sql.DB) *WorkflowRepository {
	return &WorkflowRepository{
		db: db,
	}
}

// GetAll retrieves all workflows from the database
func (r *WorkflowRepository) GetAll() ([]models.Workflow, error) {
	query := `
		SELECT id, name, description, triggers, actions, active, created_at, updated_at
		FROM workflows
		ORDER BY name ASC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var workflows []models.Workflow
	for rows.Next() {
		var workflow models.Workflow
		if err := rows.Scan(
			&workflow.ID,
			&workflow.Name,
			&workflow.Description,
			&workflow.Triggers,
			&workflow.Actions,
			&workflow.Active,
			&workflow.CreatedAt,
			&workflow.UpdatedAt,
		); err != nil {
			return nil, err
		}
		workflows = append(workflows, workflow)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return workflows, nil
}

// GetActive retrieves all active workflows
func (r *WorkflowRepository) GetActive() ([]models.Workflow, error) {
	query := `
		SELECT id, name, description, triggers, actions, active, created_at, updated_at
		FROM workflows
		WHERE active = true
		ORDER BY name ASC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var workflows []models.Workflow
	for rows.Next() {
		var workflow models.Workflow
		if err := rows.Scan(
			&workflow.ID,
			&workflow.Name,
			&workflow.Description,
			&workflow.Triggers,
			&workflow.Actions,
			&workflow.Active,
			&workflow.CreatedAt,
			&workflow.UpdatedAt,
		); err != nil {
			return nil, err
		}
		workflows = append(workflows, workflow)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return workflows, nil
}

// GetByID retrieves a workflow by ID
func (r *WorkflowRepository) GetByID(id string) (models.Workflow, error) {
	query := `
		SELECT id, name, description, triggers, actions, active, created_at, updated_at
		FROM workflows
		WHERE id = $1
	`
	
	var workflow models.Workflow
	err := r.db.QueryRow(query, id).Scan(
		&workflow.ID,
		&workflow.Name,
		&workflow.Description,
		&workflow.Triggers,
		&workflow.Actions,
		&workflow.Active,
		&workflow.CreatedAt,
		&workflow.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return workflow, fmt.Errorf("workflow not found: %s", id)
		}
		return workflow, err
	}
	
	return workflow, nil
}

// Create inserts a new workflow into the database
func (r *WorkflowRepository) Create(workflow models.Workflow) error {
	query := `
		INSERT INTO workflows (id, name, description, triggers, actions, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	now := time.Now()
	workflow.CreatedAt = now
	workflow.UpdatedAt = now
	
	_, err := r.db.Exec(
		query,
		workflow.ID,
		workflow.Name,
		workflow.Description,
		workflow.Triggers,
		workflow.Actions,
		workflow.Active,
		workflow.CreatedAt,
		workflow.UpdatedAt,
	)
	
	return err
}

// Update updates an existing workflow
func (r *WorkflowRepository) Update(workflow models.Workflow) error {
	query := `
		UPDATE workflows
		SET name = $1, description = $2, triggers = $3, actions = $4, active = $5, updated_at = $6
		WHERE id = $7
	`
	
	workflow.UpdatedAt = time.Now()
	
	result, err := r.db.Exec(
		query,
		workflow.Name,
		workflow.Description,
		workflow.Triggers,
		workflow.Actions,
		workflow.Active,
		workflow.UpdatedAt,
		workflow.ID,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("workflow not found: %s", workflow.ID)
	}
	
	return nil
}

// Delete removes a workflow from the database
func (r *WorkflowRepository) Delete(id string) error {
	query := `DELETE FROM workflows WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("workflow not found: %s", id)
	}
	
	return nil
}

// SetActive updates the active status of a workflow
func (r *WorkflowRepository) SetActive(id string, active bool) error {
	query := `
		UPDATE workflows
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
		return fmt.Errorf("workflow not found: %s", id)
	}
	
	return nil
}

// FindByTriggerChannel finds all workflows with triggers for a specific channel
func (r *WorkflowRepository) FindByTriggerChannel(channelID string) ([]models.Workflow, error) {
	query := `
		SELECT id, name, description, triggers, actions, active, created_at, updated_at
		FROM workflows
		WHERE active = true AND triggers::jsonb @> ANY(ARRAY[jsonb_build_array(jsonb_build_object('channelId', $1))])
		ORDER BY name ASC
	`
	
	rows, err := r.db.Query(query, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var workflows []models.Workflow
	for rows.Next() {
		var workflow models.Workflow
		if err := rows.Scan(
			&workflow.ID,
			&workflow.Name,
			&workflow.Description,
			&workflow.Triggers,
			&workflow.Actions,
			&workflow.Active,
			&workflow.CreatedAt,
			&workflow.UpdatedAt,
		); err != nil {
			return nil, err
		}
		workflows = append(workflows, workflow)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return workflows, nil
}
