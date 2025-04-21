package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/services"
	"github.com/zepzeper/tower/internal/services/connection"
	"github.com/zepzeper/tower/internal/services/mapping"
)

// Registry holds all the handlers for the internal API
type Registry struct {
	mappingHandler    *MappingHandler
	connectionHandler *ConnectionHandler
	authHandler       *AuthHandler
}

// NewRegistry creates a new handler registry
func NewRegistry(
	authService *services.AuthService,
	mappingService *mapping.Service,
	connectionService *connection.Service,
) *Registry {
	return &Registry{
		mappingHandler:    NewMappingHandler(mappingService),
		connectionHandler: NewConnectionHandler(connectionService),
	}
}

func (r *Registry) RegisterRoutes(router chi.Router) {
	router.Route("/api", func(routes chi.Router) {
		routes.Route("/mappings", func(routes chi.Router) {
			routes.Get("/schema", r.mappingHandler.Fetch)
			routes.Post("/test", r.mappingHandler.HandleTestMapping)

			routes.Route("/connections", func(routes chi.Router) {
				routes.Get("/", r.connectionHandler.Fetch)
				routes.Post("/test", r.connectionHandler.Test)
				routes.Post("/save", r.connectionHandler.Test)
			})
		})

	})
}
