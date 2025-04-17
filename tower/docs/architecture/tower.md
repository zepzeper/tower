# Tower Architecture Documentation

This document describes the architecture of the Tower API Middleware system, which is designed to connect and transform data between different APIs.

## Overview

Tower is built with a clean, layered architecture that separates concerns and makes the system maintainable and extensible. The system is designed to handle connections between hundreds of different APIs with automated field mapping and transformation.

## System Architecture

### High-Level Architecture

![Tower Architecture](architecture-diagram.png)

Tower is organized into the following layers:

1. **API Layer** - REST API endpoints exposed to clients
2. **Service Layer** - Business logic for all operations
3. **Core Layer** - Core domain models and interfaces
4. **Data Layer** - Database access and persistence

### Key Components

#### Connector System

The connector system provides a unified interface for connecting to various external APIs.

- **Connector Interface** - Common interface all API connectors implement
- **Base Connector** - Shared functionality for API connectors
- **Connector Registry** - Central registry of all available connectors
- **Schema System** - Schema definition and discovery tools

#### Transformation System

The transformation system handles mapping and transforming data between different API formats.

- **Transformer** - Core component that handles data transformation
- **Field Mapping** - Defines how fields map between source and target
- **AutoMapper** - Automatically generates mappings based on field similarity
- **Transformation Functions** - Built-in functions for data manipulation

#### Job Management

The job management system handles scheduled and on-demand data transfers.

- **Job Manager** - Manages scheduled jobs
- **Job Scheduler** - Schedules and tracks job execution
- **Execution Tracker** - Records execution history and results

## Directory Structure

```
tower/
├── cmd/
│   └── server/           # Application entry points
│       └── main.go
├── internal/
│   ├── api/              # API endpoints
│   │   ├── dto/          # Data Transfer Objects
│   │   ├── handlers/     # Request handlers
│   │   ├── middleware/   # HTTP middleware
│   │   ├── response/     # Response formatting
│   │   └── server.go     # HTTP server
│   ├── config/           # Configuration management
│   │   └── config.go
│   ├── core/             # Core domain logic
│   │   ├── connectors/   # Connector interfaces
│   │   ├── registry/     # Component registries
│   │   └── transformers/ # Transformation logic
│   ├── connectors/       # Connector implementations
│   │   ├── woocommerce/
│   │   ├── mirakl/
│   │   └── kaufland/
│   ├── database/         # Data access
│   │   ├── models/       # Database models
│   │   └── repositories/ # Data repositories
│   └── services/         # Business logic services
│       ├── connector_service.go
│       ├── transformer_service.go
│       ├── connection_service.go
│       └── job_service.go
├── schemas/              # JSON schema definitions
├── docs/                 # Documentation
└── web/                  # Web UI assets
```

## Data Flow

A typical data flow through the system:

1. Client makes a request to create a connection
2. Connection service validates the request
3. Connection is stored in the database
4. If scheduled, job manager creates a scheduled job
5. When executing:
   - Data is fetched from the source connector
   - Transformer applies field mappings and functions
   - Transformed data is pushed to the target connector
   - Execution results are recorded

## Key Interfaces

### Connector Interface

```go
type Connector interface {
    // Connect establishes the connection with the API
    Connect(ctx context.Context) error
    
    // Fetch retrieves data from the API
    Fetch(ctx context.Context, query map[string]interface{}) ([]DataPayload, error)
    
    // Push sends data to the API
    Push(ctx context.Context, data []DataPayload) error
    
    // GetSchema returns the data schema for this connector
    GetSchema() Schema
}
```

### Transformer Interface

```go
type Transformer struct {
    ID          string
    Name        string
    Description string
    Mappings    []FieldMapping
    Functions   []Function
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type FieldMapping struct {
    SourceField string
    TargetField string
}

type Function struct {
    Name        string
    TargetField string
    Args        []string
}

// Transform applies mappings and functions to input data
func (t *Transformer) Transform(input DataPayload) (DataPayload, error)
```

## Database Schema

The database schema consists of three main tables:

### Transformers

Stores data transformation definitions.

- `id` - Primary key
- `name` - Transformer name
- `description` - Optional description
- `mappings` - JSON array of field mappings
- `functions` - JSON array of transformation functions
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Connections

Stores connections between source and target systems.

- `id` - Primary key
- `name` - Connection name
- `description` - Optional description
- `source_id` - Source connector ID
- `target_id` - Target connector ID
- `transformer_id` - Reference to transformer
- `config` - JSON configuration including query parameters
- `schedule` - Schedule expression (e.g., "1h")
- `active` - Whether the connection is active
- `last_run` - Last execution timestamp
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Executions

Stores execution history.

- `id` - Primary key
- `connection_id` - Reference to connection
- `status` - Execution status (success, failed, in_progress)
- `start_time` - Start timestamp
- `end_time` - End timestamp (if completed)
- `source_data` - JSON data fetched from source
- `target_data` - JSON data sent to target
- `error` - Error message (if failed)
- `created_at` - Creation timestamp

## Extension Points

The Tower system is designed to be extensible. The main extension points are:

1. **New Connectors** - Implement the Connector interface to add support for new APIs
2. **Transformation Functions** - Add new functions to enhance data transformation capabilities
3. **Authentication Mechanisms** - Customize authentication for different security requirements
4. **Schedulers** - Customize job scheduling for different deployment scenarios

## Security Considerations

1. **API Authentication** - Implement appropriate authentication mechanisms
2. **Credential Storage** - Store API credentials securely, preferably using a secrets manager
3. **Data Encryption** - Encrypt sensitive data in transit and at rest
4. **Access Control** - Implement proper authorization for all operations
5. **Rate Limiting** - Add rate limiting to prevent abuse
6. **Input Validation** - Validate all input to prevent injection attacks

## Performance Considerations

1. **Database Indexing** - Ensure proper indexing for frequently queried fields
2. **Connection Pooling** - Implement connection pooling for database and API connections
3. **Caching** - Consider caching frequently accessed data
4. **Batch Processing** - Process data in batches for better throughput
5. **Monitoring** - Implement monitoring to identify performance bottlenecks
6. **Horizontal Scaling** - Design for horizontal scaling in high-load scenarios
