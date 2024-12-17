package store

import (
	"database/sql"
	"testing"
	"time"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
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
	}

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

	tests := []testCase{
		{
			name: "Scenario 1: Successfully Retrieve Comments",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"body", "user_id", "article_id",
					"author.id", "author.username", "author.email",
				}).
					AddRow(1, time.Now(), time.Now(), nil,
						"Comment 1", 1, 1,
						1, "user1", "user1@example.com").
					AddRow(2, time.Now(), time.Now(), nil,
						"Comment 2", 2, 1,
						2, "user2", "user2@example.com")

				mock.ExpectQuery(`SELECT .+ FROM "comments" .+ WHERE .+article_id = ?`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "Scenario 2: Article with No Comments",
			article: &model.Article{
				Model: gorm.Model{ID: 2},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"body", "user_id", "article_id",
					"author.id", "author.username", "author.email",
				})

				mock.ExpectQuery(`SELECT .+ FROM "comments" .+ WHERE .+article_id = ?`).
					WithArgs(2).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name: "Scenario 3: Database Error",
			article: &model.Article{
				Model: gorm.Model{ID: 3},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .+ FROM "comments" .+ WHERE .+article_id = ?`).
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:    "Scenario 6: Invalid Article Parameter",
			article: nil,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			tc.setupMock(mock)

			comments, err := store.GetComments(tc.article)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, comments, tc.expectedCount)

				if tc.expectedCount > 0 {
					for _, comment := range comments {
						assert.NotEmpty(t, comment.Author)
						assert.Equal(t, tc.article.ID, comment.ArticleID)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %s", err)
			}
		})
	}
}
