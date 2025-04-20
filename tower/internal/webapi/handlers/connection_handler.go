package handlers

import (
	"github.com/zepzeper/tower/internal/services/connection"
)

type ConnectionHandler struct {
    connectionService *connection.Service    
}

func NewConnectionHandler(connectionService *connection.Service) *ConnectionHandler {
    return &ConnectionHandler{
        connectionService: connectionService,
    }
}
