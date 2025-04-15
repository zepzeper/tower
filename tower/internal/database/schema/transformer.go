package schema

// TransformersTableSchema defines the SQL schema for the transformers table
const TransformersTableSchema = `
CREATE TABLE IF NOT EXISTS transformers (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    mappings JSONB NOT NULL DEFAULT '[]'::jsonb,
    functions JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
`
