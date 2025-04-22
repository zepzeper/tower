package schema

const CredentialsTableSchema = `
  CREATE TABLE IF NOT EXISTS credentials (
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
const CredentialsConfigSchema = `
  CREATE TABLE IF NOT EXISTS credentials_configs (
  connection_id VARCHAR(50) NOT NULL REFERENCES credentials(id) ON DELETE CASCADE,
  key VARCHAR(100) NOT NULL,
  value TEXT,
  is_secret BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY (connection_id, key)
  );

  CREATE INDEX IF NOT EXISTS idx_credentials_configs_conn_id ON credentials_configs(connection_id);
  `
