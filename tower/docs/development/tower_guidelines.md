# Tower Development Guide

This document provides guidance for developers working on the Tower API Middleware project. It covers setup, development workflows, coding standards, and how to extend the system with new connectors.

## Development Environment Setup

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL (local or Docker)
- Git

### Getting Started

1. Clone the repository:

```bash
git clone https://github.com/zepzeper/tower.git
cd tower
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables for local development (create a `.env` file):

```
# Server settings
SERVER_PORT=8080

# Database settings
DB_HOST=localhost
DB_PORT=5432
DB_NAME=tower
DB_USER=postgres
DB_PASSWORD=postgres
```

4. Start the development environment:

```bash
# Using Docker Compose for development
./scripts/start-dev.sh

# OR manually start a PostgreSQL instance
docker run -d --name tower-db \
  -e POSTGRES_DB=tower \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:15-alpine
```

5. Run the application:

```bash
# Using Air for hot reloading (install Air first)
air

# OR using the Go command
go run cmd/server/main.go
```

## Development Workflow

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/services

# Run tests with coverage
go test -cover ./...
```

### Building the Application

```bash
# Build for local platform
go build -o bin/server cmd/server/main.go

# Build for a specific platform
GOOS=linux GOARCH=amd64 go build -o bin/server-linux-amd64 cmd/server/main.go
```

### Docker Development

```bash
# Build and start all services
docker-compose up --build

# Run in background
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f
```

## Code Structure and Standards

### Project Structure

- `cmd/` - Application entry points
- `internal/` - Private application code
  - `api/` - API handlers and server
  - `config/` - Configuration management
  - `core/` - Core business logic
  - `connectors/` - Connector implementations
  - `database/` - Database access
  - `services/` - Business logic services
- `schemas/` - JSON schema definitions

### Coding Standards

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Format code with `gofmt` or `go fmt`
- Use `golint` and `go vet` to check for code issues
- Write comments for exported functions, types, and packages
- Write unit tests for business logic

### Error Handling

- Use error wrapping with `fmt.Errorf("some context: %w", err)`
- Return errors rather than using panic
- Check all errors and provide context for debugging
- Use structured logging for error reporting

### Logging

- Use the standard `log` package or a structured logging library
- Include relevant context in log messages
- Use appropriate log levels (debug, info, warn, error)
- Include timestamps and request IDs in logs

## Adding a New Connector

To add a new connector to Tower, follow these steps:

1. Create a new directory under `internal/connectors/` with your connector name
2. Implement the Connector interface
3. Register the connector in `cmd/server/main.go`

### 1. Create Directory Structure

```
internal/connectors/myconnector/
├── connector.go
└── schema.go
```

### 2. Implement the Connector Interface

#### connector.go

```go
package myconnector

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    
    "github.com/zepzeper/tower/internal/core/connectors"
)

// Connector implements the Connector interface for MyAPI
type Connector struct {
    connectors.BaseConnector
    ApiKey string
}

// NewConnector creates a new MyAPI connector
func NewConnector(config map[string]interface{}) (*Connector, error) {
    // Extract configuration
    apiURL, ok := config["api_url"].(string)
    if !ok {
        return nil, fmt.Errorf("missing or invalid api_url in configuration")
    }
    
    apiKey, ok := config["api_key"].(string)
    if !ok {
        return nil, fmt.Errorf("missing or invalid api_key in configuration")
    }
    
    // Set up headers
    headers := map[string]string{
        "Content-Type": "application/json",
        "Authorization": "Bearer " + apiKey,
    }
    
    // Create the connector
    connector := &Connector{
        BaseConnector: connectors.NewBaseConnector(apiURL, headers, GetSchema()),
        ApiKey: apiKey,
    }
    
    return connector, nil
}

// Connect implements the Connector interface
func (c *Connector) Connect(ctx context.Context) error {
    // Implement connection test
    _, err := c.MakeRequest(ctx, "GET", "ping", nil)
    return err
}

// Fetch implements the Connector interface
func (c *Connector) Fetch(ctx context.Context, query map[string]interface{}) ([]connectors.DataPayload, error) {
    // Implement data fetching logic
    // ...
    return items, nil
}

// Push implements the Connector interface
func (c *Connector) Push(ctx context.Context, data []connectors.DataPayload) error {
    // Implement data pushing logic
    // ...
    return nil
}
```

#### schema.go

```go
package myconnector

import "github.com/zepzeper/tower/internal/core/connectors"

// GetSchema returns the schema for MyAPI
func GetSchema() connectors.Schema {
    return connectors.Schema{
        EntityName: "product",
        Fields: map[string]connectors.FieldDefinition{
            "id": {
                Type:     "string",
                Required: false,
                Path:     "id",
            },
            "name": {
                Type:     "string",
                Required: true,
                Path:     "name",
            },
            "price": {
                Type:     "number",
                Required: true,
                Path:     "price",
            },
            // Add more field definitions
        },
    }
}
```

### 3. Register the Connector

In `cmd/server/main.go`, add your connector to the registration function:

```go
func registerConnectors(registry *registry.ConnectorRegistry) {
    // Existing connectors...
    
    // Register MyAPI connector
    myConnector, err := myconnector.NewConnector(map[string]interface{}{
        "api_url": "https://api.example.com/v1",
        "api_key": "default_key", // Default value, will be overridden in real use
    })
    if err == nil {
        registry.Register("myconnector", myConnector)
    }
}
```
