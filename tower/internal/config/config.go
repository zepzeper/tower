package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Auth     AuthConfig     `json:"auth"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port int `json:"port"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret     string `json:"jwt_secret"`
	TokenDuration int    `json:"token_duration"` // in hours
}

// Load loads configuration from environment variables or config file
func Load() (*Config, error) {
	// Set default configuration
	config := &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Database: DatabaseConfig{
			Driver: "postgres",
			Host:   "localhost",
			Port:   5432,
			Name:   "apimiddleware",
			User:   "postgres",
			// Password is intentionally left empty for default
		},
		Auth: AuthConfig{
			JWTSecret:     "your-secret-key",
			TokenDuration: 24,
		},
	}

	// Override with environment variables
	if err := loadFromEnv(config); err != nil {
		return nil, err
	}

	// If config file exists, override with it
	configFile := getConfigFile()
	if configFile != "" {
		if err := loadFromFile(config, configFile); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(cfg *Config) error {
	// Server configuration
	if port := os.Getenv("SERVER_PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil && p > 0 {
			cfg.Server.Port = p
		}
	}

	// Database configuration
	if driver := os.Getenv("DB_DRIVER"); driver != "" {
		cfg.Database.Driver = driver
	}
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil && p > 0 {
			cfg.Database.Port = p
		}
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		cfg.Database.Name = name
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}

	// Auth configuration
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.Auth.JWTSecret = jwtSecret
	}
	if tokenDuration := os.Getenv("TOKEN_DURATION"); tokenDuration != "" {
		var d int
		if _, err := fmt.Sscanf(tokenDuration, "%d", &d); err == nil && d > 0 {
			cfg.Auth.TokenDuration = d
		}
	}

	return nil
}

// loadFromFile loads configuration from a JSON file
func loadFromFile(cfg *Config, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("failed to decode config file: %v", err)
	}

	return nil
}

// getConfigFile returns the path to the config file, if it exists
func getConfigFile() string {
	// Check for config file path in environment
	if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
		if _, err := os.Stat(configFile); err == nil {
			return configFile
		}
	}

	// Check for default config file locations
	locations := []string{
		"./config.json",
		"./configs/config.json",
		"/etc/api-middleware/config.json",
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc
		}
	}

	return ""
}
