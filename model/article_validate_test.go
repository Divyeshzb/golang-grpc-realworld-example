package model

import (
	"strings"
	"testing"
)

func TestArticleValidate(t *testing.T) {

	tests := []struct {
		name        string
		article     Article
		expectError bool
		errorFields []string
	}{
		{
			name: "Valid Article with All Required Fields",
			article: Article{
				Title: "Test Article",
				Body:  "This is a test article body",
				Tags: []Tag{
					{Name: "test"},
				},
			},
			expectError: false,
		},
		{
			name: "Missing Title Field",
			article: Article{
				Title: "",
				Body:  "This is a test article body",
				Tags: []Tag{
					{Name: "test"},
				},
			},
			expectError: true,
			errorFields: []string{"Title"},
		},
		{
			name: "Missing Body Field",
			article: Article{
				Title: "Test Article",
				Body:  "",
				Tags: []Tag{
					{Name: "test"},
				},
			},
			expectError: true,
			errorFields: []string{"Body"},
		},
		{
			name: "Empty Tags Array",
			article: Article{
				Title: "Test Article",
				Body:  "This is a test article body",
				Tags:  []Tag{},
			},
			expectError: true,
			errorFields: []string{"Tags"},
		},
		{
			name: "Nil Tags Array",
			article: Article{
				Title: "Test Article",
				Body:  "This is a test article body",
				Tags:  nil,
			},
			expectError: true,
			errorFields: []string{"Tags"},
		},
		{
			name: "Multiple Validation Errors",
			article: Article{
				Title: "",
				Body:  "",
				Tags:  nil,
			},
			expectError: true,
			errorFields: []string{"Title", "Body", "Tags"},
		},
		{
			name: "Whitespace Only in Required Fields",
			article: Article{
				Title: "   ",
				Body:  "    ",
				Tags: []Tag{
					{Name: "test"},
				},
			},
			expectError: true,
			errorFields: []string{"Title", "Body"},
		},
		{
			name: "Valid Article with Maximum Field Values",
			article: Article{
				Title: strings.Repeat("a", 1000),
				Body:  strings.Repeat("b", 10000),
				Tags: []Tag{
					{Name: "tag1"},
					{Name: "tag2"},
					{Name: "tag3"},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.article.Validate()

			t.Logf("Test case: %s", tt.name)
			t.Logf("Article: %+v", tt.article)
			t.Logf("Expected error: %v", tt.expectError)
			if tt.expectError {
				t.Logf("Expected error fields: %v", tt.errorFields)
			}

			if tt.expectError && err == nil {
				t.Errorf("Expected validation error but got nil")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no validation error but got: %v", err)
				return
			}

			if tt.expectError {
				for _, field := range tt.errorFields {
					if !strings.Contains(err.Error(), field) {
						t.Errorf("Expected error for field %s, but it was not found in error: %v", field, err)
					}
				}
			}
		})
	}
}
