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
		input         *model.Comment
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
			input: &model.Comment{
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
			input: &model.Comment{
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
			input: &model.Comment{
				Model: gorm.Model{
					ID: 1,
				},
			},
			expectedError: sql.ErrConnDone,
		},
		{
			name: "Nil Comment Parameter",
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			input:         nil,
			expectedError: errors.New("invalid comment: nil pointer"),
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

			store := &ArticleStore{
				db: gormDB,
			}

			if tc.setupMock != nil {
				tc.setupMock(mock)
			}

			err = store.DeleteComment(tc.input)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test case '%s' completed", tc.name)
			if err != nil {
				t.Logf("Error: %v", err)
			}
		})
	}
}
