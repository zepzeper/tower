package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/api/handlers"
	"github.com/zepzeper/tower/internal/api/middleware"
)

// Server represents the external API server
type Server struct {
	router            *chi.Mux
	handlers          *handlers.Registry
}

// NewServer creates a new API server
func NewServer(
) *Server {
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Recover)
	r.Use(middleware.Logging)
	r.Use(middleware.CORS)

	// Create server
	server := &Server{
		router:            r,
	}

	// Create handler registry
	server.handlers = handlers.NewRegistry(
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
