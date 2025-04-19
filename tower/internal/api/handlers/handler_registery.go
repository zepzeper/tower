package handlers

import (
	"github.com/go-chi/chi/v5"
)

// Registry holds all the handlers for the API
type Registry struct {
}

// NewRegistry creates a new handler registry
func NewRegistry(
) *Registry {
	return &Registry{
	}
}

// RegisterRoutes registers all API routes with the provided router
func (r *Registry) RegisterRoutes(router chi.Router) {
}
