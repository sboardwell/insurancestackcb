package middleware

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const userIDKey contextKey = "userID"

// AuthMiddleware extracts user ID from X-User-ID header (simplified for demo)
func AuthMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health check
			if r.URL.Path == "/healthz" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract user ID from X-User-ID header (demo purposes)
			userID := r.Header.Get("X-User-ID")
			if userID == "" {
				userID = "cust-001" // Default for demo
			}

			// Add user ID to request context
			ctx := context.WithValue(r.Context(), userIDKey, userID)

			logger.WithField("userId", userID).Debug("User authenticated")

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts the user ID from the request context
func GetUserID(r *http.Request) string {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		return "cust-001" // Default fallback
	}
	return userID
}
