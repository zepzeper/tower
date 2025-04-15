package response

import (
	"encoding/json"
	"net/http"
)

// Response is a standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// JSON sends a JSON response
func JSON(w http.ResponseWriter, data interface{}, status int) {
	response := Response{
		Success: status >= 200 && status < 300,
		Data:    data,
	}

	jsonResponse(w, response, status)
}

// Error sends an error response
func Error(w http.ResponseWriter, message string, status int) {
	response := Response{
		Success: false,
		Error:   message,
	}

	jsonResponse(w, response, status)
}

// Paginated sends a paginated response
func Paginated(w http.ResponseWriter, data interface{}, page, limit, total int) {
	meta := map[string]interface{}{
		"page":     page,
		"limit":    limit,
		"total":    total,
		"lastPage": (total + limit - 1) / limit,
	}

	response := Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	}

	jsonResponse(w, response, http.StatusOK)
}

// jsonResponse sends a JSON response with the given status code
func jsonResponse(w http.ResponseWriter, response interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	json.NewEncoder(w).Encode(response)
}
