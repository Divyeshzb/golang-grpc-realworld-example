package model

import (
	"strings"
	"testing"
)

func TestArticleOverwrite(t *testing.T) {

	tests := []struct {
		name           string
		initialArticle Article
		inputTitle     string
		inputDesc      string
		inputBody      string
		expected       Article
	}{
		{
			name: "Update All Fields with Valid Values",
			initialArticle: Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			},
			inputTitle: "New Title",
			inputDesc:  "New Description",
			inputBody:  "New Body",
			expected: Article{
				Title:       "New Title",
				Description: "New Description",
				Body:        "New Body",
			},
		},
		{
			name: "Update No Fields with Empty Strings",
			initialArticle: Article{
				Title:       "Keep Title",
				Description: "Keep Description",
				Body:        "Keep Body",
			},
			inputTitle: "",
			inputDesc:  "",
			inputBody:  "",
			expected: Article{
				Title:       "Keep Title",
				Description: "Keep Description",
				Body:        "Keep Body",
			},
		},
		{
			name: "Partial Update - Only Title",
			initialArticle: Article{
				Title:       "Old Title",
				Description: "Keep Description",
				Body:        "Keep Body",
			},
			inputTitle: "Updated Title",
			inputDesc:  "",
			inputBody:  "",
			expected: Article{
				Title:       "Updated Title",
				Description: "Keep Description",
				Body:        "Keep Body",
			},
		},
		{
			name: "Partial Update - Only Description",
			initialArticle: Article{
				Title:       "Keep Title",
				Description: "Old Description",
				Body:        "Keep Body",
			},
			inputTitle: "",
			inputDesc:  "Updated Description",
			inputBody:  "",
			expected: Article{
				Title:       "Keep Title",
				Description: "Updated Description",
				Body:        "Keep Body",
			},
		},
		{
			name: "Partial Update - Only Body",
			initialArticle: Article{
				Title:       "Keep Title",
				Description: "Keep Description",
				Body:        "Old Body",
			},
			inputTitle: "",
			inputDesc:  "",
			inputBody:  "Updated Body",
			expected: Article{
				Title:       "Keep Title",
				Description: "Keep Description",
				Body:        "Updated Body",
			},
		},
		{
			name: "Update with Special Characters",
			initialArticle: Article{
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
			},
			inputTitle: "Title with 特殊字符 and ñ",
			inputDesc:  "Description with ♠♣♥♦ symbols",
			inputBody:  "Body with émojis 🎉🎊",
			expected: Article{
				Title:       "Title with 特殊字符 and ñ",
				Description: "Description with ♠♣♥♦ symbols",
				Body:        "Body with émojis 🎉🎊",
			},
		},
		{
			name: "Update with Maximum Length Values",
			initialArticle: Article{
				Title:       "Short Title",
				Description: "Short Description",
				Body:        "Short Body",
			},
			inputTitle: strings.Repeat("Long Title ", 100),
			inputDesc:  strings.Repeat("Long Description ", 100),
			inputBody:  strings.Repeat("Long Body Content ", 1000),
			expected: Article{
				Title:       strings.Repeat("Long Title ", 100),
				Description: strings.Repeat("Long Description ", 100),
				Body:        strings.Repeat("Long Body Content ", 1000),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			article := tt.initialArticle

			t.Logf("Testing scenario: %s", tt.name)
			t.Logf("Initial state - Title: %s, Description: %s, Body: %s",
				article.Title, article.Description, article.Body)

			article.Overwrite(tt.inputTitle, tt.inputDesc, tt.inputBody)

			if article.Title != tt.expected.Title {
				t.Errorf("Title mismatch - got: %v, want: %v", article.Title, tt.expected.Title)
			}
			if article.Description != tt.expected.Description {
				t.Errorf("Description mismatch - got: %v, want: %v", article.Description, tt.expected.Description)
			}
			if article.Body != tt.expected.Body {
				t.Errorf("Body mismatch - got: %v, want: %v", article.Body, tt.expected.Body)
			}

			t.Logf("Test completed successfully - Final state - Title: %s, Description: %s, Body: %s",
				article.Title, article.Description, article.Body)
		})
	}
}
