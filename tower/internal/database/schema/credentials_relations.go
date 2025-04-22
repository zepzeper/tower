package schema

const CredentialsRelations = `
  DO $$ BEGIN
    CREATE TYPE connection_type AS ENUM ('inbound', 'outbound', 'bidirectional');
  EXCEPTION
    WHEN duplicate_object THEN null;
  END $$;
  
  CREATE TABLE IF NOT EXISTS credentials_connections (
    initiator_id VARCHAR(50) NOT NULL,
    target_id VARCHAR(50) NOT NULL,
    connection_type connection_type NOT NULL,
    active bool NOT NULL,
    endpoint VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (initiator_id, target_id),
    FOREIGN KEY (initiator_id) REFERENCES credentials(id) ON DELETE CASCADE,
    FOREIGN KEY (target_id) REFERENCES credentials(id) ON DELETE CASCADE
  );
`

const CredentialsRelationsLogs = `
  CREATE TABLE IF NOT EXISTS credentials_connections_logs (
    id SERIAL PRIMARY KEY,
    initiator_id VARCHAR(50) NOT NULL,
    target_id VARCHAR(50) NOT NULL,
    connection_type connection_type NOT NULL,
    message VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (initiator_id) REFERENCES credentials(id) ON DELETE CASCADE,
    FOREIGN KEY (target_id) REFERENCES credentials(id) ON DELETE CASCADE
  );
`
