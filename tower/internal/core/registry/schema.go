package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SchemaRegistry manages API schemas and their mappings
type SchemaFetcher struct {
	basePath string
}

// NewSchemaFetcher creates a new StaticSchemaFetcher with preloaded schemas
func NewSchemaFetcher(basePath string) *SchemaFetcher {
	return &SchemaFetcher{
		basePath: basePath,
	}
}

// GetSchema retrieves the schema by type from the in-memory map
func (f *SchemaFetcher) GetSchema(schemaType, operation, sourceortarget string) (map[string]interface{}, error) {
	filename := filepath.Join(f.basePath, "schemas/"+sourceortarget+"_"+schemaType+"_"+operation+".json")

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open schema file %s: %w", filename, err)
	}
	defer file.Close()

	var schema map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&schema); err != nil {
		return nil, fmt.Errorf("could not decode JSON schema from %s: %w", filename, err)
	}

	return schema, nil
}

func (f *SchemaFetcher) GetApis() (map[string]interface{}, error) {
	filename := filepath.Join(f.basePath, "apis.json")
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open schema file %s: %w", filename, err)
	}
	defer file.Close()

	var schema map[string]interface{}
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&schema); err != nil {
		return nil, fmt.Errorf("could not decode JSON schema from %s: %w", filename, err)
	}

	return schema, nil
}
