// internal/server/server.go
package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/api"
	"github.com/zepzeper/tower/internal/services"
	"github.com/zepzeper/tower/internal/webapi"
)

// Server is the central server wrapper that manages both API and WebAPI
type Server struct {
	router           *chi.Mux
	api              *api.Server
	webapi           *webapi.Server
	httpServer       *http.Server
}

// NewServer creates a new central server with both API and WebAPI
func NewServer(
	connectorService *services.ConnectorService,
	transformerService *services.TransformerService,
	connectionService *services.ConnectionService,
) *Server {
	// Create main router
	router := chi.NewRouter()

	// Create API server (external API)
	apiServer := api.NewServer(
		connectorService,
		transformerService,
		connectionService,
	)

	// Create WebAPI server (internal API for web UI)
	webapiServer := webapi.NewServer(
		connectorService,
		transformerService,
		connectionService,
	)

	// Create central server
	server := &Server{
		router: router,
		api:    apiServer,
		webapi: webapiServer,
	}

	// Setup routes
	server.setupRoutes()

	return server
}

// setupRoutes configures the routes for both API and WebAPI
func (s *Server) setupRoutes() {
	// Health check
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Mount external API
	s.router.Mount("/api", s.api.Router())

	// Mount internal API for web UI
	s.router.Mount("/internal", s.webapi.Router())

	fileServer := http.FileServer(http.Dir("./web/dist"))
	s.router.Handle("/*", fileServer)
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	log.Printf("Starting server on %s", addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
