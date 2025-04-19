// internal/webapi/handlers/stats_handler.go
package handlers

import (
	"net/http"
	"time"

	"github.com/zepzeper/tower/internal/response"
	"github.com/zepzeper/tower/internal/services"
	"github.com/zepzeper/tower/internal/webapi/dto"
)

// StatsHandler handles stats-related internal API endpoints
type StatsHandler struct {
	connectorService   *services.ConnectorService
	transformerService *services.TransformerService
	connectionService  *services.ConnectionService
}

// NewStatsHandler creates a new stats handler
func NewStatsHandler(
	connectorService *services.ConnectorService,
	transformerService *services.TransformerService,
	connectionService *services.ConnectionService,
) *StatsHandler {
	return &StatsHandler{
		connectorService:   connectorService,
		transformerService: transformerService,
		connectionService:  connectionService,
	}
}

// GetConnectionStats handles GET /internal/stats/connections
func (h *StatsHandler) GetConnectionStats(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, this would get data from services
	
	// For now, we'll return mock data
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	
	stats := dto.ConnectionStatsResponse{
		TotalCount:    10,
		ActiveCount:   7,
		InactiveCount: 3,
		ConnectorDistribution: map[string]int{
			"woocommerce": 4,
			"mirakl":      3,
			"kaufland":    2,
			"shopify":     1,
		},
		RecentlyCreated: []dto.ConnectionStatSummary{
			{
				ID:        "conn1",
				Name:      "WooCommerce to Mirakl",
				SourceID:  "woocommerce",
				TargetID:  "mirakl",
				Active:    true,
				LastRun:   &yesterday,
				CreatedAt: yesterday,
			},
			{
				ID:        "conn2",
				Name:      "Shopify to Kaufland",
				SourceID:  "shopify",
				TargetID:  "kaufland",
				Active:    true,
				LastRun:   &now,
				CreatedAt: yesterday,
			},
		},
	}
	
	response.JSON(w, stats, http.StatusOK)
}

// GetExecutionStats handles GET /internal/stats/executions
func (h *StatsHandler) GetExecutionStats(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, this would get data from services
	
	// For now, we'll return mock data
	stats := dto.ExecutionStatsResponse{
		TotalCount:   165,
		SuccessCount: 150,
		FailureCount: 12,
		PendingCount: 3,
		TimeDistribution: map[string]int{
			"24h": 25,
			"3d":  78,
			"7d":  165,
		},
	}
	
	response.JSON(w, stats, http.StatusOK)
}
