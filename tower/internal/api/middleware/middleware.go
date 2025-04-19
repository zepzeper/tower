package middleware

import (
	"log"
	"net/http"
)

// CORS middleware adds Cross-Origin Resource Sharing headers
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always set these headers, even for preflight or error responses
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Logging middleware logs HTTP requests
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Recover middleware recovers from panics
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				log.Printf("Panic: %v", err)
				
				// Return an error response
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}
