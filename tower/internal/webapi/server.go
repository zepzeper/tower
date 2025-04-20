package webapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/api/middleware"
	"github.com/zepzeper/tower/internal/services"
	"github.com/zepzeper/tower/internal/services/connection"
	"github.com/zepzeper/tower/internal/services/mapping"
	"github.com/zepzeper/tower/internal/webapi/handlers"
)

// Server represents the internal API server for the web UI
type Server struct {
	router            *chi.Mux
	handlers          *handlers.Registry
  authService       *services.AuthService
}

// NewServer creates a new WebAPI server
func NewServer(
	authService *services.AuthService,
	mappingService *mapping.Service,
	connectionService *connection.Service,
) *Server {
	r := chi.NewRouter()

	// Use the same middleware as the external API
	r.Use(middleware.Recover)
	r.Use(middleware.Logging)
	r.Use(middleware.CORS)

	// Create server
	server := &Server{
		router:            r,
    authService:       authService,  // Store the auth service
	}

	// Create handler registry
	server.handlers = handlers.NewRegistry(
    authService,
    mappingService,
    connectionService,
	)

	// Setup routes
	server.setupRoutes()

	return server
}

// setupRoutes configures the routes
func (s *Server) setupRoutes() {
	// Register all internal API routes through the handler registry
	s.handlers.RegisterRoutes(s.router)
}

// Router returns the router for use in the central server
func (s *Server) Router() *chi.Mux {
	return s.router
}
