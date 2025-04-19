package transformers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// FieldDefinition represents a single mappable field
type FieldDefinition struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Path     string `json:"path"`
	Sample   string `json:"sample,omitempty"`
	Required bool   `json:"required,omitempty"`
}

// MappingDefinition represents a field mapping
type MappingDefinition struct {
	ID          string  `json:"id"`
	SourceField string  `json:"sourceField"`
	TargetField string  `json:"targetField"`
	Transform   *string `json:"transform"`
}

// MappingData holds all fields and suggested mappings
type MappingData struct {
	SourceFields []FieldDefinition   `json:"sourceFields"`
	TargetFields []FieldDefinition   `json:"targetFields"`
	Mappings     []MappingDefinition `json:"mappings"`
}

type schemaTransformOptions struct {
	Prefix         string
	IncludeRequired bool
}

// TransformSchema flattens and describes all fields in the schema
func TransformSchema(schema map[string]interface{}, schemaType string, opts *schemaTransformOptions) []FieldDefinition {
	if opts == nil {
		opts = &schemaTransformOptions{
			Prefix:         map[string]string{"source": "s", "target": "t"}[schemaType],
			IncludeRequired: true,
		}
	}
	fields := []FieldDefinition{}
	counter := 1

	var processObject func(mapData interface{}, path string, parent string)
	processObject = func(mapData interface{}, path string, parent string) {
		dataMap, ok := mapData.(map[string]interface{})
		if !ok {
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
						// Generate JSON snippet for preview
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
				fieldType = reflect.TypeOf(v).String()
				sample = fmt.Sprintf("%v", v)
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
func GenerateMappingData(sourceSchema, targetSchema map[string]interface{}) MappingData {
	sourceFields := TransformSchema(sourceSchema, "source", nil)
	targetFields := TransformSchema(targetSchema, "target", nil)

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

		if matchedSource != nil {
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
	}

	return MappingData{
		SourceFields: sourceFields,
		TargetFields: targetFields,
		Mappings:     mappings,
	}
}

// Utility
func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
