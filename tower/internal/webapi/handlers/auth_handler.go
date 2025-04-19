package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zepzeper/tower/internal/response"
	"github.com/zepzeper/tower/internal/services"
	"github.com/zepzeper/tower/internal/webapi/dto"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        response.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
        return
    }

    var req dto.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, "Failed to parse request body", http.StatusBadRequest)
        return
    }

    // Validate request (you can use validator package)
    if len(req.Email) < 8 || len(req.Password) < 8 {
        response.Error(w, "Email and password must be at least 8 characters", http.StatusNotAcceptable)
        return
    }

    // Check if user exists in db (via authService)
    exists, err := h.authService.UserExists(req.Email)
    if err != nil {
        response.Error(w, "Failed to check user existence", http.StatusInternalServerError)
        return
    }
    if exists {
        response.Error(w, "Email already registered", http.StatusConflict)
        return
    }

    // Create user via authService
    user, err := h.authService.CreateUser(req.Email, req.Password)
    if err != nil {
        response.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    // Generate token
    _, err = h.authService.GenerateToken(user.ID)
    if err != nil {
        response.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    // Return success response
    response.JSON(w, &dto.RegisterResponse{
        Message: "User registered successfully",
        User: dto.UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            CreatedAt: user.CreatedAt.Format(time.RFC3339),
        },
    }, http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        response.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
        return
    }

    var req dto.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, "Failed to parse request body", http.StatusBadRequest)
        return
    }

    // Authenticate user
    user, err := h.authService.Authenticate(req.Email, req.Password)
    if err != nil {
        response.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Generate token
    token, err := h.authService.GenerateToken(user.ID)
    if err != nil {
        response.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    // Return success response
    response.JSON(w, &dto.AuthResponse{
        Token: token,
        User: dto.UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            CreatedAt: user.CreatedAt.Format(time.RFC3339),
        },
    }, http.StatusOK)
}
