package auth

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"simple password", "password123", false},
		{"complex password", "P@ssw0rd!2023", false},
		{"long password", "ThisIsAVeryLongPasswordWithManyCharacters123!@#", false},
		{"short password", "pass", false},
		{"empty password", "", false},
		{"unicode password", "パスワード", false},
		{"special chars", "!@#$%^&*()", false},
		{"numbers only", "12345678", false},
		{"letters only", "abcdefgh", false},
		{"mixed case", "AbCdEfGh", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}
				if hash == tt.password {
					t.Error("HashPassword() returned plain password instead of hash")
				}
				if !strings.HasPrefix(hash, "$2a$") {
					t.Error("HashPassword() did not return bcrypt hash")
				}
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	// Pre-generate some hashes for testing
	validHash, _ := HashPassword("correctpassword")

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		wantErr        bool
	}{
		{"correct password", validHash, "correctpassword", false},
		{"incorrect password", validHash, "wrongpassword", true},
		{"empty password", validHash, "", true},
		{"case sensitive", validHash, "CorrectPassword", true},
		{"extra characters", validHash, "correctpassword123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyPassword(tt.hashedPassword, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	password := "testpassword"

	hash1, err1 := HashPassword(password)
	if err1 != nil {
		t.Fatalf("Failed to hash password: %v", err1)
	}

	hash2, err2 := HashPassword(password)
	if err2 != nil {
		t.Fatalf("Failed to hash password: %v", err2)
	}

	// Hashes should be different due to random salt
	if hash1 == hash2 {
		t.Error("Two hashes of the same password should be different")
	}

	// But both should verify correctly
	if err := VerifyPassword(hash1, password); err != nil {
		t.Error("First hash failed to verify")
	}
	if err := VerifyPassword(hash2, password); err != nil {
		t.Error("Second hash failed to verify")
	}
}

func TestPasswordRoundTrip(t *testing.T) {
	tests := []string{
		"simple",
		"complex!@#123",
		"VeryLongPasswordWith123Numbers!",
		"短い",
		"مرحبا",
		"😀🎉",
	}

	for _, password := range tests {
		t.Run("roundtrip_"+password, func(t *testing.T) {
			hash, err := HashPassword(password)
			if err != nil {
				t.Fatalf("HashPassword failed: %v", err)
			}

			if err := VerifyPassword(hash, password); err != nil {
				t.Errorf("VerifyPassword failed for password %q: %v", password, err)
			}
		})
	}
}

func TestVerifyPasswordWithInvalidHash(t *testing.T) {
	tests := []struct {
		name string
		hash string
	}{
		{"empty hash", ""},
		{"invalid format", "notahash"},
		{"partial hash", "$2a$10$incomplete"},
		{"wrong algorithm", "$2b$10$somehashvalue"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyPassword(tt.hash, "anypassword")
			if err == nil {
				t.Error("VerifyPassword should fail with invalid hash")
			}
		})
	}
}
