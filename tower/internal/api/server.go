package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yourusername/tower/internal/api/handlers"
	apimiddleware "github.com/yourusername/tower/internal/api/middleware"
	"github.com/yourusername/tower/internal/config"
)

// Server represents the HTTP server
type Server struct {
	server *http.Server
	router *chi.Mux
	config *config.Config
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
	r.Use(apimiddleware.CORS)

	// Setup routes
	setupRoutes(r)

	return &Server{
		router: r,
		config: cfg,
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// setupRoutes configures the API routes
func setupRoutes(r *chi.Mux) {
	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
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
	r.Handle("/*", http.StripPrefix("/", fileServer))
}
