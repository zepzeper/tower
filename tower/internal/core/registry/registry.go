package registry

import (
	"sync"
	
	"github.com/zepzeper/tower/internal/core/connectors"
)

// ConnectorRegistry manages all available connectors
type ConnectorRegistry struct {
	connectors map[string]connectors.Connector
	mu         sync.RWMutex
}

// NewConnectorRegistry creates a new connector registry
func NewConnectorRegistry() *ConnectorRegistry {
	return &ConnectorRegistry{
		connectors: make(map[string]connectors.Connector),
	}
}

// Register adds a connector to the registry
func (r *ConnectorRegistry) Register(name string, connector connectors.Connector) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.connectors[name] = connector
}

// Get retrieves a connector by name
func (r *ConnectorRegistry) Get(name string) (connectors.Connector, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, exists := r.connectors[name]
	return c, exists
}

// List returns all registered connector names
func (r *ConnectorRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	names := make([]string, 0, len(r.connectors))
	for name := range r.connectors {
		names = append(names, name)
	}
	
	return names
}
