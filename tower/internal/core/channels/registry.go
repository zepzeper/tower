package channels

import (
	"errors"
	"sync"

	"github.com/yourusername/api-middleware/internal/core/channels/models"
)

// Registry is a service that manages channel types
type Registry struct {
	channelTypes map[string]models.ChannelType
	mu           sync.RWMutex
}

// NewRegistry creates a new channel registry
func NewRegistry() *Registry {
	return &Registry{
		channelTypes: make(map[string]models.ChannelType),
	}
}

// RegisterChannelType adds a channel type to the registry
func (r *Registry) RegisterChannelType(channelType models.ChannelType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.channelTypes[channelType.ID]; exists {
		return errors.New("channel type already registered")
	}

	r.channelTypes[channelType.ID] = channelType
	return nil
}

// GetChannelType returns a channel type by ID
func (r *Registry) GetChannelType(id string) (models.ChannelType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	channelType, exists := r.channelTypes[id]
	if !exists {
		return models.ChannelType{}, errors.New("channel type not found")
	}

	return channelType, nil
}

// GetAllChannelTypes returns all registered channel types
func (r *Registry) GetAllChannelTypes() []models.ChannelType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]models.ChannelType, 0, len(r.channelTypes))
	for _, channelType := range r.channelTypes {
		types = append(types, channelType)
	}

	return types
}

// RemoveChannelType removes a channel type from the registry
func (r *Registry) RemoveChannelType(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.channelTypes[id]; !exists {
		return errors.New("channel type not found")
	}

	delete(r.channelTypes, id)
	return nil
}
