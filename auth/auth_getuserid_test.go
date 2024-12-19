package auth

import (
	"context"
	"os"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestGetUserID(t *testing.T) {

	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	testSecret := "test-secret-key"
	os.Setenv("JWT_SECRET", testSecret)

	createToken := func(userID uint, exp time.Time) string {
		claims := jwt.MapClaims{
			"user_id": userID,
			"exp":     exp.Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(testSecret))
		return tokenString
	}

	createContextWithToken := func(token string) context.Context {
		md := metadata.New(map[string]string{
			"authorization": "Token " + token,
		})
		return metadata.NewIncomingContext(context.Background(), md)
	}

	tests := []struct {
		name          string
		setupContext  func() context.Context
		expectedID    uint
		expectedError string
	}{
		{
			name: "Valid Token",
			setupContext: func() context.Context {
				token := createToken(123, time.Now().Add(time.Hour))
				return createContextWithToken(token)
			},
			expectedID:    123,
			expectedError: "",
		},
		{
			name: "Missing Token",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedID:    0,
			expectedError: "Request unauthenticated with Token",
		},
		{
			name: "Expired Token",
			setupContext: func() context.Context {
				token := createToken(123, time.Now().Add(-time.Hour))
				return createContextWithToken(token)
			},
			expectedID:    0,
			expectedError: "token expired",
		},
		{
			name: "Malformed Token",
			setupContext: func() context.Context {
				return createContextWithToken("malformed.token.string")
			},
			expectedID:    0,
			expectedError: "invalid token",
		},
		{
			name: "Invalid Token Signature",
			setupContext: func() context.Context {

				claims := jwt.MapClaims{
					"user_id": 123,
					"exp":     time.Now().Add(time.Hour).Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte("wrong-secret"))
				return createContextWithToken(tokenString)
			},
			expectedID:    0,
			expectedError: "invalid token",
		},
		{
			name: "Future Token",
			setupContext: func() context.Context {
				token := createToken(123, time.Now().Add(24*time.Hour))
				return createContextWithToken(token)
			},
			expectedID:    123,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := tt.setupContext()

			userID, err := GetUserID(ctx)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Equal(t, uint(0), userID)
				t.Logf("Expected error received: %v", err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
				t.Logf("Successfully retrieved user ID: %d", userID)
			}
		})
	}
}
