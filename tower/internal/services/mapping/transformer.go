package mapping

import (
	"fmt"
	"strconv"
	"strings"
	"encoding/json"
	"log"
	"reflect"
)

// Transformer provides transformation functions for mapping operations
type Transformer struct{}

// NewTransformer creates a new transformer
func NewTransformer() *Transformer {
	return &Transformer{}
}

// TransformSchema flattens and describes all fields in the schema
func TransformSchema(schema map[string]interface{}, schemaType string, opts *schemaTransformOptions) []FieldDefinition {
	if opts == nil {
		opts = &schemaTransformOptions{
			Prefix:          map[string]string{"source": "s", "target": "t"}[schemaType],
			IncludeRequired: true,
		}
	}

	fields := []FieldDefinition{}
	counter := 1

	var processObject func(mapData interface{}, path string, parent string)
	processObject = func(mapData interface{}, path string, parent string) {
		dataMap, ok := mapData.(map[string]interface{})
		if !ok {
			log.Printf("⚠️ Skipping non-object at path: %s (type: %T)\n", path, mapData)
			return
		}

		for key, value := range dataMap {
			fullPath := key
			if path != "" {
				fullPath = path + "." + key
			}

			displayName := key
			if parent != "" {
				displayName = parent + "." + key
			}

			fieldType := "string"
			sample := ""

			switch v := value.(type) {
			case string:
				fieldType = "string"
				sample = v
			case float64:
				fieldType = "number"
				sample = fmt.Sprintf("%v", v)
			case bool:
				fieldType = "boolean"
				sample = fmt.Sprintf("%v", v)
			case []interface{}:
				if len(v) > 0 {
					first := v[0]
					switch item := first.(type) {
					case map[string]interface{}:
						fieldType = "array.object"
						for k, val := range item {
							processObject(
								map[string]interface{}{k: val},
								fullPath+"[]."+k,
								displayName+"[]."+k,
							)
						}
						j, _ := json.Marshal(item)
						sample = truncateString(string(j), 100)
					default:
						fieldType = fmt.Sprintf("array.%s", reflect.TypeOf(item).Kind())
						j, _ := json.Marshal(v)
						sample = truncateString(string(j), 100)
					}
				} else {
					fieldType = "array"
					sample = "[]"
				}
			case map[string]interface{}:
				fieldType = "object"
				processObject(v, fullPath, displayName)
				j, _ := json.Marshal(v)
				sample = truncateString(string(j), 100)
			default:
				log.Printf("❓ Unhandled value for key '%s': %v (type: %T)\n", key, v, v)
				if v == nil {
					fieldType = "null"
					sample = "null"
				} else {
					fieldType = reflect.TypeOf(v).String()
					sample = fmt.Sprintf("%v", v)
				}
			}

			if fieldType != "object" {
				isRequired := opts.IncludeRequired && (key == "id" || key == "sku" || key == "name" || key == "price" || key == "code" || key == "description")
				field := FieldDefinition{
					ID:       fmt.Sprintf("%s%d", opts.Prefix, counter),
					Name:     key,
					Type:     fieldType,
					Path:     strings.ReplaceAll(fullPath, ".[].", "[]."),
					Sample:   sample,
					Required: isRequired,
				}
				fields = append(fields, field)
				counter++
			}
		}
	}

	processObject(schema, "", "")
	return fields
}

