package connection

import (
	"fmt"
	"github.com/zepzeper/tower/internal/database"
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

// func (s *Service) FetchConnections() ([]models.Connection, error) {
// 	conn, err := s.databaseManager.Repos.Connection().GetAll()
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed fetching connections: %s", err)
// 	}
//
// 	return conn, nil
// }

func (s *Service) GetConnection(id string) (*Connection, error) {
	conn, exists := s.apiManager.Connections[id]
	if !exists {
		return nil, fmt.Errorf("connection not found with ID: %s", id)
	}
	return conn, nil
}
