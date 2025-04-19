// internal/webapi/handlers/handler_registry.go
package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/zepzeper/tower/internal/services"
)

// Registry holds all the handlers for the internal API
type Registry struct {
	dashboardHandler *DashboardHandler
	statsHandler     *StatsHandler
  authHandler      *AuthHandler
}

// NewRegistry creates a new handler registry
func NewRegistry(
	connectorService *services.ConnectorService,
	transformerService *services.TransformerService,
	connectionService *services.ConnectionService,
	authService *services.AuthService,
) *Registry {
	return &Registry{
		dashboardHandler: NewDashboardHandler(connectionService),
		statsHandler:     NewStatsHandler(connectorService, transformerService, connectionService),
    authHandler:      NewAuthHandler(authService),
	}
}

func (r *Registry) RegisterRoutes(router chi.Router) {
    router.Route("/dashboard", func(router chi.Router) {
        router.Get("/summary", r.dashboardHandler.GetSummary)
        router.Get("/recent-activity", r.dashboardHandler.GetRecentActivity)
    })

    router.Route("/stats", func(router chi.Router) {
        router.Get("/connections", r.statsHandler.GetConnectionStats)
        router.Get("/executions", r.statsHandler.GetExecutionStats)
    })
}
