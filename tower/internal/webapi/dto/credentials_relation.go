package dto

import (
	"time"
)

type CredentialsRelationCreateRequest struct {
	InitiatorID    string `json:"initiator_id" binding:"required"`
	TargetID       string `json:"target_id" binding:"required"`
	ConnectionType string `json:"type" binding:"required,oneof=inbound outbound bidirectional"`
	Active         bool   `json:"active" binding:"required"`
	Endpoint       string `json:"endpoint" binding:"required"`
}

type CredentialsRelationUpdateRequest struct {
	InitiatorID    string `json:"initiator_id" binding:"required"`
	TargetID       string `json:"target_id" binding:"required"`
	ConnectionType string `json:"connection_type" binding:"required,oneof=inbound outbound bidirectional"`
	Active         bool   `json:"active" binding:"required"`
	Endpoint       string `json:"endpoint" binding:"required"`
}

type CredentialsRelationResponse struct {
	InitiatorID    string    `json:"initiator_id"`
	TargetID       string    `json:"target_id"`
	ConnectionType string    `json:"connection_type"`
	Active         bool      `json:"active"`
	Endpoint       string    `json:"endpoint"`
	CreatedAt      time.Time `json:"created_at"`
}

type CredentialsRelationLogsResponse struct {
	ID          string    `json:"id"`
	InitiatorID string    `json:"initiator_id"`
	TargetID    string    `json:"target_id"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
}
