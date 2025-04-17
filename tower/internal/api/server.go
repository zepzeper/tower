package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/api/handlers"
	"github.com/zepzeper/tower/internal/api/middleware"
	"github.com/zepzeper/tower/internal/services"
)

// Server represents the HTTP server
type Server struct {
	server           *http.Server
	router           *chi.Mux
	handlers         *handlers.Registry
	connectorService *services.ConnectorService
	transformerService *services.TransformerService
	connectionService *services.ConnectionService
}

// NewServer creates a new server with the given services
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
	// Health check
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	s.router.Route("/api", func(r chi.Router) {
		// Register all API routes through the handler registry
		s.handlers.RegisterRoutes(r)
	})

	// Serve static files for web UI
	fileServer := http.FileServer(http.Dir("./web/dist"))
	s.router.Handle("/web/*", http.StripPrefix("/web", fileServer))
	
	// Redirect root to web UI
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web", http.StatusSeeOther)
	})
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("Starting server on %s", addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

