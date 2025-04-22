package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/middleware"
	"github.com/zepzeper/tower/internal/services/authentication"
	"github.com/zepzeper/tower/internal/services/credentials"
	"github.com/zepzeper/tower/internal/services/mapping"
	"github.com/zepzeper/tower/internal/services/relations"
)

// Registry holds all the handlers for the internal API
type Registry struct {
	mappingHandler     *MappingHandler
	credentialsHandler *CredentialsHandler
	relationsHandler   *RelationsHandler
	authHandler        *AuthHandler
}

// NewRegistry creates a new handler registry
func NewRegistry(
	mappingService *mapping.Service,
	credentialsService *credentials.Service,
	relationsService *relations.Service,
	authService *auth.Service,
) *Registry {
	return &Registry{
		mappingHandler:     NewMappingHandler(mappingService),
		credentialsHandler: NewCredentialsHandler(credentialsService),
		relationsHandler:   NewRelationsHandler(relationsService),
		authHandler:        NewAuthHandler(authService),
	}
}

func (r *Registry) RegisterRoutes(router chi.Router) {
	router.Route("/api", func(routes chi.Router) {
		// Public auth routes
		routes.Route("/auth", func(routes chi.Router) {
			routes.With(middleware.POST, middleware.JSON).HandleFunc("/register", r.authHandler.Register)
			routes.With(middleware.POST, middleware.JSON).HandleFunc("/login", r.authHandler.Login)
			routes.With(middleware.POST).HandleFunc("/refresh", r.authHandler.RefreshToken)
			routes.With(middleware.POST).HandleFunc("/logout", r.authHandler.Logout)
		})

		// Protected routes that require authentication
		routes.Group(func(routes chi.Router) {
			// Apply auth middleware to this group
			routes.Use(r.authHandler.AuthMiddleware)

			// User routes
			routes.With(middleware.GET).HandleFunc("/user/me", r.authHandler.GetMe)

			// Connection routes
			routes.Route("/connections", func(routes chi.Router) {
				routes.With(middleware.GET).HandleFunc("/all", r.credentialsHandler.Fetch)
				routes.With(middleware.GET).HandleFunc("/{id}", r.credentialsHandler.GetByID)
				routes.With(middleware.GET, middleware.JSON).HandleFunc("/test/{id}", r.credentialsHandler.Test)
				routes.With(middleware.POST, middleware.JSON).HandleFunc("/", r.credentialsHandler.Create)
				routes.With(middleware.PATCH, middleware.JSON).HandleFunc("/patch/{id}", r.credentialsHandler.Update)
				routes.With(middleware.DELETE).HandleFunc("/delete/{id}", r.credentialsHandler.Delete)

				routes.Route("/{id}/relations", func(routes chi.Router) {
					routes.With(middleware.GET).HandleFunc("/", r.relationsHandler.Fetch)
					routes.With(middleware.GET).HandleFunc("/logs", r.relationsHandler.FetchLogs)
					routes.With(middleware.POST).HandleFunc("/create", r.relationsHandler.Create)
				})
			})

			// Mapping routes
			routes.Route("/mappings", func(routes chi.Router) {
				routes.With(middleware.GET).HandleFunc("/schema", r.mappingHandler.Fetch)
				routes.With(middleware.POST, middleware.JSON).HandleFunc("/test", r.mappingHandler.HandleTestMapping)
			})
		})
	})
}
