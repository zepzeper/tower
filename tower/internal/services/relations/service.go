package relations

import (
	"fmt"

	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/webapi/dto"
)

type Service struct {
	databaseManager *database.Manager
}

func NewService(databaseManager *database.Manager) *Service {
	return &Service{
		databaseManager: databaseManager,
	}
}

func (s *Service) Fetch(initiatorId string) ([]dto.CredentialsRelationResponse, error) {
	relations, err := s.databaseManager.Repos.Relations().GetByInitiatorID(initiatorId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch relations: %w", err)
	}

	return relations, nil
}

func (s *Service) FetchLogs(initiatorId string) ([]dto.CredentialsRelationLogsResponse, error) {
	relations, err := s.databaseManager.Repos.Relations().GetLogsByInitiatorID(initiatorId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch relations: %w", err)
	}

	return relations, nil
}

func (s *Service) Create(req dto.CredentialsRelationCreateRequest) error {
	// Validate that initiator and target are not the same
	if req.InitiatorID == req.TargetID {
		return fmt.Errorf("initiator and target cannot be the same credential")
	}

	// Create the connection
	err := s.databaseManager.Repos.Relations().Create(req)
	if err != nil {
		return fmt.Errorf("failed to create connection: %w", err)
	}

	return nil
}
