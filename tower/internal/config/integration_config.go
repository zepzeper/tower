package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// IntegrationConfig holds configuration for API integrations
type IntegrationConfig struct {
	Connections     []ConnectionConfig  `json:"connections"`
	Transformers    []TransformerConfig `json:"transformers"`
	AdapterMappings []AdapterMapping    `json:"adapterMappings"`
}

// ConnectionConfig defines a connection between two APIs
type ConnectionConfig struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	SourceID      string                 `json:"sourceId"`
	TargetID      string                 `json:"targetId"`
	TransformerID string                 `json:"transformerId"`
	Query         map[string]interface{} `json:"query"`
	Schedule      string                 `json:"schedule"`
	Active        bool                   `json:"active"`
}

// TransformerConfig defines a data transformer
type TransformerConfig struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Mappings    []FieldMapping `json:"mappings"`
	Functions   []Function     `json:"functions"`
}

// FieldMapping defines how to map fields between schemas
type FieldMapping struct {
	SourceField string `json:"sourceField"`
	TargetField string `json:"targetField"`
}

// Function defines a transformation function
type Function struct {
	Name        string   `json:"name"`
	TargetField string   `json:"targetField"`
	Args        []string `json:"args"`
}

// AdapterMapping defines field mappings for an adapter
type AdapterMapping struct {
	AdapterName string                `json:"adapterName"`
	EntityType  string                `json:"entityType"`
	Mappings    []AdapterFieldMapping `json:"mappings"`
}

// AdapterFieldMapping defines how to map a field for an adapter
type AdapterFieldMapping struct {
	SourceField        string `json:"sourceField"`
	CanonicalField     string `json:"canonicalField"`
	IsRequired         bool   `json:"isRequired"`
	DefaultValue       string `json:"defaultValue,omitempty"`
	TransformToCanon   string `json:"transformToCanon,omitempty"`
	TransformFromCanon string `json:"transformFromCanon,omitempty"`
}

// LoadIntegrationConfig loads integration configuration from JSON files
func LoadIntegrationConfig(configDir string) (*IntegrationConfig, error) {
	config := &IntegrationConfig{
		Connections:     make([]ConnectionConfig, 0),
		Transformers:    make([]TransformerConfig, 0),
		AdapterMappings: make([]AdapterMapping, 0),
	}

	// Load connections
	connectionsPath := filepath.Join(configDir, "connections.json")
	if fileExists(connectionsPath) {
		connections, err := loadConnectionsConfig(connectionsPath)
		if err != nil {
			return nil, fmt.Errorf("error loading connections config: %w", err)
		}
		config.Connections = connections
	}

	// Load transformers
	transformersPath := filepath.Join(configDir, "transformers.json")
	if fileExists(transformersPath) {
		transformers, err := loadTransformersConfig(transformersPath)
		if err != nil {
			return nil, fmt.Errorf("error loading transformers config: %w", err)
		}
		config.Transformers = transformers
	}

	// Load adapter mappings
	adapterMappingsPath := filepath.Join(configDir, "adapter_mappings.json")
	if fileExists(adapterMappingsPath) {
		adapterMappings, err := loadAdapterMappingsConfig(adapterMappingsPath)
		if err != nil {
			return nil, fmt.Errorf("error loading adapter mappings config: %w", err)
		}
		config.AdapterMappings = adapterMappings
	}

	return config, nil
}

// Helper functions to load specific config files
func loadConnectionsConfig(path string) ([]ConnectionConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var connections []ConnectionConfig
	err = json.Unmarshal(data, &connections)
	if err != nil {
		return nil, err
	}

	return connections, nil
}

func loadTransformersConfig(path string) ([]TransformerConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var transformers []TransformerConfig
	err = json.Unmarshal(data, &transformers)
	if err != nil {
		return nil, err
	}

	return transformers, nil
}

func loadAdapterMappingsConfig(path string) ([]AdapterMapping, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var adapterMappings []AdapterMapping
	err = json.Unmarshal(data, &adapterMappings)
	if err != nil {
		return nil, err
	}

	return adapterMappings, nil
}

// SaveIntegrationConfig saves integration configuration to JSON files
func SaveIntegrationConfig(configDir string, config *IntegrationConfig) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Save connections
	if len(config.Connections) > 0 {
		connectionsPath := filepath.Join(configDir, "connections.json")
		err := saveToJSONFile(connectionsPath, config.Connections)
		if err != nil {
			return fmt.Errorf("error saving connections config: %w", err)
		}
	}

	// Save transformers
	if len(config.Transformers) > 0 {
		transformersPath := filepath.Join(configDir, "transformers.json")
		err := saveToJSONFile(transformersPath, config.Transformers)
		if err != nil {
			return fmt.Errorf("error saving transformers config: %w", err)
		}
	}

	// Save adapter mappings
	if len(config.AdapterMappings) > 0 {
		adapterMappingsPath := filepath.Join(configDir, "adapter_mappings.json")
		err := saveToJSONFile(adapterMappingsPath, config.AdapterMappings)
		if err != nil {
			return fmt.Errorf("error saving adapter mappings config: %w", err)
		}
	}

	return nil
}

// Helper function to save data to JSON file
func saveToJSONFile(path string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, jsonData, 0644)
}

// Helper function to check if file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
