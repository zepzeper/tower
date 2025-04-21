package connection

import (
	"fmt"

	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/webapi/dto"
)

// Service provides mapping-related functionality
type Service struct {
	apiManager      *APIManager
	databaseManager *database.Manager
}

// NewService creates a new instance with the given fetcher
func NewService(databaseManager *database.Manager) *Service {
	return &Service{
		apiManager: &APIManager{
			Connections: make(map[string]*Connection),
		},
		databaseManager: databaseManager,
	}
}

// FetchConnections retrieves all connections from the database
func (s *Service) FetchConnections() ([]*Connection, error) {
	// In a real implementation, this would fetch from the database
	// For now, return connections from the in-memory store
	connections := make([]*Connection, 0, len(s.apiManager.Connections))
	for _, conn := range s.apiManager.Connections {
		connections = append(connections, conn)
	}

	return connections, nil
}

// GetConnection retrieves a connection by ID
func (s *Service) GetConnection(id string) (*Connection, error) {
	conn, exists := s.apiManager.Connections[id]
	if !exists {
		return nil, fmt.Errorf("connection not found with ID: %s", id)
	}
	return conn, nil
}

func (s *Service) TestConnection(req dto.ApiConnectionCreateRequest) (*Connection, error) {
	// Create a test connection instance from the request
	conn := &Connection{
		ID:          "test", // Temporary ID for testing
		Name:        req.Name,
		Description: "",
		Type:        ConnectionType(req.Type),
		Active:      true,
		Headers:     make(map[string]string),
		Endpoints:   make(map[string]Endpoint),
	}

	// If description is provided, use it
	if req.Description != nil {
		conn.Description = *req.Description
	}

	// Process the configuration values
	for _, config := range req.Configs {
		// Map configuration values to the appropriate connection fields
		// based on the config key
		switch config.Key {
		case "base_url":
			conn.BaseURL = config.Value
		case "timeout":
			// Add proper error handling in production code
			timeout := 30 // Default timeout
			fmt.Sscanf(config.Value, "%d", &timeout)
			conn.Timeout = timeout
		case "retry_attempts":
			// Add proper error handling in production code
			retries := 3 // Default retries
			fmt.Sscanf(config.Value, "%d", &retries)
			conn.RetryAttempts = retries
		case "auth_type":
			conn.Auth.Type = AuthType(config.Value)
		case "username":
			conn.Auth.Username = config.Value
		case "password":
			conn.Auth.Password = config.Value
		case "token":
			conn.Auth.Token = config.Value
		case "token_type":
			conn.Auth.TokenType = config.Value
		case "api_key_name":
			conn.Auth.APIKeyName = config.Value
		case "api_key_in_header":
			conn.Auth.APIKeyInHeader = (config.Value == "true")
		case "relation":
			conn.Relation = Relation(config.Value)
		default:
			// For custom headers or other configs, store in headers map
			if config.Key != "" {
				conn.Headers[config.Key] = config.Value
			}
		}
	}

	// Validate the connection - check required fields based on type
	if conn.BaseURL == "" {
		return nil, fmt.Errorf("base_url is required")
	}

	// Validate auth configuration based on auth type
	switch conn.Auth.Type {
	case BASIC:
		if conn.Auth.Username == "" || conn.Auth.Password == "" {
			return nil, fmt.Errorf("username and password are required for BASIC authentication")
		}
	case TOKEN:
		if conn.Auth.Token == "" {
			return nil, fmt.Errorf("token is required for TOKEN authentication")
		}
	case OAUTH:
		// Simplified validation for OAuth - actual implementation would need more checks
		if conn.Auth.OAuth2Config.ClientID == "" || conn.Auth.OAuth2Config.ClientSecret == "" {
			return nil, fmt.Errorf("client_id and client_secret are required for OAUTH authentication")
		}
	}

	// TODO: Attempt to make a test connection to the API
	// This would involve creating a test request to validate connectivity
	// For the sake of this example, we'll just return the connection object

	return conn, nil
}
