package services

import (
	"context"
  "database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
	
  "github.com/zepzeper/tower/internal/core/connectors"
	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/database/models"
)

// Job represents a data transfer job
type Job struct {
	ID            string
	ConnectionID  string
	SourceID      string
	TargetID      string
	TransformerID string
	Query         map[string]interface{}
	Schedule      string
	LastRun       time.Time
	Status        string
	Errors        []string
	cancel        context.CancelFunc
}

// JobManager manages scheduled jobs
type JobManager struct {
	dbManager         *database.Manager
	connectorService  *ConnectorService
	transformerService *TransformerService
	activeJobs        map[string]*Job
	mu                sync.RWMutex
}

// NewJobManager creates a new job manager
func NewJobManager(
	dbManager *database.Manager,
	connectorService *ConnectorService,
	transformerService *TransformerService,
) *JobManager {
	return &JobManager{
		dbManager:         dbManager,
		connectorService:  connectorService,
		transformerService: transformerService,
		activeJobs:        make(map[string]*Job),
	}
}

// ScheduleJob schedules a job for execution
func (m *JobManager) ScheduleJob(
	ctx context.Context,
	connectionID string,
	sourceID string,
	targetID string,
	transformerID string,
	query map[string]interface{},
	schedule string,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Check if job is already scheduled
	if _, exists := m.activeJobs[connectionID]; exists {
		return fmt.Errorf("job already scheduled for connection: %s", connectionID)
	}
	
	// Create a cancellable context
	jobCtx, cancel := context.WithCancel(ctx)
	
	// Create the job
	job := &Job{
		ID:            generateID(),
		ConnectionID:  connectionID,
		SourceID:      sourceID,
		TargetID:      targetID,
		TransformerID: transformerID,
		Query:         query,
		Schedule:      schedule,
		Status:        "scheduled",
		cancel:        cancel,
	}
	
	// Store the job
	m.activeJobs[connectionID] = job
	
	// Start the job in a goroutine
	go m.runJob(jobCtx, job)
	
	return nil
}

// runJob executes a job based on its schedule
func (m *JobManager) runJob(ctx context.Context, job *Job) {
	// Parse schedule
	interval, err := parseSchedule(job.Schedule)
	if err != nil {
		job.Status = "error"
		job.Errors = append(job.Errors, fmt.Sprintf("Invalid schedule: %v", err))
		return
	}
	
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	// Execute job immediately, then on schedule
	err = m.executeTransfer(ctx, job)
	if err != nil {
		job.Errors = append(job.Errors, err.Error())
	}
	
for {
		select {
		case <-ctx.Done():
			// Context cancelled, stop the job
			job.Status = "stopped"
			return
		case <-ticker.C:
			// It's time to run the job again
			err := m.executeTransfer(ctx, job)
			if err != nil {
				job.Errors = append(job.Errors, err.Error())
				// Limit error history
				if len(job.Errors) > 10 {
					job.Errors = job.Errors[len(job.Errors)-10:]
				}
			}
		}
	}
}

// executeTransfer performs the actual data transfer
func (m *JobManager) executeTransfer(ctx context.Context, job *Job) error {
	log.Printf("Executing transfer job: %s (connection: %s)", job.ID, job.ConnectionID)
	
	// Update job status
	job.Status = "running"
	job.LastRun = time.Now()
	
	// Create execution record
	execution := models.Execution{
		ID:           generateID(),
		ConnectionID: job.ConnectionID,
		Status:       "in_progress",
		StartTime:    job.LastRun,
		CreatedAt:    job.LastRun,
	}
	err := m.dbManager.Repos.Execution().Create(execution)
	if err != nil {
		return fmt.Errorf("error creating execution record: %w", err)
	}
	
	// Update connection last run time
	m.dbManager.Repos.Connection().UpdateLastRun(job.ConnectionID, job.LastRun)
	
	// 1. Fetch data from source
	sourceData, err := m.connectorService.FetchData(ctx, job.SourceID, job.Query)
	if err != nil {
		updateExecutionError(m.dbManager, execution, err)
		job.Status = "error"
		return fmt.Errorf("error fetching data from source: %w", err)
	}
	
	// Save source data in execution record
	sourceDataJSON, _ := json.Marshal(sourceData)
	execution.SourceData = sourceDataJSON
	
	// Skip if no data
	if len(sourceData) == 0 {
		updateExecutionSuccess(m.dbManager, execution, nil)
		job.Status = "success"
		return nil
	}
	
	// 2. Transform data
  transformedData := make([]connectors.DataPayload, 0, len(sourceData))
	for _, item := range sourceData {
		transformed, err := m.transformerService.TransformData(ctx, job.TransformerID, item)
		if err != nil {
			updateExecutionError(m.dbManager, execution, err)
			job.Status = "error"
			return fmt.Errorf("error transforming data: %w", err)
		}
		transformedData = append(transformedData, transformed)
	}
	
	// 3. Push to target
	err = m.connectorService.PushData(ctx, job.TargetID, transformedData)
	if err != nil {
		updateExecutionError(m.dbManager, execution, err)
		job.Status = "error"
		return fmt.Errorf("error pushing data to target: %w", err)
	}
	
	// Update execution record with success
	targetDataJSON, _ := json.Marshal(transformedData)
	updateExecutionSuccess(m.dbManager, execution, targetDataJSON)
	
	// Update job status
	job.Status = "success"
	
	log.Printf("Transfer job completed successfully: %s", job.ID)
	return nil
}

