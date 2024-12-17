package model

import (
	"testing"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)

func TestUserProtoUser(t *testing.T) {
	// Test cases structure
	type testCase struct {
		name     string
		user     *User
		token    string
		expected *pb.User
		wantErr  bool
	}

	// Test cases
	tests := []testCase{
		{
			name: "Scenario 1: Basic User Data Conversion with Valid Token",
			user: &User{
				Email:    "test@example.com",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
			},
			token: "valid.jwt.token",
			expected: &pb.User{
				Email:    "test@example.com",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
				Token:    "valid.jwt.token",
			},
			wantErr: false,
		},
		{
			name: "Scenario 2: Empty Token Handling",
			user: &User{
				Email:    "test@example.com",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
			},
			token: "",
			expected: &pb.User{
				Email:    "test@example.com",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
				Token:    "",
			},
			wantErr: false,
		},
		{
			name: "Scenario 3: User with Empty Optional Fields",
			user: &User{
				Email:    "test@example.com",
				Username: "testuser",
				Bio:      "",
				Image:    "",
			},
			token: "valid.token",
			expected: &pb.User{
				Email:    "test@example.com",
				Username: "testuser",
				Bio:      "",
				Image:    "",
				Token:    "valid.token",
			},
			wantErr: false,
		},
		{
			name: "Scenario 4: Maximum Field Length Values",
			user: &User{
				Email:    "verylongemail@verylongdomain.com",
				Username: "verylongusername123456789",
				Bio:      "Very long bio text that contains multiple sentences and paragraphs...",
				Image:    "https://very-long-domain.com/very-long-image-path/image.jpg",
			},
			token: "very.long.jwt.token.with.multiple.segments",
			expected: &pb.User{
				Email:    "verylongemail@verylongdomain.com",
				Username: "verylongusername123456789",
				Bio:      "Very long bio text that contains multiple sentences and paragraphs...",
				Image:    "https://very-long-domain.com/very-long-image-path/image.jpg",
				Token:    "very.long.jwt.token.with.multiple.segments",
			},
			wantErr: false,
		},
		{
			name: "Scenario 5: Special Characters in User Data",
			user: &User{
				Email:    "special@例子.com",
				Username: "user名_🎉",
				Bio:      "Bio with émojis 🎈 and спеціальні characters",
				Image:    "http://example.com/イメージ.jpg",
			},
			token: "valid.token",
			expected: &pb.User{
				Email:    "special@例子.com",
				Username: "user名_🎉",
				Bio:      "Bio with émojis 🎈 and спеціальні characters",
				Image:    "http://example.com/イメージ.jpg",
				Token:    "valid.token",
			},
			wantErr: false,
		},
		{
			name:     "Scenario 6: Nil User Object Handling",
			user:     nil,
			token:    "valid.token",
			expected: nil,
			wantErr:  true,
		},
		{
			name: "Scenario 7: User with Associated Data",
			user: &User{
				Email:    "test@example.com",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
				Follows: []User{
					{Username: "follower1"},
					{Username: "follower2"},
				},
				FavoriteArticles: []Article{
					{Title: "Article 1"},
					{Title: "Article 2"},
				},
			},
			token: "valid.token",
			expected: &pb.User{
				Email:    "test@example.com",
				Username: "testuser",
				Bio:      "Test bio",
				Image:    "http://example.com/image.jpg",
				Token:    "valid.token",
			},
			wantErr: false,
		},
	}

	// Execute test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Testing:", tc.name)

			// Handle nil user case
			if tc.user == nil {
				if !tc.wantErr {
					t.Error("Expected error for nil user but wantErr is false")
				}
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic for nil user but got none")
					}
				}()
			}

			// Call the function under test
			result := tc.user.ProtoUser(tc.token)

			// Assertions
			if tc.wantErr {
				assert.Nil(t, result, "Expected nil result for error case")
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.Equal(t, tc.expected.Email, result.Email, "Email mismatch")
				assert.Equal(t, tc.expected.Username, result.Username, "Username mismatch")
				assert.Equal(t, tc.expected.Bio, result.Bio, "Bio mismatch")
				assert.Equal(t, tc.expected.Image, result.Image, "Image mismatch")
				assert.Equal(t, tc.expected.Token, result.Token, "Token mismatch")
				t.Log("Test passed successfully")
			}
		})
	}
}
