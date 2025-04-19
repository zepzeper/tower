package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/zepzeper/tower/internal/config"
	"github.com/zepzeper/tower/internal/core/registry"
	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/database/repositories"
	"github.com/zepzeper/tower/internal/server"
	"github.com/zepzeper/tower/internal/services"
	"github.com/zepzeper/tower/internal/services/mapping"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	dbManager, err := database.NewManager(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbManager.Close()

	// Migrate database schema
	if err := dbManager.MigrateSchema(); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	authRepo := repositories.NewAuthRepository(dbManager.DB)

  cwd, err := os.Getwd()
  if err != nil {
    log.Fatal(err)
  }
  schemaPath := filepath.Join(cwd, "internal/core/schemas")
  schemaFetcher := registry.NewSchemaFetcher(schemaPath)

	// Create service layer
  authService := services.NewAuthService(*authRepo, "123", 24*time.Hour);
  mappingService := mapping.NewService(schemaFetcher)

	// Create central server
	server := server.NewServer(
    authService,
    mappingService,
	)

	// // Initialize job manager with existing connections
	// if err := jobManager.InitializeFromDatabase(context.Background()); err != nil {
	// 	log.Printf("Warning: Failed to initialize job manager: %v", err)
	// }

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		if err := server.Start(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped")
}
