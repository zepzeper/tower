package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/api/handlers"
	"github.com/zepzeper/tower/internal/api/middleware"
	"github.com/zepzeper/tower/internal/services"
)

// Server represents the external API server
type Server struct {
	router            *chi.Mux
	handlers          *handlers.Registry
	connectorService  *services.ConnectorService
	transformerService *services.TransformerService
	connectionService *services.ConnectionService
}

// NewServer creates a new API server
func NewServer(
	connectorService *services.ConnectorService,
	transformerService *services.TransformerService,
	connectionService *services.ConnectionService,
) *Server {
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Recover)
	r.Use(middleware.Logging)
	r.Use(middleware.CORS)

	// Create server
	server := &Server{
		router:            r,
		connectorService:  connectorService,
		transformerService: transformerService,
		connectionService: connectionService,
	}

	// Create handler registry
	server.handlers = handlers.NewRegistry(
		connectorService,
		transformerService,
		connectionService,
	)

	// Setup routes
	server.setupRoutes()

	return server
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	// Register all API routes through the handler registry
	s.handlers.RegisterRoutes(s.router)
}

// Router returns the router for use in the central server
func (s *Server) Router() *chi.Mux {
	return s.router
}
