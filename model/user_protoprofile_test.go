package model

import (
	"testing"
	pb "github.com/raahii/golang-grpc-realworld-example/proto" // TODO: Ensure this import path matches your project structure
	"github.com/stretchr/testify/assert"
)

func TestUserProtoProfile(t *testing.T) {
	// Define test cases structure
	type testCase struct {
		name      string
		user      User
		following bool
		expected  *pb.Profile
	}

	// Define test cases
	tests := []testCase{
		{
			name: "Basic Profile Conversion - Following True",
			user: User{
				Username: "testuser1",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			following: true,
			expected: &pb.Profile{
				Username:  "testuser1",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
				Following: true,
			},
		},
		{
			name: "Basic Profile Conversion - Following False",
			user: User{
				Username: "testuser2",
				Bio:      "Another bio",
				Image:    "https://example.com/image2.jpg",
			},
			following: false,
			expected: &pb.Profile{
				Username:  "testuser2",
				Bio:      "Another bio",
				Image:    "https://example.com/image2.jpg",
				Following: false,
			},
		},
		{
			name: "Empty Optional Fields",
			user: User{
				Username: "testuser3",
				Bio:      "",
				Image:    "",
			},
			following: false,
			expected: &pb.Profile{
				Username:  "testuser3",
				Bio:      "",
				Image:    "",
				Following: false,
			},
		},
		{
			name: "Special Characters in Fields",
			user: User{
				Username: "test_user_4©",
				Bio:      "Bio with émojis 🎉",
				Image:    "https://example.com/image_special_№.jpg",
			},
			following: true,
			expected: &pb.Profile{
				Username:  "test_user_4©",
				Bio:      "Bio with émojis 🎉",
				Image:    "https://example.com/image_special_№.jpg",
				Following: true,
			},
		},
		{
			name: "Whitespace-Only Content",
			user: User{
				Username: "testuser5",
				Bio:      "   ",
				Image:    "  ",
			},
			following: true,
			expected: &pb.Profile{
				Username:  "testuser5",
				Bio:      "   ",
				Image:    "  ",
				Following: true,
			},
		},
		{
			name: "Zero Value User Struct",
			user: User{},
			following: false,
			expected: &pb.Profile{
				Username:  "",
				Bio:      "",
				Image:    "",
				Following: false,
			},
		},
	}

	// Execute test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Log test case execution
			t.Logf("Executing test case: %s", tc.name)

			// Execute the method
			result := tc.user.ProtoProfile(tc.following)

			// Assert results
			assert.Equal(t, tc.expected.Username, result.Username, "Username mismatch")
			assert.Equal(t, tc.expected.Bio, result.Bio, "Bio mismatch")
			assert.Equal(t, tc.expected.Image, result.Image, "Image mismatch")
			assert.Equal(t, tc.expected.Following, result.Following, "Following status mismatch")

			// Log success
			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}
}
