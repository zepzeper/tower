package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
  "github.com/zepzeper/tower/internal/core/connectors"  // Add this
	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/database/models"
)

// ConnectionService handles connections between systems
type ConnectionService struct {
	dbManager         *database.Manager
	connectorService  *ConnectorService
	transformerService *TransformerService
	jobManager        *JobManager
}

// NewConnectionService creates a new connection service
func NewConnectionService(
	dbManager *database.Manager,
	connectorService *ConnectorService,
	transformerService *TransformerService,
	jobManager *JobManager,
) *ConnectionService {
	return &ConnectionService{
		dbManager:         dbManager,
		connectorService:  connectorService,
		transformerService: transformerService,
		jobManager:        jobManager,
	}
}

// ConnectionInfo contains information about a connection
type ConnectionInfo struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	SourceID      string                 `json:"sourceId"`
	TargetID      string                 `json:"targetId"`
	TransformerID string                 `json:"transformerId"`
	Query         map[string]interface{} `json:"query"`
	Schedule      string                 `json:"schedule"`
	Active        bool                   `json:"active"`
	LastRun       *time.Time             `json:"lastRun,omitempty"`
	Status        string                 `json:"status,omitempty"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

// ListConnections retrieves all connections
func (s *ConnectionService) ListConnections() ([]ConnectionInfo, error) {
	// Get all connections from database
	connections, err := s.dbManager.Repos.Connection().GetAll()
	if err != nil {
		return nil, fmt.Errorf("error retrieving connections: %w", err)
	}
	
	// Convert to response format
	result := make([]ConnectionInfo, len(connections))
	for i, conn := range connections {
		info, err := s.modelToConnectionInfo(conn)
		if err != nil {
			return nil, err
		}
		result[i] = *info
	}
	
	return result, nil
}

// GetConnection gets a connection by ID
func (s *ConnectionService) GetConnection(id string) (*ConnectionInfo, error) {
	// Get connection from database
	connection, err := s.dbManager.Repos.Connection().GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("connection not found: %w", err)
	}
	
	// Convert to response format
	return s.modelToConnectionInfo(connection)
}

// CreateConnection creates a new connection
func (s *ConnectionService) CreateConnection(
	ctx context.Context,
	name, description string,
	sourceID, targetID, transformerID string,
	query map[string]interface{},
	schedule string,
) (string, error) {
	// Validate source connector
	_, err := s.connectorService.GetConnector(sourceID)
	if err != nil {
		return "", fmt.Errorf("source connector not found: %w", err)
	}
	
	// Validate target connector
	_, err = s.connectorService.GetConnector(targetID)
	if err != nil {
		return "", fmt.Errorf("target connector not found: %w", err)
	}
	
	// Validate transformer
	_, err = s.transformerService.GetTransformer(transformerID)
	if err != nil {
		return "", fmt.Errorf("transformer not found: %w", err)
	}
	
	// Generate ID
	connectionID := generateID()
	
	// Create config with query
	config := map[string]interface{}{
		"query": query,
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("error serializing config: %w", err)
	}
	
	// Create connection model
	connection := models.Connection{
		ID:            connectionID,
		Name:          name,
		Description:   sqlNullString(description),
		SourceID:      sourceID,
		TargetID:      targetID,
		TransformerID: transformerID,
		Config:        configJSON,
		Schedule:      schedule,
		Active:        true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	
	// Save to database
	if err := s.dbManager.Repos.Connection().Create(connection); err != nil {
		return "", fmt.Errorf("error saving connection: %w", err)
	}
	
	// Schedule the job if needed
	if schedule != "" {
		if err := s.jobManager.ScheduleJob(ctx, connectionID, sourceID, targetID, transformerID, query, schedule); err != nil {
			return connectionID, fmt.Errorf("connection created but job scheduling failed: %w", err)
		}
	}
	
	return connectionID, nil
}

// UpdateConnection updates an existing connection
func (s *ConnectionService) UpdateConnection(
	ctx context.Context,
	id, name, description string,
	sourceID, targetID, transformerID string,
	query map[string]interface{},
	schedule string,
	active bool,
) error {
	// Get existing connection
	connection, err := s.dbManager.Repos.Connection().GetByID(id)
	if err != nil {
		return fmt.Errorf("connection not found: %w", err)
	}
	
	// Cancel existing job if active
	if connection.Active && schedule != connection.Schedule {
		s.jobManager.CancelJob(id)
	}
	
	// Validate source connector
	_, err = s.connectorService.GetConnector(sourceID)
	if err != nil {
		return fmt.Errorf("source connector not found: %w", err)
	}
	
	// Validate target connector
	_, err = s.connectorService.GetConnector(targetID)
	if err != nil {
		return fmt.Errorf("target connector not found: %w", err)
	}
	
	// Validate transformer
	_, err = s.transformerService.GetTransformer(transformerID)
	if err != nil {
		return fmt.Errorf("transformer not found: %w", err)
	}
	
	// Update config with query
	config := map[string]interface{}{
		"query": query,
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("error serializing config: %w", err)
	}
	
	// Update connection model
	connection.Name = name
	connection.Description = sqlNullString(description)
	connection.SourceID = sourceID
	connection.TargetID = targetID
	connection.TransformerID = transformerID
	connection.Config = configJSON
	connection.Schedule = schedule
	connection.Active = active
	connection.UpdatedAt = time.Now()
	
	// Save to database
	if err := s.dbManager.Repos.Connection().Update(connection); err != nil {
		return fmt.Errorf("error updating connection: %w", err)
	}
	
	// Schedule a new job if active and has schedule
	if active && schedule != "" {
		if err := s.jobManager.ScheduleJob(ctx, id, sourceID, targetID, transformerID, query, schedule); err != nil {
			return fmt.Errorf("connection updated but job scheduling failed: %w", err)
		}
	}
	
	return nil
}

// DeleteConnection deletes a connection
func (s *ConnectionService) DeleteConnection(id string) error {
	// Cancel any running job
	s.jobManager.CancelJob(id)
	
	// Delete from database
	return s.dbManager.Repos.Connection().Delete(id)
}

// ExecuteConnection executes a connection once
func (s *ConnectionService) ExecuteConnection(ctx context.Context, id string) (string, error) {
	// Get connection
	connection, err := s.dbManager.Repos.Connection().GetByID(id)
	if err != nil {
		return "", fmt.Errorf("connection not found: %w", err)
	}
	
	// Extract query from config
	var config map[string]interface{}
	err = json.Unmarshal(connection.Config, &config)
	if err != nil {
		return "", fmt.Errorf("error parsing config: %w", err)
	}
	
	query, _ := config["query"].(map[string]interface{})
	
	// Create execution record
	execution := models.Execution{
		ID:           generateID(),
		ConnectionID: id,
		Status:       "in_progress",
		StartTime:    time.Now(),
		CreatedAt:    time.Now(),
	}
	err = s.dbManager.Repos.Execution().Create(execution)
	if err != nil {
		return "", fmt.Errorf("error creating execution record: %w", err)
	}
	
	// Update connection last run time
	s.dbManager.Repos.Connection().UpdateLastRun(id, execution.StartTime)
	
	// Execute in a goroutine
	go func() {
		// Create a context with timeout
		execCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		
		// 1. Fetch data from source
		sourceData, err := s.connectorService.FetchData(execCtx, connection.SourceID, query)
		if err != nil {
			updateExecutionError(s.dbManager, execution, err)
			return
		}
		
		// Save source data in execution record
		sourceDataJSON, _ := json.Marshal(sourceData)
		execution.SourceData = sourceDataJSON
		
		// Skip if no data
		if len(sourceData) == 0 {
			updateExecutionSuccess(s.dbManager, execution, nil)
			return
		}
		
		// 2. Transform data
    transformedData := make([]connectors.DataPayload, 0, len(sourceData))
		for _, item := range sourceData {
			transformed, err := s.transformerService.TransformData(execCtx, connection.TransformerID, item)
			if err != nil {
				updateExecutionError(s.dbManager, execution, err)
				return
			}
			transformedData = append(transformedData, transformed)
		}
		
		// 3. Push to target
		err = s.connectorService.PushData(execCtx, connection.TargetID, transformedData)
		if err != nil {
			updateExecutionError(s.dbManager, execution, err)
			return
		}
		
		// Update execution record with success
		targetDataJSON, _ := json.Marshal(transformedData)
		updateExecutionSuccess(s.dbManager, execution, targetDataJSON)
	}()
	
	return execution.ID, nil
}

// GetExecutions returns executions for a connection
func (s *ConnectionService) GetExecutions(connectionID string, limit, offset int) ([]models.Execution, error) {
	if limit <= 0 {
		limit = 20
	}
	
	return s.dbManager.Repos.Execution().GetByConnectionID(connectionID, limit, offset)
}

// SetActive sets the active status of a connection
func (s *ConnectionService) SetActive(ctx context.Context, id string, active bool) error {
	// Get connection
	connection, err := s.dbManager.Repos.Connection().GetByID(id)
	if err != nil {
		return fmt.Errorf("connection not found: %w", err)
	}
	
	// No change needed if already in desired state
	if connection.Active == active {
		return nil
	}
	
	// Update in database
	err = s.dbManager.Repos.Connection().SetActive(id, active)
	if err != nil {
		return fmt.Errorf("error updating connection status: %w", err)
	}
	
	// If activating and has a schedule, start the job
	if active && connection.Schedule != "" {
		// Extract query from config
		var config map[string]interface{}
		err = json.Unmarshal(connection.Config, &config)
		if err != nil {
			return fmt.Errorf("error parsing config: %w", err)
		}
		
		query, _ := config["query"].(map[string]interface{})
		
		err = s.jobManager.ScheduleJob(
			ctx,
			id,
			connection.SourceID,
			connection.TargetID,
			connection.TransformerID,
			query,
			connection.Schedule,
		)
		if err != nil {
			return fmt.Errorf("connection activated but job scheduling failed: %w", err)
		}
	} else if !active {
		// If deactivating, cancel any running job
		s.jobManager.CancelJob(id)
	}
	
	return nil
}

// Helper function to convert model to ConnectionInfo
func (s *ConnectionService) modelToConnectionInfo(model models.Connection) (*ConnectionInfo, error) {
	// Extract query from config
	var config map[string]interface{}
	err := json.Unmarshal(model.Config, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}
	
	query, _ := config["query"].(map[string]interface{})
	
	// Get status if job is active
	status := ""
	if model.Active && model.Schedule != "" {
		status, _ = s.jobManager.GetJobStatus(model.ID)
	}
	
	info := &ConnectionInfo{
		ID:            model.ID,
		Name:          model.Name,
		Description:   model.Description.String,
		SourceID:      model.SourceID,
		TargetID:      model.TargetID,
		TransformerID: model.TransformerID,
		Query:         query,
		Schedule:      model.Schedule,
		Active:        model.Active,
		Status:        status,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
	
	if model.LastRun.Valid {
		lastRun := model.LastRun.Time
		info.LastRun = &lastRun
	}
	
	return info, nil
}
