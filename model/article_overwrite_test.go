package model

import (
	"strings"
	"testing"
)

func TestArticleOverwrite(t *testing.T) {

	type testCase struct {
		name           string
		initialArticle Article
		inputTitle     string
		inputDesc      string
		inputBody      string
		expected       Article
	}

	tests := []testCase{
		{
			name: "Scenario 1: Update All Fields with Valid Values",
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
			name: "Scenario 2: Update No Fields with Empty Strings",
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
			name: "Scenario 3: Partial Update - Only Title",
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
			name: "Scenario 4: Partial Update - Only Description",
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
			name: "Scenario 5: Partial Update - Only Body",
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
			name: "Scenario 6: Update with Special Characters",
			initialArticle: Article{
				Title:       "Original Title",
				Description: "Original Description",
				Body:        "Original Body",
			},
			inputTitle: "Title with 特殊文字 and ñ",
			inputDesc:  "Description with €$¥ symbols",
			inputBody:  "Body with 🚀 emoji and \n newline",
			expected: Article{
				Title:       "Title with 特殊文字 and ñ",
				Description: "Description with €$¥ symbols",
				Body:        "Body with 🚀 emoji and \n newline",
			},
		},
		{
			name: "Scenario 7: Update with Maximum Length Values",
			initialArticle: Article{
				Title:       "Short Title",
				Description: "Short Description",
				Body:        "Short Body",
			},
			inputTitle: strings.Repeat("Long Title ", 100),
			inputDesc:  strings.Repeat("Long Description ", 100),
			inputBody:  strings.Repeat("Long Body Content ", 100),
			expected: Article{
				Title:       strings.Repeat("Long Title ", 100),
				Description: strings.Repeat("Long Description ", 100),
				Body:        strings.Repeat("Long Body Content ", 100),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			article := tc.initialArticle

			article.Overwrite(tc.inputTitle, tc.inputDesc, tc.inputBody)

			if article.Title != tc.expected.Title {
				t.Errorf("Title mismatch\nexpected: %v\ngot: %v", tc.expected.Title, article.Title)
			}
			if article.Description != tc.expected.Description {
				t.Errorf("Description mismatch\nexpected: %v\ngot: %v", tc.expected.Description, article.Description)
			}
			if article.Body != tc.expected.Body {
				t.Errorf("Body mismatch\nexpected: %v\ngot: %v", tc.expected.Body, article.Body)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}
}
