#!/bin/bash
# Script to set up the API middleware project structure

set -e

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go first."
    exit 1
fi

# Project name
PROJECT_NAME="api-middleware"

# Create root directory
mkdir -p $PROJECT_NAME
cd $PROJECT_NAME

# Initialize Go module
go mod init github.com/yourusername/$PROJECT_NAME

# Create directory structure
mkdir -p cmd/server
mkdir -p internal/api/{handlers,middleware,routes}
mkdir -p internal/core/{channels,workflows,transformers}
mkdir -p internal/core/channels/{models,facebook,twitter,shopify}
mkdir -p internal/auth
mkdir -p internal/config
mkdir -p internal/database/{models,migrations,repositories}
mkdir -p internal/services/{cache,queue,storage,metrics}
mkdir -p pkg/{apiclient,validator,logger,errors}
mkdir -p ui/{web,assets}
mkdir -p scripts
mkdir -p deployments/{docker,kubernetes,terraform}
mkdir -p test/{integration,e2e}
mkdir -p docs/{api,development,architecture}

# Create basic files
touch cmd/server/main.go
touch internal/api/server.go
touch internal/core/channels/models/channel.go
touch internal/core/channels/registry.go
touch internal/core/workflows/executor.go
touch internal/core/workflows/triggers.go
touch internal/core/workflows/actions.go
touch internal/core/transformers/mapping.go
touch internal/core/transformers/functions.go
touch internal/auth/middleware.go
touch internal/auth/jwt.go
touch internal/config/config.go
touch internal/config/environment.go
touch pkg/logger/logger.go
touch pkg/errors/errors.go
touch Makefile
touch .gitignore
touch README.md

# Create Docker files
touch deployments/docker/Dockerfile
touch deployments/docker/docker-compose.yml

# Create basic Makefile
cat > Makefile << 'EOF'
.PHONY: build run test clean

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Test the application
test:
	go test ./...

# Clean the binary
clean:
	rm -f bin/server
EOF

# Create a basic README
cat > README.md << 'EOF'
# API Middleware

A professional API middleware service for connecting and transforming data between different APIs.

## Features

- Connect multiple API channels
- Transform data between APIs
- Set up rules for data processing
- Monitor data flow

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Installation

```bash
git clone https://github.com/yourusername/api-middleware.git
cd api-middleware
make build
```

### Running the Application

```bash
make run
```

## Project Structure

The project follows a clean architecture approach with the following structure:

- `cmd/`: Application entry points
- `internal/`: Private application code
- `pkg/`: Public libraries
- `ui/`: User interface
- `deployments/`: Deployment configurations
- `docs/`: Documentation
EOF

# Create a basic .gitignore
cat > .gitignore << 'EOF'
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# Environment variables
.env

# IDE files
.idea/
.vscode/

# Logs
*.log
EOF

# Write a basic Dockerfile
cat > deployments/docker/Dockerfile << 'EOF'
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server /app/cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/server .

EXPOSE 8080

CMD ["./server"]
EOF

# Write a basic docker-compose.yml
cat > deployments/docker/docker-compose.yml << 'EOF'
version: '3.8'

services:
  api:
    build:
      context: ../..
      dockerfile: deployments/docker/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - ENV=development
    restart: unless-stopped
EOF

echo "Project structure created successfully!"
