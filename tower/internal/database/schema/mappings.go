package schema

const MappingsTableSchema = `
  CREATE TABLE IF NOT EXISTS mappings (
  id VARCHAR(50) PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  source_id VARCHAR(50) NOT NULL,
  target_id VARCHAR(50) NOT NULL,
  transformer_id VARCHAR(50) NOT NULL,
  config JSONB NOT NULL DEFAULT '{}'::jsonb,
  schedule VARCHAR(100),
  active BOOLEAN NOT NULL DEFAULT TRUE,
  last_run TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_transformer FOREIGN KEY (transformer_id)
  REFERENCES transformers(id) ON DELETE RESTRICT
  );
  `
