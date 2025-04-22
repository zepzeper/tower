package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zepzeper/tower/internal/database/models"
	"github.com/zepzeper/tower/internal/webapi/dto"
)

type RelationsRepository struct {
	db *sql.DB
}

func NewRelationsRepository(db *sql.DB) *RelationsRepository {
	return &RelationsRepository{
		db: db,
	}
}

// Create inserts a new connection between credentials
func (r *RelationsRepository) Create(conn dto.CredentialsRelationCreateRequest) error {
	query := `
		INSERT INTO credentials_connections (initiator_id, target_id, connection_type, active, endpoint, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(
		query,
		conn.InitiatorID,
		conn.TargetID,
		conn.ConnectionType,
		conn.Active,
		conn.Endpoint,
		time.Now(),
	)

	return err
}

// GetByInitiatorID finds all connections where the specified credential is the initiator
func (r *RelationsRepository) GetByInitiatorID(initiatorID string) ([]dto.CredentialsRelationResponse, error) {
	query := `
		SELECT initiator_id, target_id, connection_type, active, endpoint, created_at
		FROM credentials_connections
		WHERE initiator_id = $1
	`

	rows, err := r.db.Query(query, initiatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations []dto.CredentialsRelationResponse

	for rows.Next() {
		var relation dto.CredentialsRelationResponse
		if err := rows.Scan(
			&relation.InitiatorID,
			&relation.TargetID,
			&relation.ConnectionType,
			&relation.Active,
			&relation.Endpoint,
			&relation.CreatedAt,
		); err != nil {
			return nil, err
		}
		relations = append(relations, relation)
	}

	return relations, rows.Err()
}

func (r *RelationsRepository) GetLogsByInitiatorID(initiatorID string) ([]dto.CredentialsRelationLogsResponse, error) {
	query := `
		SELECT id, initiator_id, target_id, message, created_at
		FROM credentials_connections_logs
		WHERE initiator_id = $1
	`

	fmt.Println(initiatorID)

	rows, err := r.db.Query(query, initiatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations []dto.CredentialsRelationLogsResponse
	for rows.Next() {
		var relation dto.CredentialsRelationLogsResponse
		if err := rows.Scan(
			&relation.ID,
			&relation.InitiatorID,
			&relation.TargetID,
			&relation.Message,
			&relation.CreatedAt,
		); err != nil {
			return nil, err
		}
		relations = append(relations, relation)
	}
	return relations, rows.Err()
}

// GetByTargetID finds all connections where the specified credential is the target
func (r *RelationsRepository) GetByTargetID(targetID string) ([]models.CredentialsRelation, error) {
	query := `
		SELECT initiator_id, target_id, connection_type, active, endpoint, created_at
		FROM credentials_connections
		WHERE target_id = $1
	`

	rows, err := r.db.Query(query, targetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations []models.CredentialsRelation
	for rows.Next() {
		var relation models.CredentialsRelation
		if err := rows.Scan(
			&relation.InitiatorID,
			&relation.TargetID,
			&relation.ConnectionType,
			&relation.Active,
			&relation.Endpoint,
			&relation.CreatedAt,
		); err != nil {
			return nil, err
		}
		relations = append(relations, relation)
	}

	return relations, rows.Err()
}

// GetAllConnections returns all connections for a credential (both as initiator and target)
func (r *RelationsRepository) GetAllConnections(credentialID string) ([]models.CredentialsRelation, error) {
	query := `
		SELECT initiator_id, target_id, connection_type, active, endpoint, created_at
		FROM credentials_connections
		WHERE initiator_id = $1 OR target_id = $1
	`

	rows, err := r.db.Query(query, credentialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations []models.CredentialsRelation
	for rows.Next() {
		var relation models.CredentialsRelation
		if err := rows.Scan(
			&relation.InitiatorID,
			&relation.TargetID,
			&relation.ConnectionType,
			&relation.Active,
			&relation.Endpoint,
			&relation.CreatedAt,
		); err != nil {
			return nil, err
		}
		relations = append(relations, relation)
	}

	return relations, rows.Err()
}

// GetConnection retrieves a specific connection between two credentials
func (r *RelationsRepository) GetConnection(initiatorID, targetID string) (models.CredentialsRelation, error) {
	query := `
		SELECT initiator_id, target_id, connection_type, active, endpoint, created_at
		FROM credentials_connections
		WHERE initiator_id = $1 AND target_id = $2
	`

	var relation models.CredentialsRelation
	err := r.db.QueryRow(query, initiatorID, targetID).Scan(
		&relation.InitiatorID,
		&relation.TargetID,
		&relation.ConnectionType,
		&relation.Active,
		&relation.Endpoint,
		&relation.CreatedAt,
	)

	return relation, err
}

// UpdateConnection updates an existing connection between credentials
func (r *RelationsRepository) UpdateConnection(conn dto.CredentialsRelationUpdateRequest) error {
	query := `
		UPDATE credentials_connections
		SET connection_type = $3, active = $4, endpoint = $5
		WHERE initiator_id = $1 AND target_id = $2
	`

	_, err := r.db.Exec(
		query,
		conn.InitiatorID,
		conn.TargetID,
		conn.ConnectionType,
		conn.Active,
		conn.Endpoint,
	)

	return err
}

// DeleteConnection removes a connection between credentials
func (r *RelationsRepository) DeleteConnection(initiatorID, targetID string) error {
	query := `
		DELETE FROM credentials_connections
		WHERE initiator_id = $1 AND target_id = $2
	`

	_, err := r.db.Exec(query, initiatorID, targetID)
	return err
}

// GetConnectionsByType retrieves all connections of a specific type
func (r *RelationsRepository) GetConnectionsByType(connectionType string) ([]models.CredentialsRelation, error) {
	query := `
		SELECT initiator_id, target_id, connection_type, active, endpoint, created_at
		FROM credentials_connections
		WHERE connection_type = $1
	`

	rows, err := r.db.Query(query, connectionType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations []models.CredentialsRelation
	for rows.Next() {
		var relation models.CredentialsRelation
		if err := rows.Scan(
			&relation.InitiatorID,
			&relation.TargetID,
			&relation.ConnectionType,
			&relation.Active,
			&relation.Endpoint,
			&relation.CreatedAt,
		); err != nil {
			return nil, err
		}
		relations = append(relations, relation)
	}

	return relations, rows.Err()
}
