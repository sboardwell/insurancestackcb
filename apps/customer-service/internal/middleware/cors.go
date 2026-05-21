package middleware

import (
	"github.com/rs/cors"
)

// NewCORS creates a new CORS middleware with appropriate settings
func NewCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // In production, specify exact origins
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-User-ID",
		},
		ExposedHeaders: []string{
			"Link",
		},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	})
}
