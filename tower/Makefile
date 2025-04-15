.PHONY: build run test test-with-env clean

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Clean the binary
clean:
	rm -f bin/server

# Enhanced test command with environment loading
test-with-env: export $(shell sed 's/=.*//' .env)
test-with-env:
	@echo "Loading environment variables from .env file..."
	@go test ./... -v

# Alternative for Unix-like systems
test-env-unix:
	@echo "Loading environment variables (Unix method)..."
	@env $$(cat .env | grep -v '^#' | xargs) go test ./... -v
