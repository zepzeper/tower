package credentials

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/database/models"
	"github.com/zepzeper/tower/internal/webapi/dto"
)

// Credentials defines the connection and its configurations
type Credentials struct {
	Credentials       *models.Credentials
	CredentialsConfig []models.CredentialsConfig
}

// Service provides connection-related functionality
type Service struct {
	databaseManager *database.Manager
}

// NewService creates a new instance with the given database manager
func NewService(databaseManager *database.Manager) *Service {
	return &Service{
		databaseManager: databaseManager,
	}
}

// FetchConnections retrieves all connections with their configurations
func (s *Service) FetchConnections() ([]*Credentials, error) {
	dbConns, err := s.databaseManager.Repos.Connection().GetAll()
	if err != nil {
		return nil, err
	}
	var result []*Credentials
	for _, dbConn := range dbConns {
		configs, err := s.databaseManager.Repos.Connection().GetConfigsByConnectionID(dbConn.ID)
		if err != nil {
			// Initialize to empty slice if error
			configs = []models.CredentialsConfig{}
		}
		conn := Credentials{
			Credentials:       &dbConn,
			CredentialsConfig: configs,
		}
		result = append(result, &conn)
	}
	return result, nil
}

// GetConnection retrieves a single connection by ID with its configurations
func (s *Service) GetConnection(id string) (*Credentials, error) {
	dbConn, err := s.databaseManager.Repos.Connection().GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("connection not found with ID: %s", id)
	}
	configs, err := s.databaseManager.Repos.Connection().GetConfigsByConnectionID(id)
	if err != nil {
		// Initialize to empty slice if error
		configs = []models.CredentialsConfig{}
	}
	conn := Credentials{
		Credentials:       dbConn,
		CredentialsConfig: configs,
	}
	return &conn, nil
}

// CreateConnection creates a new connection from the given request
func (s *Service) Create(req dto.CredentialsCreateRequest) (*Credentials, error) {
	id := uuid.New().String()

	// Create the API connection
	dbConn := &models.Credentials{
		ID:        id,
		Name:      req.Name,
		Type:      req.Type,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Handle description as sql.NullString
	if req.Description != nil {
		dbConn.Description = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	} else {
		dbConn.Description = sql.NullString{
			Valid: false,
		}
	}

	// Create configurations
	var configs []models.CredentialsConfig
	for _, config := range req.Configs {
		configs = append(configs, models.CredentialsConfig{
			ConnectionID: id,
			Key:          config.Key,
			Value:        config.Value,
		})
	}

	// Validate connection based on configs
	if err := s.validateConnectionConfigs(configs); err != nil {
		return nil, err
	}

	// Save to database
	if err := s.databaseManager.Repos.Connection().Create(dbConn, configs); err != nil {
		return nil, err
	}

	// Return the created connection
	return &Credentials{
		Credentials:       dbConn,
		CredentialsConfig: configs,
	}, nil
}

// UpdateConnection updates an existing connection
func (s *Service) UpdateConnection(req dto.CredentialsUpdateRequest) (*Credentials, error) {
	// Get the current connection
	credentials, err := s.GetConnection(req.ID)
	if err != nil {
		return nil, err
	}

	dbConn := credentials.Credentials
	configs := credentials.CredentialsConfig

	// Update basic fields
	if req.Name != nil {
		dbConn.Name = *req.Name
	}
	if req.Description != nil {
		dbConn.Description = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	}
	if req.Active != nil {
		dbConn.Active = *req.Active
	}

	// Update configs if provided
	if req.Configs != nil {
		// Create a map of existing configs for easy lookup
		configMap := make(map[string]int)
		for i, config := range configs {
			configMap[config.Key] = i
		}

		// Update or add new configs
		for _, newConfig := range req.Configs {
			if idx, exists := configMap[newConfig.Key]; exists {
				// Update existing config
				configs[idx].Value = newConfig.Value
			} else {
				// Add new config
				configs = append(configs, models.CredentialsConfig{
					ConnectionID: dbConn.ID,
					Key:          newConfig.Key,
					Value:        newConfig.Value,
				})
			}
		}
	}

	// Update timestamp
	dbConn.UpdatedAt = time.Now()

	// Validate connection based on configs
	if err := s.validateConnectionConfigs(configs); err != nil {
		return nil, err
	}

	// Save to database
	if err := s.databaseManager.Repos.Connection().Update(dbConn, configs); err != nil {
		return nil, err
	}

	// Return the updated connection
	return &Credentials{
		Credentials:       dbConn,
		CredentialsConfig: configs,
	}, nil
}

// DeleteConnection removes a connection by ID
func (s *Service) DeleteConnection(id string) error {
	return s.databaseManager.Repos.Connection().Delete(id)
}

// validateConnectionConfigs validates connection configurations
func (s *Service) validateConnectionConfigs(configs []models.CredentialsConfig) error {
	// Extract authentication type and relevant configs
	var authType string
	var username, password, token, baseURL string
	var clientID, clientSecret string

	for _, config := range configs {
		switch config.Key {
		case "auth_type":
			authType = config.Value
		case "username":
			username = config.Value
		case "password":
			password = config.Value
		case "token":
			token = config.Value
		case "base_url":
			baseURL = config.Value
		case "client_id":
			clientID = config.Value
		case "client_secret":
			clientSecret = config.Value
		}
	}

	// Validate based on auth type
	switch authType {
	case "BASIC":
		if username == "" || password == "" {
			return fmt.Errorf("username and password are required for BASIC authentication")
		}
	case "TOKEN":
		if baseURL == "" {
			return fmt.Errorf("base_url is required")
		}
		if token == "" {
			return fmt.Errorf("token is required for TOKEN authentication")
		}
	case "OAUTH":
		if clientID == "" || clientSecret == "" {
			return fmt.Errorf("client_id and client_secret are required for OAUTH authentication")
		}
	}

	return nil
}
