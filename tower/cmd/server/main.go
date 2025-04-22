package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/zepzeper/tower/internal/config"
	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/services/authentication"
	"github.com/zepzeper/tower/internal/services/credentials"
	"github.com/zepzeper/tower/internal/services/mapping"
	"github.com/zepzeper/tower/internal/services/relations"
	"github.com/zepzeper/tower/internal/webapi/handlers"
)

// Application represents the web application
type Application struct {
	config             *config.Config
	router             *chi.Mux
	databaseManager    *database.Manager
	authService        *auth.Service
	mappingService     *mapping.Service
	credentialsService *credentials.Service
	relationsService   *relations.Service
	handlerRegistry    *handlers.Registry
}

// NewApplication creates a new application instance
func NewApplication(cfg *config.Config) (*Application, error) {
	// Initialize database connection
	dbManager, err := database.NewManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Setup database schema
	if err := dbManager.MigrateSchema(); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	// Initialize services
	authService := auth.NewService(
		cfg.Auth.JWTSecret,
		cfg.Auth.AccessTokenExpiry,
		cfg.Auth.RefreshTokenExpiry,
		dbManager,
	)

	// Initialize other services
	mappingService := mapping.NewService(nil, dbManager)
	credentialsService := credentials.NewService(dbManager)
	relationsService := relations.NewService(dbManager)

	// Initialize router
	router := chi.NewRouter()
	setupGlobalMiddleware(router)

	// Initialize handler registry
	handlerRegistry := handlers.NewRegistry(mappingService, credentialsService, relationsService, authService)
	handlerRegistry.RegisterRoutes(router)

	return &Application{
		config:             cfg,
		router:             router,
		databaseManager:    dbManager,
		authService:        authService,
		mappingService:     mappingService,
		credentialsService: credentialsService,
		relationsService:   relationsService,
		handlerRegistry:    handlerRegistry,
	}, nil
}

// setupGlobalMiddleware adds global middleware to the router
func setupGlobalMiddleware(router *chi.Mux) {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)

	// CORS configuration
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ancient browsers
	}))
}

// Run starts the application server
func (a *Application) Run() error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.Server.Port),
		Handler: a.router,
	}

	log.Printf("Starting server on port %d", a.config.Server.Port)
	return server.ListenAndServe()
}

// Shutdown gracefully shuts down the application
func (a *Application) Shutdown() error {
	// Close database connection
	if a.databaseManager != nil {
		return a.databaseManager.Close()
	}
	return nil
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start the application
	app, err := NewApplication(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Ensure clean shutdown
	defer func() {
		if err := app.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	// Start the server
	log.Fatal(app.Run())
}
