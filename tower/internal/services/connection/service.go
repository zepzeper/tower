package connection

import "fmt"

// Service provides mapping-related functionality
type Service struct {
    apiManager *APIManager
}

// NewService creates a new instance with the given fetcher
func NewService() *Service {
	return &Service{
        apiManager: &APIManager{
            Connections: make(map[string]*Connection),
        },
	}
}

func (s *Service) FetchConnections() map[string]*Connection {
	return s.apiManager.Connections
}

func (s *Service) GetConnection(id string) (*Connection, error) {
	conn, exists := s.apiManager.Connections[id]
	if !exists {
        return nil, fmt.Errorf("connection not found with ID: %s", id)
	}
	return conn, nil
}
