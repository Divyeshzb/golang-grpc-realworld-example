package model

import (
	"testing"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)

func TestUserProtoUser(t *testing.T) {
	// Define test cases structure
	type testCase struct {
		name     string
		user     *User
		token    string
		expected *pb.User
	}

	// Define test cases
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
		},
		{
			name: "Scenario 4: Maximum Field Length Values",
			user: &User{
				Email:    "verylongemail@verylongdomain.com",
				Username: "verylongusername123456789",
				Bio:      "Very long bio text that contains multiple sentences and paragraphs...",
				Image:    "https://very-long-domain-name.com/very/long/path/to/image.jpg",
			},
			token: "very.long.jwt.token.with.multiple.segments",
			expected: &pb.User{
				Email:    "verylongemail@verylongdomain.com",
				Username: "verylongusername123456789",
				Bio:      "Very long bio text that contains multiple sentences and paragraphs...",
				Image:    "https://very-long-domain-name.com/very/long/path/to/image.jpg",
				Token:    "very.long.jwt.token.with.multiple.segments",
			},
		},
		{
			name: "Scenario 5: Special Characters in User Data",
			user: &User{
				Email:    "special@example.com",
				Username: "user©™",
				Bio:      "Bio with émojis 🎉 and ünicode",
				Image:    "http://example.com/image-#special$.jpg",
			},
			token: "normal.token",
			expected: &pb.User{
				Email:    "special@example.com",
				Username: "user©™",
				Bio:      "Bio with émojis 🎉 and ünicode",
				Image:    "http://example.com/image-#special$.jpg",
				Token:    "normal.token",
			},
		},
	}

	// Execute test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Testing:", tc.name)

			// Execute the function
			result := tc.user.ProtoUser(tc.token)

			// Assert the results
			assert.NotNil(t, result, "ProtoUser should not return nil")
			assert.Equal(t, tc.expected.Email, result.Email, "Email mismatch")
			assert.Equal(t, tc.expected.Username, result.Username, "Username mismatch")
			assert.Equal(t, tc.expected.Bio, result.Bio, "Bio mismatch")
			assert.Equal(t, tc.expected.Image, result.Image, "Image mismatch")
			assert.Equal(t, tc.expected.Token, result.Token, "Token mismatch")

			t.Log("Test passed successfully")
		})
	}

	// Test nil user handling
	t.Run("Scenario 6: Nil User Object Handling", func(t *testing.T) {
		t.Log("Testing nil user handling")
		
		var nilUser *User
		// Using assert.Panics to check if the function panics with nil user
		assert.Panics(t, func() {
			nilUser.ProtoUser("token")
		}, "Expected panic when calling ProtoUser on nil user")
		
		t.Log("Nil user test passed successfully")
	})

	// Test user with associated data
	t.Run("Scenario 7: User with Associated Data", func(t *testing.T) {
		t.Log("Testing user with associated data")
		
		user := &User{
			Email:    "test@example.com",
			Username: "testuser",
			Bio:      "Test bio",
			Image:    "http://example.com/image.jpg",
			Follows:  []User{{Username: "follower1"}, {Username: "follower2"}},
			FavoriteArticles: []Article{{Title: "Article 1"}, {Title: "Article 2"}},
		}
		
		result := user.ProtoUser("token")
		
		// Verify only main fields are included
		assert.Equal(t, user.Email, result.Email)
		assert.Equal(t, user.Username, result.Username)
		assert.Equal(t, user.Bio, result.Bio)
		assert.Equal(t, user.Image, result.Image)
		assert.Equal(t, "token", result.Token)
		
		t.Log("Associated data test passed successfully")
	})
}
