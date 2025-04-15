package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zepzeper/tower/internal/database/models"
)

// ExecutionRepository handles database operations for workflow executions
type ExecutionRepository struct {
	db *sql.DB
}

// NewExecutionRepository creates a new execution repository
func NewExecutionRepository(db *sql.DB) *ExecutionRepository {
	return &ExecutionRepository{
		db: db,
	}
}

// Create inserts a new execution record into the database
func (r *ExecutionRepository) Create(execution models.Execution) error {
	query := `
		INSERT INTO executions (id, workflow_id, status, start_time, end_time, result, error, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	now := time.Now()
	execution.CreatedAt = now
	
	_, err := r.db.Exec(
		query,
		execution.ID,
		execution.WorkflowID,
		execution.Status,
		execution.StartTime,
		execution.EndTime,
		execution.Result,
		execution.Error,
		execution.CreatedAt,
	)
	
	return err
}

// GetByID retrieves an execution by ID
func (r *ExecutionRepository) GetByID(id string) (models.Execution, error) {
	query := `
		SELECT id, workflow_id, status, start_time, end_time, result, error, created_at
		FROM executions
		WHERE id = $1
	`
	
	var execution models.Execution
	err := r.db.QueryRow(query, id).Scan(
		&execution.ID,
		&execution.WorkflowID,
		&execution.Status,
		&execution.StartTime,
		&execution.EndTime,
		&execution.Result,
		&execution.Error,
		&execution.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return execution, fmt.Errorf("execution not found: %s", id)
		}
		return execution, err
	}
	
	return execution, nil
}

// GetByWorkflowID retrieves all executions for a workflow
func (r *ExecutionRepository) GetByWorkflowID(workflowID string, limit, offset int) ([]models.Execution, error) {
	query := `
		SELECT id, workflow_id, status, start_time, end_time, result, error, created_at
		FROM executions
		WHERE workflow_id = $1
		ORDER BY start_time DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Query(query, workflowID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var executions []models.Execution
	for rows.Next() {
		var execution models.Execution
		if err := rows.Scan(
			&execution.ID,
			&execution.WorkflowID,
			&execution.Status,
			&execution.StartTime,
			&execution.EndTime,
			&execution.Result,
			&execution.Error,
			&execution.CreatedAt,
		); err != nil {
			return nil, err
		}
		executions = append(executions, execution)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return executions, nil
}

// UpdateStatus updates the status of an execution
func (r *ExecutionRepository) UpdateStatus(id, status string, endTime time.Time, resultData []byte, errorMsg string) error {
	query := `
		UPDATE executions
		SET status = $1, end_time = $2, result = $3, error = $4
		WHERE id = $5
	`
	
	var sqlError sql.NullString
	if errorMsg != "" {
		sqlError = sql.NullString{String: errorMsg, Valid: true}
	}
	
	dbResult, err := r.db.Exec(
		query,
		status,
		endTime,
		resultData,
		sqlError,
		id,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := dbResult.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("execution not found: %s", id)
	}
	
	return nil
}

// GetRecentExecutions retrieves recent executions across all workflows
func (r *ExecutionRepository) GetRecentExecutions(limit, offset int) ([]models.Execution, error) {
	query := `
		SELECT id, workflow_id, status, start_time, end_time, result, error, created_at
		FROM executions
		ORDER BY start_time DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var executions []models.Execution
	for rows.Next() {
		var execution models.Execution
		if err := rows.Scan(
			&execution.ID,
			&execution.WorkflowID,
			&execution.Status,
			&execution.StartTime,
			&execution.EndTime,
			&execution.Result,
			&execution.Error,
			&execution.CreatedAt,
		); err != nil {
			return nil, err
		}
		executions = append(executions, execution)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return executions, nil
}

// GetExecutionCounts retrieves counts of executions grouped by status
func (r *ExecutionRepository) GetExecutionCounts(workflowID string) (map[string]int, error) {
	query := `
		SELECT status, COUNT(*) as count
		FROM executions
		WHERE workflow_id = $1
		GROUP BY status
	`
	
	rows, err := r.db.Query(query, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	counts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return counts, nil
}
