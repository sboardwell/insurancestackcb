package middleware

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const customerIDKey contextKey = "customerID"

// AuthMiddleware extracts customer ID from X-User-ID header (simplified for demo)
func AuthMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health check
			if r.URL.Path == "/healthz" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract customer ID from X-User-ID header (demo purposes)
			customerID := r.Header.Get("X-User-ID")
			if customerID == "" {
				customerID = "cust-001" // Default for demo
			}

			// Add customer ID to request context
			ctx := context.WithValue(r.Context(), customerIDKey, customerID)

			logger.WithField("customerId", customerID).Debug("Customer authenticated")

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts the customer ID from the request context
// Named GetUserID for backwards compatibility with handlers
func GetUserID(r *http.Request) string {
	customerID, ok := r.Context().Value(customerIDKey).(string)
	if !ok {
		return "cust-001" // Default fallback
	}
	return customerID
}