// CancelJob cancels a running job
func (m *JobManager) CancelJob(connectionID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	job, exists := m.activeJobs[connectionID]
	if !exists {
		return fmt.Errorf("job not found for connection: %s", connectionID)
	}
	
	job.cancel()
	delete(m.activeJobs, connectionID)
	return nil
}

// GetJobStatus returns the status of a job
func (m *JobManager) GetJobStatus(connectionID string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	job, exists := m.activeJobs[connectionID]
	if !exists {
		return "", fmt.Errorf("job not found for connection: %s", connectionID)
	}
	
	return job.Status, nil
}

// ListActiveJobs returns all active jobs
func (m *JobManager) ListActiveJobs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	jobs := make([]string, 0, len(m.activeJobs))
	for id := range m.activeJobs {
		jobs = append(jobs, id)
	}
	
	return jobs
}

// InitializeFromDatabase loads and starts jobs from the database
func (m *JobManager) InitializeFromDatabase(ctx context.Context) error {
	// Get all active connections
	connections, err := m.dbManager.Repos.Connection().GetAll()
	if err != nil {
		return fmt.Errorf("error loading connections: %w", err)
	}
	
	// Start jobs for connections with schedules
	for _, conn := range connections {
		if !conn.Active || conn.Schedule == "" {
			continue
		}
		
		// Extract query from config
		var config map[string]interface{}
		err := json.Unmarshal(conn.Config, &config)
		if err != nil {
			log.Printf("Error parsing config for connection %s: %v", conn.ID, err)
			continue
		}
		
		query, _ := config["query"].(map[string]interface{})
		
		// Schedule the job
		err = m.ScheduleJob(
			ctx,
			conn.ID,
			conn.SourceID,
			conn.TargetID,
			conn.TransformerID,
			query,
			conn.Schedule,
		)
		
		if err != nil {
			log.Printf("Error scheduling job for connection %s: %v", conn.ID, err)
			continue
		}
		
		log.Printf("Scheduled job for connection: %s", conn.ID)
	}
	
	return nil
}

// Helper functions

func parseSchedule(schedule string) (time.Duration, error) {
	// Simple implementation - extend as needed
	return time.ParseDuration(schedule)
}

func updateExecutionError(dbManager *database.Manager, execution models.Execution, err error) {
	errorMsg := err.Error()
	execution.Status = "failed"
	execution.EndTime = sql.NullTime{Time: time.Now(), Valid: true}
	execution.Error = sql.NullString{String: errorMsg, Valid: true}
	
	dbManager.Repos.Execution().UpdateStatus(
		execution.ID,
		execution.Status,
		execution.EndTime.Time,
		nil,
		errorMsg,
	)
}

func updateExecutionSuccess(dbManager *database.Manager, execution models.Execution, targetData []byte) {
	execution.Status = "success"
	execution.EndTime = sql.NullTime{Time: time.Now(), Valid: true}
	
	dbManager.Repos.Execution().UpdateStatus(
		execution.ID,
		execution.Status,
		execution.EndTime.Time,
		targetData,
		"",
	)
}
