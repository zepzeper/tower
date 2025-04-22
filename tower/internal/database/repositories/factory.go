package repositories

import (
	"database/sql"
)

// Factory provides access to all repository instances
type Factory struct {
	db              *sql.DB
	transformerRepo *TransformerRepository
	credentialsRepo *CredentialsRepository
	relationsRepo   *RelationsRepository
	executionRepo   *ExecutionRepository
	authRepo        *AuthRepository
}

// NewFactory creates a new repository factory
func NewFactory(db *sql.DB) *Factory {
	return &Factory{
		db: db,
	}
}

// Transformer returns a transformer repository
func (f *Factory) Transformer() *TransformerRepository {
	if f.transformerRepo == nil {
		f.transformerRepo = NewTransformerRepository(f.db)
	}
	return f.transformerRepo
}

// Connection returns a connection repository
func (f *Factory) Connection() *CredentialsRepository {
	if f.credentialsRepo == nil {
		f.credentialsRepo = NewCredentialsRepository(f.db)
	}
	return f.credentialsRepo
}

// Relations return a relation repository
func (f *Factory) Relations() *RelationsRepository {
	if f.relationsRepo == nil {
		f.relationsRepo = NewRelationsRepository(f.db)
	}
	return f.relationsRepo
}

// Execution returns an execution repository
func (f *Factory) Execution() *ExecutionRepository {
	if f.executionRepo == nil {
		f.executionRepo = NewExecutionRepository(f.db)
	}
	return f.executionRepo
}

// Execution returns an execution repository
func (f *Factory) Auth() *AuthRepository {
	if f.authRepo == nil {
		f.authRepo = NewAuthRepository(f.db)
	}
	return f.authRepo
}
