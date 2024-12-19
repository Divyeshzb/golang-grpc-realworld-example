package model

import (
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

// TestCommentValidate tests the Validate method of the Comment struct
func TestCommentValidate(t *testing.T) {
	// Define test cases
	tests := []struct {
		name    string
		comment Comment
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid Comment with Non-Empty Body",
			comment: Comment{
				Body:      "This is a valid comment",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "Empty Comment Body",
			comment: Comment{
				Body:      "",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: true,
			errMsg:  "body: cannot be blank",
		},
		{
			name: "Comment Body with Only Whitespace",
			comment: Comment{
				Body:      "    \t\n",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: true,
			errMsg:  "body: cannot be blank",
		},
		{
			name: "Comment Body with Maximum Length",
			comment: Comment{
				Body:      strings.Repeat("a", 1000),
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "Comment with Zero Values for Required Fields",
			comment: Comment{
				Body: "Valid body but zero values for other fields",
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "Comment with Special Characters in Body",
			comment: Comment{
				Body:      "Special chars: !@#$%^&*()_+-=[]{}|;:'\",.<>?/",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "Comment with Unicode Characters",
			comment: Comment{
				Body:      "Unicode chars: 你好世界 こんにちは世界 안녕하세요 세계",
				UserID:    1,
				ArticleID: 1,
			},
			wantErr: false,
			errMsg:  "",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing scenario: %s", tt.name)

			err := tt.comment.Validate()

			// Check if error matches expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("Comment.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If expecting an error, verify error message
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error message '%s', but got nil", tt.errMsg)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Expected error message '%s', but got '%s'", tt.errMsg, err.Error())
					return
				}
			}

			t.Logf("Test case passed successfully")
		})
	}
}
