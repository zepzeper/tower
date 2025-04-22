package schema

// AllSchemas returns all database schema definitions in the order they should be applied
func AllSchemas() []string {
	return []string{
		TransformersTableSchema,
		CredentialsTableSchema,
		CredentialsConfigSchema,
		CredentialsRelations,
		CredentialsRelationsLogs,
		MappingsTableSchema,
		ExecutionsTableSchema,
		UsersTableSchema,
	}
}
