package store

import (
	"database/sql"
	"errors"
	"testing"
	"time"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreDeleteComment(t *testing.T) {

	type testCase struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		comment       *model.Comment
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successfully Delete Existing Comment",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			comment: &model.Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Body:      "Test comment",
				UserID:    1,
				ArticleID: 1,
			},
			expectedError: nil,
		},
		{
			name: "Delete Non-existent Comment",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 999,
				},
			},
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 1,
				},
			},
			expectedError: sql.ErrConnDone,
		},
		{
			name: "Delete Comment with Null Fields",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 1,
				},
				Body:      "",
				UserID:    0,
				ArticleID: 0,
			},
			expectedError: nil,
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
			gormDB.LogMode(true)
			defer gormDB.Close()

			store := &ArticleStore{
				db: gormDB,
			}

			tc.setupMock(mock)

			err = store.DeleteComment(tc.comment)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError) || err.Error() == tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
