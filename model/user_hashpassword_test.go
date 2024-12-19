package model

import (
	"strings"
	"testing"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

func TestUserHashPassword(t *testing.T) {
	// Test cases structure
	type testCase struct {
		name     string
		password string
		wantErr  bool
		errMsg   string
	}

	// Define test cases
	tests := []testCase{
		{
			name:     "Scenario 1: Successfully Hash a Valid Password",
			password: "validPassword123",
			wantErr:  false,
			errMsg:   "",
		},
		{
			name:     "Scenario 2: Empty Password Handling",
			password: "",
			wantErr:  true,
			errMsg:   "password should not be empty",
		},
		{
			name:     "Scenario 3: Very Long Password Handling",
			password: strings.Repeat("a", 1024*1024), // 1MB of data
			wantErr:  false,
			errMsg:   "",
		},
		{
			name:     "Scenario 4: Special Characters in Password",
			password: "!@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`",
			wantErr:  false,
			errMsg:   "",
		},
		{
			name:     "Scenario 6: Unicode Password Handling",
			password: "パスワード123アБВ",
			wantErr:  false,
			errMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			user := &User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: tt.password,
			}
			originalPassword := tt.password

			// Act
			err := user.HashPassword()

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err.Error() != tt.errMsg {
					t.Errorf("HashPassword() error message = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			// Additional validations for successful cases
			if !tt.wantErr {
				// Verify password has changed
				if user.Password == originalPassword {
					t.Error("Password was not hashed, still matches original")
					return
				}

				// Verify hash format
				if !strings.HasPrefix(user.Password, "$2a$") {
					t.Error("Generated hash does not appear to be a bcrypt hash")
					return
				}

				// Verify hash can be compared with original password
				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(originalPassword))
				if err != nil {
					t.Errorf("Failed to verify hashed password: %v", err)
					return
				}

				// Log success details
				t.Logf("Successfully hashed password of length %d", utf8.RuneCountInString(originalPassword))
			}
		})
	}

	// Scenario 5: Multiple Hash Calls on Same User
	t.Run("Scenario 5: Multiple Hash Calls on Same User", func(t *testing.T) {
		user := &User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "testPassword123",
		}
		originalPassword := user.Password

		// First hash
		err1 := user.HashPassword()
		if err1 != nil {
			t.Errorf("First HashPassword() failed: %v", err1)
			return
		}
		firstHash := user.Password

		// Second hash
		err2 := user.HashPassword()
		if err2 != nil {
			t.Errorf("Second HashPassword() failed: %v", err2)
			return
		}
		secondHash := user.Password

		// Verify hashes are different
		if firstHash == secondHash {
			t.Error("Multiple hash calls produced identical hashes")
			return
		}

		// Verify both hashes can be validated against original password
		err := bcrypt.CompareHashAndPassword([]byte(firstHash), []byte(originalPassword))
		if err != nil {
			t.Errorf("First hash validation failed: %v", err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(secondHash), []byte(originalPassword))
		if err != nil {
			t.Errorf("Second hash validation failed: %v", err)
		}

		t.Log("Successfully verified multiple hash calls")
	})
}
