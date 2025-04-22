package models

import "time"

type ConnectionType string

const (
	ConnectionTypeInbound       ConnectionType = "inbound"
	ConnectionTypeOutbound      ConnectionType = "outbound"
	ConnectionTypeBidirectional ConnectionType = "bidirectional"
)

type CredentialsRelation struct {
	InitiatorID    string         `json:"initiator_id" db:"initiator_id"`
	TargetID       string         `json:"target_id" db:"target_id"`
	Endpoint       string         `json:"endpoint" db:"endpoint"`
	Active         bool           `json:"active" db:"active"`
	ConnectionType ConnectionType `json:"connection_type" db:"connection_type"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
}
