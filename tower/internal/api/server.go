package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/zepzeper/tower/internal/api/handlers"
	towermiddleware "github.com/zepzeper/tower/internal/api/middleware"
	"github.com/zepzeper/tower/internal/config"
)

// Server represents the HTTP server
type Server struct {
	server *http.Server
	router *chi.Mux
	config *config.Config
	db     *sql.DB
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config) *Server {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(towermiddleware.CORS)

	// Connect to database
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create server
	s := &Server{
		router: r,
		config: cfg,
		db:     db,
	}

	// Setup routes
	s.setupRoutes()

	return s
}

// connectDB establishes a connection to the database
func connectDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	// Try to connect with retries
	var db *sql.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open(cfg.Database.Driver, connStr)
		if err == nil {
			// Test the connection
			err = db.Ping()
			if err == nil {
				log.Println("Database connection established successfully")
				return db, nil
			}
		}
		
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Second * 2)
		}
	}
	
	return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
}

// initDB initializes the database schema
func (s *Server) initDB() error {
	// Create tables if they don't exist
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS channels (
			id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			type VARCHAR(50) NOT NULL,
			description TEXT,
			config JSONB
		);
		
		CREATE TABLE IF NOT EXISTS workflows (
			id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			triggers JSONB,
			actions JSONB,
			active BOOLEAN DEFAULT TRUE
		);
		
		CREATE TABLE IF NOT EXISTS transformers (
			id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			mappings JSONB,
			functions JSONB
		);
	`)
	
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}
	
	log.Println("Database schema initialized successfully")
	return nil
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
		// V1 API
		r.Route("/v1", func(r chi.Router) {
			// Channels
			r.Route("/channels", func(r chi.Router) {
				r.Get("/", handlers.ListChannels)
				r.Post("/", handlers.CreateChannel)
				r.Route("/{channelID}", func(r chi.Router) {
					r.Get("/", handlers.GetChannel)
					r.Put("/", handlers.UpdateChannel)
					r.Delete("/", handlers.DeleteChannel)
				})
			})

			// Workflows
			r.Route("/workflows", func(r chi.Router) {
				r.Get("/", handlers.ListWorkflows)
				r.Post("/", handlers.CreateWorkflow)
				r.Route("/{workflowID}", func(r chi.Router) {
					r.Get("/", handlers.GetWorkflow)
					r.Put("/", handlers.UpdateWorkflow)
					r.Delete("/", handlers.DeleteWorkflow)
					r.Post("/execute", handlers.ExecuteWorkflow)
				})
			})

			// Transformers
			r.Route("/transformers", func(r chi.Router) {
				r.Get("/", handlers.ListTransformers)
				r.Post("/", handlers.CreateTransformer)
				r.Route("/{transformerID}", func(r chi.Router) {
					r.Get("/", handlers.GetTransformer)
					r.Put("/", handlers.UpdateTransformer)
					r.Delete("/", handlers.DeleteTransformer)
				})
			})
		})
	})

	// Static files (for the UI)
	fileServer := http.FileServer(http.Dir("./ui/web"))
	s.router.Handle("/*", http.StripPrefix("/", fileServer))
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	// Initialize database schema
	if err := s.initDB(); err != nil {
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
