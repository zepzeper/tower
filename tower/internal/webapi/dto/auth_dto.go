package dto

type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email,min=8"`
    Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
    Token string `json:"token"`
    User  UserResponse `json:"user"`
}

type UserResponse struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at,omitempty"`
}

type RegisterResponse struct {
    Message string       `json:"message"`
    User    UserResponse `json:"user"`
}
