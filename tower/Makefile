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
