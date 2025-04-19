package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/services"
	"github.com/zepzeper/tower/internal/services/mapping"
)

// Registry holds all the handlers for the internal API
type Registry struct {
	mappingHandler *MappingHandler
	authHandler *AuthHandler
}

// NewRegistry creates a new handler registry
func NewRegistry(
	authService *services.AuthService,
	mappingService *mapping.Service,
) *Registry {
	return &Registry{
		mappingHandler: NewMappingHandler(mappingService),
	}
}

func (r *Registry) RegisterRoutes(router chi.Router) {
	router.Route("/api/mappings", func(routes chi.Router) {
		routes.Get("/schema", r.mappingHandler.Fetch)
		routes.Post("/test", r.mappingHandler.HandleTestMapping)
	})
}
