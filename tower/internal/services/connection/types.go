package connection

import "time"

// ConnectionType represents the protocol used for the connection
type ConnectionType string

const (
	REST      ConnectionType = "REST"
	GRPC      ConnectionType = "gRPC"
	WEBSOCKET ConnectionType = "WebSocket"
	GRAPHQL   ConnectionType = "GraphQL"
)

// Relation defines how this connection relates to other systems
type Relation string

const (
	SOURCE        Relation = "SOURCE"        // This API provides data
	DESTINATION   Relation = "DESTINATION"   // This API receives data
	BIDIRECTIONAL Relation = "BIDIRECTIONAL" // Data flows both ways
)

// AuthType defines authentication mechanism
type AuthType string

const (
	BASIC AuthType = "BASIC"
	TOKEN AuthType = "TOKEN"
	OAUTH AuthType = "OAUTH2"
)

// Auth holds authentication details
type Auth struct {
	Type           AuthType
	Username       string
	Password       string
	Token          string
	TokenType      string
	APIKeyName     string
	APIKeyInHeader bool
	OAuth2Config   OAuth2Config
}

type OAuth2Config struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
	AuthURL      string
	RedirectURL  string
	Scope        string
	RefreshToken string
	ExpiresAt    time.Time
}

// Endpoint defines an API endpoint
type Endpoint struct {
	URL        string
	Method     string // GET, POST, etc. for REST
	Parameters map[string]string
}

// Connection represents a connection to another API
type Connection struct {
	ID             string
	Name           string
	Description    string
	Type           ConnectionType
	Relation       Relation
	BaseURL        string
	Endpoints      map[string]Endpoint
	Auth           Auth
	Timeout        int // in seconds
	RetryAttempts  int
	Active         bool
	Headers        map[string]string
	DefaultPayload interface{}
}

// APIManager handles multiple connections
type APIManager struct {
	Connections map[string]*Connection
}
