package handlers

import (
	"net/http"
	"time"

	"github.com/zepzeper/tower/internal/response"
	"github.com/zepzeper/tower/internal/services"
	"github.com/zepzeper/tower/internal/webapi/dto"
)

// DashboardHandler handles dashboard-related internal API endpoints
type DashboardHandler struct {
	connectionService *services.ConnectionService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(connectionService *services.ConnectionService) *DashboardHandler {
	return &DashboardHandler{
		connectionService: connectionService,
	}
}

// GetSummary handles GET /internal/dashboard/summary
func (h *DashboardHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	
	// For now, we'll return mock data
	summary := dto.SummaryResponse{
		TotalConnections:    10,
		ActiveConnections:   7,
		SuccessfulTransfers: 150,
		FailedTransfers:     12,
		PendingTransfers:    3,
	}
	
	response.JSON(w, summary, http.StatusOK)
}

// GetRecentActivity handles GET /internal/dashboard/recent-activity
func (h *DashboardHandler) GetRecentActivity(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, get recent executions from the service
	
	// For now, return mock data
	now := time.Now()
	tenMinutesAgo := now.Add(-10 * time.Minute)
	fiveMinutesAgo := now.Add(-5 * time.Minute)
	
	activity := dto.RecentActivityResponse{
		Executions: []dto.ExecutionSummary{
			{
				ID:             "exec123",
				ConnectionID:   "conn1",
				ConnectionName: "WooCommerce to Mirakl",
				Status:         "success",
				StartTime:      tenMinutesAgo,
				EndTime:        &fiveMinutesAgo,
			},
			{
				ID:             "exec124",
				ConnectionID:   "conn2",
				ConnectionName: "Shopify to Kaufland",
				Status:         "failed",
				StartTime:      fiveMinutesAgo,
				EndTime:        &now,
			},
			{
				ID:             "exec125",
				ConnectionID:   "conn3",
				ConnectionName: "Magento to WooCommerce",
				Status:         "in_progress",
				StartTime:      now,
				EndTime:        nil,
			},
		},
	}
	
	response.JSON(w, activity, http.StatusOK)
}
