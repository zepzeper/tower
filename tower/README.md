# Tower API Middleware

A professional API middleware service for connecting and transforming data between different APIs.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose

### Running with Docker Compose

The easiest way to run the application is using Docker Compose, which will set up both the API server and the PostgreSQL database:

```bash
# Start the services
docker-compose up -d

# Check the logs
docker-compose logs -f
```

The server will be available at http://localhost:8080

### Running Locally

If you prefer to run the application locally:

1. Start a PostgreSQL instance:

```bash
docker run -d --name tower-db \
  -e POSTGRES_DB=tower \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:15-alpine
```

2. Update the `.env` file with local database configuration:

```bash
# Database configuration
DB_HOST=localhost
DB_PORT=5432
```

3. Build and run the application:

```bash
make build
make run
```

## Testing the API

You can use the provided script to test the API endpoints:

```bash
chmod +x scripts/test-api.sh
./scripts/test-api.sh
```

Alternatively, you can use curl commands manually:

```bash
# Health check
curl -X GET http://localhost:8080/health

# List channels
curl -X GET http://localhost:8080/api/v1/channels

# Create a new channel
curl -X POST http://localhost:8080/api/v1/channels \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Channel",
    "type": "webhook",
    "description": "A new webhook channel",
    "config": {
      "url": "https://example.com/webhook"
    }
  }'
```

## Project Structure

The project follows a clean architecture approach with the following structure:

- `cmd/`: Application entry points
- `internal/`: Private application code
  - `api/`: API server and handlers
  - `config/`: Configuration management
  - `core/`: Core business logic
- `ui/`: User interface
- `deployments/`: Deployment configurations
- `scripts/`: Utility scripts

## Docker Development

The Docker setup includes:

1. API service: The Go server application
2. PostgreSQL database: For data storage

To rebuild and restart services:

```bash
docker-compose down
docker-compose up --build -d
```
