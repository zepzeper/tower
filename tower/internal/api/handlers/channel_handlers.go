package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/zepzeper/tower/internal/api/response"
	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/database/models"
)

// channelRequest represents the data expected in channel API requests
type channelRequest struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
}

// ChannelHandler handles channel-related API endpoints
type ChannelHandler struct {
	db *db.Manager
}

// NewChannelHandler creates a new channel handler
func NewChannelHandler(dbManager *db.Manager) *ChannelHandler {
	return &ChannelHandler{
		db: dbManager,
	}
}

// ListChannels handles GET /api/v1/channels
func (h *ChannelHandler) ListChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := h.db.Repos.Channel().GetAll()
	if err != nil {
		response.Error(w, "Failed to retrieve channels", http.StatusInternalServerError)
		return
	}
	
	// Convert to API response format
	result := make([]interface{}, len(channels))
	for i, channel := range channels {
		result[i] = channel.ToAPIChannel()
	}
	
	response.JSON(w, result, http.StatusOK)
}

// GetChannel handles GET /api/v1/channels/{channelID}
func (h *ChannelHandler) GetChannel(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	
	channel, err := h.db.Repos.Channel().GetByID(channelID)
	if err != nil {
		response.Error(w, "Channel not found", http.StatusNotFound)
		return
	}
	
	response.JSON(w, channel.ToAPIChannel(), http.StatusOK)
}

// CreateChannel handles POST /api/v1/channels
func (h *ChannelHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	var req channelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.Name == "" || req.Type == "" {
		response.Error(w, "Name and type are required", http.StatusBadRequest)
		return
	}
	
	// Convert config to JSON
	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		response.Error(w, "Invalid config format", http.StatusBadRequest)
		return
	}
	
	// Create channel model
	channel := models.Channel{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Type:      req.Type,
		Config:    configJSON,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Set description if provided
	if req.Description != "" {
		channel.Description.String = req.Description
		channel.Description.Valid = true
	}
	
	// Save to database
	if err := h.db.Repos.Channel().Create(channel); err != nil {
		response.Error(w, "Failed to create channel", http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, channel.ToAPIChannel(), http.StatusCreated)
}

// UpdateChannel handles PUT /api/v1/channels/{channelID}
func (h *ChannelHandler) UpdateChannel(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	
	// Check if channel exists
	channel, err := h.db.Repos.Channel().GetByID(channelID)
	if err != nil {
		response.Error(w, "Channel not found", http.StatusNotFound)
		return
	}
	
	// Parse request body
	var req channelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.Name == "" || req.Type == "" {
		response.Error(w, "Name and type are required", http.StatusBadRequest)
		return
	}
	
	// Convert config to JSON
	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		response.Error(w, "Invalid config format", http.StatusBadRequest)
		return
	}
	
	// Update channel fields
	channel.Name = req.Name
	channel.Type = req.Type
	channel.Config = configJSON
	channel.UpdatedAt = time.Now()
	
	// Update description if provided
	if req.Description != "" {
		channel.Description.String = req.Description
		channel.Description.Valid = true
	} else {
		channel.Description.Valid = false
	}
	
	// Save to database
	if err := h.db.Repos.Channel().Update(channel); err != nil {
		response.Error(w, "Failed to update channel", http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, channel.ToAPIChannel(), http.StatusOK)
}

// DeleteChannel handles DELETE /api/v1/channels/{channelID}
func (h *ChannelHandler) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	
	// Delete from database
	if err := h.db.Repos.Channel().Delete(channelID); err != nil {
		response.Error(w, "Failed to delete channel", http.StatusInternalServerError)
		return
	}
	
	response.JSON(w, map[string]interface{}{
		"success": true,
		"message": "Channel deleted successfully",
	}, http.StatusOK)
}
