// cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zepzeper/tower/internal/config"
	"github.com/zepzeper/tower/internal/connectors/woocommerce"
	"github.com/zepzeper/tower/internal/core/registry"
	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/database/repositories"
	"github.com/zepzeper/tower/internal/server"
	"github.com/zepzeper/tower/internal/services"
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

	// Initialize connector registry
	connectorRegistry := registry.NewConnectorRegistry()

	authRepo := repositories.NewAuthRepository(dbManager.DB)

	// Register built-in connectors
	registerConnectors(connectorRegistry)

	// Create service layer
	connectorService := services.NewConnectorService(connectorRegistry)
	transformerService := services.NewTransformerService(dbManager, connectorRegistry)
	jobManager := services.NewJobManager(dbManager, connectorService, transformerService)
	connectionService := services.NewConnectionService(dbManager, connectorService, transformerService, jobManager)
  authService := services.NewAuthService(*authRepo, "123", 24*time.Hour);

	// Create central server
	server := server.NewServer(
		connectorService,
		transformerService,
		connectionService,
    authService,
	)

	// Initialize job manager with existing connections
	if err := jobManager.InitializeFromDatabase(context.Background()); err != nil {
		log.Printf("Warning: Failed to initialize job manager: %v", err)
	}

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

// registerConnectors registers built-in connector implementations
func registerConnectors(registry *registry.ConnectorRegistry) {
	// WooCommerce connector (example configuration)
	wooCommerceConnector, err := woocommerce.NewConnector(map[string]interface{}{
		"api_url":        "https://example.com/wp-json/wc/v3",
		"consumer_key":   "your_consumer_key",
		"consumer_secret": "your_consumer_secret",
	})
	if err == nil {
		registry.Register("woocommerce", wooCommerceConnector)
	}

	// // Mirakl connector (example configuration)
	// miraklConnector, err := mirakl.NewConnector(map[string]interface{}{
	// 	"marketplace": "example",
	// 	"shop_id":     "your_shop_id",
	// 	"api_key":     "your_api_key",
	// })
	// if err == nil {
	// 	registry.Register("mirakl", miraklConnector)
	// }
	//
	// // Kaufland connector (example configuration)
	// kauflandConnector, err := kaufland.NewConnector(map[string]interface{}{
	// 	"public_key":  "your_public_key",
	// 	"private_key": "your_private_key",
	// })
	// if err == nil {
	// 	registry.Register("kaufland", kauflandConnector)
	// }
}
