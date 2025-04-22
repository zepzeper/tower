package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/zepzeper/tower/internal/config"
	"github.com/zepzeper/tower/internal/database/repositories"
	"github.com/zepzeper/tower/internal/database/schema"
)

// Manager manages database connections and operations
type Manager struct {
	DB    *sql.DB
	Repos *repositories.Factory
}

// NewManager creates a new database manager
func NewManager(cfg *config.Config) (*Manager, error) {
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
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			// Test the connection
			err = db.Ping()
			if err == nil {
				log.Println("Database connection established successfully")
				manager := &Manager{
					DB: db,
				}

				// Create repository factory
				manager.Repos = repositories.NewFactory(db)

				return manager, nil
			}
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Second * 2)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
}

// Close closes the database connection
func (m *Manager) Close() error {
	if m.DB != nil {
		return m.DB.Close()
	}
	return nil
}

// MigrateSchema creates all necessary database tables if they don't exist
func (m *Manager) MigrateSchema() error {
	// Apply each schema in the correct order
	for _, schemaSQL := range schema.AllSchemas() {
		_, err := m.DB.Exec(schemaSQL)
		if err != nil {
			return fmt.Errorf("failed to apply database schema: %v", err)
		}
	}

	log.Println("Database schema initialized successfully")
	return nil

}

func (m *Manager) Transaction(fn func(*sql.Tx) error) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	// If the function returns an error, rollback the transaction
	// Otherwise, commit it
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
