package workflows

// Workflow represents a data processing workflow
type Workflow struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Triggers    []Trigger  `json:"triggers"`
	Actions     []Action   `json:"actions"`
	Active      bool       `json:"active"`
}

// Trigger represents an event that can start a workflow
type Trigger struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	ChannelID string                 `json:"channelId"`
	Event     string                 `json:"event"`
	Config    map[string]interface{} `json:"config"`
}

// Action represents an operation in a workflow
type Action struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	ChannelID     string                 `json:"channelId"`
	Operation     string                 `json:"operation"`
	Config        map[string]interface{} `json:"config"`
	TransformerID string                 `json:"transformerId,omitempty"`
}

// ActionHandler defines the interface for executing workflow actions
type ActionHandler interface {
	Execute(action Action, data WorkflowData) (WorkflowData, error)
}

// ActionRegistry manages action handlers for different channel types
type ActionRegistry struct {
	handlers map[string]ActionHandler
}

// NewActionRegistry creates a new action registry
func NewActionRegistry() *ActionRegistry {
	return &ActionRegistry{
		handlers: make(map[string]ActionHandler),
	}
}

// RegisterHandler registers an action handler for a channel type
func (r *ActionRegistry) RegisterHandler(channelType string, handler ActionHandler) {
	r.handlers[channelType] = handler
}

// GetHandler returns an action handler for a channel type
func (r *ActionRegistry) GetHandler(channelType string) (ActionHandler, bool) {
	handler, exists := r.handlers[channelType]
	return handler, exists
}

// DefaultActionHandler is a basic implementation that simulates action execution
type DefaultActionHandler struct{}

// Execute simulates executing an action
func (h *DefaultActionHandler) Execute(action Action, data WorkflowData) (WorkflowData, error) {
	// In a real implementation, this would perform the actual operation
	// For now, we'll just echo the input with an added field
	result := WorkflowData{}
	
	// Copy the input data
	for k, v := range data {
		result[k] = v
	}
	
	// Add an action identifier
	result["_action"] = action.ID
	result["_operation"] = action.Operation
	
	return result, nil
}
