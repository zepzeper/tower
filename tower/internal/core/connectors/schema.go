package connectors

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Schema represents the field structure of a connector
type Schema struct {
	EntityName string
	Fields     map[string]FieldDefinition
}

// FieldDefinition describes a field in a connector's schema
type FieldDefinition struct {
	Type        string   // string, number, boolean, object, array
	Required    bool
	Path        string   // JSON path for nested fields
	EnumValues  []string // Possible values for enum fields
	Description string
}

// SchemaDiscoverer helps discover schema from sample data
type SchemaDiscoverer struct {
	MaxSamples int
}

// NewSchemaDiscoverer creates a new schema discoverer
func NewSchemaDiscoverer(maxSamples int) *SchemaDiscoverer {
	return &SchemaDiscoverer{
		MaxSamples: maxSamples,
	}
}

// DiscoverSchema analyzes samples to generate a schema
func (sd *SchemaDiscoverer) DiscoverSchema(entityName string, samples []DataPayload) Schema {
	schema := Schema{
		EntityName: entityName,
		Fields:     make(map[string]FieldDefinition),
	}
	
	// Process each sample
	for i, sample := range samples {
		if i >= sd.MaxSamples {
			break
		}
		
		sd.analyzeFields("", sample, schema.Fields)
	}
	
	return schema
}

// analyzeFields recursively analyzes fields in a sample
func (sd *SchemaDiscoverer) analyzeFields(prefix string, data interface{}, fields map[string]FieldDefinition) {
	if data == nil {
		return
	}
	
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		// For maps, we go through each key
		mapData, ok := data.(map[string]interface{})
		if !ok {
			// Try another type of map
			mapDataAlt, ok := data.(DataPayload)
			if !ok {
				return
			}
			mapData = mapDataAlt
		}
		
		for key, value := range mapData {
			field := key
			if prefix != "" {
				field = prefix + "." + key
			}
			
			// Add or update the field
			sd.updateField(field, value, fields)
			
			// Recursively analyze nested fields
			sd.analyzeFields(field, value, fields)
		}
		
	case reflect.Slice, reflect.Array:
		// For arrays, we analyze the first element
		sliceVal := reflect.ValueOf(data)
		if sliceVal.Len() > 0 {
			firstElem := sliceVal.Index(0).Interface()
			
			// Add or update the field as an array
			sd.updateField(prefix, data, fields)
			
			// Recursively analyze the first element
			fieldName := prefix + "[0]"
			sd.analyzeFields(fieldName, firstElem, fields)
		}
	default:
		// For other types, we just add/update the field
		sd.updateField(prefix, data, fields)
	}
}

// updateField updates field information based on a value
func (sd *SchemaDiscoverer) updateField(field string, value interface{}, fields map[string]FieldDefinition) {
	if field == "" {
		return
	}
	
	// Skip array index notation
	if strings.Contains(field, "[") {
		return
	}
	
	// Get existing field or create new one
	def, exists := fields[field]
	if !exists {
		def = FieldDefinition{
			Path:     field,
			Required: true, // Assume required by default
		}
	}
	
	// Determine type based on the value
	valueType := reflect.TypeOf(value)
	if valueType == nil {
		return
	}
	
	switch valueType.Kind() {
	case reflect.String:
		def.Type = "string"
	case reflect.Bool:
		def.Type = "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		def.Type = "number"
	case reflect.Map:
		def.Type = "object"
	case reflect.Slice, reflect.Array:
		def.Type = "array"
	default:
		def.Type = "string" // Default to string for unknown types
	}
	
	// Update the field definition
	fields[field] = def
}

// GenerateSchemaFromJSON generates a schema from a JSON sample
func GenerateSchemaFromJSON(entityName, jsonSample string) (Schema, error) {
	var sample DataPayload
	err := json.Unmarshal([]byte(jsonSample), &sample)
	if err != nil {
		return Schema{}, err
	}
	
	discoverer := NewSchemaDiscoverer(1)
	return discoverer.DiscoverSchema(entityName, []DataPayload{sample}), nil
}
