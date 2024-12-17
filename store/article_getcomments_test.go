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

func TestArticleStoreGetComments(t *testing.T) {

	type testCase struct {
		name          string
		article       *model.Article
		setupMock     func(sqlmock.Sqlmock)
		expectedCount int
		expectError   bool
		errorMessage  string
	}

	tests := []testCase{
		{
			name: "Successfully Retrieve Comments",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Comment 1", 1, 1).
					AddRow(2, time.Now(), time.Now(), nil, "Comment 2", 1, 1)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email", "bio", "image"}).
					AddRow(1, "testuser", "test@example.com", "bio", "image.jpg")

				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE article_id = \? AND "comments"\."deleted_at" IS NULL`).
					WithArgs(1).
					WillReturnRows(rows)

				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WillReturnRows(authorRows)
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "Article With No Comments",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id"})
				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE article_id = \? AND "comments"\."deleted_at" IS NULL`).
					WithArgs(2).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name: "Database Error",
			article: &model.Article{
				Model: gorm.Model{ID: 3},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE article_id = \? AND "comments"\."deleted_at" IS NULL`).
					WithArgs(3).
					WillReturnError(errors.New("database error"))
			},
			expectedCount: 0,
			expectError:   true,
			errorMessage:  "database error",
		},
		{
			name:          "Nil Article Parameter",
			article:       nil,
			setupMock:     func(mock sqlmock.Sqlmock) {},
			expectedCount: 0,
			expectError:   true,
			errorMessage:  "invalid article parameter",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to open gorm connection: %v", err)
			}
			defer gormDB.Close()

			store := &ArticleStore{
				db: gormDB,
			}

			if tc.article != nil {
				tc.setupMock(mock)
			}

			comments, err := store.GetComments(tc.article)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMessage != "" {
					assert.Contains(t, err.Error(), tc.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, comments, tc.expectedCount)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %s", err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}
}
