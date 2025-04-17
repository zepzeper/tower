package transformers

import (
	"errors"
	"fmt"
	"strings"
	"time"
	
	"github.com/zepzeper/tower/internal/core/connectors"
)

// FieldMapping defines how data should be mapped between fields
type FieldMapping struct {
	SourceField string `json:"sourceField"`
	TargetField string `json:"targetField"`
}

// Function represents a transformation function to apply
type Function struct {
	Name        string   `json:"name"`
	TargetField string   `json:"targetField"`
	Args        []string `json:"args"`
}

// Transformer holds the configuration for data transformation
type Transformer struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Mappings    []FieldMapping `json:"mappings"`
	Functions   []Function    `json:"functions"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

// Transform applies the transformer's mappings and functions to input data
func (t *Transformer) Transform(input connectors.DataPayload) (connectors.DataPayload, error) {
	// Create the output map
	output := connectors.DataPayload{}

	// Apply field mappings
	for _, mapping := range t.Mappings {
		value, err := getNestedValue(input, mapping.SourceField)
		if err != nil {
			// Skip this mapping if source field doesn't exist
			continue
		}
		
		// Set the value in the output
		err = setNestedValue(output, mapping.TargetField, value)
		if err != nil {
			return nil, fmt.Errorf("mapping error for %s → %s: %v", 
				mapping.SourceField, mapping.TargetField, err)
		}
	}

	// Apply transformation functions
	for _, fn := range t.Functions {
		// Get function arguments (resolve field references)
		args := make([]interface{}, 0, len(fn.Args))
		for _, arg := range fn.Args {
			if strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'") {
				// Literal string value
				args = append(args, arg[1:len(arg)-1])
			} else {
				// Field reference
				value, err := getNestedValue(input, arg)
				if err != nil {
					// Use nil for missing fields
					args = append(args, nil)
				} else {
					args = append(args, value)
				}
			}
		}
		
		// Apply the function
		result, err := applyFunction(fn.Name, args)
		if err != nil {
			return nil, fmt.Errorf("function error for %s: %v", fn.Name, err)
		}
		
		// Store the result
		err = setNestedValue(output, fn.TargetField, result)
		if err != nil {
			return nil, fmt.Errorf("error setting function result for %s: %v", 
				fn.TargetField, err)
		}
	}

	return output, nil
}

// getNestedValue retrieves a potentially nested value from a DataPayload
func getNestedValue(data connectors.DataPayload, path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	current := data
	
	// Navigate through all but the last part
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		next, ok := current[part]
		if !ok {
			return nil, fmt.Errorf("path %s not found at %s", path, part)
		}
		
		// Try to convert to DataPayload
		nextMap, ok := next.(map[string]interface{})
		if !ok {
			// Try another type of map
			nextMapAlt, ok := next.(connectors.DataPayload)
			if !ok {
				return nil, fmt.Errorf("path %s blocked at %s: not a map", path, part)
			}
			current = nextMapAlt
			continue
		}
		
		// Convert to DataPayload for consistency
		current = connectors.DataPayload(nextMap)
	}
	
	// Get the final value
	value, ok := current[parts[len(parts)-1]]
	if !ok {
		return nil, fmt.Errorf("path %s final key %s not found", path, parts[len(parts)-1])
	}
	
	return value, nil
}

// setNestedValue sets a potentially nested value in a DataPayload
func setNestedValue(data connectors.DataPayload, path string, value interface{}) error {
	parts := strings.Split(path, ".")
	current := data
	
	// Create nested maps as needed for all but the last part
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		
		// Check if this level exists
		next, exists := current[part]
		if !exists {
			// Create a new map
			newMap := connectors.DataPayload{}
			current[part] = newMap
			current = newMap
			continue
		}
		
		// Try to convert existing value to DataPayload
		nextMap, ok := next.(map[string]interface{})
		if !ok {
			// Try another type of map
			nextMapAlt, ok := next.(connectors.DataPayload)
			if !ok {
				return fmt.Errorf("path %s blocked at %s: existing value is not a map", path, part)
			}
			current = nextMapAlt
			continue
		}
		
		// Convert to DataPayload for consistency
		current = connectors.DataPayload(nextMap)
	}
	
	// Set the final value
	current[parts[len(parts)-1]] = value
	return nil
}

// applyFunction applies a named function to the given arguments
func applyFunction(name string, args []interface{}) (interface{}, error) {
	switch name {
	case "concatenate":
		return functionConcatenate(args)
	case "uppercase":
		return functionUppercase(args)
	case "lowercase":
		return functionLowercase(args)
	case "trim":
		return functionTrim(args)
	// Add more functions as needed
	default:
		return nil, fmt.Errorf("unknown function: %s", name)
	}
}

// Function implementations

func functionConcatenate(args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return "", nil
	}
	
	var result strings.Builder
	for _, arg := range args {
		if arg != nil {
			result.WriteString(fmt.Sprintf("%v", arg))
		}
	}
	
	return result.String(), nil
}

func functionUppercase(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("uppercase function requires exactly 1 argument")
	}
	
	if args[0] == nil {
		return "", nil
	}
	
	return strings.ToUpper(fmt.Sprintf("%v", args[0])), nil
}

func functionLowercase(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("lowercase function requires exactly 1 argument")
	}
	
	if args[0] == nil {
		return "", nil
	}
	
	return strings.ToLower(fmt.Sprintf("%v", args[0])), nil
}

func functionTrim(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("trim function requires exactly 1 argument")
	}
	
	if args[0] == nil {
		return "", nil
	}
	
	return strings.TrimSpace(fmt.Sprintf("%v", args[0])), nil
}
