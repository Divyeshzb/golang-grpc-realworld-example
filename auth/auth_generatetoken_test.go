package auth

import (
	"github.com/dgrijalva/jwt-go"
	"math"
	"os"
	"sync"
	"testing"
	"time"
)

func TestgenerateToken(t *testing.T) {

	type testCase struct {
		name        string
		userID      uint
		currentTime time.Time
		setupEnv    func()
		wantErr     bool
		validate    func(t *testing.T, token string, err error)
	}

	parseToken := func(tokenString string) (*jwt.Token, error) {
		return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
	}

	tests := []testCase{
		{
			name:        "Successful Token Generation",
			userID:      1,
			currentTime: time.Now(),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test-secret")
				jwtSecret = []byte(os.Getenv("JWT_SECRET"))
			},
			wantErr: false,
			validate: func(t *testing.T, token string, err error) {
				if token == "" {
					t.Error("expected non-empty token")
				}
				parsedToken, err := parseToken(token)
				if err != nil || !parsedToken.Valid {
					t.Errorf("generated invalid token: %v", err)
				}
			},
		},
		{
			name:        "Token Expiration Validation",
			userID:      1,
			currentTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test-secret")
				jwtSecret = []byte(os.Getenv("JWT_SECRET"))
			},
			wantErr: false,
			validate: func(t *testing.T, token string, err error) {
				parsedToken, _ := parseToken(token)
				claims := parsedToken.Claims.(jwt.MapClaims)
				expectedExp := time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC).Unix()
				if int64(claims["exp"].(float64)) != expectedExp {
					t.Errorf("incorrect expiration time, got: %v, want: %v", claims["exp"], expectedExp)
				}
			},
		},
		{
			name:        "Missing JWT Secret",
			userID:      1,
			currentTime: time.Now(),
			setupEnv: func() {
				os.Unsetenv("JWT_SECRET")
				jwtSecret = []byte{}
			},
			wantErr: true,
			validate: func(t *testing.T, token string, err error) {
				if err == nil {
					t.Error("expected error with missing JWT secret")
				}
			},
		},
		{
			name:        "Zero User ID",
			userID:      0,
			currentTime: time.Now(),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test-secret")
				jwtSecret = []byte(os.Getenv("JWT_SECRET"))
			},
			wantErr: false,
			validate: func(t *testing.T, token string, err error) {
				parsedToken, _ := parseToken(token)
				claims := parsedToken.Claims.(jwt.MapClaims)
				if claims["id"] != float64(0) {
					t.Errorf("incorrect user ID in claims, got: %v, want: 0", claims["id"])
				}
			},
		},
		{
			name:        "Maximum uint Value",
			userID:      math.MaxUint32,
			currentTime: time.Now(),
			setupEnv: func() {
				os.Setenv("JWT_SECRET", "test-secret")
				jwtSecret = []byte(os.Getenv("JWT_SECRET"))
			},
			wantErr: false,
			validate: func(t *testing.T, token string, err error) {
				parsedToken, _ := parseToken(token)
				claims := parsedToken.Claims.(jwt.MapClaims)
				if claims["id"] != float64(math.MaxUint32) {
					t.Errorf("incorrect user ID in claims, got: %v, want: %v", claims["id"], math.MaxUint32)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			tc.setupEnv()

			token, err := generateToken(tc.userID, tc.currentTime)

			if (err != nil) != tc.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			tc.validate(t, token, err)
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
			if err != nil {
				t.Errorf("goroutine %d failed to generate token: %v", i, err)
			}
			if tokens[i] == "" {
				t.Errorf("goroutine %d generated empty token", i)
			}
		}
	})

	t.Run("Token Uniqueness", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))

		token1, err1 := generateToken(1, time.Now())
		token2, err2 := generateToken(1, time.Now().Add(time.Second))

		if err1 != nil || err2 != nil {
			t.Errorf("failed to generate tokens: %v, %v", err1, err2)
		}

		if token1 == token2 {
			t.Error("tokens should be unique even for same user ID")
		}
	})
}
