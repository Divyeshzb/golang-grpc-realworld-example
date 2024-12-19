package model

import (
	"testing"
	"golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/gorm"
	"time"
)

func TestUserCheckPassword(t *testing.T) {
	// Helper function to create bcrypt hash
	createHash := func(password string) string {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("Failed to create hash: %v", err)
		}
		return string(hash)
	}

	// Test cases
	tests := []struct {
		name           string
		storedHash    string
		inputPassword string
		expected      bool
		description   string
	}{
		{
			name:           "Valid Password Match",
			storedHash:    createHash("correctPassword123"),
			inputPassword: "correctPassword123",
			expected:      true,
			description:   "Testing correct password validation",
		},
		{
			name:           "Invalid Password Mismatch",
			storedHash:    createHash("correctPassword123"),
			inputPassword: "wrongPassword123",
			expected:      false,
			description:   "Testing incorrect password rejection",
		},
		{
			name:           "Empty Password Check",
			storedHash:    createHash("somePassword"),
			inputPassword: "",
			expected:      false,
			description:   "Testing empty password handling",
		},
		{
			name:           "Empty Stored Hash",
			storedHash:    "",
			inputPassword: "anyPassword",
			expected:      false,
			description:   "Testing empty stored hash handling",
		},
		{
			name:           "Very Long Password Input",
			storedHash:    createHash("normal"),
			inputPassword: string(make([]byte, 1024*1024)), // 1MB string
			expected:      false,
			description:   "Testing very long password input",
		},
		{
			name:           "Special Characters Password",
			storedHash:    createHash("!@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`"),
			inputPassword: "!@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`",
			expected:      true,
			description:   "Testing special characters in password",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.description)

			// Create user instance with test data
			user := &User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "testuser",
				Email:    "test@example.com",
				Password: tt.storedHash,
				Bio:      "Test bio",
				Image:    "test.jpg",
			}

			// Execute test
			result := user.CheckPassword(tt.inputPassword)

			// Assert result
			if result != tt.expected {
				t.Errorf("CheckPassword() = %v, want %v", result, tt.expected)
			}

			// Log test outcome
			if result {
				t.Log("Password validation successful")
			} else {
				t.Log("Password validation failed as expected")
			}
		})
	}
}