// GenerateMappingData attempts to auto-map fields based on name similarity
func GenerateMappingData(sourceSchema, targetSchema map[string]interface{}) (MappingData, error) {
	sourceFields := TransformSchema(sourceSchema, "source", nil)
	targetFields := TransformSchema(targetSchema, "target", nil)

	log.Printf("✅ Parsed %d source fields, %d target fields\n", len(sourceFields), len(targetFields))

	mappings := []MappingDefinition{}
	mappingCounter := 1

	for _, target := range targetFields {
		targetName := strings.ToLower(target.Name)
		var matchedSource *FieldDefinition

		for _, source := range sourceFields {
			sourceName := strings.ToLower(source.Name)

			if sourceName == targetName ||
				strings.Contains(sourceName, targetName) ||
				strings.Contains(targetName, sourceName) ||
				(targetName == "sku" && sourceName == "id") ||
				(targetName == "description" && sourceName == "short_description") ||
				(targetName == "price" && sourceName == "regular_price") ||
				(targetName == "image" && sourceName == "images") {
				matchedSource = &source
				break
			}
		}

		if matchedSource == nil {
			log.Printf("⚠️ No match for target field: %s\n", target.Name)
			continue
		}

		var transform *string
		if matchedSource.Type != target.Type {
			if matchedSource.Type == "string" && target.Type == "number" {
				tr := "parseFloat"
				transform = &tr
			} else if matchedSource.Type == "number" && target.Type == "string" {
				tr := "toString"
				transform = &tr
			} else if strings.HasPrefix(matchedSource.Type, "array") && target.Type == "string" {
				tr := "splitFirst"
				transform = &tr
			}
		}

		mappings = append(mappings, MappingDefinition{
			ID:          fmt.Sprintf("m%d", mappingCounter),
			SourceField: matchedSource.ID,
			TargetField: target.ID,
			Transform:   transform,
		})
		mappingCounter++
	}

	return MappingData{
		SourceFields: sourceFields,
		TargetFields: targetFields,
		Mappings:     mappings,
	}, nil
}

// Utility
func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

// ExtractValueFromPath extracts a value from a nested path in a map
func (t *Transformer) ExtractValueFromPath(data map[string]interface{}, path string) (interface{}, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {
		// Handle array access using bracket notation [index]
		if strings.Contains(part, "[") && strings.Contains(part, "]") {
			beforeBracket := part[:strings.Index(part, "[")]
			
			// Handle object access into array element
			if c, ok := current.(map[string]interface{}); ok {
				if arr, ok := c[beforeBracket].([]interface{}); ok && len(arr) > 0 {
					// For simplicity in testing, just take the first element
					current = arr[0]
					continue
				}
			}
			
			// If we can't properly handle the array access, just continue with nil
			return nil, fmt.Errorf("couldn't access array at path %s", path)
		}

		// Handle regular field access
		if c, ok := current.(map[string]interface{}); ok {
			if val, ok := c[part]; ok {
				current = val
			} else {
				return nil, fmt.Errorf("key not found: %s in path %s", part, path)
			}
		} else {
			return nil, fmt.Errorf("cannot access path %s on non-map value", path)
		}
	}

	return current, nil
}

// GetFieldNameFromPath extracts the field name from a path
func (t *Transformer) GetFieldNameFromPath(path string) string {
	if path == "" {
		return ""
	}

	// Get the last part of the path
	parts := strings.Split(path, ".")
	lastPart := parts[len(parts)-1]

	// Remove array notation if present
	if idx := strings.Index(lastPart, "["); idx != -1 {
		lastPart = lastPart[:idx]
	}

	return lastPart
}

// ApplyTransformation applies the specified transformation with awareness of target field type
func (t *Transformer) ApplyTransformation(value interface{}, transform *string, targetFieldType string) interface{} {
	// Handle null values
	if value == nil {
		return nil
	}

	// If there's an explicit transformation, apply it first
	if transform != nil && *transform != "identity" {
		value = t.applyExplicitTransform(value, *transform)
	}

	// Then handle special format conversions based on target type
	return t.adaptToTargetType(value, targetFieldType)
}

// applyExplicitTransform applies the transformation specified in the mapping
func (t *Transformer) applyExplicitTransform(value interface{}, transform string) interface{} {
	switch transform {
	case "parseFloat":
		switch v := value.(type) {
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err == nil {
				return f
			}
		case float64, float32, int, int64:
			return v
		}
	case "toString":
		return fmt.Sprintf("%v", value)
	case "trim":
		if str, ok := value.(string); ok {
			return strings.TrimSpace(str)
		}
	case "uppercase":
		if str, ok := value.(string); ok {
			return strings.ToUpper(str)
		}
	case "lowercase":
		if str, ok := value.(string); ok {
			return strings.ToLower(str)
		}
	case "round":
		switch v := value.(type) {
		case float64:
			return int(v + 0.5)
		case float32:
			return int(v + 0.5)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return int(f + 0.5)
			}
		}
	case "splitFirst":
		switch v := value.(type) {
		case []interface{}:
			if len(v) > 0 {
				return v[0]
			}
		case string:
			parts := strings.Split(v, " ")
			if len(parts) > 0 {
				return parts[0]
			}
		}
	}

	return value
}

