package handlers

import (
	"net/http"

	"github.com/zepzeper/tower/internal/response"
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

func (h *ConnectionHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

  h.connectionService.FetchConnections()


}
