FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Fetch dependencies if they're not in go.sum
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd/server/main.go

# Create final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/bin/server ./
# Copy UI files
COPY --from=builder /app/ui /app/ui
# Ensure the directory exists
RUN mkdir -p /app/ui/web

# Set timezone
ENV TZ=UTC

EXPOSE 8080

CMD ["./server"]
