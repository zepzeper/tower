package transformers

import (
	"strings"
	
	"github.com/zepzeper/tower/internal/core/connectors"
)

// AutoMapper attempts to create mappings between two schemas
type AutoMapper struct {
	similarityThreshold float64
}

// NewAutoMapper creates a new automapper
func NewAutoMapper(threshold float64) *AutoMapper {
	return &AutoMapper{similarityThreshold: threshold}
}

// GenerateMappings creates mappings between source and target schemas
func (am *AutoMapper) GenerateMappings(source, target connectors.Schema) []FieldMapping {
    var mappings []FieldMapping
    
    // For each target field, find the best matching source field
    for targetField, _ := range target.Fields {
        bestMatch := ""
        bestScore := 0.0
        
        for sourceField, _ := range source.Fields {
            score := calculateSimilarity(sourceField, targetField)
            if score > am.similarityThreshold && score > bestScore {
                bestMatch = sourceField
                bestScore = score
            }
        }
        
        if bestMatch != "" {
            mappings = append(mappings, FieldMapping{
                SourceField: bestMatch,
                TargetField: targetField,
            })
        }
    }
    
    return mappings
}

// calculateSimilarity returns a similarity score between two field names
func calculateSimilarity(a, b string) float64 {
	// Normalize names
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	
	// Remove common prefixes/suffixes
	a = strings.TrimPrefix(a, "id_")
	b = strings.TrimPrefix(b, "id_")
	
	// More sophisticated similarity calculations would go here
	// For now, a simple approach of checking for direct substrings
	if a == b {
		return 1.0
	}
	if strings.Contains(a, b) || strings.Contains(b, a) {
		return 0.8
	}
	
	// Check for common words
	aWords := strings.Split(a, "_")
	bWords := strings.Split(b, "_")
	
	matchCount := 0
	for _, aWord := range aWords {
		for _, bWord := range bWords {
			if aWord == bWord && len(aWord) > 2 { // Ignore very short words
				matchCount++
			}
		}
	}
	
	if matchCount > 0 {
		return float64(matchCount) / float64(max(len(aWords), len(bWords)))
	}
	
	return 0.0
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
