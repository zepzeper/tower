package response

import (
	"encoding/json"
	"net/http"
)

// Standard response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSON sends a JSON response with the provided data and status code
func JSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// If data is already a Response type, use it directly
	if resp, ok := data.(Response); ok {
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Otherwise, wrap it in a success response
	response := Response{
		Success: true,
		Data:    data,
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Success sends a success response with a message
func Success(w http.ResponseWriter, message string, status int) {
	JSON(w, Response{
		Success: true,
		Message: message,
	}, status)
}

// Error sends an error response with the provided message and status code
func Error(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: false,
		Error:   message,
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Created is a helper for 201 Created responses
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, data, http.StatusCreated)
}

// OK is a helper for 200 OK responses
func OK(w http.ResponseWriter, data interface{}) {
	JSON(w, data, http.StatusOK)
}

// NoContent is a helper for 204 No Content responses
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// BadRequest is a helper for 400 Bad Request responses
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, message, http.StatusBadRequest)
}

// Unauthorized is a helper for 401 Unauthorized responses
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, message, http.StatusUnauthorized)
}

// Forbidden is a helper for 403 Forbidden responses
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, message, http.StatusForbidden)
}

// NotFound is a helper for 404 Not Found responses
func NotFound(w http.ResponseWriter, message string) {
	Error(w, message, http.StatusNotFound)
}

// InternalServerError is a helper for 500 Internal Server Error responses
func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, message, http.StatusInternalServerError)
}
