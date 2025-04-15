package workflows

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// TriggerManager manages workflow triggers
type TriggerManager struct {
	workflows    map[string]*Workflow
	executor     *Executor
	subscriptions map[string][]string // Maps trigger types to workflow IDs
	mu           sync.RWMutex
}

// NewTriggerManager creates a new trigger manager
func NewTriggerManager(executor *Executor) *TriggerManager {
	return &TriggerManager{
		workflows:     make(map[string]*Workflow),
		executor:      executor,
		subscriptions: make(map[string][]string),
	}
}

// RegisterWorkflow registers a workflow with the trigger manager
func (tm *TriggerManager) RegisterWorkflow(workflow *Workflow) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Store the workflow
	tm.workflows[workflow.ID] = workflow

	// Register all triggers
	for _, trigger := range workflow.Triggers {
		triggerKey := fmt.Sprintf("%s:%s", trigger.ChannelID, trigger.Event)
		
		// Create the subscription list if it doesn't exist
		if _, exists := tm.subscriptions[triggerKey]; !exists {
			tm.subscriptions[triggerKey] = []string{}
		}
		
		// Add this workflow to the subscriptions
		tm.subscriptions[triggerKey] = append(tm.subscriptions[triggerKey], workflow.ID)
	}

	log.Printf("Registered workflow %s with %d triggers", workflow.ID, len(workflow.Triggers))
}

// UnregisterWorkflow removes a workflow from the trigger manager
func (tm *TriggerManager) UnregisterWorkflow(workflowID string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	workflow, exists := tm.workflows[workflowID]
	if !exists {
		return
	}

	// Remove all trigger subscriptions
	for _, trigger := range workflow.Triggers {
		triggerKey := fmt.Sprintf("%s:%s", trigger.ChannelID, trigger.Event)
		
		if workflows, exists := tm.subscriptions[triggerKey]; exists {
			// Filter out this workflow ID
			updatedWorkflows := []string{}
			for _, id := range workflows {
				if id != workflowID {
					updatedWorkflows = append(updatedWorkflows, id)
				}
			}
			
			if len(updatedWorkflows) > 0 {
				tm.subscriptions[triggerKey] = updatedWorkflows
			} else {
				delete(tm.subscriptions, triggerKey)
			}
		}
	}

	// Remove the workflow
	delete(tm.workflows, workflowID)
	log.Printf("Unregistered workflow %s", workflowID)
}

// HandleTrigger processes a trigger event and executes corresponding workflows
func (tm *TriggerManager) HandleTrigger(ctx context.Context, channelID, event string, data WorkflowData) {
	tm.mu.RLock()
	triggerKey := fmt.Sprintf("%s:%s", channelID, event)
	workflowIDs, exists := tm.subscriptions[triggerKey]
	
	if !exists || len(workflowIDs) == 0 {
		tm.mu.RUnlock()
		log.Printf("No workflows subscribed to trigger %s", triggerKey)
		return
	}
	
	// Get the workflows we need to execute
	workflows := make([]*Workflow, 0, len(workflowIDs))
	for _, id := range workflowIDs {
		if workflow, ok := tm.workflows[id]; ok && workflow.Active {
			workflows = append(workflows, workflow)
		}
	}
	tm.mu.RUnlock()

	// Execute each workflow
	for _, workflow := range workflows {
		go func(wf *Workflow) {
			log.Printf("Executing workflow %s in response to trigger %s", wf.ID, triggerKey)
			result, err := tm.executor.ExecuteWorkflow(ctx, wf, data)
			if err != nil {
				log.Printf("Error executing workflow %s: %v", wf.ID, err)
				return
			}
			
			if result.Success {
				log.Printf("Workflow %s executed successfully", wf.ID)
			} else {
				log.Printf("Workflow %s failed: %s", wf.ID, result.Error)
			}
		}(workflow)
	}
}
