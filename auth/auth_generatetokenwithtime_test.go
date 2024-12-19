package auth

import (
	"os"
	"testing"
	"time"
)

func TestGenerateTokenWithTime(t *testing.T) {

	type testCase struct {
		name        string
		id          uint
		timestamp   time.Time
		setupEnv    func()
		cleanupEnv  func()
		expectError bool
		errorMsg    string
	}

	originalSecret := os.Getenv("JWT_SECRET")

	setJWTSecret := func(secret string) func() {
		return func() {
			os.Setenv("JWT_SECRET", secret)
			jwtSecret = []byte(secret)
		}
	}

	restoreJWTSecret := func() {
		os.Setenv("JWT_SECRET", originalSecret)
		jwtSecret = []byte(originalSecret)
	}

	tests := []testCase{
		{
			name:        "Successful Token Generation",
			id:          1,
			timestamp:   time.Now(),
			setupEnv:    setJWTSecret("test-secret"),
			cleanupEnv:  restoreJWTSecret,
			expectError: false,
		},
		{
			name:        "Zero ID",
			id:          0,
			timestamp:   time.Now(),
			setupEnv:    setJWTSecret("test-secret"),
			cleanupEnv:  restoreJWTSecret,
			expectError: true,
			errorMsg:    "invalid user ID",
		},
		{
			name:        "Future Timestamp",
			id:          1,
			timestamp:   time.Now().Add(24 * time.Hour),
			setupEnv:    setJWTSecret("test-secret"),
			cleanupEnv:  restoreJWTSecret,
			expectError: false,
		},
		{
			name:      "Missing JWT Secret",
			id:        1,
			timestamp: time.Now(),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "")
				jwtSecret = []byte("")
			},
			cleanupEnv:  restoreJWTSecret,
			expectError: true,
			errorMsg:    "JWT secret not configured",
		},
		{
			name:        "Past Timestamp",
			id:          1,
			timestamp:   time.Now().Add(-24 * time.Hour),
			setupEnv:    setJWTSecret("test-secret"),
			cleanupEnv:  restoreJWTSecret,
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			if tc.setupEnv != nil {
				tc.setupEnv()
			}

			if tc.cleanupEnv != nil {
				defer tc.cleanupEnv()
			}

			token, err := GenerateTokenWithTime(tc.id, tc.timestamp)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tc.errorMsg != "" && err.Error() != tc.errorMsg {
					t.Errorf("expected error message %q but got %q", tc.errorMsg, err.Error())
				}
				t.Logf("Successfully caught expected error: %v", err)
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if token == "" {
					t.Error("expected non-empty token but got empty string")
					return
				}
				t.Logf("Successfully generated token: %v", token)
			}
		})
	}

	t.Run("Multiple Sequential Tokens", func(t *testing.T) {
		setJWTSecret("test-secret")()
		defer restoreJWTSecret()

		tokens := make(map[string]bool)
		now := time.Now()

		for i := uint(1); i <= 3; i++ {
			token, err := GenerateTokenWithTime(i, now)
			if err != nil {
				t.Errorf("failed to generate token %d: %v", i, err)
				continue
			}
			if tokens[token] {
				t.Errorf("duplicate token generated: %s", token)
			}
			tokens[token] = true
			t.Logf("Generated unique token %d: %s", i, token)
		}
	})
}
