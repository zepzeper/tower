package credentials

// AuthType defines authentication mechanism
type AuthType string

const (
	BASIC AuthType = "BASIC"
	TOKEN AuthType = "TOKEN"
	OAUTH AuthType = "OAUTH2"
)
