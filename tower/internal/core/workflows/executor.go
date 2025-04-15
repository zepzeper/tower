package workflows

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// WorkflowData represents the data being processed in a workflow
type WorkflowData map[string]interface{}

// ExecutionResult represents the result of executing a workflow
type ExecutionResult struct {
	WorkflowID  string      `json:"workflowId"`
	Success     bool        `json:"success"`
	StartTime   time.Time   `json:"startTime"`
	EndTime     time.Time   `json:"endTime"`
	Steps       []StepResult `json:"steps"`
	Error       string      `json:"error,omitempty"`
	OutputData  WorkflowData `json:"outputData,omitempty"`
}

// StepResult represents the result of executing a single step in a workflow
type StepResult struct {
	StepID      string      `json:"stepId"`
	StepType    string      `json:"stepType"` // "trigger", "action", "transformer"
	Name        string      `json:"name"`
	Success     bool        `json:"success"`
	StartTime   time.Time   `json:"startTime"`
	EndTime     time.Time   `json:"endTime"`
	InputData   WorkflowData `json:"inputData,omitempty"`
	OutputData  WorkflowData `json:"outputData,omitempty"`
	Error       string      `json:"error,omitempty"`
}

// Executor handles workflow execution
type Executor struct {
	activeExecutions map[string]context.CancelFunc
	mu               sync.Mutex
}

// NewExecutor creates a new workflow executor
func NewExecutor() *Executor {
	return &Executor{
		activeExecutions: make(map[string]context.CancelFunc),
	}
}

// ExecuteWorkflow executes a workflow with the given data
func (e *Executor) ExecuteWorkflow(ctx context.Context, workflow *Workflow, data WorkflowData) (*ExecutionResult, error) {
	// Create a new execution context that can be cancelled
	execCtx, cancel := context.WithCancel(ctx)
	
	// Store the cancel function for potential cancellation
	executionID := fmt.Sprintf("%s-%d", workflow.ID, time.Now().UnixNano())
	e.mu.Lock()
	e.activeExecutions[executionID] = cancel
	e.mu.Unlock()

	// Ensure we clean up when done
	defer func() {
		e.mu.Lock()
		delete(e.activeExecutions, executionID)
		e.mu.Unlock()
		cancel()
	}()

	// Start tracking execution time
	result := &ExecutionResult{
		WorkflowID: workflow.ID,
		StartTime:  time.Now(),
		Steps:      []StepResult{},
	}

	// We'll use this to track the data as it flows through the workflow
	currentData := data

	// Execute each action in sequence
	for _, action := range workflow.Actions {
		// Create a step result to track this step
		stepResult := StepResult{
			StepID:    action.ID,
			StepType:  "action",
			Name:      action.Name,
			StartTime: time.Now(),
			InputData: currentData,
		}

		// Check if we've been cancelled
		select {
		case <-execCtx.Done():
			stepResult.Success = false
			stepResult.Error = "execution cancelled"
			stepResult.EndTime = time.Now()
			result.Steps = append(result.Steps, stepResult)
			result.Success = false
			result.Error = "execution cancelled"
			result.EndTime = time.Now()
			return result, errors.New("execution cancelled")
		default:
			// Continue execution
		}

		// In a real implementation, we would:
		// 1. Apply any transformations 
		// 2. Execute the action on the associated channel
		// 3. Process the results
		
		// For this example, we'll simulate successful execution
		log.Printf("Executing action: %s", action.Name)
		time.Sleep(100 * time.Millisecond) // Simulate work

		// Update the output data for the next step
		outputData := WorkflowData{
			"result": fmt.Sprintf("Processed by %s", action.Name),
			"timestamp": time.Now().Unix(),
		}
		currentData = outputData

		// Update the step result
		stepResult.Success = true
		stepResult.OutputData = outputData
		stepResult.EndTime = time.Now()
		result.Steps = append(result.Steps, stepResult)
	}

	// Workflow completed successfully
	result.Success = true
	result.EndTime = time.Now()
	result.OutputData = currentData

	return result, nil
}

// CancelExecution cancels an active workflow execution
func (e *Executor) CancelExecution(executionID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	cancel, exists := e.activeExecutions[executionID]
	if !exists {
		return errors.New("execution not found")
	}

	cancel()
	delete(e.activeExecutions, executionID)
	return nil
}

// GetActiveExecutionsCount returns the number of currently active executions
func (e *Executor) GetActiveExecutionsCount() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.activeExecutions)
}
