package middleware

import (
	"net/http"

	"github.com/zepzeper/tower/internal/response"
)

// MethodOnly restricts requests to specified HTTP methods
func MethodOnly(methods ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the request method is allowed
			allowed := false
			for _, method := range methods {
				if r.Method == method {
					allowed = true
					break
				}
			}

			// If method is not allowed, return error
			if !allowed {
				allowedMethods := ""
				for i, method := range methods {
					if i > 0 {
						allowedMethods += ", "
					}
					allowedMethods += method
				}
				w.Header().Set("Allow", allowedMethods)
				response.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// GET restricts requests to GET method
func GET(next http.Handler) http.Handler {
	return MethodOnly(http.MethodGet)(next)
}

// POST restricts requests to POST method
func POST(next http.Handler) http.Handler {
	return MethodOnly(http.MethodPost)(next)
}

// PUT restricts requests to PUT method
func PUT(next http.Handler) http.Handler {
	return MethodOnly(http.MethodPut)(next)
}

// PATCH restricts requests to PATCH method
func PATCH(next http.Handler) http.Handler {
	return MethodOnly(http.MethodPatch)(next)
}

// DELETE restricts requests to DELETE method
func DELETE(next http.Handler) http.Handler {
	return MethodOnly(http.MethodDelete)(next)
}

// ContentType ensures the request has a specific Content-Type
func ContentType(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip empty body requests
			if r.Method == http.MethodGet || r.Method == http.MethodHead {
				next.ServeHTTP(w, r)
				return
			}

			// Check Content-Type
			if r.Header.Get("Content-Type") != contentType {
				response.BadRequest(w, "Content-Type must be "+contentType)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// JSON ensures the request has Content-Type: application/json
func JSON(next http.Handler) http.Handler {
	return ContentType("application/json")(next)
}
