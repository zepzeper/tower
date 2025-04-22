package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/database/models"
	"golang.org/x/crypto/bcrypt"
)

// Claims represents the JWT claims data structure
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// Service provides authentication-related functionality
type Service struct {
	jwtSecret       string
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
	databaseManager *database.Manager
}

// NewService creates a new auth service
func NewService(jwtSecret string, accessTokenExp, refreshTokenExp time.Duration, databaseManager *database.Manager) *Service {
	return &Service{
		jwtSecret:       jwtSecret,
		accessTokenExp:  accessTokenExp,
		refreshTokenExp: refreshTokenExp,
		databaseManager: databaseManager,
	}
}

// CreateUser creates a new user with the provided details
func (s *Service) CreateUser(email, password, name string) (*models.User, error) {
	// Hash the password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user in the database
	return s.databaseManager.Repos.Auth().CreateUser(email, string(hashedPassword), name)
}

// Authenticate verifies user credentials and returns the user if valid
func (s *Service) Authenticate(email, password string) (*models.User, error) {
	// Get the user by email
	user, err := s.databaseManager.Repos.Auth().GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GenerateTokenPair generates both access and refresh tokens
func (s *Service) GenerateTokenPair(user *models.User) (accessToken, refreshToken string, err error) {
	// Generate access token
	accessToken, err = s.generateAccessToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err = s.generateRefreshToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// generateAccessToken creates a short-lived JWT access token
func (s *Service) generateAccessToken(user *models.User) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenExp)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "tower-api",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// generateRefreshToken creates a long-lived JWT refresh token
func (s *Service) generateRefreshToken(user *models.User) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenExp)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "tower-api",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// RefreshToken validates a refresh token and issues a new access token
func (s *Service) RefreshToken(refreshToken string) (string, error) {
	// Parse and validate the refresh token
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get the user from the database
	user, err := s.databaseManager.Repos.Auth().GetUserByID(claims.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	// Generate a new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

// ValidateToken parses and validates a JWT token
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(id string) (*models.User, error) {
	return s.databaseManager.Repos.Auth().GetUserByID(id)
}

// UserExists checks if a user with the given email exists
func (s *Service) UserExists(email string) (bool, error) {
	return s.databaseManager.Repos.Auth().UserExists(email)
}
