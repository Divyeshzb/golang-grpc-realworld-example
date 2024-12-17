package store

import (
	"errors"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreGetFeedArticles(t *testing.T) {

	type testCase struct {
		name      string
		userIDs   []uint
		limit     int64
		offset    int64
		mockSetup func(sqlmock.Sqlmock)
		expected  struct {
			articles []model.Article
			err      error
		}
	}

	tests := []testCase{
		{
			name:    "Successfully retrieve feed articles for multiple users",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",

					"author.id", "author.created_at", "author.updated_at", "author.deleted_at",
					"author.username", "author.email", "author.password", "author.bio", "author.image",
				}).AddRow(
					1, time.Now(), time.Now(), nil,
					"Test Title", "Test Description", "Test Body", 1, 0,
					1, time.Now(), time.Now(), nil,
					"testuser", "test@test.com", "password", "bio", "image",
				)

				mock.ExpectQuery(`SELECT .+ FROM "articles"`).
					WithArgs(1, 2).
					WillReturnRows(rows)
			},
			expected: struct {
				articles []model.Article
				err      error
			}{
				articles: []model.Article{{
					Model:       gorm.Model{ID: 1},
					Title:       "Test Title",
					Description: "Test Description",
					Body:        "Test Body",
					UserID:      1,
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
						Email:    "test@test.com",
						Password: "password",
						Bio:      "bio",
						Image:    "image",
					},
				}},
				err: nil,
			},
		},
		{
			name:    "Empty result with valid but non-existent user IDs",
			userIDs: []uint{999},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})

				mock.ExpectQuery(`SELECT .+ FROM "articles"`).
					WithArgs(999).
					WillReturnRows(rows)
			},
			expected: struct {
				articles []model.Article
				err      error
			}{
				articles: []model.Article{},
				err:      nil,
			},
		},
		{
			name:    "Database connection error",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .+ FROM "articles"`).
					WithArgs(1).
					WillReturnError(errors.New("database connection error"))
			},
			expected: struct {
				articles []model.Article
				err      error
			}{
				articles: nil,
				err:      errors.New("database connection error"),
			},
		},
		{
			name:    "Empty userIDs slice handling",
			userIDs: []uint{},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})

				mock.ExpectQuery(`SELECT .+ FROM "articles"`).
					WillReturnRows(rows)
			},
			expected: struct {
				articles []model.Article
				err      error
			}{
				articles: []model.Article{},
				err:      nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.mockSetup(mock)

			store := &ArticleStore{db: gormDB}

			articles, err := store.GetFeedArticles(tc.userIDs, tc.limit, tc.offset)

			if tc.expected.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expected.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, len(tc.expected.articles), len(articles))
			if len(tc.expected.articles) > 0 {
				assert.Equal(t, tc.expected.articles[0].Title, articles[0].Title)
				assert.Equal(t, tc.expected.articles[0].UserID, articles[0].UserID)
				assert.Equal(t, tc.expected.articles[0].Author.Username, articles[0].Author.Username)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
