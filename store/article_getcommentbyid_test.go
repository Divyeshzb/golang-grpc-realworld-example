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

func TestArticleStoreGetCommentByID(t *testing.T) {

	type testCase struct {
		name            string
		commentID       uint
		setupMock       func(sqlmock.Sqlmock)
		expectedError   error
		expectedComment *model.Comment
	}

	tests := []testCase{
		{
			name:      "Successfully retrieve existing comment",
			commentID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id"}
				mock.ExpectQuery(`SELECT \* FROM "comments"`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, time.Now(), time.Now(), nil, "Test comment", 1, 1))
			},
			expectedComment: &model.Comment{
				Model:     gorm.Model{ID: 1},
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			expectedError: nil,
		},
		{
			name:      "Non-existent comment",
			commentID: 99999,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "comments"`).
					WithArgs(99999).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectedComment: nil,
			expectedError:   gorm.ErrRecordNotFound,
		},
		{
			name:      "Database connection error",
			commentID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "comments"`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedComment: nil,
			expectedError:   sql.ErrConnDone,
		},
		{
			name:      "Zero ID input",
			commentID: 0,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "comments"`).
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedComment: nil,
			expectedError:   gorm.ErrRecordNotFound,
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

			tc.setupMock(mock)

			store := &ArticleStore{db: gormDB}

			comment, err := store.GetCommentByID(tc.commentID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, comment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
				assert.Equal(t, tc.expectedComment.ID, comment.ID)
				assert.Equal(t, tc.expectedComment.Body, comment.Body)
				assert.Equal(t, tc.expectedComment.UserID, comment.UserID)
				assert.Equal(t, tc.expectedComment.ArticleID, comment.ArticleID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
