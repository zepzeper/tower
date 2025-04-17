package schema

// ExecutionsTableSchema defines the SQL schema for the executions table
const ExecutionsTableSchema = `
    CREATE TABLE IF NOT EXISTS executions (
        id VARCHAR(50) PRIMARY KEY,
        connection_id VARCHAR(50) NOT NULL REFERENCES connections(id) ON DELETE CASCADE,
        status VARCHAR(20) NOT NULL CHECK (status IN ('success', 'failed', 'in_progress')),
        start_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
        end_time TIMESTAMP WITH TIME ZONE,
        source_data JSONB,
        target_data JSONB,
        error TEXT,
        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

    CREATE INDEX IF NOT EXISTS idx_executions_connection_id ON executions(connection_id);
    CREATE INDEX IF NOT EXISTS idx_executions_status ON executions(status);
`
