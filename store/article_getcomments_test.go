package store

import (
	"database/sql"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreGetComments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name          string
		article       *model.Article
		mockSetup     func(sqlmock.Sqlmock)
		expectedCount int
		expectError   bool
	}{
		{
			name: "Successfully retrieve comments",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"body", "article_id", "user_id",
					"author.id", "author.username", "author.email",
				}).
					AddRow(1, time.Now(), time.Now(), nil,
						"Comment 1", 1, 1,
						1, "user1", "user1@example.com").
					AddRow(2, time.Now(), time.Now(), nil,
						"Comment 2", 1, 2,
						2, "user2", "user2@example.com")

				mock.ExpectQuery(`SELECT .+ FROM "comments" LEFT JOIN "users" ON .+ WHERE \("article_id" = \?\)`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "Article with no comments",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
				Title: "Empty Article",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"body", "article_id", "user_id",
					"author.id", "author.username", "author.email",
				})

				mock.ExpectQuery(`SELECT .+ FROM "comments" LEFT JOIN "users" ON .+ WHERE \("article_id" = \?\)`).
					WithArgs(2).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name: "Database error",
			article: &model.Article{
				Model: gorm.Model{ID: 3},
				Title: "Error Article",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .+ FROM "comments" LEFT JOIN "users" ON .+ WHERE \("article_id" = \?\)`).
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectError:   true,
		},
		{
			name: "Invalid article ID",
			article: &model.Article{
				Model: gorm.Model{ID: 0},
				Title: "Invalid Article",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .+ FROM "comments" LEFT JOIN "users" ON .+ WHERE \("article_id" = \?\)`).
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			comments, err := store.GetComments(tt.article)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(comments))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
