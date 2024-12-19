package model

import (
	"testing"
	"time"
	"github.com/jinzhu/gorm"  // Add this import to resolve the gorm.Model undefined error
)

// TestUserValidate implements table-driven tests for User.Validate()
func TestUserValidate(t *testing.T) {
	// Define test cases
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid User Data",
			user: User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "validuser123",
				Email:    "valid@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "Missing Username",
			user: User{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "",
				Email:    "valid@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: true,
			errMsg:  "Username: cannot be blank.",
		},
		{
			name: "Invalid Username Format",
			user: User{
				Model: gorm.Model{
					ID:        3,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "invalid@user",
				Email:    "valid@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: true,
			errMsg:  "Username: must be in a valid format.",
		},
		{
			name: "Missing Email",
			user: User{
				Model: gorm.Model{
					ID:        4,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "validuser123",
				Email:    "",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: true,
			errMsg:  "Email: cannot be blank.",
		},
		{
			name: "Invalid Email Format",
			user: User{
				Model: gorm.Model{
					ID:        5,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "validuser123",
				Email:    "invalidemail@",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: true,
			errMsg:  "Email: must be a valid email address.",
		},
		{
			name: "Missing Password",
			user: User{
				Model: gorm.Model{
					ID:        6,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "validuser123",
				Email:    "valid@example.com",
				Password: "",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: true,
			errMsg:  "Password: cannot be blank.",
		},
		{
			name: "Multiple Validation Errors",
			user: User{
				Model: gorm.Model{
					ID:        7,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "",
				Email:    "invalidemail@",
				Password: "",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: true,
			errMsg:  "multiple validation errors",
		},
		{
			name: "Whitespace-Only Values",
			user: User{
				Model: gorm.Model{
					ID:        8,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "   ",
				Email:    "   ",
				Password: "   ",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: true,
			errMsg:  "multiple validation errors",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing scenario: %s", tt.name)
			
			err := tt.user.Validate()
			
			// Check if error was expected
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If error was expected, verify error message contains expected content
			if tt.wantErr && err != nil {
				if tt.errMsg != "multiple validation errors" && err.Error() != tt.errMsg {
					t.Errorf("User.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
				t.Logf("Validation failed as expected with error: %v", err)
			} else if !tt.wantErr {
				t.Log("Validation passed as expected")
			}
		})
	}
}
