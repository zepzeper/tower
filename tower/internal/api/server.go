package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zepzeper/tower/internal/api/handlers"
	towermiddleware "github.com/zepzeper/tower/internal/api/middleware"
	"github.com/zepzeper/tower/internal/config"
	"github.com/zepzeper/tower/internal/database"
)

// Server represents the HTTP server
type Server struct {
	server  *http.Server
	router  *chi.Mux
	config  *config.Config
	db      *db.Manager
	handlers *handlers.Registry
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config) (*Server, error) {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(towermiddleware.CORS)

	// Connect to database
	dbManager, err := db.NewManager(cfg)
	if err != nil {
		return nil, err
	}

	// Create server
	s := &Server{
		router: r,
		config: cfg,
		db:     dbManager,
	}

	// Create handler registry
	s.handlers = handlers.NewRegistry(dbManager)

	// Setup routes
	s.setupRoutes()

	return s, nil
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

	// Static files (for the UI)
	fileServer := http.FileServer(http.Dir("./ui/web"))
	s.router.Handle("/*", http.StripPrefix("/", fileServer))
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	// Initialize database schema
	if err := s.db.MigrateSchema(); err != nil {
		return err
	}

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("Server starting on %s", addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Close database connection
	if s.db != nil {
		s.db.Close()
	}

	return s.server.Shutdown(ctx)
}

// DB returns the database manager
func (s *Server) DB() *db.Manager {
	return s.db
}
