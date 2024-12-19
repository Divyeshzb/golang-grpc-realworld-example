package auth

import (
	"os"
	"testing"
	"time"
	"math"
	"sync"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)


func TestgenerateToken(t *testing.T) {
	type testCase struct {
		name        string
		userID      uint
		currentTime time.Time
		setupEnv    func()
		wantErr     bool
		errMsg      string
	}

	validateToken := func(t *testing.T, tokenString string, expectedID uint, expectedTime time.Time) bool {
		token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil {
			return false
		}

		if claims, ok := token.Claims.(*claims); ok {
			expectedExp := expectedTime.Add(time.Hour * 72).Unix()
			return claims.Id == expectedID && claims.ExpiresAt == expectedExp
		}
		return false
	}

	tests := []testCase{
		{
			name:        "Successful Token Generation",
			userID:      1,
			currentTime: time.Now(),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test-secret")
			},
			wantErr: false,
		},
		{
			name:        "Missing JWT Secret",
			userID:      1,
			currentTime: time.Now(),
			setupEnv: func() {
				os.Unsetenv("JWT_SECRET")
			},
			wantErr: true,
			errMsg:  "key is required",
		},
		{
			name:        "Zero UserID",
			userID:      0,
			currentTime: time.Now(),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test-secret")
			},
			wantErr: false,
		},
		{
			name:        "Maximum UserID",
			userID:      math.MaxUint32,
			currentTime: time.Now(),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test-secret")
			},
			wantErr: false,
		},
		{
			name:        "Zero Time",
			userID:      1,
			currentTime: time.Time{},
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test-secret")
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupEnv()
			jwtSecret = []byte(os.Getenv("JWT_SECRET"))

			token, err := generateToken(tc.userID, tc.currentTime)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.errMsg != "" {
					assert.Contains(t, err.Error(), tc.errMsg)
				}
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.True(t, validateToken(t, token, tc.userID, tc.currentTime))
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}

	t.Run("Concurrent Token Generation", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))

		var wg sync.WaitGroup
		numGoroutines := 10
		tokens := make([]string, numGoroutines)
		errors := make([]error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				token, err := generateToken(uint(index), time.Now())
				tokens[index] = token
				errors[index] = err
			}(i)
		}

		wg.Wait()

		for i, err := range errors {
			assert.NoError(t, err)
			assert.NotEmpty(t, tokens[i])
		}

		tokenMap := make(map[string]bool)
		for _, token := range tokens {
			assert.False(t, tokenMap[token], "Duplicate token found")
			tokenMap[token] = true
		}
	})

	t.Run("Token Uniqueness", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))

		token1, err1 := generateToken(1, time.Now())
		time.Sleep(time.Second)
		token2, err2 := generateToken(1, time.Now())

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, token1, token2, "Tokens should be unique even for same user ID")
	})
}
