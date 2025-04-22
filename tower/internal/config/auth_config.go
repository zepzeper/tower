package config

import (
	"time"
)

// AuthConfig represents the configuration for the authentication service
type AuthConfig struct {
	JWTSecret          string        `mapstructure:"jwt_secret"`
	AccessTokenExpiry  time.Duration `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry"`
}

// DefaultAuthConfig returns the default configuration for the auth service
func DefaultAuthConfig() AuthConfig {
	return AuthConfig{
		JWTSecret:          "your-secret-key-change-in-production",
		AccessTokenExpiry:  15 * time.Minute,    // 15 minutes
		RefreshTokenExpiry: 30 * 24 * time.Hour, // 30 days
	}
}
