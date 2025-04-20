package schema

const ConnectionsTableSchema = `
  CREATE TABLE IF NOT EXISTS connections (
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

const APIConnectionsTableSchema = `
  CREATE TABLE IF NOT EXISTS api_connections (
  id VARCHAR(50) PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  type VARCHAR(20) NOT NULL,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
  );
  `

// APIConnectionConfigSchema defines the SQL schema for connection configurations
const APIConnectionConfigSchema = `
  CREATE TABLE IF NOT EXISTS api_connection_configs (
  connection_id VARCHAR(50) NOT NULL REFERENCES api_connections(id) ON DELETE CASCADE,
  key VARCHAR(100) NOT NULL,
  value TEXT,
  is_secret BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY (connection_id, key)
  );

  CREATE INDEX IF NOT EXISTS idx_api_connection_configs_conn_id ON api_connection_configs(connection_id);
  `
