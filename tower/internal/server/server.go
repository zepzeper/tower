// internal/server/server.go
package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/api"
	"github.com/zepzeper/tower/internal/services"
	"github.com/zepzeper/tower/internal/webapi"
)

// Server is the central server wrapper that manages both API and WebAPI
type Server struct {
	router            *chi.Mux
	api               *api.Server
	webapi            *webapi.Server
	httpServer        *http.Server
}

// NewServer creates a new central server with both API and WebAPI
func NewServer(
	connectorService *services.ConnectorService,
	transformerService *services.TransformerService,
	connectionService *services.ConnectionService,
	authService *services.AuthService,
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
    authService,
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

	// Host-based routing
	s.router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		
		// Check if this is the portal subdomain
		if strings.HasPrefix(host, "portal.") {
			// Serve dashboard application
			dashboardFileServer := http.FileServer(http.Dir("./portal/dist"))
			dashboardFileServer.ServeHTTP(w, r)
			return
		}
		
		// Serve main application
		frontendFileServer := http.FileServer(http.Dir("./web/dist"))
		frontendFileServer.ServeHTTP(w, r)
	})

	// Mount external API
	s.router.Mount("/api", s.api.Router())

	// Mount internal API for web UI
	s.router.Mount("/internal", s.webapi.Router())
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
