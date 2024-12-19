package model

import (
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)

func TestArticleProtoArticle(t *testing.T) {
	tests := []struct {
		name      string
		article   Article
		favorited bool
		want      func(*testing.T, *pb.Article)
	}{
		{
			name: "Scenario 1: Basic Article Conversion with Minimal Data",
			article: Article{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Title:       "Test Title",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
				Tags:        []Tag{},
			},
			favorited: false,
			want: func(t *testing.T, got *pb.Article) {
				assert.Equal(t, "1", got.Slug)
				assert.Equal(t, "Test Title", got.Title)
				assert.Equal(t, "Test Description", got.Description)
				assert.Equal(t, "Test Body", got.Body)
				assert.Empty(t, got.TagList)
				assert.False(t, got.Favorited)
			},
		},
		{
			name: "Scenario 2: Article Conversion with Complete Data Set",
			article: Article{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Title:          "Full Article",
				Description:    "Complete Description",
				Body:           "Full Body Content",
				UserID:         1,
				FavoritesCount: 10,
				Tags: []Tag{
					{Model: gorm.Model{}, Name: "tag1"},
					{Model: gorm.Model{}, Name: "tag2"},
				},
			},
			favorited: true,
			want: func(t *testing.T, got *pb.Article) {
				assert.Equal(t, "2", got.Slug)
				assert.Equal(t, int32(10), got.FavoritesCount)
				assert.True(t, got.Favorited)
				assert.Len(t, got.TagList, 2)
				assert.Contains(t, got.TagList, "tag1")
				assert.Contains(t, got.TagList, "tag2")
			},
		},
		{
			name: "Scenario 3: Article Conversion with Zero Tags",
			article: Article{
				Model: gorm.Model{
					ID:        3,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Title:       "No Tags Article",
				Description: "Description",
				Body:        "Body",
				UserID:      1,
				Tags:        nil,
			},
			favorited: false,
			want: func(t *testing.T, got *pb.Article) {
				assert.NotNil(t, got.TagList)
				assert.Empty(t, got.TagList)
			},
		},
		{
			name: "Scenario 4: Time Format Verification",
			article: Article{
				Model: gorm.Model{
					ID:        4,
					CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				},
				Title:       "Time Test",
				Description: "Description",
				Body:        "Body",
				UserID:      1,
			},
			favorited: false,
			want: func(t *testing.T, got *pb.Article) {
				assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}-\d{4}Z`, got.CreatedAt)
				assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}-\d{4}Z`, got.UpdatedAt)
			},
		},
		{
			name: "Scenario 8: Unicode Content Handling",
			article: Article{
				Model: gorm.Model{
					ID:        8,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Title:       "测试标题",
				Description: "描述 🌟",
				Body:        "содержание текста",
				UserID:      1,
			},
			favorited: false,
			want: func(t *testing.T, got *pb.Article) {
				assert.Equal(t, "测试标题", got.Title)
				assert.Equal(t, "描述 🌟", got.Description)
				assert.Equal(t, "содержание текста", got.Body)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Running:", tt.name)

			got := tt.article.ProtoArticle(tt.favorited)

			assert.NotNil(t, got, "ProtoArticle should not return nil")
			tt.want(t, got)

			t.Log("Test completed successfully")
		})
	}
}
