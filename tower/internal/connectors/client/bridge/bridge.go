package bridge

import (
	"fmt"
	"reflect"
)

// DataMapper defines how to map data between different client models
type DataMapper interface {
	// Map converts source data to target data
	Map(source interface{}) (interface{}, error)
}

// Bridge connects two API clients and facilitates data transfer
type Bridge[S, T any] struct {
	SourceClient S
	TargetClient T
	Mappers      map[string]DataMapper
}

// NewBridge creates a new bridge between two API clients
func NewBridge[S, T any](sourceClient S, targetClient T) *Bridge[S, T] {
	return &Bridge[S, T]{
		SourceClient: sourceClient,
		TargetClient: targetClient,
		Mappers:      make(map[string]DataMapper),
	}
}

// RegisterMapper registers a data mapper for a specific data type
func (b *Bridge[S, T]) RegisterMapper(sourceType string, mapper DataMapper) {
	b.Mappers[sourceType] = mapper
}

// Transfer transfers data from source client to target client
func (b *Bridge[S, T]) Transfer(sourceMethod, targetMethod string, args ...interface{}) error {
	// Get the source method by reflection
	sourceClientValue := reflect.ValueOf(b.SourceClient)
	sourceMethodValue := sourceClientValue.MethodByName(sourceMethod)
	
	if !sourceMethodValue.IsValid() {
		return fmt.Errorf("source method %s not found", sourceMethod)
	}
	
	// Convert args to reflect.Values
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		reflectArgs[i] = reflect.ValueOf(arg)
	}
	
	// Call the source method
	sourceResults := sourceMethodValue.Call(reflectArgs)
	
	// Check for error (assuming the last return value is an error)
	if len(sourceResults) > 1 {
		lastResult := sourceResults[len(sourceResults)-1].Interface()
		if lastResult != nil {
			if err, ok := lastResult.(error); ok && err != nil {
				return fmt.Errorf("error calling source method: %w", err)
			}
		}
	}
	
	// Get the source data (assuming the first return value is the data)
	sourceData := sourceResults[0].Interface()
	
	// Determine the source type
	sourceType := reflect.TypeOf(sourceData).String()
	
	// Find the mapper for this source type
	mapper, ok := b.Mappers[sourceType]
	if !ok {
		return fmt.Errorf("no mapper registered for source type %s", sourceType)
	}
	
	// Map the source data to target data
	targetData, err := mapper.Map(sourceData)
	if err != nil {
		return fmt.Errorf("error mapping data: %w", err)
	}
	
	// Get the target method by reflection
	targetClientValue := reflect.ValueOf(b.TargetClient)
	targetMethodValue := targetClientValue.MethodByName(targetMethod)
	
	if !targetMethodValue.IsValid() {
		return fmt.Errorf("target method %s not found", targetMethod)
	}
	
	// Call the target method with the mapped data
	targetArgs := []reflect.Value{reflect.ValueOf(targetData)}
	targetResults := targetMethodValue.Call(targetArgs)
	
	// Check for error (assuming the last return value is an error)
	if len(targetResults) > 0 {
		lastResult := targetResults[len(targetResults)-1].Interface()
		if lastResult != nil {
			if err, ok := lastResult.(error); ok && err != nil {
				return fmt.Errorf("error calling target method: %w", err)
			}
		}
	}
	
	return nil
}
