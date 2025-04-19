package bridge

import (
	"fmt"
)

// CreateGenericMapper creates a mapper that uses a mapping function
func CreateGenericMapper[S, T any](mapFn func(S) (T, error)) DataMapper {
	return &genericMapper[S, T]{mapFn: mapFn}
}

// genericMapper is a generic implementation of DataMapper
type genericMapper[S, T any] struct {
	mapFn func(S) (T, error)
}

// Map converts source data to target data using the mapping function
func (m *genericMapper[S, T]) Map(source interface{}) (interface{}, error) {
	sourceCast, ok := source.(S)
	if !ok {
		var zero T
		return zero, fmt.Errorf("source is not of expected type")
	}
	
	return m.mapFn(sourceCast)
}

// BatchMap maps a slice of source items to a slice of target items
func BatchMap[S, T any](sources []S, mapFn func(S) (T, error)) ([]T, error) {
	targets := make([]T, len(sources))
	
	for i, source := range sources {
		target, err := mapFn(source)
		if err != nil {
			return nil, fmt.Errorf("error mapping item at index %d: %w", i, err)
		}
		targets[i] = target
	}
	
	return targets, nil
}