// adaptToTargetType converts the value to match the target field type
func (t *Transformer) adaptToTargetType(value interface{}, targetType string) interface{} {
	// If target type is empty, return as is
	if targetType == "" {
		return value
	}

	switch targetType {
	case "string":
		// If value is array of image objects, extract URL
		if arr, ok := value.([]interface{}); ok && len(arr) > 0 {
			if firstItem, ok := arr[0].(map[string]interface{}); ok {
				// Try to find image URL in common fields
				for _, field := range []string{"src", "url", "link", "href"} {
					if url, ok := firstItem[field].(string); ok {
						return url
					}
				}
			}
		}
		
		// If value is an image object, extract URL
		if obj, ok := value.(map[string]interface{}); ok {
			for _, field := range []string{"src", "url", "link", "href"} {
				if url, ok := obj[field].(string); ok {
					return url
				}
			}
		}
		
		// For other types, convert to string
		if !isString(value) {
			return fmt.Sprintf("%v", value)
		}
		
	case "number", "float":
		return convertToNumber(value)
		
	case "integer", "int":
		num := convertToNumber(value)
		if f, ok := num.(float64); ok {
			return int(f)
		}
		return num
		
	case "boolean", "bool":
		return convertToBoolean(value)
		
	case "array":
		// If it's not already an array, wrap in array
		if _, ok := value.([]interface{}); !ok {
			return []interface{}{value}
		}
	}
	
	return value
}

// isString checks if a value is a string
func isString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

// convertToNumber attempts to convert a value to a number
func convertToNumber(value interface{}) interface{} {
	switch v := value.(type) {
	case float64, float32, int, int64:
		return v
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	case bool:
		if v {
			return 1
		}
		return 0
	}
	return value
}

// convertToBoolean attempts to convert a value to boolean
func convertToBoolean(value interface{}) interface{} {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		lower := strings.ToLower(v)
		return lower == "true" || lower == "yes" || lower == "1"
	case float64, float32, int, int64:
		// Any non-zero number is true
		return v != 0
	}
	return value
}

// TransformFields applies mappings to convert source data to target format
func (t *Transformer) TransformFields(sourceData map[string]interface{}, targetSchema map[string]interface{}, mappings []MappingMetadata) map[string]interface{} {
	result := make(map[string]interface{})
	targetTypes := t.extractTargetFieldTypes(targetSchema)

	for _, mapping := range mappings {
		// Skip mappings without paths
		if mapping.SourcePath == "" || mapping.TargetPath == "" {
			continue
		}

		// Extract value from source path
		sourceValue, err := t.ExtractValueFromPath(sourceData, mapping.SourcePath)
		if err != nil {
			// Log error but continue with other mappings
			fmt.Printf("Error extracting value from path %s: %v\n", mapping.SourcePath, err)
			continue
		}

		// Get target field name and type
		targetName := t.GetFieldNameFromPath(mapping.TargetPath)
		if targetName == "" {
			continue
		}
		
		targetType := targetTypes[targetName]

		// Apply transformation with target type awareness
		transformedValue := t.ApplyTransformation(sourceValue, mapping.Transform, targetType)

		// Store in result
		result[targetName] = transformedValue
	}

	return result
}

// extractTargetFieldTypes extracts field types from a target schema
func (t *Transformer) extractTargetFieldTypes(schema map[string]interface{}) map[string]string {
	types := make(map[string]string)
	
	// Recursively process all fields
	var processObject func(obj map[string]interface{}, prefix string)
	processObject = func(obj map[string]interface{}, prefix string) {
		for key, value := range obj {
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}
			
			// If it's a string, it might be a type definition
			if typeStr, ok := value.(string); ok {
				types[key] = typeStr
			} else if childObj, ok := value.(map[string]interface{}); ok {
				// If it's an object, process recursively
				processObject(childObj, path)
			}
		}
	}
	
	processObject(schema, "")
	return types
}
