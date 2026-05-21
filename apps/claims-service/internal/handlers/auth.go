package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/CB-InsuranceStack/InsuranceStack/apps/claims-service/internal/auth"
	"github.com/sirupsen/logrus"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	jwtManager *auth.JWTManager
	logger     *logrus.Logger
	// In production, these would come from a database
	validUsername string
	validPassword string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(logger *logrus.Logger) *AuthHandler {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-change-in-production"
		logger.Warn("JWT_SECRET not set, using default (not secure for production)")
	}

	// Get credentials from environment
	username := os.Getenv("AUTH_USERNAME")
	if username == "" {
		username = "demo@insurancestack.com"
	}

	password := os.Getenv("AUTH_PASSWORD")
	if password == "" {
		password = "demo123"
	}

	// Hash the password for comparison
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		logger.WithError(err).Fatal("Failed to hash password")
	}

	return &AuthHandler{
		jwtManager:    auth.NewJWTManager(jwtSecret, 24*time.Hour),
		logger:        logger,
		validUsername: username,
		validPassword: hashedPassword,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"` // seconds
	User      User   `json:"user"`
}

// User represents basic user info
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode login request")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate credentials
	if req.Username != h.validUsername {
		h.logger.WithField("username", req.Username).Warn("Invalid username")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := auth.VerifyPassword(h.validPassword, req.Password); err != nil {
		h.logger.WithField("username", req.Username).Warn("Invalid password")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.Generate("user-001", req.Username)
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate token")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token:     token,
		ExpiresIn: 86400, // 24 hours in seconds
		User: User{
			ID:    "user-001",
			Email: req.Username,
			Name:  "Demo User",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	h.logger.WithField("username", req.Username).Info("User logged in successfully")
}
