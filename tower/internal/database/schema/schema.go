package schema

// AllSchemas returns all database schema definitions in the order they should be applied
func AllSchemas() []string {
	return []string{
		UsersTableSchema,         // Users must be created first due to foreign key references
		APIKeysTableSchema,       // API keys reference users
		ChannelsTableSchema,      // Independent table
		TransformersTableSchema,  // Independent table
		WorkflowsTableSchema,     // Independent table
		ExecutionsTableSchema,    // Executions reference workflows
	}
}
