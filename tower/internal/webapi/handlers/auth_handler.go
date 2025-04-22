package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/zepzeper/tower/internal/response"
	"github.com/zepzeper/tower/internal/services/authentication"
	"github.com/zepzeper/tower/internal/webapi/dto"
)

type AuthHandler struct {
	authService *auth.Service
}

func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Failed to parse request body")
		return
	}

	// Validate request
	if err := validateRegisterRequest(req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Check if user exists in db
	exists, err := h.authService.UserExists(req.Email)
	if err != nil {
		response.InternalServerError(w, "Failed to check user existence")
		return
	}
	if exists {
		response.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Create user
	user, err := h.authService.CreateUser(req.Email, req.Password, req.Name)
	if err != nil {
		response.InternalServerError(w, "Failed to create user")
		return
	}

	// Generate token pair
	accessToken, refreshToken, err := h.authService.GenerateTokenPair(user)
	if err != nil {
		response.InternalServerError(w, "Failed to generate token")
		return
	}

	// Set refresh token as HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil, // Secure if using HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(30 * 24 * 60 * 60), // 30 days in seconds
	})

	// Return success response
	response.Created(w, &dto.AuthResponse{
		Token: accessToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Failed to parse request body")
		return
	}

	// Validate request
	if err := validateLoginRequest(req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Authenticate user
	user, err := h.authService.Authenticate(req.Email, req.Password)
	if err != nil {
		response.Unauthorized(w, "Invalid email or password")
		return
	}

	// Generate token pair
	accessToken, refreshToken, err := h.authService.GenerateTokenPair(user)
	if err != nil {
		response.InternalServerError(w, "Failed to generate token")
		return
	}

	// Set refresh token as HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil, // Secure if using HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(30 * 24 * 60 * 60), // 30 days in seconds
	})

	// Return success response with access token
	response.OK(w, &dto.AuthResponse{
		Token: accessToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	})
}

// RefreshToken handles token refreshing
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		response.Unauthorized(w, "Refresh token not found")
		return
	}

	refreshToken := cookie.Value

	// Refresh the token
	accessToken, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		response.Unauthorized(w, "Invalid refresh token")
		return
	}

	// Return new access token
	response.OK(w, &dto.RefreshTokenResponse{
		Token: accessToken,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Expire immediately
	})

	response.Success(w, "Logged out successfully", http.StatusOK)
}

// GetMe returns the current authenticated user's details
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Unauthorized(w, "User not authenticated")
		return
	}

	// Get user from database
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		response.InternalServerError(w, "Failed to get user details")
		return
	}

	if user == nil {
		response.NotFound(w, "User not found")
		return
	}

	// Return user details
	response.OK(w, &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	})
}

// AuthMiddleware verifies the JWT token and adds user context
func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(w, "Authorization header is required")
			return
		}

		// Check if the Authorization header has the correct format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(w, "Authorization header must be in format: Bearer {token}")
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := h.authService.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(w, "Invalid or expired token")
			return
		}

		// Add user ID to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper functions for validation
func validateRegisterRequest(req dto.RegisterRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	return nil
}

func validateLoginRequest(req dto.LoginRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func isValidEmail(email string) bool {
	// Simple email validation for example purposes
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
