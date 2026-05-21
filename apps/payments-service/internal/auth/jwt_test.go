package auth

import (
	"testing"
	"time"
)

func TestNewJWTManager(t *testing.T) {
	tests := []struct {
		name          string
		secretKey     string
		tokenDuration time.Duration
	}{
		{"default duration", "testsecret", 24 * time.Hour},
		{"short duration", "secret123", 1 * time.Hour},
		{"long duration", "longsecret", 7 * 24 * time.Hour},
		{"very short duration", "quicksecret", 5 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewJWTManager(tt.secretKey, tt.tokenDuration)
			if manager == nil {
				t.Error("NewJWTManager returned nil")
			}
			if manager.secretKey != tt.secretKey {
				t.Errorf("Secret key mismatch: got %v, want %v", manager.secretKey, tt.secretKey)
			}
			if manager.tokenDuration != tt.tokenDuration {
				t.Errorf("Token duration mismatch: got %v, want %v", manager.tokenDuration, tt.tokenDuration)
			}
		})
	}
}

func TestJWTManagerGenerate(t *testing.T) {
	manager := NewJWTManager("test-secret-key", 24*time.Hour)

	tests := []struct {
		name   string
		userID string
		email  string
	}{
		{"user 1", "user-001", "user1@example.com"},
		{"user 2", "user-002", "user2@example.com"},
		{"user 3", "user-003", "user3@test.org"},
		{"user with long id", "user-123456789", "longid@example.com"},
		{"user with special email", "user-005", "special+tag@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := manager.Generate(tt.userID, tt.email)
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}
			if token == "" {
				t.Error("Generated token is empty")
			}
		})
	}
}

func TestJWTManagerVerify(t *testing.T) {
	secretKey := "test-secret-key"
	manager := NewJWTManager(secretKey, 24*time.Hour)

	// Generate a valid token
	userID := "test-user"
	email := "test@example.com"
	token, err := manager.Generate(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Test verification
	claims, err := manager.Verify(token)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID mismatch: got %v, want %v", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("Email mismatch: got %v, want %v", claims.Email, email)
	}
}

func TestJWTManagerVerifyInvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", 24*time.Hour)

	tests := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"invalid format", "notavalidtoken"},
		{"malformed jwt", "header.payload"},
		{"random string", "xyz123abc456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := manager.Verify(tt.token)
			if err == nil {
				t.Error("Verify should fail for invalid token")
			}
		})
	}
}

func TestJWTManagerVerifyExpiredToken(t *testing.T) {
	// Create manager with very short duration
	manager := NewJWTManager("test-secret-key", -1*time.Hour) // Already expired

	token, err := manager.Generate("user-001", "user@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	_, err = manager.Verify(token)
	if err == nil {
		t.Error("Verify should fail for expired token")
	}
}

func TestJWTManagerVerifyDifferentSecret(t *testing.T) {
	manager1 := NewJWTManager("secret1", 24*time.Hour)
	manager2 := NewJWTManager("secret2", 24*time.Hour)

	token, err := manager1.Generate("user-001", "user@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to verify with different secret
	_, err = manager2.Verify(token)
	if err == nil {
		t.Error("Verify should fail when using different secret key")
	}
}

func TestJWTTokenLifecycle(t *testing.T) {
	manager := NewJWTManager("test-secret", 1*time.Hour)

	tests := []struct {
		name   string
		userID string
		email  string
	}{
		{"lifecycle user 1", "u1", "u1@test.com"},
		{"lifecycle user 2", "u2", "u2@test.com"},
		{"lifecycle user 3", "u3", "u3@test.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate
			token, err := manager.Generate(tt.userID, tt.email)
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}

			// Verify
			claims, err := manager.Verify(token)
			if err != nil {
				t.Fatalf("Verify failed: %v", err)
			}

			// Check claims
			if claims.UserID != tt.userID {
				t.Errorf("UserID mismatch: got %v, want %v", claims.UserID, tt.userID)
			}
			if claims.Email != tt.email {
				t.Errorf("Email mismatch: got %v, want %v", claims.Email, tt.email)
			}

			// Check timestamps
			if claims.ExpiresAt == nil {
				t.Error("ExpiresAt should not be nil")
			}
			if claims.IssuedAt == nil {
				t.Error("IssuedAt should not be nil")
			}
			if claims.NotBefore == nil {
				t.Error("NotBefore should not be nil")
			}
		})
	}
}
