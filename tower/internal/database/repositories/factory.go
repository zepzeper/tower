package repositories

import (
	"database/sql"
)

// Factory provides access to all repository instances
type Factory struct {
	db               *sql.DB
	channelRepo      *ChannelRepository
	workflowRepo     *WorkflowRepository
	transformerRepo  *TransformerRepository
	executionRepo    *ExecutionRepository
	userRepo         *UserRepository
	apiKeyRepo       *APIKeyRepository
}

// NewFactory creates a new repository factory
func NewFactory(db *sql.DB) *Factory {
	return &Factory{
		db: db,
	}
}

// Channel returns a channel repository
func (f *Factory) Channel() *ChannelRepository {
	if f.channelRepo == nil {
		f.channelRepo = NewChannelRepository(f.db)
	}
	return f.channelRepo
}

// Workflow returns a workflow repository
func (f *Factory) Workflow() *WorkflowRepository {
	if f.workflowRepo == nil {
		f.workflowRepo = NewWorkflowRepository(f.db)
	}
	return f.workflowRepo
}

// Transformer returns a transformer repository
func (f *Factory) Transformer() *TransformerRepository {
	if f.transformerRepo == nil {
		f.transformerRepo = NewTransformerRepository(f.db)
	}
	return f.transformerRepo
}

// Execution returns an execution repository
func (f *Factory) Execution() *ExecutionRepository {
	if f.executionRepo == nil {
		f.executionRepo = NewExecutionRepository(f.db)
	}
	return f.executionRepo
}

// User returns a user repository
func (f *Factory) User() *UserRepository {
	if f.userRepo == nil {
		f.userRepo = NewUserRepository(f.db)
	}
	return f.userRepo
}

// APIKey returns an API key repository
func (f *Factory) APIKey() *APIKeyRepository {
	if f.apiKeyRepo == nil {
		f.apiKeyRepo = NewAPIKeyRepository(f.db)
	}
	return f.apiKeyRepo
}
